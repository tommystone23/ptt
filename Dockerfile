FROM kalilinux/kali-rolling

LABEL maintainer="PTT"
LABEL description="Dev environment for PTT"

ENV DEBIAN_FRONTEND=noninteractive

# Install Base Dependencies
RUN apt-get update && apt-get install -y \
    curl unzip git build-essential ca-certificates

# Install Go
ENV GO_VERSION=1.24.1
RUN curl -LO https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go
ENV PATH="${GOPATH}/bin:${PATH}"

# Install Node.js and pnpm
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g pnpm

# Install Go template, air, gRPC stuff
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install protobuf
ENV PROTOC_ZIP=protoc-23.4-linux-x86_64.zip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v23.4/${PROTOC_ZIP} && \
    unzip -o ${PROTOC_ZIP} -d /usr/local bin/protoc && \
    unzip -o ${PROTOC_ZIP} -d /usr/local 'include/*' && \
    rm -f ${PROTOC_ZIP}

WORKDIR /app
COPY package.json pnpm-lock.yaml ./

# Install JS and CSS Dependencies
RUN pnpm install

# Copy rest of the application
COPY . .

# Build Go binary
RUN make proto tailwind templ

ENTRYPOINT ["air"]

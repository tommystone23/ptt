#!/bin/bash

set -e

# Required versions
GO_VERSION="1.24.1"
PROTOC_VERSION="23.4"
PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-x86_64.zip"

sudo apt update
sudo apt install -y curl unzip git build-essential make ca-certificates

# Install Go 
install_go() {
    GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"

    curl -LO "https://go.dev/dl/${GO_TAR}"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "${GO_TAR}"
    rm "${GO_TAR}"

    # Add Go binary path if not already present
    if ! grep -q '/usr/local/go/bin' ~/.bashrc; then
      echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
    fi
    export PATH=$PATH:/usr/local/go/bin

    echo "Go installed: $(go version)"
}

if command -v go &>/dev/null; then
    INSTALLED_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [[ "$(printf '%s\n' "$GO_VERSION" "$INSTALLED_GO_VERSION" | sort -V | head -n1)" != "$GO_VERSION" ]]; then
        echo "Go version is less than required (${GO_VERSION}). Updating..."
        install_go
    else
        echo "Go version is sufficient: ${INSTALLED_GO_VERSION}"
        export PATH=$PATH:/usr/local/go/bin
    fi
else
    install_go
fi

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt install -y nodejs
echo "Node.js installed: $(node -v)"

# Install pnpm 
sudo npm install -g pnpm
echo "pnpm installed: $(pnpm -v)"

# Install Tailwind CSS CLI & JS dependencies
pnpm install
echo "Tailwind installed: $(npx tailwindcss -v)"

# Install Go template engine
go install github.com/a-h/templ/cmd/templ@latest

# Install air
go install github.com/air-verse/air@latest

# Install protoc
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}
sudo unzip -o ${PROTOC_ZIP} -d /usr/local bin/protoc
sudo unzip -o ${PROTOC_ZIP} -d /usr/local 'include/*'
rm -f ${PROTOC_ZIP}
echo "protoc installed: $(protoc --version)"

# Install protoc Go stuff
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo "protoc-gen-go installed: $(which protoc-gen-go)"
echo "protoc-gen-go-grpc installed: $(which protoc-gen-go-grpc)"

# Final PATH export for current session and persist for future shells
GOPATH_BIN=$(go env GOPATH)/bin
if ! grep -q "$GOPATH_BIN" ~/.bashrc; then
  echo "export PATH=\$PATH:${GOPATH_BIN}" >> ~/.bashrc
fi
export PATH=$PATH:$GOPATH_BIN

echo "Dependencies installed successfully"

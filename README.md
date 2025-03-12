# PTT (Penetration Testing Toolkit)

## Options (Use Environment Variables)

- `JSON=FALSE`
  - disables logging as JSON
- `LOG_LEVEL=2`
  - sets logging level (1-5)
  - see `hclog.Level`
- `ENV=DEV`
  - sets environment to dev
  - hosts static assets from file system instead of from embedded file system

## Dev Requirements

- Go >= 1.24
- Node.JS (LTS)
  - JS runtime for pnpm & tailwindcss
- pnpm
  - JS package manager
- tailwindcss CLI >= 4.0
  - CSS framework
  - requires Node.JS
  - install Tailwind CLI via pnpm
- github.com/a-h/templ
  - template engine for Go
- github.com/air-verse/air
  - code live reloading
- make
  - build automation
- grpc & protoc-compiler
  - rpc communication handler

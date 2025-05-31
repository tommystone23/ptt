#!/bin/bash

set -e

# Uninstall Go
if [ -d "/usr/local/go" ]; then
    sudo rm -rf /usr/local/go
fi

rm -rf "$(go env GOPATH)/bin/protoc-gen-go" || true
rm -rf "$(go env GOPATH)/bin/protoc-gen-go-grpc" || true
rm -rf "$(go env GOPATH)/bin/templ" || true
rm -rf "$(go env GOPATH)/bin/air" || true

# Unset Go path from .bashrc
sed -i '/\/usr\/local\/go\/bin/d' ~/.bashrc
sed -i '/GOPATH\/bin/d' ~/.bashrc

# Uninstall Node.js and pnpm
sudo apt remove -y nodejs || true
sudo npm uninstall -g pnpm || true

# Remove Node + pnpm configs and caches
rm -rf ~/.npm ~/.node-gyp ~/.pnpm-store || true

# Uninstall Tailwind CLI (via pnpm) 
rm -rf node_modules/

# Remove protoc and includes
sudo rm -f /usr/local/bin/protoc
sudo rm -rf /usr/local/include/google

# (Optional) Remove base dev packages (use with caution)
read -p "Remove base packages like build-essential, git, curl? [y/N] " CONFIRM
if [[ "$CONFIRM" =~ ^[Yy]$ ]]; then
    sudo apt purge -y curl unzip git build-essential make ca-certificates
    sudo apt autoremove -y
fi

echo "Successfully uninstalled PTT dependencies"

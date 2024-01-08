#!/bin/bash
# For AMD64 / x86_64
[ $(uname -m) = x86_64 ] && curl -Lo https://dl.k8s.io/release/v1.28.5/bin/linux/amd64/kubectl
# For ARM64
[ $(uname -m) = aarch64 ] && curl -Lo https://dl.k8s.io/release/v1.28.5/bin/linux/arm64/kubectl
chmod +x kubectl
mkdir -p ~/.local/bin
mv ./kubectl ~/.local/bin/kubectl
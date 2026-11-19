#!/bin/bash
set -euo pipefail

KUTTL_VERSION=0.22.0
KUTTL_SHA256=874bcea7b8820592018c660ce71f7f6dbfbdde35ee1f73cad7f560d1380e3868
KUTTL_BIN=$(mktemp)
trap 'rm -f "${KUTTL_BIN}"' EXIT

curl -fsSLo "${KUTTL_BIN}" "https://github.com/kudobuilder/kuttl/releases/download/v${KUTTL_VERSION}/kubectl-kuttl_${KUTTL_VERSION}_linux_x86_64"
echo "${KUTTL_SHA256}  ${KUTTL_BIN}" | sha256sum --check --quiet
sudo install -m 0755 "${KUTTL_BIN}" /usr/local/bin/kubectl-kuttl
export PATH=$PATH:/usr/local/bin

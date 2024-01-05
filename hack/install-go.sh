#!/bin/bash
rm -rf /usr/local/go
curl -Lo /tmp/ https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
tar -C /usr/local -xzf /tmp/go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
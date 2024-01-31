# Build the manager binary
FROM golang:1.21 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

COPY vendor/ vendor/
COPY internal/ internal/
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY clientcmd/ clientcmd/
COPY pkg/ pkg/
COPY scripts/ scripts/
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/scripts .
COPY --from=builder /workspace/manager .

USER 65532:65532

ENTRYPOINT ["/manager"]

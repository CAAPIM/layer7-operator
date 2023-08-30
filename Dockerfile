# Build the manager binary
FROM golang:1.20 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

COPY vendor/ vendor/
COPY internal/ internal/
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY clientcmd/ clientcmd/
COPY pkg/ pkg/
COPY scripts/ scripts/
# Build
RUN go env
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/scripts .
COPY --from=builder /workspace/manager .

USER 65532:65532

ENTRYPOINT ["/manager"]

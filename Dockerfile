FROM --platform=$BUILDPLATFORM GO_BUILD_IMG AS builder

ARG AUTHOR=Layer7
ARG VENDOR="Broadcom Inc."
ARG TITLE
ARG VERSION
ARG CREATED
ARG COPYRIGHT

LABEL org.opencontainers.image.created=${CREATED}
LABEL org.opencontainers.image.authors=${AUTHOR}
LABEL org.opencontainers.image.title=${TITLE}
LABEL org.opencontainers.image.version=${VERSION}
LABEL org.opencontainers.image.vendor=${VENDOR}
LABEL com.broadcom.copyright=${COPYRIGHT}

ARG GOPROXY
WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

COPY internal/ internal/
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY pkg/ pkg/
COPY scripts/ scripts/
ENV GOPROXY=${GOPROXY}
RUN go mod download

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -a -o manager cmd/main.go

FROM DISTROLESS_IMG

ARG AUTHOR=Layer7
ARG VENDOR="Broadcom Inc."
ARG TITLE
ARG VERSION
ARG CREATED
ARG COPYRIGHT

LABEL org.opencontainers.image.created=${CREATED}
LABEL org.opencontainers.image.authors=${AUTHOR}
LABEL org.opencontainers.image.title=${TITLE}
LABEL org.opencontainers.image.version=${VERSION}
LABEL org.opencontainers.image.vendor=${VENDOR}
LABEL com.broadcom.copyright=${COPYRIGHT}


WORKDIR /
COPY --from=builder /workspace/scripts .
COPY --from=builder /workspace/manager .

USER 65532:65532

ENTRYPOINT ["/manager"]

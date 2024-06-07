FROM GO_BUILD_IMG as builder

ARG AUTHOR=layer7
ARG VERSION
ARG CREATED
ARG COPYRIGHT

LABEL com.broadcom.ims.label.author=${AUTHOR}
LABEL com.broadcom.ims.label.version=${VERSION}
LABEL com.broadcom.ims.label.created=${CREATED}
LABEL com.broadcom.ims.label.copyright=${COPYRIGHT}

ARG GOPROXY
ARG TARGETARCH
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
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} GO111MODULE=on go build -a -o manager cmd/main.go


FROM DISTROLESS_IMG

ARG AUTHOR=layer7
ARG VERSION
ARG CREATED
ARG COPYRIGHT

LABEL com.broadcom.ims.label.author=${AUTHOR}
LABEL com.broadcom.ims.label.version=${VERSION}
LABEL com.broadcom.ims.label.created=${CREATED}
LABEL com.broadcom.ims.label.copyright=${COPYRIGHT}

WORKDIR /
COPY --from=builder /workspace/scripts .
COPY --from=builder /workspace/manager .

USER 65532:65532

ENTRYPOINT ["/manager"]

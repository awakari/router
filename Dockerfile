FROM golang:1.20.2-alpine3.17 AS builder
WORKDIR /go/src/router
COPY . .
RUN \
    apk add protoc protobuf-dev make git && \
    make build

FROM scratch
COPY --from=builder /go/src/router/router /bin/router
ENTRYPOINT ["/bin/router"]

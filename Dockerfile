FROM golang:1.10-alpine as build

# Install SSL certificates
RUN apk update && apk add --no-cache git ca-certificates gcc musl-dev

# Build static binary
RUN mkdir -p /go/src/github.com/uc-cdis/authproxy
WORKDIR /go/src/github.com/uc-cdis/authproxy
ADD . .
RUN mkdir -p bin && cd authProxyServer && go build -ldflags "-linkmode external -extldflags -static" -o ../bin/authProxyServer

# Set up small scratch image, and copy necessary things over
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/uc-cdis/authproxy/bin/authProxyServer /authProxyServer

ENTRYPOINT ["/authProxyServer"]
CMD ["--run"]

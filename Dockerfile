FROM golang:1.21.5 as build
RUN go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest

WORKDIR /go/src/shardrouter
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go .
RUN CGO_ENABLED=0 xcaddy build --with caddyshardrouter=.

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/src/shardrouter/caddy /
COPY Caddyfile /
COPY cert /cert
CMD ["/caddy", "run"]

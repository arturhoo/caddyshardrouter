FROM golang:1.24.6 as build

WORKDIR /go/src/shardrouter
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go .
COPY cmd/*.go cmd/
RUN CGO_ENABLED=0 go build -cover -o caddy cmd/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/src/shardrouter/caddy /
COPY Caddyfile /
COPY cert /cert
CMD ["/caddy", "run"]

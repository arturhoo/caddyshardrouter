FROM golang:1.25.4 as build

WORKDIR /go/src/upstream
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/upstream

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/upstream /
CMD ["/upstream"]

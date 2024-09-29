FROM golang:1.23.1 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o pg-to-gcs .

FROM postgres:17
COPY --from=builder /build/pg-to-gcs /usr/bin/pg-to-gcs

ENTRYPOINT ["/usr/bin/pg-to-gcs"]

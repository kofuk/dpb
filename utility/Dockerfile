FROM golang:1.26.1@sha256:e2ddb153f786ee6210bf8c40f7f35490b3ff7d38be70d1a0d358ba64225f6428 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o fixed-response .

FROM scratch
COPY --from=builder /build/fixed-response /fixed-response
ENTRYPOINT ["/fixed-response"]

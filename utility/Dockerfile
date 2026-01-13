FROM golang:1.25.5@sha256:0f406d34b7cb7255d0700af02ec28a2c88f1e00701055f4c282aa4c3ec0b3245 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o fixed-response .

FROM scratch
COPY --from=builder /build/fixed-response /fixed-response
ENTRYPOINT ["/fixed-response"]

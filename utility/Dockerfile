FROM golang:1.25.6@sha256:06d1251c59a75761ce4ebc8b299030576233d7437c886a68b43464bad62d4bb1 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o fixed-response .

FROM scratch
COPY --from=builder /build/fixed-response /fixed-response
ENTRYPOINT ["/fixed-response"]

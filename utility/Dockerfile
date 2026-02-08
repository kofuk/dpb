FROM golang:1.25.7@sha256:cc737435e2742bd6da3b7d575623968683609a3d2e0695f9d85bee84071c08e6 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o fixed-response .

FROM scratch
COPY --from=builder /build/fixed-response /fixed-response
ENTRYPOINT ["/fixed-response"]

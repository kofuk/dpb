FROM golang:1.25.4@sha256:3976069ea4797af8ebc1b78d56e0c6b42e85ec31703c83e0a1ffa34cb127a248 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o fixed-response .

FROM scratch
COPY --from=builder /build/fixed-response /fixed-response
ENTRYPOINT ["/fixed-response"]

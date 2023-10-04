FROM golang:1.21.0 AS builder
WORKDIR /build
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go


FROM alpine AS server
WORKDIR /Balancer
COPY --from=builder /build/app .
CMD ["./app"]
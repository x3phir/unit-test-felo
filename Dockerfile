FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /build/felo-demo ./cmd/felo-demo

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
COPY --from=builder /build/felo-demo /usr/local/bin/felo-demo

EXPOSE 50051 50052 50053 50054 50055

ENTRYPOINT ["felo-demo"]
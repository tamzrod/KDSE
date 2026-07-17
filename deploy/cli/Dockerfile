# KDSE CLI Runtime Dockerfile
# Isolated CLI deployment
# Version: 2.0.0

# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Build CLI binary only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kdse ./cmd/kdse/

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache git ca-certificates

RUN adduser -D -u 1000 kdse

WORKDIR /app

COPY --from=builder /build/kdse .
RUN chown -R kdse:kdse /app

USER kdse

CMD ["./kdse"]

# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.26.2-alpine AS builder

WORKDIR /src

# Download dependencies first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/container-runtime-mcp .

# Runtime stage
FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/container-runtime-mcp /usr/local/bin/container-runtime-mcp

EXPOSE 8080

ENTRYPOINT ["container-runtime-mcp"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install CA certs (needed if your app makes HTTPS calls)
RUN apk add --no-cache ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app

# -------- Runtime stage --------
FROM golang:1.25-alpine

WORKDIR /app

# Copy CA certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/app /app/app
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/static /app/static

EXPOSE 8000

ENTRYPOINT ["/app/app"]

# Dockerfile
FROM golang:1.24.2-alpine

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o wallet-backend ./cmd/main.go

# Run it
CMD ["./wallet-backend"]


# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o lark-skill cmd/server/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/lark-skill .

# Expose port
EXPOSE 8000

# Run the binary
CMD ["./lark-skill"]

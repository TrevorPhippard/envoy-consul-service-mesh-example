# Use Go 1.23+ to match module requirement
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the binary
RUN go build -o service-hello main.go

# ---- Final lightweight image ----
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/service-hello .

# Expose the app port
EXPOSE 8080

# Run the binary
CMD ["./service-hello"]

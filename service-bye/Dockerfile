# Use Go 1.23+ to match module requirement
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o service-bye ./cmd/main.go

# ---- Final lightweight image ----
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/service-bye .

EXPOSE 8080

CMD ["./service-bye"]

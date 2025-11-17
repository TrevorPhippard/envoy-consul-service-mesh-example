FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o service-bye ./cmd/service-bye

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/service-bye .

EXPOSE 8080

CMD ["./service-bye"]
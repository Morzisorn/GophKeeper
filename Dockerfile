FROM golang:1.23.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gophkeeper cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o agent cmd/agent/main.go

# Server image
FROM alpine:latest AS server
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/gophkeeper .

RUN mkdir -p ./keys
COPY --from=builder /app/*.pem ./keys/

EXPOSE 8080
CMD ["./gophkeeper"]

# Agent image
FROM alpine:latest AS agent
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/agent .

RUN mkdir -p ./keys
COPY --from=builder /app/*.pem ./keys/

CMD ["./agent"]

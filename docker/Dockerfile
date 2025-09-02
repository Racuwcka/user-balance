FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .

# Запуск приложения
CMD ["./app"]
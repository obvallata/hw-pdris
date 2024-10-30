# Строим приложение
FROM golang:1.22.7 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM arm64v8/debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]
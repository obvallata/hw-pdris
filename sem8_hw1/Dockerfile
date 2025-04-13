FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY app/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/logs

CMD ["./main"] 

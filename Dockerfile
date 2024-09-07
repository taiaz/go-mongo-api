FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o go-mongo-api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/go-mongo-api .

EXPOSE 8080

CMD ["./go-mongo-api"]

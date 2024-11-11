FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o golang-server cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/golang-server .
COPY --from=builder /app/.env .

CMD ["./golang-server"]

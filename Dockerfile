FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/url-shortener .

EXPOSE 8080 50051

CMD [ "./url-shortener", "--storage=postgres", "--postgres-conn=postgres://postgres:postgres@postgres/shortener?sslmode=disable" ]
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
COPY --from=builder /app/.env .

EXPOSE ${HTTP_PORT} ${GRPC_PORT}

CMD [ "./url-shortener" ]
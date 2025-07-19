FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM alpine:3.21

WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /app/app .
RUN chown -R appuser:appuser /app
EXPOSE 8080
USER appuser
CMD ["./app"]


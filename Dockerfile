FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o /bin/wareman ./cmd/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/wareman /bin/wareman
RUN apk add --no-cache libpq
EXPOSE 8080
CMD ["/bin/wareman"]

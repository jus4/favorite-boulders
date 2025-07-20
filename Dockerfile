FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the application code
COPY . .

RUN go build -o main ./cmd

FROM alpine:latest AS app

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/main .

# Copy static assets
COPY --from=builder /app/static ./static

# Expose port (optional, based on your app)
EXPOSE 8080

# Command to run the app
CMD ["./main"]

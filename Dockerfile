# Step 1: Build the Go binary
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

# Step 2: Run in a small image
FROM debian:bookworm-slim

# Install CA certificates (needed to make http request)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 7070

CMD ["./server"]

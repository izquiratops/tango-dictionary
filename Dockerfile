# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /root

# Set environment for cross-compilation
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy src
COPY . .

WORKDIR client
RUN go mod download
RUN go build -ldflags "-w" -o run main.go

# Stage 2: Create the Client image
FROM alpine:latest

WORKDIR /root/client
COPY --from=builder /root/client/run ./run

EXPOSE 8080

ENTRYPOINT ["./run"]
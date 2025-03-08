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

WORKDIR /
COPY --from=builder /root/client/run /client/run
COPY --from=builder /root/client/static /client/static
COPY --from=builder /root/client/template /client/template

EXPOSE 8080

ENTRYPOINT ["/client/run"]
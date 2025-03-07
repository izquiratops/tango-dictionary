# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /client

# Set environment for cross-compilation
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY client/go.mod client/go.sum ./
RUN go mod download

# Copy the rest of the source code and build it
COPY client .
RUN go build -ldflags "-w" -o run main.go

# Stage 2: Create the Client image
FROM alpine:latest

WORKDIR /client

COPY --from=builder /client/run /client
COPY --from=builder /client/client/static /client/static
COPY --from=builder /client/client/template /client/template

EXPOSE 8080

ENTRYPOINT ["/client/run"]
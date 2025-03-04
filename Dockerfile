# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Set environment for cross-compilation
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code and build it
COPY . .
RUN go build -ldflags "-w" -o tango main.go

# Stage 2: Create a lightweight image with the Go binary
FROM alpine:latest

# Create app directory
WORKDIR /app

# Copy only necessary files from builder
# - Executable
COPY --from=builder /app/tango /app
# - Static files
COPY --from=builder /app/server/static /app/server/static
COPY --from=builder /app/server/template /app/server/template

EXPOSE 8080

ENTRYPOINT ["/app/tango"]
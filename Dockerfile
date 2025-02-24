# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code and build it
COPY . .
RUN go build -o tango main.go

# Stage 2: Create a lightweight image with the Go binary
FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/tango .
COPY --from=builder /app/server/static ./server/static
COPY --from=builder /app/server/template ./server/template

EXPOSE 8080

# Run the Go binary with the specified arguments
ENTRYPOINT ["./tango"]
CMD ["-v", "$DB_VERSION"]
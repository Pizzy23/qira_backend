# Build stage
FROM golang:1.21.4-alpine as builder

# Set environment variables for Go modules
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install build dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application code
COPY . .

# Generate Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy the built binary and any other necessary files from the builder stage
COPY --from=builder /build/main .
COPY --from=builder /build/docs ./docs

# Set a non-root user and switch to it
RUN adduser -D appuser
USER appuser

# Command to run the binary
CMD ["./main"]

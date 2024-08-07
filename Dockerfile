# Build stage
FROM golang:1.21.4-alpine as builder

# Install git and other necessary tools
RUN apk add --no-cache git

# Set the working directory
WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates git

# Set environment variables
ENV GOROOT=/usr/local/go \
    GOPATH=/go \
    PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Set the working directory in the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /build/main .

# Copy your init script and swagger docs
COPY --from=builder /build/init.sh /app/init.sh
COPY --from=builder /build/docs ./docs

# Make the init script executable
RUN chmod +x /app/init.sh

# Command to run the init script
CMD ["./init.sh"]

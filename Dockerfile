# Build stage
FROM golang:1.21.4-alpine as builder

# Install git and other necessary dependencies
RUN apk add --no-cache git gcc g++

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Final stage to setup runtime container
FROM alpine:latest

# Install CA certificates, necessary for making HTTPS requests (e.g., git pull)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary and swagger files from builder stage
COPY --from=builder /build/main .
COPY --from=builder /build/docs ./docs
COPY --from=builder /build/init.sh /app/init.sh
RUN chmod +x /app/init.sh

# Set up the command to run when starting the container
ENTRYPOINT ["/app/init.sh"]

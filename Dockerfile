# Build stage
FROM golang:1.21.4-alpine as builder

# Set necessary environmet variables needed for our image
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

# Run stage
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary and swagger files from /build to /app
COPY --from=builder /build/main .
COPY --from=builder /build/docs ./docs

# Copy any other necessary files
COPY --from=builder /app/init.sh /app/init.sh
RUN chmod +x /app/init.sh

# Command to run
ENTRYPOINT ["/app/init.sh"]

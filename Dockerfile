# Build stage
FROM golang:1.21.4-alpine as builder

# Install git and other necessary tools for building
RUN apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Final stage
FROM alpine:latest

# Install runtime dependencies including Git and tools for swag and go
RUN apk --no-cache add ca-certificates git
ENV GOROOT=/usr/local/go \
    GOPATH=/go \
    PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Install Go and Swag in the runtime image if you need to run these at startup
COPY --from=builder /usr/local/go /usr/local/go
COPY --from=builder /go/bin/swag /usr/local/bin/swag

WORKDIR /app
COPY --from=builder /build/main .
COPY --from=builder /build/docs ./docs
COPY --from=builder /build/init.sh /app/init.sh 
RUN chmod +x /app/init.sh

# Command to run the init script
CMD ["/app/init.sh"]

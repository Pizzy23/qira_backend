# Dockerfile
FROM golang:1.21.4-alpine

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN apk add --no-cache sudo git openrc

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN chmod +x /app/init.sh

# Run the init script
CMD ["/app/init.sh"]

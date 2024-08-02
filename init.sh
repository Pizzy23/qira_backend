#!/bin/sh

# Set environment variables
export GOROOT=/usr/local/go
export GOPATH=/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Change to the app directory
cd /app

# Git pull with credentials
git pull https://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@github.com/Pizzy23/qira_backend.git

# Run swag init
swag init

# Build the Go application
go build -o main .

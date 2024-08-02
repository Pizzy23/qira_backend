#!/bin/sh

export GOROOT=/usr/local/go
export GOPATH=/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

cd /app

git pull https://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@github.com/Pizzy23/qira_backend.git

swag init

go build -o main .

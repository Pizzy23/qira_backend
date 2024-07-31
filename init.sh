#!/bin/sh

# Configurações do git para autenticação
git config --global credential.helper store

# Salva as credenciais do git
echo "https://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:@github.com" > ~/.git-credentials

# Executa o git pull
git pull

# Executa o swag init
swag init

# Executa o build da aplicação Go
go build -o main .

# Executa o daemon-reload
sudo systemctl daemon-reload

# Executa a aplicação Go
./main

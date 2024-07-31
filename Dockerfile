# Use uma imagem oficial do Go como base
FROM golang:1.21.4-alpine

# Instala o sudo, git, openrc (substituto para systemd no Alpine)
RUN apk add --no-cache sudo git openrc

# Configurações do ambiente Go
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Diretório de trabalho dentro do container
WORKDIR /app

# Instala o Swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copia o script de inicialização
COPY init.sh /app/init.sh

# Permite que o script de inicialização seja executado
RUN chmod +x /app/init.sh

# Executa o container como root
USER root

# Comando para rodar o script de inicialização
CMD ["/app/init.sh"]

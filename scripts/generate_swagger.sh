#!/bin/bash

# Instala o swag se não estiver instalado
if ! command -v swag &> /dev/null; then
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Gera a documentação Swagger
swag init -g main.go -o docs 
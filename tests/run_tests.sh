#!/bin/bash

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Função para executar os testes
run_tests() {
    echo "Executando testes..."
    echo "-------------------"

    # Teste de autenticação
    echo -n "Teste de autenticação: "
    if go test -run TestAuthFlow; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}FALHA${NC}"
        exit 1
    fi

    # Teste de usuário
    echo -n "Teste de usuário: "
    if go test -run TestUserFlow; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}FALHA${NC}"
        exit 1
    fi

    # Teste de perfil
    echo -n "Teste de perfil: "
    if go test -run TestProfileFlow; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}FALHA${NC}"
        exit 1
    fi

    # Cobertura de testes
    echo -e "\nCobertura de testes:"
    go test -cover ./...
}

# Verifica se a API está rodando
check_api() {
    if [ -z "$API_URL" ]; then
        API_URL="http://localhost:8080/api/v1"
    fi

    echo "Verificando se a API está rodando em $API_URL..."
    if curl -s "$API_URL/health" > /dev/null; then
        echo -e "${GREEN}API está rodando${NC}"
    else
        echo -e "${RED}API não está rodando${NC}"
        echo "Por favor, inicie a API ou configure a URL correta através da variável de ambiente API_URL"
        exit 1
    fi
}

# Executa os testes
check_api
run_tests 
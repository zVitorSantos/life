#!/bin/bash

# Script para executar testes localmente
set -e

echo "🚀 Iniciando testes locais..."

# Verifica se o Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado"
    exit 1
fi

# Verifica se o PostgreSQL está rodando
if ! pg_isready -h localhost -p 5432 &> /dev/null; then
    echo "⚠️  PostgreSQL não está rodando. Iniciando com Docker..."
    docker-compose up -d postgres
    sleep 5
fi

# Compila a API
echo "🔨 Compilando API..."
go build -o api

# Configura variáveis de ambiente para teste
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=life_test
export DB_PORT=5432
export JWT_SECRET=test_secret
export JWT_REFRESH_SECRET=test_refresh_secret
export PORT=8080
export ENV=test

# Cria arquivo .env temporário para os testes
cat > .env << EOF
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=life_test
DB_PORT=5432
JWT_SECRET=test_secret
JWT_REFRESH_SECRET=test_refresh_secret
PORT=8080
ENV=test
EOF

# Inicia a API em background
echo "🌐 Iniciando API..."
./api &
API_PID=$!

# Função para limpar ao sair
cleanup() {
    echo "🧹 Limpando..."
    kill $API_PID 2>/dev/null || true
    rm -f api coverage.txt coverage.html .env
}
trap cleanup EXIT

# Aguarda a API ficar disponível
echo "⏳ Aguardando API ficar disponível..."
for i in {1..30}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ API está rodando!"
        break
    fi
    echo "Aguardando... ($i/30)"
    sleep 2
done

# Verifica se a API está rodando
if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ Falha ao iniciar a API"
    exit 1
fi

# Executa os testes
echo "🧪 Executando testes..."
export API_URL=http://localhost:8080/api/v1
go test -v -coverprofile=coverage.txt -covermode=atomic ./tests/...

# Gera relatório de cobertura HTML
echo "📊 Gerando relatório de cobertura..."
go tool cover -html=coverage.txt -o coverage.html

echo "✅ Testes concluídos!"
echo "📊 Relatório de cobertura: coverage.html"
echo "📄 Arquivo de cobertura: coverage.txt" 
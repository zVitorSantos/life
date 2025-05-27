#!/bin/bash

# Configura as variáveis de ambiente
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=life_test
export JWT_SECRET=test_secret
export PORT=8080

# Inicia o banco de dados PostgreSQL
docker run -d \
  --name life-test-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=life_test \
  -p 5432:5432 \
  postgres:15

# Aguarda o banco de dados iniciar
echo "Aguardando banco de dados iniciar..."
sleep 5

# Inicia a API em background
echo "Iniciando API..."
go run main.go &
API_PID=$!

# Aguarda a API iniciar
echo "Aguardando API iniciar..."
sleep 5

# Executa os testes
echo "Executando testes..."
go test -v -coverprofile=coverage.txt -covermode=atomic ./tests/...

# Mata a API
kill $API_PID

# Remove o container do banco de dados
docker rm -f life-test-db 
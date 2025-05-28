# Script para rodar testes com cobertura
Write-Host "Iniciando testes com cobertura..." -ForegroundColor Green

# Verifica se o Docker está rodando
$dockerRunning = docker ps 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Docker não está rodando. Iniciando serviços..." -ForegroundColor Yellow
    docker-compose up -d
    Start-Sleep 10
}

# Roda os testes com cobertura
Write-Host "Executando testes..." -ForegroundColor Blue
go test -v ./tests/... -coverprofile=coverage.txt -covermode=atomic

# Verifica se o arquivo de cobertura foi gerado
if (Test-Path "coverage.txt") {
    Write-Host "Arquivo de cobertura gerado com sucesso!" -ForegroundColor Green
    
    # Mostra o relatório de cobertura
    Write-Host "Relatório de cobertura:" -ForegroundColor Blue
    go tool cover -func=coverage.txt
    
    # Gera relatório HTML (opcional)
    go tool cover -html=coverage.txt -o coverage.html
    Write-Host "Relatório HTML gerado: coverage.html" -ForegroundColor Green
} else {
    Write-Host "Erro: Arquivo de cobertura não foi gerado" -ForegroundColor Red
    exit 1
}

Write-Host "Testes concluídos!" -ForegroundColor Green 
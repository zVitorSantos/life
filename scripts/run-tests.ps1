# Script para executar testes localmente no Windows
param(
    [switch]$SkipBuild,
    [switch]$SkipCleanup
)

$ErrorActionPreference = "Stop"

Write-Host "Iniciando testes locais..." -ForegroundColor Green

# Verifica se o Go esta instalado
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Go nao esta instalado" -ForegroundColor Red
    exit 1
}

# Verifica se o PostgreSQL esta rodando (via Docker)
try {
    $dockerStatus = docker-compose ps postgres 2>$null
    if (-not $dockerStatus -or $dockerStatus -notmatch "Up") {
        Write-Host "PostgreSQL nao esta rodando. Iniciando com Docker..." -ForegroundColor Yellow
        docker-compose up -d postgres
        Start-Sleep 5
    }
} catch {
    Write-Host "Docker nao disponivel. Certifique-se de que o PostgreSQL esta rodando na porta 5432" -ForegroundColor Yellow
}

# Compila a API
if (-not $SkipBuild) {
    Write-Host "Compilando API..." -ForegroundColor Blue
    go build -o api.exe
}

# Configura variaveis de ambiente para teste
$env:DB_HOST = "localhost"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgres"
$env:DB_NAME = "life_test"
$env:DB_PORT = "5432"
$env:JWT_SECRET = "test_secret"
$env:JWT_REFRESH_SECRET = "test_refresh_secret"
$env:PORT = "8080"
$env:ENV = "test"

# Cria arquivo .env temporario para os testes
@"
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=life_test
DB_PORT=5432
JWT_SECRET=test_secret
JWT_REFRESH_SECRET=test_refresh_secret
PORT=8080
ENV=test
"@ | Out-File -FilePath ".env" -Encoding UTF8

# Inicia a API em background
Write-Host "Iniciando API..." -ForegroundColor Blue
$apiProcess = Start-Process -FilePath ".\api.exe" -PassThru -WindowStyle Hidden

# Funcao para limpar ao sair
function Cleanup {
    Write-Host "Limpando..." -ForegroundColor Yellow
    if ($apiProcess -and -not $apiProcess.HasExited) {
        Stop-Process -Id $apiProcess.Id -Force -ErrorAction SilentlyContinue
    }
    if (-not $SkipCleanup) {
        Remove-Item -Path "api.exe", "coverage.txt", "coverage.html", ".env" -ErrorAction SilentlyContinue
    }
}

# Registra funcao de limpeza para execucao ao sair
Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action { Cleanup }

try {
    # Aguarda a API ficar disponivel
    Write-Host "Aguardando API ficar disponivel..." -ForegroundColor Blue
    $maxAttempts = 30
    $attempt = 0
    $apiReady = $false

    do {
        $attempt++
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 2 -ErrorAction Stop
            if ($response.StatusCode -eq 200) {
                Write-Host "API esta rodando!" -ForegroundColor Green
                $apiReady = $true
                break
            }
        } catch {
            Write-Host "Aguardando... ($attempt/$maxAttempts)" -ForegroundColor Yellow
            Start-Sleep 2
        }
    } while ($attempt -lt $maxAttempts)

    if (-not $apiReady) {
        Write-Host "Falha ao iniciar a API" -ForegroundColor Red
        exit 1
    }

    # Executa os testes
    Write-Host "Executando testes..." -ForegroundColor Blue
    $env:API_URL = "http://localhost:8080/api"
    go test -v -coverprofile=coverage.txt -covermode=atomic ./tests/...

    if ($LASTEXITCODE -eq 0) {
        # Gera relatorio de cobertura HTML
        Write-Host "Gerando relatorio de cobertura..." -ForegroundColor Blue
        go tool cover -html=coverage.txt -o coverage.html

        Write-Host "Testes concluidos com sucesso!" -ForegroundColor Green
        Write-Host "Relatorio de cobertura: coverage.html" -ForegroundColor Cyan
        Write-Host "Arquivo de cobertura: coverage.txt" -ForegroundColor Cyan
    } else {
        Write-Host "Testes falharam!" -ForegroundColor Red
        exit $LASTEXITCODE
    }

} finally {
    Cleanup
} 
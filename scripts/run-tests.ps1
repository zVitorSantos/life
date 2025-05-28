# Script para executar testes localmente no Windows
param(
    [switch]$SkipBuild,
    [switch]$SkipCleanup
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 Iniciando testes locais..." -ForegroundColor Green

# Verifica se o Go está instalado
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go não está instalado" -ForegroundColor Red
    exit 1
}

# Verifica se o PostgreSQL está rodando (via Docker)
try {
    $dockerStatus = docker-compose ps postgres 2>$null
    if (-not $dockerStatus -or $dockerStatus -notmatch "Up") {
        Write-Host "⚠️  PostgreSQL não está rodando. Iniciando com Docker..." -ForegroundColor Yellow
        docker-compose up -d postgres
        Start-Sleep 5
    }
} catch {
    Write-Host "⚠️  Docker não disponível. Certifique-se de que o PostgreSQL está rodando na porta 5432" -ForegroundColor Yellow
}

# Compila a API
if (-not $SkipBuild) {
    Write-Host "🔨 Compilando API..." -ForegroundColor Blue
    go build -o api.exe
}

# Configura variáveis de ambiente para teste
$env:DB_HOST = "localhost"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgres"
$env:DB_NAME = "life_test"
$env:DB_PORT = "5432"
$env:JWT_SECRET = "test_secret"
$env:JWT_REFRESH_SECRET = "test_refresh_secret"
$env:PORT = "8080"
$env:ENV = "test"

# Inicia a API em background
Write-Host "🌐 Iniciando API..." -ForegroundColor Blue
$apiProcess = Start-Process -FilePath ".\api.exe" -PassThru -WindowStyle Hidden

# Função para limpar ao sair
function Cleanup {
    Write-Host "🧹 Limpando..." -ForegroundColor Yellow
    if ($apiProcess -and -not $apiProcess.HasExited) {
        Stop-Process -Id $apiProcess.Id -Force -ErrorAction SilentlyContinue
    }
    if (-not $SkipCleanup) {
        Remove-Item -Path "api.exe", "coverage.txt", "coverage.html" -ErrorAction SilentlyContinue
    }
}

# Registra função de limpeza para execução ao sair
Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action { Cleanup }

try {
    # Aguarda a API ficar disponível
    Write-Host "⏳ Aguardando API ficar disponível..." -ForegroundColor Blue
    $maxAttempts = 30
    $attempt = 0
    $apiReady = $false

    do {
        $attempt++
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 2 -ErrorAction Stop
            if ($response.StatusCode -eq 200) {
                Write-Host "✅ API está rodando!" -ForegroundColor Green
                $apiReady = $true
                break
            }
        } catch {
            Write-Host "Aguardando... ($attempt/$maxAttempts)" -ForegroundColor Yellow
            Start-Sleep 2
        }
    } while ($attempt -lt $maxAttempts)

    if (-not $apiReady) {
        Write-Host "❌ Falha ao iniciar a API" -ForegroundColor Red
        exit 1
    }

    # Executa os testes
    Write-Host "🧪 Executando testes..." -ForegroundColor Blue
    $env:API_URL = "http://localhost:8080/api"
    go test -v -coverprofile=coverage.txt -covermode=atomic ./tests/...

    if ($LASTEXITCODE -eq 0) {
        # Gera relatório de cobertura HTML
        Write-Host "📊 Gerando relatório de cobertura..." -ForegroundColor Blue
        go tool cover -html=coverage.txt -o coverage.html

        Write-Host "✅ Testes concluídos com sucesso!" -ForegroundColor Green
        Write-Host "📊 Relatório de cobertura: coverage.html" -ForegroundColor Cyan
        Write-Host "📄 Arquivo de cobertura: coverage.txt" -ForegroundColor Cyan
    } else {
        Write-Host "❌ Testes falharam!" -ForegroundColor Red
        exit $LASTEXITCODE
    }

} finally {
    Cleanup
} 
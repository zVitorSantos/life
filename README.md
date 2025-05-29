# Life Game API

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/zVitorSantos/life/ci.yml?branch=main&style=for-the-badge&logo=github-actions)](https://github.com/zVitorSantos/life/actions)
[![codecov](https://img.shields.io/codecov/c/github/zVitorSantos/life?style=for-the-badge&logo=codecov)](https://codecov.io/gh/zVitorSantos/life)
[![License](https://img.shields.io/github/license/zVitorSantos/life?style=for-the-badge)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/zVitorSantos/life?style=for-the-badge)](https://goreportcard.com/report/github.com/zVitorSantos/life)

API RESTful para o jogo Life, desenvolvida em Go com sistema completo de jogo, economia virtual e progressão de jogadores.

## 🎮 Funcionalidades

### Sistema de Autenticação
- Autenticação com JWT e Refresh Tokens
- Gerenciamento de usuários
- Sistema de API Keys

### Sistema de Jogo
- **GameProfile**: Perfis de jogo com sistema de levels e XP
- **Wallet**: Sistema de carteira multi-moeda (Coins, Gems, Tokens)
- **Transaction**: Sistema completo de transações com auditoria
- **GameSession**: Controle de sessões ativas dos jogadores
- **Leaderboard**: Ranking de jogadores por XP

### Recursos Técnicos
- Health Checks
- Documentação Swagger
- Logging estruturado
- Validação de dados
- Tratamento de erros personalizado
- Testes automatizados com 70%+ de cobertura

## 📋 Pré-requisitos

- Go 1.21 ou superior
- PostgreSQL 12 ou superior
- Docker e Docker Compose (opcional)

## 🔧 Instalação

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/life.git
cd life
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure o arquivo `.env`:
```env
# Configurações do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=life_game

# Configurações da API
PORT=8080
ENV=development

# Configurações de Segurança
JWT_SECRET=sua_chave_secreta

# Configurações de Log
LOG_LEVEL=debug
LOG_FORMAT=json
```

4. Execute as migrações:
```bash
go run main.go
```

## 🚀 Executando o projeto

### Usando Go
```bash
go run main.go
```

### Usando Docker
```bash
docker-compose up -d
```

## 📚 Documentação da API

A documentação completa da API está disponível via Swagger UI em:
```
http://localhost:8080/swagger/index.html
```

### Endpoints Principais

#### Autenticação
- `POST /api/v1/register` - Registra um novo usuário
- `POST /api/v1/login` - Realiza login e retorna tokens
- `POST /api/v1/refresh` - Atualiza o access token
- `POST /api/v1/logout` - Revoga um refresh token

#### Usuários
- `GET /api/v1/profile` - Obtém perfil do usuário
- `PUT /api/v1/profile` - Atualiza perfil do usuário

#### GameProfile (Sistema de Jogo)
- `POST /api/v1/game-profile` - Cria perfil de jogo
- `GET /api/v1/game-profile` - Obtém perfil de jogo
- `PUT /api/v1/game-profile` - Atualiza perfil de jogo
- `POST /api/v1/game-profile/xp` - Adiciona XP (com level up automático)
- `GET /api/v1/game-profile/stats` - Obtém estatísticas do jogador
- `PUT /api/v1/game-profile/stats` - Atualiza estatísticas
- `PUT /api/v1/game-profile/last-login` - Atualiza último login

#### Wallet (Sistema de Carteira)
- `POST /api/v1/wallet` - Cria carteira
- `GET /api/v1/wallet` - Obtém dados da carteira
- `GET /api/v1/wallet/balance/{currency}` - Obtém saldo específico
- `GET /api/v1/wallet/balances` - Obtém todos os saldos
- `POST /api/v1/wallet/lock` - Bloqueia carteira
- `POST /api/v1/wallet/unlock` - Desbloqueia carteira
- `GET /api/v1/wallet/status` - Obtém status da carteira
- `GET /api/v1/wallet/history` - Obtém histórico da carteira

#### Transaction (Sistema de Transações)
- `POST /api/v1/transactions/add` - Adiciona dinheiro
- `POST /api/v1/transactions/spend` - Gasta dinheiro
- `POST /api/v1/transactions/transfer` - Transfere entre carteiras
- `GET /api/v1/transactions/history` - Obtém histórico de transações
- `GET /api/v1/transactions/{id}` - Obtém transação específica

#### Leaderboard
- `GET /api/v1/leaderboard` - Obtém ranking de jogadores

#### API Keys
- `POST /api/v1/api-keys` - Cria uma nova API key
- `GET /api/v1/api-keys` - Lista API keys do usuário
- `PUT /api/v1/api-keys/{id}` - Atualiza uma API key
- `DELETE /api/v1/api-keys/{id}` - Remove uma API key

#### Health Checks
- `GET /health` - Verifica a saúde da aplicação
- `GET /ready` - Verifica se a aplicação está pronta
- `GET /live` - Verifica se a aplicação está viva

## 🧪 Testes

### Executando Testes Localmente

#### Usando Scripts Automatizados

**Linux/macOS:**
```bash
# Executa todos os testes com API rodando
./scripts/run-tests.sh
```

**Windows (PowerShell):**
```powershell
# Executa todos os testes com API rodando
.\scripts\run-tests.ps1

# Opções disponíveis:
.\scripts\run-tests.ps1 -SkipBuild    # Pula a compilação
.\scripts\run-tests.ps1 -SkipCleanup  # Mantém arquivos temporários
```

#### Executando Manualmente

```bash
# 1. Inicie o PostgreSQL (via Docker)
docker-compose up -d postgres

# 2. Configure as variáveis de ambiente
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=life_test
export DB_PORT=5432
export JWT_SECRET=test_secret
export JWT_REFRESH_SECRET=test_refresh_secret
export PORT=8080
export ENV=test

# 3. Compile e inicie a API
go build -o api
./api &

# 4. Execute os testes
export API_URL=http://localhost:8080/api
go test -v -coverprofile=coverage.txt -covermode=atomic ./tests/...

# 5. Gere relatório de cobertura
go tool cover -html=coverage.txt -o coverage.html
```

### Comandos de Teste Rápidos

```bash
# Apenas testes unitários (sem API)
go test ./tests/ -run TestUserModel

# Testes com verbose
go test ./tests/... -v

# Testes com cobertura detalhada
go test ./tests/... -coverprofile=coverage.txt -covermode=atomic -coverpkg=./...

# Para executar todos os testes
go test -v ./...

# Para visualizar o relatório de cobertura
go tool cover -html=coverage.txt
```

### Cobertura de Código

O projeto utiliza [Codecov](https://codecov.io/gh/zVitorSantos/life) para monitoramento de cobertura de código:

- **Meta de cobertura**: 70%
- **Threshold**: 2%
- **Relatórios automáticos**: Gerados a cada push via GitHub Actions

### Tipos de Teste

- **Testes de Unidade**: Testam componentes individuais
- **Testes de Integração**: Testam fluxos completos da API
- **Testes de Modelo**: Validam estruturas de dados

### Estrutura de Testes

```
tests/
├── auth_test.go           # Testes de autenticação
├── profile_test.go        # Testes de perfil
├── user_test.go           # Testes de modelo de usuário
├── game_profile_test.go   # Testes do sistema de perfil de jogo
├── wallet_test.go         # Testes do sistema de carteira
├── transaction_test.go    # Testes do sistema de transações
└── config.go              # Configuração dos testes
```

## 📦 Estrutura do Projeto

```
.
├── config/         # Configurações da aplicação
├── docs/          # Documentação Swagger
├── errors/        # Erros personalizados
├── handlers/      # Handlers HTTP
├── logger/        # Configuração de logging
├── middleware/    # Middlewares
├── models/        # Modelos de dados
├── routes/        # Rotas da API
├── scripts/       # Scripts utilitários
├── tests/         # Testes
├── validator/     # Validação de dados
├── .env           # Variáveis de ambiente
├── .gitignore     # Arquivos ignorados pelo git
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```

## 🔐 Segurança

- Autenticação JWT com refresh tokens
- Validação robusta de dados
- Sanitização de inputs
- Rate limiting
- Headers de segurança

## 📈 Monitoramento

- Health checks
- Logging estruturado
- Métricas de performance

## 🤝 Contribuindo

1. Faça o fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## ✨ Próximos Passos

- [ ] Implementar cache com Redis
- [ ] Adicionar sistema de pontuação
- [ ] Implementar sistema de níveis
- [ ] Adicionar sistema de conquistas
- [ ] Implementar sistema de amigos
- [ ] Adicionar sistema de chat
- [ ] Implementar WebSocket para real-time
- [ ] Adicionar suporte a múltiplos idiomas 
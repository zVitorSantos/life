# Sistema de Jogo - Life API

Este documento descreve o sistema completo de jogo implementado na Life API, incluindo perfis de jogador, economia virtual e sistema de transações.

## 📋 Visão Geral

O sistema de jogo é composto por 4 componentes principais:

1. **GameProfile** - Perfis de jogo dos usuários
2. **Wallet** - Sistema de carteira multi-moeda
3. **Transaction** - Sistema de transações com auditoria
4. **GameSession** - Controle de sessões ativas

## 🎮 GameProfile

### Descrição
O GameProfile representa o perfil de jogo de um usuário, separado dos dados pessoais. Inclui sistema de progressão, estatísticas e configurações.

### Campos Principais
- **Level**: Nível atual do jogador (inicia em 1)
- **XP**: Pontos de experiência acumulados
- **Stats**: Estatísticas flexíveis (JSONB)
- **Settings**: Configurações do jogador (JSONB)
- **IsActive**: Status ativo/inativo
- **LastLogin**: Último acesso

### Sistema de Progressão
- **Fórmula XP**: `level * 1000` XP necessário para próximo nível
- **Level Up Automático**: Ao atingir XP suficiente
- **Progresso**: Cálculo de % para próximo nível

### Métodos Úteis
```go
// Adiciona XP e verifica level up
profile.AddXP(500)

// Calcula XP necessário para próximo nível
nextLevelXP := profile.GetXPForNextLevel()

// Obtém progresso atual (0-100%)
progress := profile.GetXPProgress()

// Manipula estatísticas
profile.SetStat("games_played", 10)
gamesPlayed := profile.GetStat("games_played")

// Manipula configurações
profile.SetSetting("theme", "dark")
theme := profile.GetSetting("theme")
```

### Endpoints
```
POST   /api/v1/game-profile           # Criar perfil
GET    /api/v1/game-profile           # Obter perfil
PUT    /api/v1/game-profile           # Atualizar perfil
POST   /api/v1/game-profile/xp        # Adicionar XP
GET    /api/v1/game-profile/stats     # Obter estatísticas
PUT    /api/v1/game-profile/stats     # Atualizar estatísticas
PUT    /api/v1/game-profile/last-login # Atualizar último login
```

## 💰 Wallet (Sistema de Carteira)

### Descrição
Sistema de carteira multi-moeda com controle de segurança e auditoria completa.

### Tipos de Moeda
- **Coins**: Moeda principal do jogo
- **Gems**: Moeda premium (1 gem = 100 coins)
- **Tokens**: Moeda especial/evento (1 token = 10 coins)

### Campos Principais
- **CoinsBalance**: Saldo em coins
- **GemsBalance**: Saldo em gems
- **TokensBalance**: Saldo em tokens
- **IsLocked**: Status de bloqueio
- **LockReason**: Motivo do bloqueio

### Recursos de Segurança
- **Bloqueio/Desbloqueio**: Com motivo registrado
- **Validação de Saldo**: Antes de gastos
- **Permissões**: Verificação antes de transações

### Métodos Úteis
```go
// Obter saldo específico
balance := wallet.GetBalance(models.CurrencyCoins)

// Definir saldo
wallet.SetBalance(models.CurrencyGems, 100)

// Adicionar/subtrair saldo
newBalance := wallet.AddBalance(models.CurrencyCoins, 500)

// Verificar saldo suficiente
canSpend := wallet.HasSufficientBalance(models.CurrencyCoins, 1000)

// Verificar se pode gastar (não bloqueada + saldo)
canSpend := wallet.CanSpend(models.CurrencyCoins, 500)

// Bloquear/desbloquear
wallet.Lock("Suspeita de fraude")
wallet.Unlock()

// Valor total para ranking
totalValue := wallet.GetTotalValue()
```

### Endpoints
```
POST   /api/v1/wallet                 # Criar carteira
GET    /api/v1/wallet                 # Obter carteira
GET    /api/v1/wallet/balance/{currency} # Saldo específico
GET    /api/v1/wallet/balances        # Todos os saldos
POST   /api/v1/wallet/lock            # Bloquear carteira
POST   /api/v1/wallet/unlock          # Desbloquear carteira
GET    /api/v1/wallet/status          # Status da carteira
GET    /api/v1/wallet/history         # Histórico da carteira
```

## 📊 Transaction (Sistema de Transações)

### Descrição
Sistema completo de transações com auditoria, reversão e metadados flexíveis.

### Tipos de Transação
- **earn**: Ganhou dinheiro
- **spend**: Gastou dinheiro
- **transfer**: Transferência entre usuários
- **reward**: Recompensa do sistema
- **penalty**: Penalidade/multa
- **refund**: Reembolso

### Status de Transação
- **pending**: Pendente
- **completed**: Concluída
- **failed**: Falhou
- **cancelled**: Cancelada
- **reversed**: Revertida

### Campos Principais
- **Type**: Tipo da transação
- **Status**: Status atual
- **Currency**: Tipo de moeda
- **Amount**: Valor (positivo/negativo)
- **BalanceBefore/After**: Saldos antes/depois
- **Description**: Descrição da transação
- **Metadata**: Dados adicionais (JSONB)

### Recursos Avançados
- **Auditoria Completa**: Saldos antes/depois
- **Reversão**: Transações podem ser revertidas
- **Transferências**: Entre carteiras diferentes
- **Metadados**: Informações flexíveis

### Métodos Úteis
```go
// Verificar se pode ser revertida
canReverse := transaction.CanBeReversed()

// Marcar como concluída
transaction.Complete()

// Marcar como falhada
transaction.Fail()

// Cancelar transação
transaction.Cancel()

// Reverter transação
transaction.Reverse(reversalTransaction)

// Manipular metadados
transaction.SetMetadata("source", "quest_reward")
source := transaction.GetMetadata("source")

// Verificar se é transferência
isTransfer := transaction.IsTransfer()

// Valor absoluto
absAmount := transaction.GetAbsoluteAmount()
```

### Endpoints
```
POST   /api/v1/transactions/add       # Adicionar dinheiro
POST   /api/v1/transactions/spend     # Gastar dinheiro
POST   /api/v1/transactions/transfer  # Transferir dinheiro
GET    /api/v1/transactions/history   # Histórico de transações
GET    /api/v1/transactions/{id}      # Transação específica
```

## 🏆 Leaderboard

### Descrição
Sistema de ranking baseado no valor total da carteira e XP dos jogadores.

### Cálculo de Ranking
- **Valor Total**: Coins + (Gems × 100) + (Tokens × 10)
- **XP**: Pontos de experiência
- **Ordenação**: Por valor total decrescente

### Endpoint
```
GET    /api/v1/leaderboard?limit=10&offset=0
```

## 🎯 GameSession (Controle de Sessões)

### Descrição
Controle de sessões ativas dos jogadores com informações técnicas e estatísticas.

### Campos Principais
- **Status**: online, away, idle, offline
- **IPAddress**: IP da sessão
- **UserAgent**: Navegador/cliente
- **Platform**: Plataforma (web, mobile, etc.)
- **ExpiresAt**: Expiração da sessão
- **LastActivity**: Última atividade

### Status de Sessão
- **active**: Sessão ativa
- **inactive**: Sessão inativa
- **expired**: Sessão expirada
- **terminated**: Sessão terminada

## 🔄 Fluxo de Uso Típico

### 1. Criação de Perfil de Jogo
```bash
POST /api/v1/game-profile
{
  "stats": {"games_played": 0, "wins": 0},
  "settings": {"theme": "dark", "sound": true}
}
```

### 2. Criação de Carteira
```bash
POST /api/v1/wallet
{
  "initial_coins": 1000,
  "initial_gems": 50,
  "initial_tokens": 10
}
```

### 3. Adição de XP (com Level Up)
```bash
POST /api/v1/game-profile/xp
{
  "xp": 500,
  "reason": "Quest completed"
}
```

### 4. Transação de Gasto
```bash
POST /api/v1/transactions/spend
{
  "currency": "coins",
  "amount": 100,
  "description": "Bought item",
  "metadata": {"item_id": "sword_001"}
}
```

### 5. Consulta de Leaderboard
```bash
GET /api/v1/leaderboard?limit=10
```

## 🛡️ Segurança e Validações

### Validações Implementadas
- **Saldo Suficiente**: Antes de qualquer gasto
- **Carteira Desbloqueada**: Para transações
- **Valores Positivos**: Para adições
- **Tipos Válidos**: Moedas e transações
- **Autenticação**: JWT obrigatório

### Auditoria
- **Log Completo**: Todas as transações
- **Saldos Históricos**: Antes/depois
- **Metadados**: Contexto adicional
- **Reversibilidade**: Transações podem ser desfeitas

## 📈 Métricas e Estatísticas

### Dados Coletados
- **Progressão**: Levels e XP
- **Economia**: Saldos e transações
- **Atividade**: Sessões e último login
- **Performance**: Estatísticas personalizadas

### Relatórios Disponíveis
- **Leaderboard**: Ranking de jogadores
- **Histórico**: Transações e atividades
- **Estatísticas**: Dados personalizados
- **Sessões**: Atividade dos jogadores

## 🔧 Configurações Flexíveis

### GameProfile Stats (Exemplos)
```json
{
  "games_played": 150,
  "wins": 89,
  "losses": 61,
  "win_rate": 59.3,
  "total_playtime": 12500,
  "achievements": ["first_win", "level_10", "rich_player"]
}
```

### GameProfile Settings (Exemplos)
```json
{
  "theme": "dark",
  "sound": true,
  "notifications": true,
  "language": "pt-BR",
  "auto_save": true,
  "difficulty": "normal"
}
```

### Transaction Metadata (Exemplos)
```json
{
  "source": "quest_reward",
  "quest_id": "daily_001",
  "item_purchased": "sword_legendary",
  "shop_category": "weapons",
  "promotion_code": "WELCOME50"
}
```

## 🚀 Próximos Passos

### Funcionalidades Planejadas
- **Inventário**: Sistema de itens
- **Quests**: Sistema de missões
- **Guilds**: Sistema de grupos
- **PvP**: Batalhas entre jogadores
- **Events**: Eventos temporários
- **Achievements**: Sistema de conquistas

### Melhorias Técnicas
- **Cache**: Redis para performance
- **Websockets**: Atualizações em tempo real
- **Analytics**: Métricas avançadas
- **Backup**: Sistema de backup automático 
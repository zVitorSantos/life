# Testes da API

Este diretório contém os testes automatizados para a API.

## Estrutura dos Testes

- `auth_test.go`: Testes das rotas de autenticação (registro, login, refresh token e logout)
- `user_test.go`: Testes das rotas de usuário (obter, atualizar e listar usuários)
- `profile_test.go`: Testes das rotas de perfil (obter e atualizar perfil)
- `config.go`: Configuração do ambiente de teste

## Como Executar os Testes

1. Certifique-se de que a API está rodando localmente na porta 8080 ou configure a URL da API através da variável de ambiente `API_URL`:

```bash
# Para usar a API local
go test ./...

# Para usar uma API específica
API_URL=http://sua-api.com/api/v1 go test ./...
```

2. Para executar um teste específico:

```bash
# Teste de autenticação
go test -run TestAuthFlow

# Teste de usuário
go test -run TestUserFlow

# Teste de perfil
go test -run TestProfileFlow
```

3. Para ver a cobertura de testes:

```bash
go test -cover ./...
```

## Observações

- Os testes são independentes e podem ser executados em qualquer ordem
- Cada teste cria seus próprios dados e limpa após a execução
- O timeout padrão para as requisições é de 10 segundos
- Os testes assumem que a API está rodando e acessível 
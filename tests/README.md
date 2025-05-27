# Testes da API Life

Este diretório contém os testes automatizados da API Life.

## Estrutura dos Testes

- `auth_test.go`: Testes de autenticação (registro, login, refresh token, logout)
- `profile_test.go`: Testes de perfil (obter e atualizar perfil)
- `user_test.go`: Testes de usuário (modelo, CRUD de usuários)
- `config.go`: Configurações compartilhadas entre os testes

## Executando os Testes

### Localmente

1. Certifique-se de que o PostgreSQL está rodando:
```bash
docker-compose up -d postgres
```

2. Configure as variáveis de ambiente:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=life_test
export JWT_SECRET=test_secret
export JWT_REFRESH_SECRET=test_refresh_secret
export JWT_EXPIRATION=3600
export GIN_MODE=test
export API_PORT=8080
```

3. Execute os testes:
```bash
go test -v ./tests/...
```

### Com Cobertura

```bash
go test -v -coverprofile=coverage.txt -covermode=atomic ./tests/...
go tool cover -html=coverage.txt
```

### No CI/CD

Os testes são executados automaticamente no GitHub Actions quando:
- Um push é feito para as branches `main` ou `develop`
- Um Pull Request é criado para as branches `main` ou `develop`

O workflow:
1. Configura o ambiente de teste
2. Inicia o PostgreSQL
3. Compila e inicia a API
4. Executa os testes
5. Gera o relatório de cobertura
6. Envia a cobertura para o Codecov

## Padrões de Teste

1. **Estrutura**:
   - Cada arquivo de teste testa um módulo específico
   - Funções de teste são nomeadas com prefixo `Test`
   - Funções auxiliares são nomeadas com prefixo `test`

2. **Logs**:
   - Todas as requisições são logadas
   - Todas as respostas são logadas
   - Erros são logados com detalhes

3. **Assertions**:
   - Verificação de status code
   - Verificação de corpo da resposta
   - Verificação de headers quando necessário

4. **Cleanup**:
   - Dados de teste são limpos após cada teste
   - Conexões são fechadas adequadamente
   - Recursos são liberados

## Cobertura de Código

A cobertura de código é monitorada pelo Codecov. O relatório pode ser visualizado em:
https://codecov.io/gh/seu-usuario/life

## Contribuindo

1. Adicione testes para novas funcionalidades
2. Mantenha a cobertura de código acima de 80%
3. Siga os padrões de teste existentes
4. Documente casos de teste complexos 
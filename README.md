# Life Game API

API RESTful para o jogo Life, desenvolvida em Go com Docker.

## Requisitos

- Docker
- Docker Compose
- Go 1.21 ou superior

## Configuração

1. Clone o repositório
2. Copie o arquivo `.env.example` para `.env` e configure as variáveis de ambiente:
   ```bash
   cp .env.example .env
   ```
3. Edite o arquivo `.env` com suas configurações:
   - `DB_PASSWORD`: Senha do banco de dados
   - `JWT_SECRET`: Chave secreta para JWT (use uma string aleatória segura)
4. Execute o projeto com Docker Compose:
   ```bash
   docker-compose up --build
   ```

## Documentação

A documentação da API está disponível através do Swagger UI em:
```
http://localhost:8080/swagger/index.html
```

Para gerar a documentação Swagger:
```bash
swag init -g main.go -o docs
```

## Testes

Para executar os testes:
```bash
go test ./tests/...
```

## Logging

O sistema utiliza logging estruturado com zerolog. Os logs incluem:
- Método HTTP
- Path
- Status code
- Latência
- IP do cliente
- User Agent
- Erros (se houver)

## Endpoints

### Públicos

- `POST /api/register` - Registro de novo usuário
- `POST /api/login` - Login de usuário

### Protegidos (requer autenticação)

- `GET /api/profile` - Obter perfil do usuário
- `PUT /api/profile` - Atualizar perfil do usuário

## Estrutura do Projeto

```
.
├── config/         # Configurações do projeto
├── docs/          # Documentação Swagger
├── logger/        # Sistema de logging
├── middleware/    # Middlewares da aplicação
├── models/        # Modelos de dados
├── scripts/       # Scripts utilitários
├── tests/         # Testes unitários e de integração
├── .env           # Variáveis de ambiente (não versionado)
├── .env.example   # Exemplo de variáveis de ambiente
├── docker-compose.yml
├── Dockerfile
└── main.go        # Arquivo principal
```

## Segurança

- Autenticação JWT
- Senhas criptografadas
- Validação de dados
- Proteção contra SQL Injection (GORM)
- Headers de segurança
- Logging estruturado
- Documentação Swagger

## Arquivos Sensíveis

Os seguintes arquivos NÃO devem ser versionados:
- `.env` (contém senhas e tokens)
- `postgres_data/` (dados do banco)

## Próximos Passos

- [ ] Implementar rate limiting
- [ ] Adicionar validação de dados
- [ ] Implementar recuperação de senha
- [ ] Adicionar mais testes unitários
- [ ] Implementar testes de integração
- [ ] Adicionar monitoramento 
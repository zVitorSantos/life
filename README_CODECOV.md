# Configuração do Codecov

Este documento explica como configurar e usar o Codecov no projeto Life.

## Configuração Inicial

### 1. Obter Token do Codecov

1. Acesse [codecov.io](https://codecov.io)
2. Faça login com sua conta GitHub
3. Adicione o repositório `life`
4. Copie o token fornecido

### 2. Configurar GitHub Actions

Adicione o token como secret no GitHub:

1. Vá para Settings > Secrets and variables > Actions
2. Clique em "New repository secret"
3. Nome: `CODECOV_TOKEN`
4. Valor: Cole o token obtido do Codecov

### 3. Arquivo de Configuração

O arquivo `codecov.yml` já está configurado com:

- **Target de cobertura**: 80%
- **Threshold**: 1%
- **Arquivos ignorados**: Testes, docs, README

## Executando Testes Localmente

### Usando o Script PowerShell

```powershell
.\scripts\test-coverage.ps1
```

### Manualmente

```bash
# Iniciar serviços
docker-compose up -d

# Rodar testes com cobertura
go test -v ./tests/... -coverprofile=coverage.txt -covermode=atomic

# Ver relatório
go tool cover -func=coverage.txt

# Gerar HTML
go tool cover -html=coverage.txt -o coverage.html
```

## Enviando para Codecov

### Via GitHub Actions (Automático)

Os dados são enviados automaticamente quando você faz push para qualquer branch.

### Manualmente (Local)

```bash
# Instalar codecov CLI
go install github.com/codecov/codecov-cli@latest

# Enviar cobertura
codecov upload --file coverage.txt --token YOUR_TOKEN
```

## Interpretando Resultados

### Status Badges

Adicione ao README principal:

```markdown
[![codecov](https://codecov.io/gh/zVitorSantos/life/branch/main/graph/badge.svg)](https://codecov.io/gh/zVitorSantos/life)
```

### Métricas Importantes

- **Project Coverage**: Cobertura total do projeto
- **Patch Coverage**: Cobertura das mudanças no PR
- **Files Changed**: Arquivos modificados e sua cobertura

## Configuração do codecov.yml

```yaml
coverage:
  status:
    project:
      default:
        target: 80%        # Meta de cobertura
        threshold: 1%      # Tolerância de queda
    patch:
      default:
        target: 80%        # Meta para patches
        threshold: 1%      # Tolerância para patches

comment:
  layout: "reach, diff, flags, files"
  behavior: default
  require_changes: false

ignore:
  - "**/*_test.go"        # Ignora arquivos de teste
  - "**/testdata/*"       # Ignora dados de teste
  - "**/tests/*"          # Ignora diretório de testes
  - "docs/*"              # Ignora documentação
  - "*.md"                # Ignora arquivos markdown
```

## Troubleshooting

### Cobertura 0%

Se a cobertura aparecer como 0%, verifique:

1. Os testes estão passando?
2. O arquivo `coverage.txt` foi gerado?
3. Os testes estão testando o código principal?

### Token Inválido

Se o upload falhar:

1. Verifique se o token está correto
2. Confirme se o repositório foi adicionado no Codecov
3. Verifique se o secret está configurado no GitHub

### Conflitos de Merge

Para evitar conflitos ao fazer merge:

1. Sempre faça rebase da develop antes de criar PR
2. Use `--no-ff` para preservar histórico de features
3. Resolva conflitos localmente antes do push

## GitFlow Recomendado

```bash
# Criar feature
git checkout develop
git pull origin develop
git checkout -b feature/nova-funcionalidade

# Desenvolver e testar
# ... fazer alterações ...
git add .
git commit -m "feat: adiciona nova funcionalidade"

# Finalizar feature
git checkout develop
git pull origin develop
git checkout feature/nova-funcionalidade
git rebase develop
git checkout develop
git merge feature/nova-funcionalidade --no-ff
git push origin develop

# Limpar branch
git branch -d feature/nova-funcionalidade
git push origin --delete feature/nova-funcionalidade
``` 
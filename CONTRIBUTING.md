# Guia de Contribuição

Obrigado por considerar contribuir com o projeto Life Game API! Este documento fornece diretrizes e instruções para contribuir.

## 📋 Código de Conduta

Este projeto e todos que participam dele estão comprometidos com um ambiente amigável e seguro para todos. Por favor, leia nosso [Código de Conduta](CODE_OF_CONDUCT.md) para manter um ambiente respeitoso e inclusivo.

## 🤝 Como Contribuir

### 1. Configuração do Ambiente

1. Faça um fork do projeto
2. Clone seu fork:
```bash
git clone https://github.com/seu-usuario/life.git
cd life
```

3. Adicione o repositório original como upstream:
```bash
git remote add upstream https://github.com/original-usuario/life.git
```

4. Instale as dependências:
```bash
go mod download
```

5. Configure o ambiente de desenvolvimento:
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

### 2. Fluxo de Trabalho

1. Crie uma branch para sua feature:
```bash
git flow feature start nome-da-feature
```

2. Faça suas alterações seguindo as convenções de código

3. Execute os testes:
```bash
go test ./...
```

4. Commit suas alterações:
```bash
git add .
git commit -m "feat: descrição da feature"
```

5. Push para sua branch:
```bash
git push origin feature/nome-da-feature
```

6. Abra um Pull Request

### 3. Convenções de Código

- Use `gofmt` para formatação
- Siga as convenções de nomenclatura do Go
- Documente funções e tipos públicos
- Adicione testes para novas funcionalidades
- Mantenha a cobertura de testes acima de 80%

### 4. Estrutura de Commits

Siga o padrão [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` para novas funcionalidades
- `fix:` para correções de bugs
- `docs:` para alterações na documentação
- `style:` para formatação, ponto e vírgula faltando, etc
- `refactor:` para refatoração de código
- `test:` para adicionar ou modificar testes
- `chore:` para tarefas de manutenção

### 5. Pull Requests

1. Atualize sua branch com a develop:
```bash
git checkout develop
git pull upstream develop
git checkout feature/nome-da-feature
git rebase develop
```

2. Resolva conflitos se houver

3. Certifique-se que todos os testes passam

4. Atualize a documentação se necessário

5. Abra o PR com uma descrição clara das mudanças

### 6. Revisão de Código

- Responda aos comentários da revisão
- Faça as alterações necessárias
- Mantenha o histórico de commits limpo
- Force push apenas se necessário

## 🧪 Testes

### Testes Unitários
```bash
go test ./...
```

### Testes de Integração
```bash
go test ./... -tags=integration
```

### Cobertura de Testes
```bash
go test ./... -cover
```

## 📚 Documentação

- Atualize o README.md se necessário
- Documente novas funcionalidades
- Atualize a documentação Swagger
- Adicione exemplos de uso

## 🔍 Checklist do Pull Request

- [ ] Código segue as convenções
- [ ] Testes adicionados/atualizados
- [ ] Documentação atualizada
- [ ] Commits seguem o padrão
- [ ] Branch atualizada com develop
- [ ] Todos os testes passam
- [ ] Cobertura de testes mantida
- [ ] PR tem uma descrição clara

## 📝 Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a mesma licença do projeto. 
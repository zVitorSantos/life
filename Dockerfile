# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar dependências de build
RUN apk add --no-cache git

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Instalar dependências de runtime
RUN apk add --no-cache ca-certificates tzdata

# Copiar o binário compilado
COPY --from=builder /app/main .

# Expor a porta da aplicação
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 
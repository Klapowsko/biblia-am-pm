# Script de População do Banco de Dados

Este script popula o banco de dados com o plano de leitura anual (365 dias).

## Uso

### Desenvolvimento (Docker)

```bash
# Popular o banco
docker-compose exec backend go run cmd/populate/main.go

# Limpar e recriar todos os planos
docker-compose exec backend go run cmd/populate/main.go -clear
```

### Local

```bash
# Certifique-se de que as variáveis de ambiente estão configuradas
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=biblia_db

# Popular o banco
go run cmd/populate/main.go

# Limpar e recriar
go run cmd/populate/main.go -clear
```

## Nota

O script atual cria um plano básico de exemplo. Para um plano completo baseado no "Bíblia 365", você precisará substituir a lógica de geração com as referências bíblicas reais do plano.


# Populate Catechism

Comando CLI para popular o banco de dados com as 107 perguntas do Catecismo Menor de Westminster.

## Uso

### Via Makefile (recomendado)

```bash
# Desenvolvimento
make populate-catechism

# Limpar e popular novamente
make populate-catechism-clear

# Produção
make populate-catechism-prod
```

### Diretamente

```bash
cd backend/cmd/populate-catechism
go run . [flags]
```

### Flags

- `-clear`: Limpa todas as perguntas existentes antes de popular
- `-url`: URL customizada para buscar o catecismo (opcional)

### Exemplos

```bash
# Popular normalmente
go run .

# Limpar e popular
go run . -clear

# Usar URL customizada
go run . -url "https://sua-url.com/catechism.json"
```

## Fonte dos Dados

Por padrão, o comando busca os dados de:
- https://raw.githubusercontent.com/ReformedWiki/westminster-shorter-catechism/master/data/catechism.json

## Requisitos

- Banco de dados configurado e acessível
- Variáveis de ambiente do banco de dados configuradas (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)


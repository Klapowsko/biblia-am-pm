# Bíblia AM/PM - Aplicação de Controle de Leitura Bíblica

Aplicação full-stack para controle de leituras bíblicas diárias, com detecção automática de horário (manhã/noite) e exibição das leituras correspondentes do plano anual.

## Tecnologias

- **Backend**: Go (Golang) com `database/sql` padrão
- **Frontend**: React (JavaScript)
- **Banco de Dados**: PostgreSQL
- **Containerização**: Docker e Docker Compose

## Estrutura do Projeto

```
biblia-am-pm/
├── backend/          # API Go
├── frontend/         # Aplicação React
├── docker-compose.yml
├── .env.development
└── .env.production
```

## Como Executar

### Desenvolvimento

1. Copie o arquivo `.env.development` para `.env`:
```bash
cp .env.development .env
```

2. Inicie os serviços com Docker Compose (profile dev):
```bash
docker compose --profile dev up --build
```

3. Acesse a aplicação:
   - Frontend: http://localhost:3001
   - Backend API: http://localhost:8081

### Produção

1. Configure o arquivo `.env.production` com suas credenciais:
```bash
cp .env.production .env
# Edite o arquivo .env com suas configurações de produção
```

2. Inicie os serviços com Docker Compose (profile prod):
```bash
docker compose --profile prod up --build
```

## Popular Banco de Dados

Para popular o banco de dados com o plano de leitura anual (365 dias):

```bash
# Dentro do container do backend
docker compose exec backend go run cmd/populate/main.go

# Ou para limpar e recriar
docker compose exec backend go run cmd/populate/main.go -clear
```

## Endpoints da API

### Autenticação
- `POST /api/auth/register` - Registrar novo usuário
- `POST /api/auth/login` - Login

### Leituras (requer autenticação)
- `GET /api/readings/today` - Buscar leituras do dia atual
- `POST /api/readings/mark-completed` - Marcar leitura como concluída
- `GET /api/progress` - Obter progresso do usuário

## Funcionalidades

- **Detecção automática de horário**: A aplicação detecta se é manhã (6h-12h) ou noite (18h-23h) e exibe as leituras correspondentes
- **Plano anual**: Sistema baseado no dia do ano (1-365), permitindo começar em qualquer data
- **Controle de progresso**: Marque leituras de manhã e noite como concluídas
- **Visualização de progresso**: Acompanhe seu histórico de leituras

## Lógica de Horário

- **Manhã (6h-12h)**: Exibe Antigo Testamento + Salmos
- **Noite (18h-23h)**: Exibe Novo Testamento + Provérbios
- **Outros horários**: Exibe todas as leituras do dia


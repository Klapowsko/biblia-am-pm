.PHONY: help dev prod build-dev build-prod down logs populate clean

# Variáveis
DEV_PROFILE = --profile dev
PROD_PROFILE = --profile prod

# Cores para output
GREEN = \033[0;32m
YELLOW = \033[0;33m
NC = \033[0m # No Color

help: ## Mostra esta mensagem de ajuda
	@echo "$(GREEN)Comandos disponíveis:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

# Desenvolvimento
dev: setup-env-dev build-dev up-dev ## Builda e sobe os containers em modo desenvolvimento

setup-env-dev: ## Configura arquivos .env para desenvolvimento
	@echo "$(GREEN)Setting up development environment files...$(NC)"
	@cp backend/.env.development backend/.env 2>/dev/null || true
	@cp frontend/.env.development frontend/.env 2>/dev/null || true
	@chmod -R u+rwX backend frontend 2>/dev/null || true

build-dev: setup-env-dev ## Builda as imagens Docker para desenvolvimento
	@echo "$(GREEN)Building development images...$(NC)"
	docker compose $(DEV_PROFILE) build

up-dev: setup-env-dev ## Sobe os containers em modo desenvolvimento
	@echo "$(GREEN)Starting development containers...$(NC)"
	docker compose $(DEV_PROFILE) up -d
	@echo "$(GREEN)Containers started!$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3001$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8081$(NC)"

# Produção
prod: setup-env-prod build-prod up-prod ## Builda e sobe os containers em modo produção

setup-env-prod: ## Configura arquivos .env para produção
	@echo "$(GREEN)Setting up production environment files...$(NC)"
	@cp backend/.env.production backend/.env 2>/dev/null || true
	@cp frontend/.env.production frontend/.env 2>/dev/null || true
	@chmod -R u+rwX backend frontend 2>/dev/null || true

build-prod: setup-env-prod ## Builda as imagens Docker para produção
	@echo "$(GREEN)Building production images...$(NC)"
	docker compose $(PROD_PROFILE) build

up-prod: setup-env-prod ## Sobe os containers em modo produção
	@echo "$(GREEN)Starting production containers...$(NC)"
	docker compose $(PROD_PROFILE) up -d
	@echo "$(GREEN)Containers started!$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3001$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8081$(NC)"

# Gerenciamento
down: ## Para e remove os containers
	@echo "$(GREEN)Stopping and removing containers...$(NC)"
	docker compose $(DEV_PROFILE) down
	docker compose $(PROD_PROFILE) down

logs: ## Mostra logs de todos os containers
	docker compose $(DEV_PROFILE) logs -f

# Banco de dados
populate: ## Popula o banco de dados com o plano de leitura (desenvolvimento)
	@echo "$(GREEN)Populating database with reading plan...$(NC)"
	docker compose $(DEV_PROFILE) exec backend go run cmd/populate/main.go

populate-prod: ## Popula o banco de dados com o plano de leitura (produção)
	@echo "$(GREEN)Populating production database with reading plan...$(NC)"
	@NETWORK=$$(docker inspect biblia_postgres --format '{{range $$k, $$v := .NetworkSettings.Networks}}{{$$k}}{{end}}' 2>/dev/null | head -1 || docker network ls --filter name=biblia --format '{{.Name}}' | grep network | head -1 || echo "biblia-am-pm_biblia-network"); \
	echo "$(YELLOW)Using network: $$NETWORK$(NC)"; \
	if ! docker network inspect $$NETWORK >/dev/null 2>&1; then \
		echo "$(YELLOW)Network $$NETWORK not found. Make sure containers are running with 'make prod'$(NC)"; \
		exit 1; \
	fi; \
	docker run --rm \
		--network $$NETWORK \
		-v $$(pwd)/backend:/app \
		-w /app \
		--env-file backend/.env 2>/dev/null || true \
		-e DB_HOST=postgres \
		-e DB_PORT=$${DB_PORT:-5432} \
		-e DB_USER=$${DB_USER:-postgres} \
		-e DB_PASSWORD=$${DB_PASSWORD:-postgres} \
		-e DB_NAME=$${DB_NAME:-biblia_db} \
		golang:alpine \
		sh -c "go mod download && go run cmd/populate/main.go"

# Limpeza
clean: ## Remove containers, volumes e imagens não utilizadas
	@echo "$(GREEN)Cleaning up...$(NC)"
	docker compose $(DEV_PROFILE) down -v
	docker compose $(PROD_PROFILE) down -v
	docker image prune -f

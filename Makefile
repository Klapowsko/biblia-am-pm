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
dev: build-dev up-dev ## Builda e sobe os containers em modo desenvolvimento

build-dev: ## Builda as imagens Docker para desenvolvimento
	@echo "$(GREEN)Building development images...$(NC)"
	docker-compose $(DEV_PROFILE) build

up-dev: ## Sobe os containers em modo desenvolvimento
	@echo "$(GREEN)Starting development containers...$(NC)"
	docker-compose $(DEV_PROFILE) up -d
	@echo "$(GREEN)Containers started!$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3000$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8080$(NC)"

# Produção
prod: build-prod up-prod ## Builda e sobe os containers em modo produção

build-prod: ## Builda as imagens Docker para produção
	@echo "$(GREEN)Building production images...$(NC)"
	docker-compose $(PROD_PROFILE) build

up-prod: ## Sobe os containers em modo produção
	@echo "$(GREEN)Starting production containers...$(NC)"
	docker-compose $(PROD_PROFILE) up -d
	@echo "$(GREEN)Containers started!$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3000$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8080$(NC)"

# Gerenciamento
down: ## Para e remove os containers
	@echo "$(GREEN)Stopping and removing containers...$(NC)"
	docker-compose $(DEV_PROFILE) down
	docker-compose $(PROD_PROFILE) down

logs: ## Mostra logs de todos os containers
	docker-compose $(DEV_PROFILE) logs -f

# Banco de dados
populate: ## Popula o banco de dados com o plano de leitura
	@echo "$(GREEN)Populating database with reading plan...$(NC)"
	docker-compose $(DEV_PROFILE) exec backend go run cmd/populate/main.go

# Limpeza
clean: ## Remove containers, volumes e imagens não utilizadas
	@echo "$(GREEN)Cleaning up...$(NC)"
	docker-compose $(DEV_PROFILE) down -v
	docker-compose $(PROD_PROFILE) down -v
	docker image prune -f

.PHONY: help setup start stop restart status logs clean build test health db-shell qdrant-ui dev

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "$(BLUE)TestPilot AI - Available Commands$(NC)"
	@echo "=================================="
	@echo ""
	@echo "$(GREEN)Setup & Start:$(NC)"
	@echo "  make setup          - First time setup (creates .env, checks prerequisites)"
	@echo "  make start          - Start all services"
	@echo "  make start-infra    - Start only infrastructure (PostgreSQL + Qdrant)"
	@echo "  make start-dev      - Start infrastructure for development"
	@echo ""
	@echo "$(GREEN)Stop & Clean:$(NC)"
	@echo "  make stop           - Stop all services"
	@echo "  make restart        - Restart all services"
	@echo "  make clean          - Stop and remove containers"
	@echo "  make clean-volumes  - Stop and remove containers + volumes (deletes data!)"
	@echo ""
	@echo "$(GREEN)Status & Logs:$(NC)"
	@echo "  make status         - Show status of all services"
	@echo "  make logs           - View logs for all services"
	@echo "  make logs-postgres  - View PostgreSQL logs"
	@echo "  make logs-qdrant    - View Qdrant logs"
	@echo "  make logs-ingestion - View Ingestion service logs"
	@echo "  make health         - Check health of all services"
	@echo ""
	@echo "$(GREEN)Build:$(NC)"
	@echo "  make build          - Build all services"
	@echo "  make build-ingestion- Build ingestion service only"
	@echo "  make rebuild        - Rebuild and restart all services"
	@echo ""
	@echo "$(GREEN)Testing:$(NC)"
	@echo "  make test           - Run all tests"
	@echo "  make test-ingestion - Test ingestion service endpoints"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  make db-shell       - Access PostgreSQL shell"
	@echo "  make db-tables      - Show database tables"
	@echo "  make qdrant-ui      - Open Qdrant dashboard in browser"
	@echo "  make dev            - Start development environment"
	@echo ""
	@echo "$(GREEN)Quick Actions:$(NC)"
	@echo "  make quickstart     - Complete setup + start infrastructure"
	@echo "  make reset          - Clean everything and start fresh"
	@echo ""

## setup: First time setup
setup:
	@echo "$(BLUE)TestPilot AI - First Time Setup$(NC)"
	@echo "================================"
	@echo ""
	@echo "$(YELLOW)1. Checking prerequisites...$(NC)"
	@which docker > /dev/null || (echo "$(RED)❌ Docker not found. Please install Docker first.$(NC)" && exit 1)
	@echo "$(GREEN)✅ Docker found$(NC)"
	@which docker-compose > /dev/null || (echo "$(RED)❌ Docker Compose not found. Please install Docker Compose first.$(NC)" && exit 1)
	@echo "$(GREEN)✅ Docker Compose found$(NC)"
	@echo ""
	@echo "$(YELLOW)2. Setting up environment...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN)✅ Created .env file from template$(NC)"; \
		echo "$(RED)⚠️  IMPORTANT: Edit .env and add your LLM API key!$(NC)"; \
		echo "$(RED)   Required: OPENAI_API_KEY or ANTHROPIC_API_KEY$(NC)"; \
	else \
		echo "$(GREEN)✅ .env file already exists$(NC)"; \
	fi
	@echo ""
	@echo "$(GREEN)Setup complete!$(NC)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Edit .env and add your API key"
	@echo "  2. Run: make quickstart"

## quickstart: Complete setup and start infrastructure
quickstart: setup
	@echo ""
	@echo "$(BLUE)Starting infrastructure services...$(NC)"
	@docker-compose up -d postgres qdrant
	@echo ""
	@echo "$(YELLOW)Waiting for services to be ready (30 seconds)...$(NC)"
	@sleep 30
	@echo ""
	@make status
	@echo ""
	@echo "$(GREEN)🎉 Quick start complete!$(NC)"
	@echo ""
	@echo "Infrastructure is running. Next steps:"
	@echo "  - View logs: make logs-postgres"
	@echo "  - Check status: make status"
	@echo "  - Build service: make build-ingestion"
	@echo "  - Start service: docker-compose up -d ingestion"
	@echo ""

## start: Start all services
start:
	@echo "$(BLUE)Starting all services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✅ All services started$(NC)"
	@echo ""
	@make status

## start-infra: Start only infrastructure
start-infra:
	@echo "$(BLUE)Starting infrastructure (PostgreSQL + Qdrant)...$(NC)"
	@docker-compose up -d postgres qdrant
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	@sleep 10
	@echo "$(GREEN)✅ Infrastructure started$(NC)"

## start-dev: Start development environment
start-dev: start-infra
	@echo ""
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "Infrastructure services running:"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Qdrant: localhost:6333"
	@echo ""
	@echo "Build and start your service:"
	@echo "  make build-ingestion && docker-compose up -d ingestion"

## stop: Stop all services
stop:
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@docker-compose stop
	@echo "$(GREEN)✅ All services stopped$(NC)"

## restart: Restart all services
restart: stop start

## clean: Stop and remove containers
clean:
	@echo "$(YELLOW)Stopping and removing containers...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✅ Containers removed$(NC)"

## clean-volumes: Stop and remove containers + volumes (DANGER: Deletes data!)
clean-volumes:
	@echo "$(RED)⚠️  WARNING: This will delete all data!$(NC)"
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		docker-compose down -v; \
		echo "$(GREEN)✅ Containers and volumes removed$(NC)"; \
	else \
		echo "$(YELLOW)Cancelled.$(NC)"; \
	fi

## reset: Clean everything and start fresh
reset: clean-volumes quickstart

## status: Show status of all services
status:
	@echo "$(BLUE)TestPilot AI - Service Status$(NC)"
	@echo "=============================="
	@echo ""
	@echo "$(YELLOW)Docker Containers:$(NC)"
	@docker-compose ps
	@echo ""
	@echo "$(YELLOW)Port Usage:$(NC)"
	@echo "  Frontend:     localhost:3000"
	@echo "  API Gateway:  localhost:8000"
	@echo "  Ingestion:    localhost:8001"
	@echo "  LLM:          localhost:8002"
	@echo "  Execution:    localhost:8003"
	@echo "  Validation:   localhost:8004"
	@echo "  Query:        localhost:8005"
	@echo "  PostgreSQL:   localhost:5432"
	@echo "  Qdrant:       localhost:6333"

## health: Check health of all services
health:
	@echo "$(BLUE)Health Checks$(NC)"
	@echo "=============="
	@echo ""
	@echo "$(YELLOW)Ingestion Service:$(NC)"
	@curl -s -f http://localhost:8001/health > /dev/null 2>&1 && echo "$(GREEN)✅ Healthy$(NC)" || echo "$(RED)❌ Not responding$(NC)"
	@echo ""
	@echo "$(YELLOW)Qdrant:$(NC)"
	@curl -s -f http://localhost:6333/healthz > /dev/null 2>&1 && echo "$(GREEN)✅ Healthy$(NC)" || echo "$(RED)❌ Not responding$(NC)"
	@echo ""
	@echo "$(YELLOW)PostgreSQL:$(NC)"
	@docker exec testpilot-postgres pg_isready -U testpilot > /dev/null 2>&1 && echo "$(GREEN)✅ Ready$(NC)" || echo "$(RED)❌ Not ready$(NC)"

## logs: View logs for all services
logs:
	@docker-compose logs -f

## logs-postgres: View PostgreSQL logs
logs-postgres:
	@docker-compose logs -f postgres

## logs-qdrant: View Qdrant logs
logs-qdrant:
	@docker-compose logs -f qdrant

## logs-ingestion: View Ingestion service logs
logs-ingestion:
	@docker-compose logs -f ingestion

## logs-llm: View LLM service logs
logs-llm:
	@docker-compose logs -f llm

## logs-gateway: View Gateway logs
logs-gateway:
	@docker-compose logs -f gateway

## build: Build all services
build:
	@echo "$(BLUE)Building all services...$(NC)"
	@docker-compose build
	@echo "$(GREEN)✅ Build complete$(NC)"

## build-ingestion: Build ingestion service
build-ingestion:
	@echo "$(BLUE)Building ingestion service...$(NC)"
	@docker-compose build ingestion
	@echo "$(GREEN)✅ Ingestion service built$(NC)"

## build-llm: Build LLM service
build-llm:
	@echo "$(BLUE)Building LLM service...$(NC)"
	@docker-compose build llm
	@echo "$(GREEN)✅ LLM service built$(NC)"

## rebuild: Rebuild and restart all services
rebuild:
	@echo "$(BLUE)Rebuilding all services...$(NC)"
	@docker-compose up -d --build
	@echo "$(GREEN)✅ All services rebuilt and restarted$(NC)"

## test: Run all tests
test: test-ingestion

## test-ingestion: Test ingestion service endpoints
test-ingestion:
	@echo "$(BLUE)Testing Ingestion Service$(NC)"
	@echo "=========================="
	@echo ""
	@echo "$(YELLOW)1. Health check...$(NC)"
	@curl -s http://localhost:8001/health | python3 -m json.tool && echo "$(GREEN)✅ Passed$(NC)" || echo "$(RED)❌ Failed$(NC)"
	@echo ""
	@echo "$(YELLOW)2. List APIs...$(NC)"
	@curl -s http://localhost:8001/api/v1/apis | python3 -m json.tool | head -20 && echo "$(GREEN)✅ Passed$(NC)" || echo "$(RED)❌ Failed$(NC)"
	@echo ""
	@echo "$(YELLOW)3. Ingestion status...$(NC)"
	@curl -s http://localhost:8001/api/v1/ingest/status | python3 -m json.tool | head -30 && echo "$(GREEN)✅ Passed$(NC)" || echo "$(RED)❌ Failed$(NC)"

## db-shell: Access PostgreSQL shell
db-shell:
	@echo "$(BLUE)Connecting to PostgreSQL...$(NC)"
	@docker exec -it testpilot-postgres psql -U testpilot -d testpilot

## db-tables: Show database tables
db-tables:
	@echo "$(BLUE)Database Tables:$(NC)"
	@docker exec testpilot-postgres psql -U testpilot -d testpilot -c "\dt"

## db-users: Show database users
db-users:
	@echo "$(BLUE)Database Users:$(NC)"
	@docker exec testpilot-postgres psql -U testpilot -d testpilot -c "SELECT username, role, created_at FROM users;"

## qdrant-ui: Open Qdrant dashboard
qdrant-ui:
	@echo "$(BLUE)Opening Qdrant dashboard...$(NC)"
	@open http://localhost:6333/dashboard || xdg-open http://localhost:6333/dashboard || echo "Open http://localhost:6333/dashboard in your browser"

## qdrant-collections: Show Qdrant collections
qdrant-collections:
	@echo "$(BLUE)Qdrant Collections:$(NC)"
	@curl -s http://localhost:6333/collections | python3 -m json.tool

## dev: Interactive development menu
dev:
	@echo "$(BLUE)TestPilot AI - Development Menu$(NC)"
	@echo "================================"
	@echo ""
	@echo "Quick commands:"
	@echo "  make start-dev      - Start infrastructure for development"
	@echo "  make build-ingestion- Build ingestion service"
	@echo "  make logs-ingestion - View ingestion logs"
	@echo "  make test-ingestion - Test ingestion endpoints"
	@echo "  make db-shell       - Access database"
	@echo "  make qdrant-ui      - View Qdrant dashboard"
	@echo ""

## check-ports: Check if required ports are available
check-ports:
	@echo "$(BLUE)Checking port availability...$(NC)"
	@echo ""
	@for port in 3000 5432 6333 8000 8001 8002 8003 8004 8005; do \
		if lsof -Pi :$$port -sTCP:LISTEN -t >/dev/null 2>&1; then \
			echo "$(RED)❌ Port $$port is in use$(NC)"; \
		else \
			echo "$(GREEN)✅ Port $$port is available$(NC)"; \
		fi; \
	done


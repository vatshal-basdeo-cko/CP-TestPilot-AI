# TestPilot AI - Makefile Commands

This project uses a Makefile for all common operations. All commands are documented and easy to use.

## Quick Start

```bash
# First time setup
make setup

# Edit .env and add your API key, then:
make quickstart

# Check everything is running
make status
```

## Available Commands

Run `make help` to see all available commands with descriptions.

### Setup & Start

- `make setup` - First time setup (creates .env, checks prerequisites)
- `make quickstart` - Complete setup + start infrastructure
- `make start` - Start all services
- `make start-infra` - Start only PostgreSQL + Qdrant
- `make start-dev` - Start development environment

### Stop & Clean

- `make stop` - Stop all services
- `make restart` - Restart all services
- `make clean` - Stop and remove containers
- `make clean-volumes` - Remove containers + volumes (deletes data!)
- `make reset` - Clean everything and start fresh

### Status & Logs

- `make status` - Show status of all services
- `make health` - Check health of all services
- `make logs` - View logs for all services
- `make logs-postgres` - View PostgreSQL logs only
- `make logs-qdrant` - View Qdrant logs only
- `make logs-ingestion` - View Ingestion service logs
- `make logs-llm` - View LLM service logs
- `make logs-gateway` - View Gateway logs

### Build

- `make build` - Build all services
- `make build-ingestion` - Build ingestion service only
- `make build-llm` - Build LLM service only
- `make rebuild` - Rebuild and restart all services

### Testing

- `make test` - Run all tests
- `make test-ingestion` - Test ingestion service endpoints

### Development

- `make dev` - Show development menu
- `make db-shell` - Access PostgreSQL shell
- `make db-tables` - Show database tables
- `make db-users` - Show users in database
- `make qdrant-ui` - Open Qdrant dashboard
- `make qdrant-collections` - Show Qdrant collections
- `make check-ports` - Check if required ports are available

## Common Workflows

### First Time Setup

```bash
make setup              # Creates .env, checks Docker
# Edit .env to add your API key
make quickstart         # Starts infrastructure
make status             # Verify everything is running
```

### Daily Development

```bash
make start-dev          # Start infrastructure
make build-ingestion    # Build your service
docker-compose up -d ingestion
make logs-ingestion     # Watch logs
make test-ingestion     # Test endpoints
```

### Troubleshooting

```bash
make health             # Check service health
make logs-postgres      # View database logs
make db-tables          # Check database tables
make check-ports        # Verify ports are available
```

### Clean Restart

```bash
make stop               # Stop services
make clean              # Remove containers
make quickstart         # Start fresh
```

### Complete Reset (Deletes Data!)

```bash
make clean-volumes      # Remove everything including data
make quickstart         # Start from scratch
```

## Environment Variables

Required in `.env`:

```bash
# LLM Provider (at least one required)
OPENAI_API_KEY=sk-...
# OR
ANTHROPIC_API_KEY=sk-ant-...

# Database
POSTGRES_PASSWORD=changeme_in_production

# JWT
JWT_SECRET=change_this_to_a_random_secret
```

## Port Reference

- `3000` - Frontend
- `5432` - PostgreSQL
- `6333` - Qdrant
- `8000` - API Gateway
- `8001` - Ingestion Service
- `8002` - LLM Service
- `8003` - Execution Service
- `8004` - Validation Service
- `8005` - Query Service

## Tips

- Run `make` or `make help` to see all available commands
- All commands have colored output for better readability
- Use `Ctrl+C` to exit log viewing
- Commands check prerequisites before running
- Destructive operations ask for confirmation


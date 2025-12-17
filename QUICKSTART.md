# 🚀 Quick Start Guide

This guide will get TestPilot AI running in 5 minutes.

## Prerequisites

- Docker & Docker Compose installed
- Git installed
- 8GB RAM minimum
- Ports 8000-8005, 5432, 6333, 3000 available

## Step 1: Clone & Configure

```bash
# Clone the repository
git clone <your-repo-url>
cd CP-TestPilot-AI

# Copy environment template
cp .env.example .env

# Edit .env and add your API keys
nano .env  # or use your preferred editor
```

**Required:**
- Set `OPENAI_API_KEY` or `ANTHROPIC_API_KEY`
- Update `POSTGRES_PASSWORD` (production)
- Update `JWT_SECRET` (production)

## Step 2: Start Everything

```bash
# Option 1: Using Makefile (Recommended)
make quickstart

# Option 2: Using Docker Compose directly
docker-compose up -d

# Option 3: Build and start
make build
make start
```

## Step 3: Verify Services

```bash
# Check all services are healthy
make health

# Or check manually
curl http://localhost:8000/health/all
```

**Expected Output:**
```json
{
  "status": "healthy",
  "services": {
    "ingestion": "healthy",
    "llm": "healthy",
    "execution": "healthy",
    "validation": "healthy",
    "query": "healthy"
  }
}
```

## Step 4: Test the API

### Login
```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**Save the token from response!**

### Ingest Sample APIs
```bash
export TOKEN="your-jwt-token-here"

curl -X POST http://localhost:8000/api/v1/ingest/folder \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "folder_path": "/app/api_configs"
  }'
```

### List Ingested APIs
```bash
curl http://localhost:8000/api/v1/ingest/apis \
  -H "Authorization: Bearer $TOKEN"
```

## Step 5: Access Services

| Service | URL | Purpose |
|---------|-----|---------|
| API Gateway | http://localhost:8000 | Main entry point |
| Frontend | http://localhost:3000 | Web UI (when ready) |
| Ingestion | http://localhost:8001 | API config management |
| LLM | http://localhost:8002 | Natural language processing |
| Execution | http://localhost:8003 | API call execution |
| Validation | http://localhost:8004 | Response validation |
| Query | http://localhost:8005 | History & analytics |
| PostgreSQL | localhost:5432 | Database |
| Qdrant | http://localhost:6333 | Vector database |

## Common Commands

```bash
# View logs
make logs                    # All services
make logs SERVICE=ingestion  # Specific service

# Restart a service
make restart-llm

# Stop everything
make stop

# Clean and restart
make clean
make quickstart

# Check service status
make status
```

## Troubleshooting

### Services not starting?
```bash
# Check logs
docker-compose logs <service-name>

# Restart specific service
docker-compose restart <service-name>

# Rebuild if needed
docker-compose build <service-name>
docker-compose up -d <service-name>
```

### Database connection issues?
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check database connection
docker-compose exec postgres psql -U testpilot -d testpilot -c "SELECT 1;"
```

### Port conflicts?
```bash
# Check what's using ports
lsof -i :8000
lsof -i :5432

# Or kill the process
kill -9 <PID>
```

### Out of memory?
```bash
# Increase Docker memory to 8GB minimum
# Docker Desktop → Settings → Resources → Memory

# Or run services individually
docker-compose up postgres qdrant gateway -d
```

## Next Steps

✅ Services running  
✅ APIs ingested  
✅ Authentication working  

**Now you can:**
1. 🧪 Test API execution
2. 📊 View test history
3. 🎨 Build frontend
4. 🤖 Add more APIs

## Need Help?

- Check logs: `make logs`
- Health status: `make health`
- Service status: `make status`
- Read full docs: `README.md`
- Makefile commands: `make help`

---

**Happy Testing! 🚀**





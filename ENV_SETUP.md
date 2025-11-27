# Environment Configuration Guide

## Quick Setup

1. **Copy the template:**
   ```bash
   cp env.example .env
   ```

2. **Edit `.env` with your values:**
   ```bash
   nano .env  # or use your preferred editor
   ```

3. **Required values to update:**
   - `OPENAI_API_KEY` or `ANTHROPIC_API_KEY` (choose your LLM provider)
   - `POSTGRES_PASSWORD` (for production)
   - `JWT_SECRET` (for production)

## Complete Environment Variables

### Database Configuration

```bash
# PostgreSQL
POSTGRES_DB=testpilot
POSTGRES_USER=testpilot
POSTGRES_PASSWORD=changeme_in_production

# Or use full connection string
DATABASE_URL=postgres://testpilot:changeme_in_production@postgres:5432/testpilot?sslmode=disable
```

### Vector Database

```bash
QDRANT_HOST=qdrant
QDRANT_PORT=6333
```

### API Gateway

```bash
GATEWAY_SERVICE_PORT=8000
JWT_SECRET=change_this_to_a_random_secret_key_at_least_32_characters_long
JWT_EXPIRATION_HOURS=24
CORS_ORIGINS=http://localhost:3000,http://localhost:8000
RATE_LIMIT_PER_MINUTE=60
```

### Microservices Ports

```bash
INGESTION_SERVICE_PORT=8001
LLM_SERVICE_PORT=8002
EXECUTION_SERVICE_PORT=8003
VALIDATION_SERVICE_PORT=8004
QUERY_SERVICE_PORT=8005
FRONTEND_PORT=3000
```

### LLM Provider Configuration

**Choose ONE provider and set the appropriate API key:**

```bash
# Provider selection (openai, anthropic, or azure)
DEFAULT_LLM_PROVIDER=openai

# OpenAI
OPENAI_API_KEY=sk-your-openai-api-key-here

# Anthropic Claude
ANTHROPIC_API_KEY=sk-ant-your-anthropic-key-here

# Azure OpenAI
AZURE_OPENAI_API_KEY=your-azure-key-here
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
```

### System Configuration

```bash
# Learning threshold (number of successful tests before learning)
LEARNING_THRESHOLD=5

# History retention (days)
HISTORY_RETENTION_DAYS=90

# Default environment
DEFAULT_ENVIRONMENT=QA1

# API configs path
API_CONFIGS_PATH=/app/api_configs

# Embedding model
EMBEDDING_MODEL=all-MiniLM-L6-v2
```

### Execution Service

```bash
DEFAULT_TIMEOUT=30
MAX_RETRIES=3
REQUEST_TIMEOUT_SECONDS=30
```

### Logging

```bash
# Options: DEBUG, INFO, WARNING, ERROR, CRITICAL
LOG_LEVEL=INFO
LOG_FORMAT=json
```

### Admin User

```bash
# Default admin credentials (change in production!)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123
ADMIN_EMAIL=admin@testpilot.local
```

### Frontend

```bash
VITE_API_GATEWAY_URL=http://localhost:8000
VITE_APP_NAME=TestPilot AI
```

### Security (Production)

```bash
# Session configuration
SESSION_TIMEOUT_MINUTES=30
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION_MINUTES=15

# Generate a secure JWT secret (32+ characters)
JWT_SECRET=$(openssl rand -base64 32)
```

### Development vs Production

```bash
# Options: development, staging, production
ENVIRONMENT=development
DEBUG=false
```

## Getting API Keys

### OpenAI
1. Go to https://platform.openai.com/api-keys
2. Create new secret key
3. Copy and paste into `OPENAI_API_KEY`

### Anthropic Claude
1. Go to https://console.anthropic.com/settings/keys
2. Create new API key
3. Copy and paste into `ANTHROPIC_API_KEY`

### Azure OpenAI
1. Go to Azure Portal
2. Navigate to your Azure OpenAI resource
3. Get keys and endpoint from Keys and Endpoint section
4. Set `AZURE_OPENAI_API_KEY` and `AZURE_OPENAI_ENDPOINT`

## Security Best Practices

### 🔒 For Production

**MUST CHANGE:**
1. `POSTGRES_PASSWORD` - Use strong password (20+ characters)
2. `JWT_SECRET` - Use random string (32+ characters)
3. `ADMIN_PASSWORD` - Change default immediately

**Generate secure secrets:**
```bash
# Generate secure password
openssl rand -base64 32

# Or use this one-liner to update .env
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
```

### 🚫 Never Commit

- `.env` file (contains secrets)
- API keys
- Passwords
- Database credentials

The `.env` file is already in `.gitignore` to prevent accidental commits.

## Validation

After setting up, validate your configuration:

```bash
# Check all required variables are set
make validate-env

# Or manually check
docker-compose config
```

## Troubleshooting

### Missing API Key Error
```
Error: OPENAI_API_KEY not found
```
**Solution:** Set your LLM provider API key in `.env`

### Database Connection Failed
```
Error: could not connect to postgres
```
**Solution:** Check `DATABASE_URL` or PostgreSQL credentials

### Port Already in Use
```
Error: port 8000 already allocated
```
**Solution:** Change port in `.env` or kill process using that port:
```bash
lsof -i :8000
kill -9 <PID>
```

## Example Configurations

### Development (Local Testing)
```bash
ENVIRONMENT=development
DEBUG=true
LOG_LEVEL=DEBUG
DEFAULT_LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-dev-key
```

### Staging
```bash
ENVIRONMENT=staging
DEBUG=false
LOG_LEVEL=INFO
POSTGRES_PASSWORD=strong-staging-password
JWT_SECRET=generated-secret-key-here
```

### Production
```bash
ENVIRONMENT=production
DEBUG=false
LOG_LEVEL=WARNING
POSTGRES_PASSWORD=very-strong-production-password
JWT_SECRET=generated-production-secret-key
ENABLE_RATE_LIMITING=true
RATE_LIMIT_PER_MINUTE=30
```

## Need Help?

- Check `env.example` for all available variables
- See `QUICKSTART.md` for setup instructions
- Read `README.md` for detailed documentation


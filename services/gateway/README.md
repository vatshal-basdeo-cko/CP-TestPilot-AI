# API Gateway

Central API Gateway for TestPilot AI - handles authentication, routing, and request proxying.

## Features

- JWT authentication
- Request routing to backend services
- CORS handling
- Request logging
- Health check aggregation
- Password hashing with bcrypt

## Architecture

The gateway routes requests to backend services:
- `/api/v1/ingest/*` → Ingestion Service (port 8001)
- `/api/v1/parse/*`, `/api/v1/construct/*` → LLM Service (port 8002)
- `/api/v1/execute/*`, `/api/v1/environments/*` → Execution Service (port 8003)
- `/api/v1/validate/*`, `/api/v1/rules/*` → Validation Service (port 8004)
- `/api/v1/history/*`, `/api/v1/analytics/*` → Query Service (port 8005)

## Endpoints

### Authentication (Public)
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### Protected Routes
- `GET /api/v1/auth/me` - Get current user info
- All other `/api/v1/*` routes require JWT token

### Health Checks
- `GET /health` - Gateway health
- `GET /health/all` - All services health

## Usage

### Login
```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Use Protected Endpoints
```bash
curl -X GET http://localhost:8000/api/v1/history \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Environment Variables

- `SERVER_PORT` - Server port (default: 8000)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret (production)


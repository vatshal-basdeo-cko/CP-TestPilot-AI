# Query Service

Query and analytics service for TestPilot AI - provides test history and statistics.

## Features

- Test execution history with filtering
- Search by natural language request
- Analytics and statistics
- Per-API metrics
- Pagination support

## Endpoints

### History
- `GET /api/v1/history` - List test executions (with filters)
- `GET /api/v1/history/:id` - Get specific execution

### Analytics
- `GET /api/v1/analytics/overview` - Overall statistics
- `GET /api/v1/analytics/by-api/:id` - Per-API metrics

### Health
- `GET /health` - Health check

## Query Parameters

**History endpoint:**
- `user_id` - Filter by user
- `api_spec_id` - Filter by API
- `status` - Filter by status (success/failed)
- `search` - Search in natural language requests

**Analytics overview:**
- `start_date` - Start date (RFC3339)
- `end_date` - End date (RFC3339)

## Development

```bash
go run main.go
```


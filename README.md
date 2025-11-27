# TestPilot AI

**Enterprise AI-powered API testing platform that converts natural language requests into executable API calls against QA environments.**

![Architecture](https://img.shields.io/badge/Architecture-Microservices-blue)
![Python](https://img.shields.io/badge/Python-3.11+-green)
![Go](https://img.shields.io/badge/Go-1.21+-cyan)
![C%23](https://img.shields.io/badge/C%23-.NET%208-purple)
![React](https://img.shields.io/badge/React-18-61DAFB)

---

## 🎯 What is TestPilot AI?

TestPilot AI is an intelligent API testing platform designed for development and QA teams. It allows you to:

- **Test APIs using natural language** - "Test Mastercard PTC Authorisation with amount 200"
- **Automatic API request construction** - AI understands your API specifications and builds valid requests
- **Smart validation** - Automatic response validation against schemas and custom rules
- **Learning from success** - System learns from successful tests to improve over time
- **Admin-friendly** - Manage API configurations, users, and settings through a web UI

---

## 🏗️ Architecture Overview

### Polyglot Microservices Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend (React + TypeScript)        │
│                    Test UI │ Admin Panel │ History           │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                    API Gateway (Go)                          │
│              Auth │ Routing │ Rate Limiting                  │
└──┬────────┬────────┬────────┬────────┬────────────────────┬─┘
   │        │        │        │        │                    │
   ▼        ▼        ▼        ▼        ▼                    ▼
┌──────┐┌──────┐┌──────┐┌──────┐┌──────────┐    ┌──────────────┐
│Ingest││ LLM  ││Exec  ││Valid ││  Query   │    │   Qdrant     │
│(Py)  ││ (Py) ││ (Go) ││ (C#) ││  (Go)    │◄───┤Vector Database│
└──┬───┘└──┬───┘└──┬───┘└──┬───┘└────┬─────┘    └──────────────┘
   │       │       │       │         │
   └───────┴───────┴───────┴─────────┘
                   │
          ┌────────▼────────┐
          │   PostgreSQL    │
          │   Database      │
          └─────────────────┘
```

### Services

| Service | Language | Purpose |
|---------|----------|---------|
| **Ingestion** | Python | Ingest API configs from files, Postman collections, Git repos |
| **LLM** | Python | Parse natural language, RAG pipeline, construct API requests |
| **Execution** | Go | Execute API calls to QA environments |
| **Validation** | C# | Validate responses with JSON schema and custom rules |
| **Query** | Go | Test history, search, analytics |
| **Gateway** | Go | API gateway, authentication, routing |
| **Frontend** | React + TS | Web UI with test interface and admin panel |

### Infrastructure

- **Vector Database**: Qdrant (stores API embeddings and learned patterns)
- **SQL Database**: PostgreSQL 15+ (test history, users, configurations)
- **Auth**: JWT + bcrypt
- **Containers**: Docker Compose

---

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- LLM API Key (OpenAI, Anthropic, or Azure OpenAI)
- 8GB RAM minimum
- Ports available: 3000, 8000, 8001-8005, 5432, 6333

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/CP-TestPilot-AI.git
   cd CP-TestPilot-AI
   ```

2. **Quick start (recommended)**
   ```bash
   make quickstart
   ```
   
   This will:
   - Check prerequisites (Docker, Docker Compose)
   - Create .env file
   - Start infrastructure services
   - Initialize databases
   
3. **Or manual setup:**
   ```bash
   # Setup
   make setup
   
   # Edit .env and add your LLM API key
   nano .env
   
   # Start services
   make start
   ```

4. **Verify everything is running**
   ```bash
   make status
   make health
   ```

5. **Access the application**
   - Frontend: http://localhost:3000
   - API Gateway: http://localhost:8000
   - Default credentials: `admin` / `admin123`

**See [MAKEFILE_GUIDE.md](MAKEFILE_GUIDE.md) for all available commands.**

### First Test

1. Log in to http://localhost:3000
2. Navigate to Test Execution
3. Type: `Test Mastercard PTC authorization with amount 200`
4. Click "Execute Test"
5. View constructed request, response, and validation results

---

## 📁 Project Structure

```
CP-TestPilot-AI/
├── services/                   # Microservices
│   ├── ingestion/             # Python - Ingestion service
│   ├── llm/                   # Python - LLM service
│   ├── execution/             # Go - Execution service
│   ├── validation/            # C# - Validation service
│   ├── query/                 # Go - Query service
│   └── gateway/               # Go - API Gateway
├── frontend/                  # React + TypeScript
├── shared/                    # Shared contracts and schemas
│   └── contracts/             # OpenAPI specifications
├── infrastructure/            # Infrastructure configuration
│   ├── postgres/              # Database init scripts
│   └── qdrant/                # Vector DB configuration
├── api_configs/               # Sample API configurations
│   ├── mastercard_ptc.yaml
│   └── payment_api.yaml
├── docker-compose.yml         # Development orchestration
└── README.md                  # This file
```

---

## 🔧 Configuration

### Environment Variables

Key environment variables (see `.env.example` for complete list):

```bash
# Database
POSTGRES_PASSWORD=changeme_in_production

# LLM Provider
OPENAI_API_KEY=sk-...
DEFAULT_LLM_PROVIDER=openai

# JWT
JWT_SECRET=change_this_to_a_random_secret

# System
DEFAULT_LEARNING_THRESHOLD=5
DEFAULT_HISTORY_RETENTION_DAYS=90
```

### Adding API Configurations

Place YAML files in `api_configs/` directory:

```yaml
name: "My API"
version: "1.0.0"
description: "API description"
base_url: "${ENV_BASE_URL}"
endpoints:
  - name: "endpoint_name"
    path: "/api/v1/resource"
    method: "POST"
    # ... see api_configs/README.md for full schema
```

---

## 💡 Usage Examples

### Natural Language Test Examples

```
✅ "Test Mastercard PTC authorization with amount 200 and currency USD"
✅ "Call the payment API to authorize $150"
✅ "Test refund endpoint with transaction ID txn_123 and amount 50"
✅ "Authorize a payment of 1000 EUR using test card"
```

### Admin Tasks

- **Ingest API Configs**: Admin Panel → API Management → Upload YAML or Postman
- **Manage Users**: Admin Panel → User Management → Add User
- **Configure Learning**: Admin Panel → System Config → Set learning threshold
- **View All Tests**: Admin Panel → Test Executions → Filter and search

---

## 🎓 Key Features

### 1. Natural Language Testing
Type test requests in plain English - no need to remember API endpoints, parameters, or formats.

### 2. RAG-Powered Intelligence
Uses Retrieval Augmented Generation to understand your API specifications and construct valid requests.

### 3. Multi-Source Ingestion
- YAML configuration files
- Postman collections
- Git repositories (coming soon)

### 4. Learning System
After N successful tests (configurable), the system learns patterns and improves future request construction.

### 5. Automatic Validation
- JSON schema validation
- Status code validation
- Custom business rule validation

### 6. Complete Admin UI
- Manage API configurations
- User management
- View all test executions
- Configure system settings

---

## 🔒 Security

- JWT-based authentication
- Bcrypt password hashing
- Role-based access control (admin/user)
- Encrypted environment credentials
- Rate limiting on API Gateway

**⚠️ IMPORTANT**: Change default admin password immediately!

---

## 📊 Database Schema

### Key Tables

- `users` - User accounts and roles
- `api_specifications` - Ingested API configurations
- `test_executions` - Test execution history
- `environments` - QA/staging environment configs
- `validation_rules` - Custom validation rules
- `learned_patterns` - AI-learned test patterns
- `system_config` - System configuration

See `infrastructure/postgres/init.sql` for complete schema.

---

## 🛠️ Development

### Running Individual Services

Each service can be run independently for development:

```bash
# Ingestion Service (Python)
cd services/ingestion
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn api.main:app --reload --port 8001

# Execution Service (Go)
cd services/execution
go mod download
go run main.go

# Frontend
cd frontend
npm install
npm run dev
```

### Adding a New Service

1. Create service directory in `services/`
2. Implement clean architecture (domain, application, infrastructure, api)
3. Add Dockerfile
4. Add service to `docker-compose.yml`
5. Update API Gateway routing
6. Document in OpenAPI spec

---

## 📚 Documentation

- **Architecture Deep-Dive**: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- **API Configuration Guide**: [docs/API_CONFIGURATION.md](docs/API_CONFIGURATION.md)
- **Development Guide**: [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)
- **Deployment Guide**: [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)
- **API Reference**: [docs/API_REFERENCE.md](docs/API_REFERENCE.md)
- **Troubleshooting**: [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)

---

## 🧪 Testing

```bash
# Run all tests
docker-compose exec ingestion pytest
docker-compose exec execution go test ./...
docker-compose exec validation dotnet test

# Integration tests
./scripts/run-integration-tests.sh
```

---

## 🗺️ Roadmap

### ✅ MVP (Current)
- Natural language API testing
- File + Postman ingestion
- Basic validation
- Admin UI
- Learning from successful tests

### 🚧 Coming Soon
- Git repository sync
- Multi-step test flows
- Scheduled tests
- CI/CD integration
- Advanced analytics
- Slack/email notifications
- Test comparison and regression detection

---

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines first.

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

---

## 📄 License

[Your License Here]

---

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/your-org/CP-TestPilot-AI/issues)
- **Documentation**: [docs/](docs/)
- **Email**: support@testpilot-ai.com

---

## 🙏 Acknowledgments

Built with:
- [FastAPI](https://fastapi.tiangolo.com/) - Python web framework
- [LangChain](https://python.langchain.com/) - LLM orchestration
- [Gin](https://gin-gonic.com/) - Go web framework
- [ASP.NET Core](https://dotnet.microsoft.com/apps/aspnet) - C# web framework
- [React](https://react.dev/) - Frontend framework
- [Qdrant](https://qdrant.tech/) - Vector database
- [PostgreSQL](https://www.postgresql.org/) - SQL database

---

**Made with ❤️ by the TestPilot AI Team**


"""FastAPI main application."""
import logging
import os
from datetime import datetime
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from api.controllers.ingestion_controller import router as ingestion_router
from api.dtos.ingestion_dtos import HealthResponse
from infrastructure.adapters import (
    FileReaderAdapter,
    PostmanParser,
    EmbeddingService,
    QdrantAdapter,
    PostgresRepository,
)
from application.use_cases import (
    IngestFromFileUseCase,
    IngestFromFolderUseCase,
    IngestPostmanCollectionUseCase,
    GetIngestionStatusUseCase,
)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Global instances (simplified dependency injection)
repository = None
file_use_case = None
folder_use_case = None
postman_use_case = None
status_use_case = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager for startup/shutdown."""
    # Startup
    logger.info("Starting Ingestion Service...")
    
    # Initialize infrastructure
    global repository, file_use_case, folder_use_case, postman_use_case, status_use_case
    
    db_url = os.getenv(
        "DATABASE_URL",
        f"postgresql+asyncpg://{os.getenv('POSTGRES_USER')}:{os.getenv('POSTGRES_PASSWORD')}"
        f"@{os.getenv('POSTGRES_HOST')}:{os.getenv('POSTGRES_PORT')}/{os.getenv('POSTGRES_DB')}"
    )
    
    repository = PostgresRepository(db_url)
    file_reader = FileReaderAdapter()
    postman_parser = PostmanParser()
    embedding_service = EmbeddingService()
    
    qdrant_adapter = QdrantAdapter(
        host=os.getenv("QDRANT_HOST", "localhost"),
        port=int(os.getenv("QDRANT_PORT", "6333"))
    )
    
    # Ensure collections exist
    qdrant_adapter.ensure_collection("api-knowledge")
    
    # Initialize use cases
    file_use_case = IngestFromFileUseCase(
        repository, file_reader, embedding_service, qdrant_adapter
    )
    folder_use_case = IngestFromFolderUseCase(file_use_case)
    postman_use_case = IngestPostmanCollectionUseCase(
        repository, postman_parser, embedding_service, qdrant_adapter
    )
    status_use_case = GetIngestionStatusUseCase(repository)
    
    logger.info("Ingestion Service started successfully")
    
    yield
    
    # Shutdown
    logger.info("Shutting down Ingestion Service...")


# Create FastAPI app
app = FastAPI(
    title="TestPilot AI - Ingestion Service",
    description="API configuration ingestion service",
    version="1.0.0",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=os.getenv("CORS_ORIGINS", "*").split(","),
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(ingestion_router)


@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint."""
    return HealthResponse(
        status="healthy",
        service="ingestion",
        timestamp=datetime.utcnow()
    )


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "service": "TestPilot AI - Ingestion Service",
        "version": "1.0.0",
        "status": "running"
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)


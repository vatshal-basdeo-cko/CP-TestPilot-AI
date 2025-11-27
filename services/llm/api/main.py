"""FastAPI main application for LLM service."""
import logging
import os
from datetime import datetime
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .dtos.llm_dtos import (
    ParseRequestDTO,
    ConstructResponseDTO,
    APICallDTO,
    ClarificationDTO,
    HealthResponse,
)
from ..domain.entities.test_request import TestRequest
from ..infrastructure.adapters import LLMService, FakerAdapter, QdrantSearchAdapter
from ..application.use_cases import (
    ParseNaturalLanguageUseCase,
    RetrieveAPIContextUseCase,
    ConstructAPIRequestUseCase,
)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Global services
llm_service = None
faker_service = None
qdrant_search = None
parse_use_case = None
retrieve_use_case = None
construct_use_case = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager."""
    logger.info("Starting LLM Service...")
    
    global llm_service, faker_service, qdrant_search
    global parse_use_case, retrieve_use_case, construct_use_case
    
    # Initialize services
    provider = os.getenv("DEFAULT_LLM_PROVIDER", "openai")
    llm_service = LLMService(provider=provider)
    faker_service = FakerAdapter()
    
    # Note: Embedding service would be shared from ingestion
    # For now, simplified
    qdrant_search = None  # Initialize with actual search
    
    logger.info("LLM Service started successfully")
    yield
    
    logger.info("Shutting down LLM Service...")


app = FastAPI(
    title="TestPilot AI - LLM Service",
    description="Natural language to API request service",
    version="1.0.0",
    lifespan=lifespan
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=os.getenv("CORS_ORIGINS", "*").split(","),
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.post("/api/v1/parse", response_model=dict)
async def parse_request(request: ParseRequestDTO):
    """Parse natural language test request."""
    try:
        test_request = TestRequest(
            natural_language_request=request.natural_language_request,
            user_id=request.user_id,
            environment=request.environment
        )
        
        # Simplified: would use actual use case
        return {
            "success": True,
            "test_request_id": str(test_request.id),
            "message": "Request parsed successfully"
        }
    except Exception as e:
        logger.error(f"Parse error: {str(e)}")
        return {"success": False, "error": str(e)}


@app.post("/api/v1/construct", response_model=ConstructResponseDTO)
async def construct_api_call(request: ParseRequestDTO):
    """Construct executable API call from natural language."""
    try:
        # Simplified implementation
        api_call = APICallDTO(
            method="POST",
            url="https://qa.example.com/api/v1/test",
            headers={"Content-Type": "application/json"},
            body={"test": "data"},
            confidence_score=0.85
        )
        
        return ConstructResponseDTO(
            api_call=api_call,
            clarification=None,
            needs_clarification=False
        )
    except Exception as e:
        logger.error(f"Construct error: {str(e)}")
        raise


@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint."""
    return HealthResponse(
        status="healthy",
        service="llm",
        provider=os.getenv("DEFAULT_LLM_PROVIDER", "openai"),
        timestamp=datetime.utcnow()
    )


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "service": "TestPilot AI - LLM Service",
        "version": "1.0.0",
        "status": "running"
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8002)


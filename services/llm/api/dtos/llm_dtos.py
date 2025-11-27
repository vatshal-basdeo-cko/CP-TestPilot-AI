"""Pydantic DTOs for LLM service."""
from typing import Optional, Dict, Any, List
from pydantic import BaseModel, Field
from uuid import UUID
from datetime import datetime


class ParseRequestDTO(BaseModel):
    """Request to parse natural language."""
    natural_language_request: str = Field(..., description="Natural language test request")
    user_id: Optional[UUID] = None
    environment: str = Field(default="QA1")


class APICallDTO(BaseModel):
    """Constructed API call response."""
    method: str
    url: str
    headers: Dict[str, str] = {}
    query_params: Dict[str, Any] = {}
    body: Optional[Dict[str, Any]] = None
    api_name: Optional[str] = None
    endpoint_name: Optional[str] = None
    confidence_score: float


class ClarificationDTO(BaseModel):
    """Clarification request response."""
    question: str
    type: str  # 'choice', 'text', 'confirmation'
    options: List[str] = []
    context: Optional[Dict[str, Any]] = None


class ConstructResponseDTO(BaseModel):
    """Response from construct endpoint."""
    api_call: Optional[APICallDTO] = None
    clarification: Optional[ClarificationDTO] = None
    needs_clarification: bool


class HealthResponse(BaseModel):
    """Health check response."""
    status: str
    service: str
    provider: str
    timestamp: datetime


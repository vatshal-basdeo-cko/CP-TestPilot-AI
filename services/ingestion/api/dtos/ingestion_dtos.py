"""Pydantic DTOs for API requests and responses."""
from typing import List, Optional, Dict, Any
from pydantic import BaseModel, Field
from uuid import UUID
from datetime import datetime


class IngestFileRequest(BaseModel):
    """Request to ingest a single file."""
    file_path: str = Field(..., description="Path to configuration file")


class IngestFolderRequest(BaseModel):
    """Request to ingest a folder."""
    folder_path: str = Field(..., description="Path to folder containing configs")
    recursive: bool = Field(default=False, description="Search subdirectories")


class IngestPostmanRequest(BaseModel):
    """Request to ingest Postman collection."""
    collection: Dict[str, Any] = Field(..., description="Postman collection JSON")
    source_path: str = Field(..., description="Source identifier")


class IngestionResultResponse(BaseModel):
    """Response for ingestion operations."""
    id: UUID
    source_type: str
    source_path: str
    status: str
    apis_ingested: int
    api_ids: List[UUID]
    error_message: Optional[str]
    created_at: datetime


class APISpecificationResponse(BaseModel):
    """Response for API specification."""
    id: UUID
    name: str
    version: str
    source_type: str
    source_path: str
    metadata: Dict[str, Any]
    created_at: datetime
    updated_at: datetime


class IngestionStatusResponse(BaseModel):
    """Response for ingestion status."""
    recent_ingestions: List[IngestionResultResponse]
    summary: Dict[str, Any]


class HealthResponse(BaseModel):
    """Health check response."""
    status: str
    service: str
    timestamp: datetime


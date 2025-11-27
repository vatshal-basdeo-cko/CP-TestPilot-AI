"""FastAPI controllers for ingestion endpoints."""
import logging
from datetime import datetime
from typing import List
from fastapi import APIRouter, HTTPException, Depends, UploadFile, File
from uuid import UUID
import json

from ..dtos.ingestion_dtos import (
    IngestFileRequest,
    IngestFolderRequest,
    IngestPostmanRequest,
    IngestionResultResponse,
    APISpecificationResponse,
    IngestionStatusResponse,
)
from ...application.use_cases import (
    IngestFromFileUseCase,
    IngestFromFolderUseCase,
    IngestPostmanCollectionUseCase,
    GetIngestionStatusUseCase,
)
from ...infrastructure.adapters import PostgresRepository

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/v1", tags=["ingestion"])


# Dependency injection (simplified)
def get_repository():
    """Get repository instance."""
    # This will be properly initialized in main.py
    pass


@router.post("/ingest/file", response_model=IngestionResultResponse)
async def ingest_file(request: IngestFileRequest):
    """Ingest a single configuration file."""
    try:
        # Use case will be injected
        use_case = None  # Injected from container
        result = await use_case.execute(request.file_path)
        return IngestionResultResponse(
            id=result.id,
            source_type=result.source_type,
            source_path=result.source_path,
            status=result.status,
            apis_ingested=result.apis_ingested,
            api_ids=result.api_ids,
            error_message=result.error_message,
            created_at=result.created_at
        )
    except Exception as e:
        logger.error(f"Error ingesting file: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/ingest/folder", response_model=IngestionResultResponse)
async def ingest_folder(request: IngestFolderRequest):
    """Ingest all configuration files from a folder."""
    try:
        use_case = None  # Injected
        result = await use_case.execute(request.folder_path, request.recursive)
        return IngestionResultResponse(
            id=result.id,
            source_type=result.source_type,
            source_path=result.source_path,
            status=result.status,
            apis_ingested=result.apis_ingested,
            api_ids=result.api_ids,
            error_message=result.error_message,
            created_at=result.created_at
        )
    except Exception as e:
        logger.error(f"Error ingesting folder: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/ingest/postman", response_model=IngestionResultResponse)
async def ingest_postman(file: UploadFile = File(...)):
    """Ingest a Postman collection file."""
    try:
        content = await file.read()
        collection = json.loads(content)
        
        use_case = None  # Injected
        result = await use_case.execute(
            collection_data=collection,
            source_path=file.filename
        )
        return IngestionResultResponse(
            id=result.id,
            source_type=result.source_type,
            source_path=result.source_path,
            status=result.status,
            apis_ingested=result.apis_ingested,
            api_ids=result.api_ids,
            error_message=result.error_message,
            created_at=result.created_at
        )
    except json.JSONDecodeError:
        raise HTTPException(status_code=400, detail="Invalid JSON file")
    except Exception as e:
        logger.error(f"Error ingesting Postman collection: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/ingest/status", response_model=IngestionStatusResponse)
async def get_ingestion_status():
    """Get ingestion status and history."""
    try:
        use_case = None  # Injected
        recent = await use_case.get_recent_ingestions(limit=20)
        summary = await use_case.get_ingestion_summary()
        
        return IngestionStatusResponse(
            recent_ingestions=[
                IngestionResultResponse(**r) for r in recent
            ],
            summary=summary
        )
    except Exception as e:
        logger.error(f"Error getting status: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/apis", response_model=List[APISpecificationResponse])
async def list_apis(limit: int = 100, offset: int = 0):
    """List all ingested API specifications."""
    try:
        repository = None  # Injected
        specs = await repository.list_all(limit=limit, offset=offset)
        return [
            APISpecificationResponse(
                id=spec.id,
                name=spec.name,
                version=str(spec.version),
                source_type=spec.source_type.value,
                source_path=spec.source_path,
                metadata=spec.metadata,
                created_at=spec.created_at,
                updated_at=spec.updated_at
            )
            for spec in specs
        ]
    except Exception as e:
        logger.error(f"Error listing APIs: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/apis/{api_id}", response_model=APISpecificationResponse)
async def get_api(api_id: UUID):
    """Get a specific API specification."""
    try:
        repository = None  # Injected
        spec = await repository.find_by_id(api_id)
        if not spec:
            raise HTTPException(status_code=404, detail="API not found")
        
        return APISpecificationResponse(
            id=spec.id,
            name=spec.name,
            version=str(spec.version),
            source_type=spec.source_type.value,
            source_path=spec.source_path,
            metadata=spec.metadata,
            created_at=spec.created_at,
            updated_at=spec.updated_at
        )
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error getting API: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@router.delete("/apis/{api_id}")
async def delete_api(api_id: UUID):
    """Delete an API specification."""
    try:
        repository = None  # Injected
        deleted = await repository.delete_api_specification(api_id)
        if not deleted:
            raise HTTPException(status_code=404, detail="API not found")
        return {"message": "API deleted successfully"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error deleting API: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


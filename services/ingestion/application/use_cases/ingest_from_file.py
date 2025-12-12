"""Use case for ingesting API configuration from a single file."""
import hashlib
import logging
from pathlib import Path
from typing import Optional
from uuid import UUID

from domain.entities.api_specification import APISpecification
from domain.entities.ingestion_result import IngestionResult
from domain.value_objects.source_type import SourceType
from domain.value_objects.version import Version
from domain.repositories.ingestion_repository import IngestionRepository

logger = logging.getLogger(__name__)


class IngestFromFileUseCase:
    """Use case for ingesting an API configuration from a file."""
    
    def __init__(
        self,
        repository: IngestionRepository,
        file_reader,  # FileReaderAdapter
        embedding_service,  # EmbeddingService
        vector_store,  # QdrantAdapter
    ):
        self.repository = repository
        self.file_reader = file_reader
        self.embedding_service = embedding_service
        self.vector_store = vector_store
    
    async def execute(
        self,
        file_path: str,
        created_by: Optional[UUID] = None
    ) -> IngestionResult:
        """
        Ingest an API configuration from a file.
        
        Args:
            file_path: Path to the configuration file
            created_by: UUID of the user creating this specification
            
        Returns:
            IngestionResult with status and details
        """
        try:
            logger.info(f"Starting ingestion from file: {file_path}")
            
            # Read file content
            content = await self.file_reader.read_file(file_path)
            
            # Calculate content hash
            content_hash = self._calculate_hash(content)
            
            # Check if already ingested
            existing = await self.repository.find_by_content_hash(content_hash)
            if existing:
                logger.info(f"File {file_path} already ingested (hash match). Skipping.")
                return IngestionResult(
                    source_type=SourceType.FILE.value,
                    source_path=file_path,
                    status="success",
                    apis_ingested=0,
                    api_ids=[existing.id],
                    error_message="File already ingested (no changes detected)"
                )
            
            # Parse configuration
            config = await self.file_reader.parse_config(content, file_path)
            
            # Create API specification
            version = Version.parse(config.get("version", "1.0.0"))
            api_spec = APISpecification(
                name=config["name"],
                version=version,
                source_type=SourceType.FILE,
                source_path=file_path,
                content_hash=content_hash,
                metadata=config,
                created_by=created_by
            )
            
            # Save to database
            saved_spec = await self.repository.save_api_specification(api_spec)
            logger.info(f"Saved API specification: {saved_spec.name} v{saved_spec.version}")
            
            # Generate and store embeddings
            await self._store_embeddings(saved_spec, config)
            
            # Create and save ingestion result
            result = IngestionResult(
                source_type=SourceType.FILE.value,
                source_path=file_path,
                status="success",
                apis_ingested=1,
                api_ids=[saved_spec.id]
            )
            await self.repository.save_ingestion_result(result)
            
            logger.info(f"Successfully ingested API from {file_path}")
            return result
            
        except Exception as e:
            logger.error(f"Error ingesting file {file_path}: {str(e)}", exc_info=True)
            result = IngestionResult(
                source_type=SourceType.FILE.value,
                source_path=file_path,
                status="failed",
                apis_ingested=0,
                error_message=str(e)
            )
            await self.repository.save_ingestion_result(result)
            return result
    
    def _calculate_hash(self, content: str) -> str:
        """Calculate SHA-256 hash of content."""
        return hashlib.sha256(content.encode('utf-8')).hexdigest()
    
    async def _store_embeddings(self, api_spec: APISpecification, config: dict):
        """Generate and store embeddings in vector database."""
        try:
            # Create text representation for embedding
            text_parts = [
                f"API Name: {config['name']}",
                f"Version: {config.get('version', '1.0.0')}",
                f"Description: {config.get('description', '')}",
            ]
            
            # Add endpoint information
            if 'endpoints' in config:
                for endpoint in config['endpoints']:
                    ep_text = f"Endpoint: {endpoint.get('method', 'GET')} {endpoint.get('path', '')}"
                    if 'description' in endpoint:
                        ep_text += f" - {endpoint['description']}"
                    text_parts.append(ep_text)
            
            full_text = "\n".join(text_parts)
            
            # Generate embedding
            embedding = await self.embedding_service.generate_embedding(full_text)
            
            # Store in vector database
            await self.vector_store.store_embedding(
                collection_name="api-knowledge",
                point_id=str(api_spec.id),
                embedding=embedding,
                metadata={
                    "api_name": api_spec.name,
                    "version": str(api_spec.version),
                    "source_type": api_spec.source_type.value,
                    "source_path": api_spec.source_path,
                }
            )
            logger.info(f"Stored embeddings for API: {api_spec.name}")
            
        except Exception as e:
            logger.error(f"Error storing embeddings: {str(e)}")
            # Don't fail the entire ingestion if embeddings fail
            pass


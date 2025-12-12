"""Use case for ingesting Postman collection."""
import hashlib
import json
import logging
from typing import Optional
from uuid import UUID

from domain.entities.api_specification import APISpecification
from domain.entities.ingestion_result import IngestionResult
from domain.value_objects.source_type import SourceType
from domain.value_objects.version import Version
from domain.repositories.ingestion_repository import IngestionRepository

logger = logging.getLogger(__name__)


class IngestPostmanCollectionUseCase:
    """Use case for ingesting a Postman collection."""
    
    def __init__(
        self,
        repository: IngestionRepository,
        postman_parser,  # PostmanParser
        embedding_service,  # EmbeddingService
        vector_store,  # QdrantAdapter
    ):
        self.repository = repository
        self.postman_parser = postman_parser
        self.embedding_service = embedding_service
        self.vector_store = vector_store
    
    async def execute(
        self,
        collection_data: dict,
        source_path: str,
        created_by: Optional[UUID] = None
    ) -> IngestionResult:
        """
        Ingest a Postman collection.
        
        Args:
            collection_data: Parsed Postman collection JSON
            source_path: Source identifier (filename or URL)
            created_by: UUID of the user creating this specification
            
        Returns:
            IngestionResult with status and details
        """
        try:
            logger.info(f"Starting Postman collection ingestion: {source_path}")
            
            # Parse Postman collection to our format
            api_config = await self.postman_parser.parse(collection_data)
            
            # Calculate content hash
            content_str = json.dumps(collection_data, sort_keys=True)
            content_hash = hashlib.sha256(content_str.encode('utf-8')).hexdigest()
            
            # Check if already ingested
            existing = await self.repository.find_by_content_hash(content_hash)
            if existing:
                logger.info(f"Postman collection {source_path} already ingested. Skipping.")
                return IngestionResult(
                    source_type=SourceType.POSTMAN.value,
                    source_path=source_path,
                    status="success",
                    apis_ingested=0,
                    api_ids=[existing.id],
                    error_message="Collection already ingested (no changes detected)"
                )
            
            # Create API specification
            version = Version.parse(api_config.get("version", "1.0.0"))
            api_spec = APISpecification(
                name=api_config["name"],
                version=version,
                source_type=SourceType.POSTMAN,
                source_path=source_path,
                content_hash=content_hash,
                metadata=api_config,
                created_by=created_by
            )
            
            # Save to database
            saved_spec = await self.repository.save_api_specification(api_spec)
            logger.info(f"Saved Postman API: {saved_spec.name} v{saved_spec.version}")
            
            # Generate and store embeddings
            await self._store_embeddings(saved_spec, api_config)
            
            # Create and save ingestion result
            result = IngestionResult(
                source_type=SourceType.POSTMAN.value,
                source_path=source_path,
                status="success",
                apis_ingested=1,
                api_ids=[saved_spec.id]
            )
            await self.repository.save_ingestion_result(result)
            
            logger.info(f"Successfully ingested Postman collection: {source_path}")
            return result
            
        except Exception as e:
            logger.error(f"Error ingesting Postman collection {source_path}: {str(e)}", exc_info=True)
            result = IngestionResult(
                source_type=SourceType.POSTMAN.value,
                source_path=source_path,
                status="failed",
                apis_ingested=0,
                error_message=str(e)
            )
            await self.repository.save_ingestion_result(result)
            return result
    
    async def _store_embeddings(self, api_spec: APISpecification, config: dict):
        """Generate and store embeddings."""
        try:
            # Create text representation
            text_parts = [
                f"API Name: {config['name']}",
                f"Description: {config.get('description', '')}",
            ]
            
            if 'endpoints' in config:
                for endpoint in config['endpoints']:
                    ep_text = f"Endpoint: {endpoint.get('method', 'GET')} {endpoint.get('path', '')}"
                    if 'description' in endpoint:
                        ep_text += f" - {endpoint['description']}"
                    text_parts.append(ep_text)
            
            full_text = "\n".join(text_parts)
            
            # Generate and store embedding
            embedding = await self.embedding_service.generate_embedding(full_text)
            await self.vector_store.store_embedding(
                collection_name="api-knowledge",
                point_id=str(api_spec.id),
                embedding=embedding,
                metadata={
                    "api_name": api_spec.name,
                    "version": str(api_spec.version),
                    "source_type": api_spec.source_type.value,
                }
            )
            logger.info(f"Stored embeddings for Postman API: {api_spec.name}")
            
        except Exception as e:
            logger.error(f"Error storing embeddings: {str(e)}")


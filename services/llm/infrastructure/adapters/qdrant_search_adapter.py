"""Qdrant search adapter for LLM service."""
import logging
from typing import List, Dict, Any
from qdrant_client import QdrantClient

logger = logging.getLogger(__name__)


class QdrantSearchAdapter:
    """Adapter for searching Qdrant vector database."""
    
    def __init__(self, host: str, port: int, embedding_service):
        self.client = QdrantClient(host=host, port=port)
        self.embedding_service = embedding_service
        logger.info(f"Initialized Qdrant search at {host}:{port}")
    
    async def search_apis(
        self,
        query_text: str,
        limit: int = 5,
        filters: Dict[str, Any] = None
    ) -> List[Dict[str, Any]]:
        """Search for relevant APIs."""
        try:
            # Generate embedding for query
            query_embedding = await self.embedding_service.generate_embedding(query_text)
            
            # Search in api-knowledge collection
            results = self.client.search(
                collection_name="api-knowledge",
                query_vector=query_embedding,
                limit=limit,
                query_filter=filters
            )
            
            return [
                {
                    "id": hit.id,
                    "score": hit.score,
                    "metadata": hit.payload
                }
                for hit in results
            ]
        except Exception as e:
            logger.error(f"Qdrant search error: {str(e)}")
            return []
    
    async def search_learned_patterns(
        self,
        api_name: str,
        query_text: str,
        limit: int = 3
    ) -> List[Dict[str, Any]]:
        """Search for learned patterns."""
        try:
            query_embedding = await self.embedding_service.generate_embedding(query_text)
            
            results = self.client.search(
                collection_name="learned-patterns",
                query_vector=query_embedding,
                limit=limit,
                query_filter={"must": [{"key": "api_name", "match": {"value": api_name}}]}
            )
            
            return [
                {
                    "id": hit.id,
                    "score": hit.score,
                    "pattern": hit.payload
                }
                for hit in results
            ]
        except Exception as e:
            logger.error(f"Pattern search error: {str(e)}")
            return []


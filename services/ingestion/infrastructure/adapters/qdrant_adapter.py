"""Qdrant vector database adapter."""
import logging
from typing import List, Dict, Any, Optional
from qdrant_client import QdrantClient
from qdrant_client.models import PointStruct, Distance, VectorParams

logger = logging.getLogger(__name__)


class QdrantAdapter:
    """Adapter for Qdrant vector database operations."""
    
    def __init__(self, host: str = "localhost", port: int = 6333):
        self.client = QdrantClient(host=host, port=port)
        logger.info(f"Connected to Qdrant at {host}:{port}")
    
    async def store_embedding(
        self,
        collection_name: str,
        point_id: str,
        embedding: List[float],
        metadata: Dict[str, Any]
    ) -> None:
        """Store an embedding vector with metadata."""
        try:
            point = PointStruct(
                id=point_id,
                vector=embedding,
                payload=metadata
            )
            self.client.upsert(collection_name=collection_name, points=[point])
            logger.debug(f"Stored embedding {point_id} in {collection_name}")
        except Exception as e:
            logger.error(f"Error storing embedding: {str(e)}")
            raise
    
    async def search(
        self,
        collection_name: str,
        query_vector: List[float],
        limit: int = 5,
        filter_conditions: Optional[Dict] = None
    ) -> List[Dict[str, Any]]:
        """Search for similar vectors."""
        try:
            results = self.client.search(
                collection_name=collection_name,
                query_vector=query_vector,
                limit=limit,
                query_filter=filter_conditions
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
            logger.error(f"Error searching: {str(e)}")
            return []
    
    def ensure_collection(self, collection_name: str, vector_size: int = 384):
        """Ensure collection exists."""
        try:
            collections = self.client.get_collections().collections
            if not any(c.name == collection_name for c in collections):
                self.client.create_collection(
                    collection_name=collection_name,
                    vectors_config=VectorParams(size=vector_size, distance=Distance.COSINE)
                )
                logger.info(f"Created collection: {collection_name}")
        except Exception as e:
            logger.warning(f"Collection setup: {str(e)}")


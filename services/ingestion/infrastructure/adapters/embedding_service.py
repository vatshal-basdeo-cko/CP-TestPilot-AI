"""Embedding service using OpenAI API."""
import logging
import os
from typing import List
from openai import OpenAI

logger = logging.getLogger(__name__)


class EmbeddingService:
    """Service for generating text embeddings using OpenAI API."""
    
    def __init__(self, model_name: str = "text-embedding-3-small"):
        """Initialize with OpenAI embedding model."""
        api_key = os.getenv("OPENAI_API_KEY")
        if not api_key:
            logger.warning("OPENAI_API_KEY not set, embeddings will fail")
        self.client = OpenAI(api_key=api_key) if api_key else None
        self.model = model_name
        logger.info(f"Initialized OpenAI embedding service with model: {model_name}")
    
    async def generate_embedding(self, text: str) -> List[float]:
        """Generate embedding vector for text using OpenAI API."""
        if not self.client:
            logger.error("OpenAI client not initialized - missing API key")
            return [0.0] * 1536  # Return zero vector as fallback
        
        try:
            response = self.client.embeddings.create(
                model=self.model,
                input=text
            )
            return response.data[0].embedding
        except Exception as e:
            logger.error(f"Error generating embedding: {str(e)}")
            return [0.0] * 1536  # Return zero vector as fallback
    
    async def generate_embeddings(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings for multiple texts using OpenAI API."""
        if not self.client:
            logger.error("OpenAI client not initialized - missing API key")
            return [[0.0] * 1536 for _ in texts]
        
        try:
            response = self.client.embeddings.create(
                model=self.model,
                input=texts
            )
            return [item.embedding for item in response.data]
        except Exception as e:
            logger.error(f"Error generating embeddings: {str(e)}")
            return [[0.0] * 1536 for _ in texts]

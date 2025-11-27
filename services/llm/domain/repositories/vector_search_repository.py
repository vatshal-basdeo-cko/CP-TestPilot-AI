"""Repository interface for vector search operations."""
from abc import ABC, abstractmethod
from typing import List, Dict, Any


class VectorSearchRepository(ABC):
    """
    Abstract repository interface for vector search operations.
    
    This interface defines the contract for searching API knowledge
    in the vector database.
    """
    
    @abstractmethod
    async def search_apis(
        self,
        query_text: str,
        limit: int = 5,
        filters: Dict[str, Any] = None
    ) -> List[Dict[str, Any]]:
        """
        Search for relevant APIs using semantic search.
        
        Args:
            query_text: Natural language query
            limit: Maximum number of results
            filters: Optional metadata filters
            
        Returns:
            List of search results with relevance scores
        """
        pass
    
    @abstractmethod
    async def search_learned_patterns(
        self,
        api_name: str,
        query_text: str,
        limit: int = 3
    ) -> List[Dict[str, Any]]:
        """
        Search for learned patterns for a specific API.
        
        Args:
            api_name: Name of the API
            query_text: Natural language query
            limit: Maximum number of results
            
        Returns:
            List of learned patterns
        """
        pass


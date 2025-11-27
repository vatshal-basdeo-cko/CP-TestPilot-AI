"""Repository interface for ingestion operations."""
from abc import ABC, abstractmethod
from typing import List, Optional
from uuid import UUID

from ..entities.api_specification import APISpecification
from ..entities.ingestion_result import IngestionResult


class IngestionRepository(ABC):
    """
    Abstract repository interface for ingestion operations.
    
    This interface defines the contract for persisting and retrieving
    API specifications and ingestion results. Implementations will handle
    the actual database operations.
    """
    
    @abstractmethod
    async def save_api_specification(
        self, 
        api_spec: APISpecification
    ) -> APISpecification:
        """
        Save an API specification to the repository.
        
        Args:
            api_spec: API specification to save
            
        Returns:
            Saved API specification with any auto-generated fields
        """
        pass
    
    @abstractmethod
    async def find_by_name_and_version(
        self, 
        name: str, 
        version: str
    ) -> Optional[APISpecification]:
        """
        Find an API specification by name and version.
        
        Args:
            name: API name
            version: API version
            
        Returns:
            API specification if found, None otherwise
        """
        pass
    
    @abstractmethod
    async def find_by_content_hash(
        self, 
        content_hash: str
    ) -> Optional[APISpecification]:
        """
        Find an API specification by content hash.
        
        Args:
            content_hash: Content hash to search for
            
        Returns:
            API specification if found, None otherwise
        """
        pass
    
    @abstractmethod
    async def find_by_id(self, api_spec_id: UUID) -> Optional[APISpecification]:
        """
        Find an API specification by ID.
        
        Args:
            api_spec_id: API specification UUID
            
        Returns:
            API specification if found, None otherwise
        """
        pass
    
    @abstractmethod
    async def list_all(
        self, 
        limit: int = 100, 
        offset: int = 0
    ) -> List[APISpecification]:
        """
        List all API specifications with pagination.
        
        Args:
            limit: Maximum number of results
            offset: Offset for pagination
            
        Returns:
            List of API specifications
        """
        pass
    
    @abstractmethod
    async def update_api_specification(
        self, 
        api_spec: APISpecification
    ) -> APISpecification:
        """
        Update an existing API specification.
        
        Args:
            api_spec: API specification with updated values
            
        Returns:
            Updated API specification
        """
        pass
    
    @abstractmethod
    async def delete_api_specification(self, api_spec_id: UUID) -> bool:
        """
        Delete an API specification.
        
        Args:
            api_spec_id: API specification UUID to delete
            
        Returns:
            True if deleted, False if not found
        """
        pass
    
    @abstractmethod
    async def save_ingestion_result(
        self, 
        result: IngestionResult
    ) -> IngestionResult:
        """
        Save an ingestion result to track ingestion history.
        
        Args:
            result: Ingestion result to save
            
        Returns:
            Saved ingestion result
        """
        pass
    
    @abstractmethod
    async def get_recent_ingestions(
        self, 
        limit: int = 10
    ) -> List[IngestionResult]:
        """
        Get recent ingestion results.
        
        Args:
            limit: Maximum number of results
            
        Returns:
            List of recent ingestion results
        """
        pass


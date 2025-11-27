"""Repository interface for LLM operations."""
from abc import ABC, abstractmethod
from typing import Optional
from uuid import UUID

from ..entities.test_request import TestRequest


class LLMRepository(ABC):
    """
    Abstract repository interface for LLM operations.
    
    This interface defines the contract for persisting test requests
    and related data.
    """
    
    @abstractmethod
    async def save_test_request(self, test_request: TestRequest) -> TestRequest:
        """Save a test request."""
        pass
    
    @abstractmethod
    async def find_test_request_by_id(self, request_id: UUID) -> Optional[TestRequest]:
        """Find a test request by ID."""
        pass
    
    @abstractmethod
    async def increment_api_success_count(self, api_spec_id: UUID) -> int:
        """Increment success count for an API and return new count."""
        pass
    
    @abstractmethod
    async def get_learning_threshold(self) -> int:
        """Get the learning threshold from system config."""
        pass


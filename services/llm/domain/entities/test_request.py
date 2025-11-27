"""Domain entity for Test Request."""
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional, Dict, Any
from uuid import UUID, uuid4


@dataclass
class TestRequest:
    """
    Test request entity representing a user's natural language test request.
    
    This is the core domain model for test requests that need to be parsed
    and converted into executable API calls.
    """
    
    natural_language_request: str
    user_id: Optional[UUID] = None
    api_name: Optional[str] = None
    environment: str = "QA1"
    id: UUID = field(default_factory=uuid4)
    created_at: datetime = field(default_factory=datetime.utcnow)
    
    def __post_init__(self):
        """Validate the entity after initialization."""
        if not self.natural_language_request or not self.natural_language_request.strip():
            raise ValueError("Natural language request cannot be empty")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert entity to dictionary representation."""
        return {
            "id": str(self.id),
            "natural_language_request": self.natural_language_request,
            "user_id": str(self.user_id) if self.user_id else None,
            "api_name": self.api_name,
            "environment": self.environment,
            "created_at": self.created_at.isoformat(),
        }


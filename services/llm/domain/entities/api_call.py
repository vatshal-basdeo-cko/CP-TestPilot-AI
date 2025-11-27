"""Domain entity for API Call."""
from dataclasses import dataclass, field
from typing import Dict, Any, Optional
from uuid import UUID, uuid4


@dataclass
class APICall:
    """
    API Call entity representing a constructed API request.
    
    This entity contains all the information needed to execute
    an API call against a target environment.
    """
    
    method: str
    url: str
    headers: Dict[str, str] = field(default_factory=dict)
    query_params: Dict[str, Any] = field(default_factory=dict)
    body: Optional[Dict[str, Any]] = None
    api_spec_id: Optional[UUID] = None
    api_name: Optional[str] = None
    endpoint_name: Optional[str] = None
    confidence_score: float = 1.0
    id: UUID = field(default_factory=uuid4)
    
    def __post_init__(self):
        """Validate the entity after initialization."""
        valid_methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']
        if self.method.upper() not in valid_methods:
            raise ValueError(f"Method must be one of {valid_methods}")
        
        if not self.url:
            raise ValueError("URL cannot be empty")
        
        self.method = self.method.upper()
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert entity to dictionary representation."""
        return {
            "id": str(self.id),
            "method": self.method,
            "url": self.url,
            "headers": self.headers,
            "query_params": self.query_params,
            "body": self.body,
            "api_spec_id": str(self.api_spec_id) if self.api_spec_id else None,
            "api_name": self.api_name,
            "endpoint_name": self.endpoint_name,
            "confidence_score": self.confidence_score,
        }


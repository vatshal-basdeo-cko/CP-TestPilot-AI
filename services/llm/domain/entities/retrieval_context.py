"""Domain entity for Retrieval Context."""
from dataclasses import dataclass, field
from typing import List, Dict, Any
from uuid import UUID


@dataclass
class RetrievalContext:
    """
    Retrieval context entity representing API knowledge retrieved from vector database.
    
    This entity contains the API specifications and examples retrieved
    through semantic search for constructing API requests.
    """
    
    api_spec_id: UUID
    api_name: str
    api_version: str
    relevance_score: float
    api_config: Dict[str, Any]
    matched_endpoints: List[Dict[str, Any]] = field(default_factory=list)
    examples: List[Dict[str, Any]] = field(default_factory=list)
    
    def __post_init__(self):
        """Validate the entity."""
        if self.relevance_score < 0 or self.relevance_score > 1:
            raise ValueError("Relevance score must be between 0 and 1")
    
    def get_endpoint_by_name(self, name: str) -> Dict[str, Any]:
        """Get endpoint configuration by name."""
        for endpoint in self.api_config.get('endpoints', []):
            if endpoint.get('name') == name:
                return endpoint
        return {}
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary representation."""
        return {
            "api_spec_id": str(self.api_spec_id),
            "api_name": self.api_name,
            "api_version": self.api_version,
            "relevance_score": self.relevance_score,
            "api_config": self.api_config,
            "matched_endpoints": self.matched_endpoints,
            "examples": self.examples,
        }


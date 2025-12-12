"""Domain entity for API Specification."""
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional, Dict, Any
from uuid import UUID, uuid4

from domain.value_objects.version import Version
from domain.value_objects.source_type import SourceType


@dataclass
class APISpecification:
    """
    API Specification entity representing an API configuration.
    
    This is the core domain model for API configurations ingested from
    various sources (files, Postman collections, Git repos, etc.).
    """
    
    name: str
    version: Version
    source_type: SourceType
    source_path: str
    content_hash: str
    metadata: Dict[str, Any]
    id: UUID = field(default_factory=uuid4)
    created_at: datetime = field(default_factory=datetime.utcnow)
    updated_at: datetime = field(default_factory=datetime.utcnow)
    created_by: Optional[UUID] = None
    
    def __post_init__(self):
        """Validate the entity after initialization."""
        if not self.name:
            raise ValueError("API name cannot be empty")
        if not self.source_path:
            raise ValueError("Source path cannot be empty")
        if not self.content_hash:
            raise ValueError("Content hash cannot be empty")
    
    def update_metadata(self, new_metadata: Dict[str, Any]) -> None:
        """Update the metadata and refresh updated_at timestamp."""
        self.metadata.update(new_metadata)
        self.updated_at = datetime.utcnow()
    
    def update_version(self, new_version: Version) -> None:
        """Update the version and refresh updated_at timestamp."""
        self.version = new_version
        self.updated_at = datetime.utcnow()
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert entity to dictionary representation."""
        return {
            "id": str(self.id),
            "name": self.name,
            "version": str(self.version),
            "source_type": self.source_type.value,
            "source_path": self.source_path,
            "content_hash": self.content_hash,
            "metadata": self.metadata,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
            "created_by": str(self.created_by) if self.created_by else None,
        }
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "APISpecification":
        """Create entity from dictionary representation."""
        return cls(
            id=UUID(data["id"]) if isinstance(data["id"], str) else data["id"],
            name=data["name"],
            version=Version.parse(data["version"]),
            source_type=SourceType(data["source_type"]),
            source_path=data["source_path"],
            content_hash=data["content_hash"],
            metadata=data["metadata"],
            created_at=datetime.fromisoformat(data["created_at"]) if isinstance(data["created_at"], str) else data["created_at"],
            updated_at=datetime.fromisoformat(data["updated_at"]) if isinstance(data["updated_at"], str) else data["updated_at"],
            created_by=UUID(data["created_by"]) if data.get("created_by") else None,
        )


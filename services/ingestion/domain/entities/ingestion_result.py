"""Domain entity for Ingestion Result."""
from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional
from uuid import UUID, uuid4


@dataclass
class IngestionResult:
    """
    Represents the result of an ingestion operation.
    
    Tracks success/failure status, ingested APIs, and any errors encountered.
    """
    
    source_type: str
    source_path: str
    status: str  # 'success', 'failed', 'partial'
    apis_ingested: int = 0
    api_ids: List[UUID] = field(default_factory=list)
    error_message: Optional[str] = None
    id: UUID = field(default_factory=uuid4)
    created_at: datetime = field(default_factory=datetime.utcnow)
    
    def __post_init__(self):
        """Validate status value."""
        valid_statuses = ['success', 'failed', 'partial']
        if self.status not in valid_statuses:
            raise ValueError(f"Status must be one of {valid_statuses}")
    
    def is_successful(self) -> bool:
        """Check if ingestion was fully successful."""
        return self.status == 'success'
    
    def has_errors(self) -> bool:
        """Check if ingestion had errors."""
        return self.error_message is not None
    
    def to_dict(self):
        """Convert to dictionary representation."""
        return {
            "id": str(self.id),
            "source_type": self.source_type,
            "source_path": self.source_path,
            "status": self.status,
            "apis_ingested": self.apis_ingested,
            "api_ids": [str(api_id) for api_id in self.api_ids],
            "error_message": self.error_message,
            "created_at": self.created_at.isoformat(),
        }


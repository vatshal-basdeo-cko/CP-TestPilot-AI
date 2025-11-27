"""Value object for Source Type."""
from enum import Enum


class SourceType(str, Enum):
    """
    Enumeration of supported ingestion source types.
    
    Represents where API configurations can be ingested from.
    """
    
    FILE = "file"
    POSTMAN = "postman"
    GIT = "git"
    URL = "url"
    
    def __str__(self) -> str:
        """String representation."""
        return self.value
    
    @classmethod
    def from_string(cls, value: str) -> "SourceType":
        """Create SourceType from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            valid_types = [t.value for t in cls]
            raise ValueError(f"Invalid source type: {value}. Must be one of {valid_types}")
    
    def is_file_based(self) -> bool:
        """Check if source type is file-based."""
        return self in [SourceType.FILE, SourceType.POSTMAN]
    
    def is_remote(self) -> bool:
        """Check if source type is remote."""
        return self in [SourceType.GIT, SourceType.URL]


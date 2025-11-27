"""Domain entity for Clarification."""
from dataclasses import dataclass, field
from typing import List, Optional, Dict, Any
from uuid import UUID, uuid4


@dataclass
class Clarification:
    """
    Clarification entity representing a request for user input.
    
    Used when the LLM needs additional information to construct
    a valid API request.
    """
    
    question: str
    clarification_type: str  # 'choice', 'text', 'confirmation'
    options: List[str] = field(default_factory=list)
    context: Optional[Dict[str, Any]] = None
    id: UUID = field(default_factory=uuid4)
    
    def __post_init__(self):
        """Validate the entity."""
        valid_types = ['choice', 'text', 'confirmation']
        if self.clarification_type not in valid_types:
            raise ValueError(f"Type must be one of {valid_types}")
        
        if self.clarification_type == 'choice' and not self.options:
            raise ValueError("Options required for choice type clarification")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary representation."""
        return {
            "id": str(self.id),
            "question": self.question,
            "type": self.clarification_type,
            "options": self.options,
            "context": self.context,
        }


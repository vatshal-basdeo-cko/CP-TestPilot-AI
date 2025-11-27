"""Value object for Version."""
import re
from dataclasses import dataclass
from typing import Tuple


@dataclass(frozen=True)
class Version:
    """
    Semantic version value object.
    
    Immutable representation of a version following semantic versioning (major.minor.patch).
    """
    
    major: int
    minor: int
    patch: int
    
    def __post_init__(self):
        """Validate version components."""
        if self.major < 0 or self.minor < 0 or self.patch < 0:
            raise ValueError("Version components must be non-negative")
    
    def __str__(self) -> str:
        """String representation in semantic versioning format."""
        return f"{self.major}.{self.minor}.{self.patch}"
    
    def __lt__(self, other: "Version") -> bool:
        """Compare versions for ordering."""
        if not isinstance(other, Version):
            return NotImplemented
        return self.to_tuple() < other.to_tuple()
    
    def __le__(self, other: "Version") -> bool:
        """Less than or equal comparison."""
        return self == other or self < other
    
    def __gt__(self, other: "Version") -> bool:
        """Greater than comparison."""
        if not isinstance(other, Version):
            return NotImplemented
        return self.to_tuple() > other.to_tuple()
    
    def __ge__(self, other: "Version") -> bool:
        """Greater than or equal comparison."""
        return self == other or self > other
    
    def to_tuple(self) -> Tuple[int, int, int]:
        """Convert to tuple for comparison."""
        return (self.major, self.minor, self.patch)
    
    @classmethod
    def parse(cls, version_string: str) -> "Version":
        """
        Parse a version string into a Version object.
        
        Supports formats:
        - "1.0.0" (semantic version)
        - "1.0" (assumes patch=0)
        - "1" (assumes minor=0, patch=0)
        - "v1.0.0" (with 'v' prefix)
        
        Args:
            version_string: Version string to parse
            
        Returns:
            Version object
            
        Raises:
            ValueError: If version string is invalid
        """
        if not version_string:
            raise ValueError("Version string cannot be empty")
        
        # Remove 'v' prefix if present
        version_string = version_string.strip()
        if version_string.lower().startswith('v'):
            version_string = version_string[1:]
        
        # Parse version components
        pattern = r'^(\d+)(?:\.(\d+))?(?:\.(\d+))?$'
        match = re.match(pattern, version_string)
        
        if not match:
            raise ValueError(f"Invalid version format: {version_string}")
        
        major = int(match.group(1))
        minor = int(match.group(2)) if match.group(2) else 0
        patch = int(match.group(3)) if match.group(3) else 0
        
        return cls(major=major, minor=minor, patch=patch)
    
    def increment_major(self) -> "Version":
        """Create new version with incremented major version."""
        return Version(major=self.major + 1, minor=0, patch=0)
    
    def increment_minor(self) -> "Version":
        """Create new version with incremented minor version."""
        return Version(major=self.major, minor=self.minor + 1, patch=0)
    
    def increment_patch(self) -> "Version":
        """Create new version with incremented patch version."""
        return Version(major=self.major, minor=self.minor, patch=self.patch + 1)


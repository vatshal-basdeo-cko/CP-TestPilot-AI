"""Use case for detecting changes in API configurations."""
import hashlib
import logging
from typing import List, Tuple

from domain.entities.api_specification import APISpecification
from domain.repositories.ingestion_repository import IngestionRepository

logger = logging.getLogger(__name__)


class DetectChangesUseCase:
    """Use case for detecting changes in API configurations."""
    
    def __init__(self, repository: IngestionRepository):
        self.repository = repository
    
    async def has_changed(
        self, 
        name: str, 
        version: str, 
        new_content: str
    ) -> Tuple[bool, APISpecification | None]:
        """
        Check if content has changed compared to existing specification.
        
        Args:
            name: API name
            version: API version
            new_content: New content to compare
            
        Returns:
            Tuple of (has_changed, existing_spec_or_None)
        """
        # Calculate hash of new content
        new_hash = self._calculate_hash(new_content)
        
        # Find existing specification
        existing = await self.repository.find_by_name_and_version(name, version)
        
        if not existing:
            # No existing spec, this is new
            return (True, None)
        
        # Compare hashes
        changed = existing.content_hash != new_hash
        
        if changed:
            logger.info(f"Detected changes in API {name} v{version}")
        else:
            logger.debug(f"No changes detected in API {name} v{version}")
        
        return (changed, existing)
    
    async def find_changed_files(
        self, 
        file_contents: List[Tuple[str, str, str]]
    ) -> List[Tuple[str, bool, APISpecification | None]]:
        """
        Check multiple files for changes.
        
        Args:
            file_contents: List of (name, version, content) tuples
            
        Returns:
            List of (name, has_changed, existing_spec) tuples
        """
        results = []
        
        for name, version, content in file_contents:
            try:
                changed, existing = await self.has_changed(name, version, content)
                results.append((name, changed, existing))
            except Exception as e:
                logger.error(f"Error checking changes for {name}: {str(e)}")
                results.append((name, True, None))  # Assume changed on error
        
        return results
    
    def _calculate_hash(self, content: str) -> str:
        """Calculate SHA-256 hash of content."""
        return hashlib.sha256(content.encode('utf-8')).hexdigest()


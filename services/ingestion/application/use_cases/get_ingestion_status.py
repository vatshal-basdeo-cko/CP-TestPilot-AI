"""Use case for getting ingestion status."""
import logging
from typing import List, Dict, Any

from domain.entities.ingestion_result import IngestionResult
from domain.repositories.ingestion_repository import IngestionRepository

logger = logging.getLogger(__name__)


class GetIngestionStatusUseCase:
    """Use case for retrieving ingestion status and history."""
    
    def __init__(self, repository: IngestionRepository):
        self.repository = repository
    
    async def get_recent_ingestions(self, limit: int = 10) -> List[Dict[str, Any]]:
        """
        Get recent ingestion results.
        
        Args:
            limit: Maximum number of results to return
            
        Returns:
            List of ingestion result dictionaries
        """
        try:
            results = await self.repository.get_recent_ingestions(limit)
            return [result.to_dict() for result in results]
        except Exception as e:
            logger.error(f"Error retrieving ingestion history: {str(e)}")
            return []
    
    async def get_ingestion_summary(self) -> Dict[str, Any]:
        """
        Get summary statistics for ingestions.
        
        Returns:
            Dictionary with summary statistics
        """
        try:
            recent = await self.repository.get_recent_ingestions(100)
            
            total = len(recent)
            successful = sum(1 for r in recent if r.is_successful())
            failed = sum(1 for r in recent if r.status == 'failed')
            partial = sum(1 for r in recent if r.status == 'partial')
            total_apis = sum(r.apis_ingested for r in recent)
            
            return {
                "total_ingestions": total,
                "successful": successful,
                "failed": failed,
                "partial": partial,
                "total_apis_ingested": total_apis,
                "success_rate": (successful / total * 100) if total > 0 else 0
            }
        except Exception as e:
            logger.error(f"Error generating ingestion summary: {str(e)}")
            return {
                "total_ingestions": 0,
                "successful": 0,
                "failed": 0,
                "partial": 0,
                "total_apis_ingested": 0,
                "success_rate": 0
            }


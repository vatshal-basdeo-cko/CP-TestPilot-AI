"""Use case for parsing natural language requests."""
import logging
from typing import Dict, Any

from domain.entities.test_request import TestRequest

logger = logging.getLogger(__name__)


class ParseNaturalLanguageUseCase:
    """Use case for parsing natural language test requests."""
    
    def __init__(self, llm_service, vector_search):
        self.llm_service = llm_service
        self.vector_search = vector_search
    
    async def execute(self, test_request: TestRequest) -> Dict[str, Any]:
        """
        Parse natural language request and extract intent.
        
        Returns:
            Dictionary with extracted information
        """
        try:
            # Search for relevant APIs
            search_results = await self.vector_search.search_apis(
                query_text=test_request.natural_language_request,
                limit=3
            )
            
            if not search_results:
                return {
                    "success": False,
                    "error": "No matching APIs found"
                }
            
            # Use LLM to parse request
            parsed = await self.llm_service.parse_request(
                natural_language=test_request.natural_language_request,
                api_context=search_results
            )
            
            return {
                "success": True,
                "parsed_request": parsed,
                "matched_apis": search_results
            }
            
        except Exception as e:
            logger.error(f"Error parsing request: {str(e)}")
            return {
                "success": False,
                "error": str(e)
            }


"""Use case for retrieving API context via RAG."""
import logging
from typing import List
from uuid import UUID

from ...domain.entities.retrieval_context import RetrievalContext

logger = logging.getLogger(__name__)


class RetrieveAPIContextUseCase:
    """Use case for retrieving API context using RAG."""
    
    def __init__(self, vector_search, api_repository):
        self.vector_search = vector_search
        self.api_repository = api_repository
    
    async def execute(
        self,
        query_text: str,
        limit: int = 5
    ) -> List[RetrievalContext]:
        """
        Retrieve relevant API context using semantic search.
        
        Args:
            query_text: Natural language query
            limit: Maximum number of results
            
        Returns:
            List of retrieval contexts
        """
        try:
            # Search vector database
            results = await self.vector_search.search_apis(
                query_text=query_text,
                limit=limit
            )
            
            # Convert to RetrievalContext entities
            contexts = []
            for result in results:
                # Get full API spec from database
                api_spec = await self.api_repository.find_by_id(
                    UUID(result['id'])
                )
                
                if api_spec:
                    context = RetrievalContext(
                        api_spec_id=api_spec.id,
                        api_name=api_spec.name,
                        api_version=str(api_spec.version),
                        relevance_score=result['score'],
                        api_config=api_spec.metadata,
                        matched_endpoints=api_spec.metadata.get('endpoints', []),
                        examples=self._extract_examples(api_spec.metadata)
                    )
                    contexts.append(context)
            
            return contexts
            
        except Exception as e:
            logger.error(f"Error retrieving context: {str(e)}")
            return []
    
    def _extract_examples(self, config):
        """Extract examples from API config."""
        examples = []
        for endpoint in config.get('endpoints', []):
            examples.extend(endpoint.get('examples', []))
        return examples


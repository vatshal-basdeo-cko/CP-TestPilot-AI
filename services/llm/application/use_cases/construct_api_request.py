"""Use case for constructing API requests."""
import logging
from typing import Optional

from domain.entities.test_request import TestRequest
from domain.entities.api_call import APICall
from domain.entities.retrieval_context import RetrievalContext
from domain.entities.clarification import Clarification

logger = logging.getLogger(__name__)


class ConstructAPIRequestUseCase:
    """Use case for constructing executable API requests."""
    
    def __init__(self, llm_service, faker_service):
        self.llm_service = llm_service
        self.faker_service = faker_service
    
    async def execute(
        self,
        test_request: TestRequest,
        context: RetrievalContext
    ) -> tuple[Optional[APICall], Optional[Clarification]]:
        """
        Construct API call from natural language and context.
        
        Returns:
            Tuple of (APICall, Clarification) - one will be None
        """
        try:
            # Use LLM to construct request
            result = await self.llm_service.construct_request(
                natural_language=test_request.natural_language_request,
                api_spec=context.api_config,
                examples=context.examples
            )
            
            # Check if clarification needed
            if result.get('needs_clarification'):
                clarification = Clarification(
                    question=result['clarification_question'],
                    clarification_type=result.get('clarification_type', 'choice'),
                    options=result.get('options', []),
                    context={'test_request_id': str(test_request.id)}
                )
                return (None, clarification)
            
            # Generate missing required data
            if result.get('missing_required'):
                result = await self._generate_missing_data(result, context)
            
            # Build API call entity
            api_call = APICall(
                method=result['method'],
                url=self._build_url(context, result['path']),
                headers=result.get('headers', {}),
                query_params=result.get('query_params', {}),
                body=result.get('body'),
                api_spec_id=context.api_spec_id,
                api_name=context.api_name,
                endpoint_name=result.get('endpoint'),
                confidence_score=result.get('confidence', 0.8)
            )
            
            return (api_call, None)
            
        except Exception as e:
            logger.error(f"Error constructing request: {str(e)}")
            return (None, None)
    
    async def _generate_missing_data(self, result, context):
        """Generate realistic test data for missing required fields."""
        for field in result.get('missing_required', []):
            field_type = self._get_field_type(field, context)
            result['parameters'][field] = await self.faker_service.generate(field_type, field)
        return result
    
    def _get_field_type(self, field_name, context):
        """Infer field type from schema."""
        # Simplified type inference
        if 'email' in field_name.lower():
            return 'email'
        elif 'card' in field_name.lower():
            return 'card_number'
        elif 'amount' in field_name.lower():
            return 'amount'
        return 'string'
    
    def _build_url(self, context, path):
        """Build complete URL from base and path."""
        base_url = context.api_config.get('base_url', '')
        return f"{base_url}{path}"


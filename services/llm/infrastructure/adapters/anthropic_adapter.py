"""Anthropic LLM adapter."""
import json
import logging
from typing import Dict, Any
from anthropic import AsyncAnthropic

logger = logging.getLogger(__name__)


class AnthropicAdapter:
    """Adapter for Anthropic Claude API."""
    
    def __init__(self, api_key: str, model: str = "claude-3-sonnet-20240229"):
        self.client = AsyncAnthropic(api_key=api_key)
        self.model = model
        logger.info(f"Initialized Anthropic adapter with model: {model}")
    
    async def generate(
        self,
        system_prompt: str,
        user_prompt: str,
        temperature: float = 0.7
    ) -> Dict[str, Any]:
        """Generate completion from Claude."""
        try:
            response = await self.client.messages.create(
                model=self.model,
                max_tokens=4096,
                system=system_prompt,
                messages=[
                    {"role": "user", "content": user_prompt}
                ],
                temperature=temperature
            )
            
            content = response.content[0].text
            # Parse JSON from response
            return json.loads(content)
            
        except Exception as e:
            logger.error(f"Anthropic API error: {str(e)}")
            raise


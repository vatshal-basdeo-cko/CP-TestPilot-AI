"""OpenAI LLM adapter."""
import json
import logging
from typing import Dict, Any
from openai import AsyncOpenAI

logger = logging.getLogger(__name__)


class OpenAIAdapter:
    """Adapter for OpenAI API."""
    
    def __init__(self, api_key: str, model: str = "gpt-4"):
        self.client = AsyncOpenAI(api_key=api_key)
        self.model = model
        logger.info(f"Initialized OpenAI adapter with model: {model}")
    
    async def generate(
        self,
        system_prompt: str,
        user_prompt: str,
        temperature: float = 0.7
    ) -> Dict[str, Any]:
        """Generate completion from OpenAI."""
        try:
            response = await self.client.chat.completions.create(
                model=self.model,
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_prompt}
                ],
                temperature=temperature,
                response_format={"type": "json_object"}
            )
            
            content = response.choices[0].message.content
            return json.loads(content)
            
        except Exception as e:
            logger.error(f"OpenAI API error: {str(e)}")
            raise


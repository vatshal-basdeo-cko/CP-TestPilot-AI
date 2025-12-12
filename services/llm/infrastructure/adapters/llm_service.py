"""LLM service with multi-provider support."""
import logging
import os
from pathlib import Path
from typing import Dict, Any

logger = logging.getLogger(__name__)


class LLMService:
    """Service for LLM operations with multi-provider support."""
    
    def __init__(self, provider: str = "openai"):
        self.provider = provider
        self.llm_adapter = self._create_adapter(provider)
        self.prompts = self._load_prompts()
    
    def _create_adapter(self, provider: str):
        """Create LLM adapter based on provider."""
        if provider == "openai":
            from .openai_adapter import OpenAIAdapter
            api_key = os.getenv("OPENAI_API_KEY")
            return OpenAIAdapter(api_key)
        elif provider == "anthropic":
            from .anthropic_adapter import AnthropicAdapter
            api_key = os.getenv("ANTHROPIC_API_KEY")
            return AnthropicAdapter(api_key)
        elif provider == "gemini":
            from .gemini_adapter import GeminiAdapter
            api_key = os.getenv("GEMINI_API_KEY")
            return GeminiAdapter(api_key)
        else:
            raise ValueError(f"Unsupported provider: {provider}. Supported: openai, anthropic, gemini")
    
    def _load_prompts(self):
        """Load prompt templates."""
        prompts_dir = Path(__file__).parent.parent.parent / "prompts"
        prompts = {}
        for prompt_file in prompts_dir.glob("*.txt"):
            with open(prompt_file, 'r') as f:
                prompts[prompt_file.stem] = f.read()
        return prompts
    
    async def parse_request(
        self,
        natural_language: str,
        api_context: list
    ) -> Dict[str, Any]:
        """Parse natural language request."""
        system_prompt = self.prompts['system_prompt']
        user_prompt = self.prompts['parse_request'].format(
            natural_language_request=natural_language,
            api_context=str(api_context)
        )
        
        return await self.llm_adapter.generate(system_prompt, user_prompt)
    
    async def construct_request(
        self,
        natural_language: str,
        api_spec: Dict[str, Any],
        examples: list
    ) -> Dict[str, Any]:
        """Construct API request."""
        system_prompt = self.prompts['system_prompt']
        user_prompt = self.prompts['construct_api_call'].format(
            natural_language_request=natural_language,
            api_spec=str(api_spec),
            endpoint=api_spec.get('endpoints', [{}])[0],
            parameters={}
        )
        
        return await self.llm_adapter.generate(system_prompt, user_prompt, temperature=0.3)


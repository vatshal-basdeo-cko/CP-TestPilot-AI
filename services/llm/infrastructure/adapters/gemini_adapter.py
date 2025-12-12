"""Google Gemini Pro LLM adapter."""
import json
import logging
from typing import Dict, Any
import google.generativeai as genai

logger = logging.getLogger(__name__)


class GeminiAdapter:
    """Adapter for Google Gemini Pro API."""
    
    def __init__(self, api_key: str, model: str = "gemini-pro"):
        genai.configure(api_key=api_key)
        self.model = genai.GenerativeModel(model)
        logger.info(f"Initialized Gemini adapter with model: {model}")
    
    async def generate(
        self,
        system_prompt: str,
        user_prompt: str,
        temperature: float = 0.7
    ) -> Dict[str, Any]:
        """Generate completion from Gemini Pro."""
        try:
            # Combine system and user prompts for Gemini
            full_prompt = f"""You are an AI assistant. Follow these instructions:

{system_prompt}

---

User Request:
{user_prompt}

---

IMPORTANT: Respond with valid JSON only, no additional text or markdown."""

            # Configure generation
            generation_config = genai.GenerationConfig(
                temperature=temperature,
                top_p=0.95,
                top_k=40,
            )
            
            # Generate response
            response = self.model.generate_content(
                full_prompt,
                generation_config=generation_config
            )
            
            # Extract text and parse JSON
            content = response.text.strip()
            
            # Remove markdown code blocks if present
            if content.startswith("```json"):
                content = content[7:]
            if content.startswith("```"):
                content = content[3:]
            if content.endswith("```"):
                content = content[:-3]
            content = content.strip()
            
            return json.loads(content)
            
        except json.JSONDecodeError as e:
            logger.error(f"Failed to parse Gemini response as JSON: {content[:200]}")
            raise ValueError(f"Invalid JSON response from Gemini: {str(e)}")
        except Exception as e:
            logger.error(f"Gemini API error: {str(e)}")
            raise


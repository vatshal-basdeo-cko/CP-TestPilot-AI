"""File reader adapter for reading and parsing configuration files."""
import json
import logging
from pathlib import Path
from typing import Dict, Any
import yaml

logger = logging.getLogger(__name__)


class FileReaderAdapter:
    """Adapter for reading and parsing YAML/JSON configuration files."""
    
    async def read_file(self, file_path: str) -> str:
        """
        Read file content as string.
        
        Args:
            file_path: Path to the file
            
        Returns:
            File content as string
            
        Raises:
            FileNotFoundError: If file doesn't exist
            PermissionError: If file can't be read
        """
        path = Path(file_path)
        
        if not path.exists():
            raise FileNotFoundError(f"File not found: {file_path}")
        
        if not path.is_file():
            raise ValueError(f"Path is not a file: {file_path}")
        
        try:
            with open(path, 'r', encoding='utf-8') as f:
                content = f.read()
            logger.debug(f"Read file: {file_path} ({len(content)} bytes)")
            return content
        except Exception as e:
            logger.error(f"Error reading file {file_path}: {str(e)}")
            raise
    
    async def parse_config(self, content: str, file_path: str) -> Dict[str, Any]:
        """
        Parse configuration content based on file extension.
        
        Args:
            content: File content string
            file_path: File path (used to determine format)
            
        Returns:
            Parsed configuration as dictionary
            
        Raises:
            ValueError: If format is unsupported or parsing fails
        """
        path = Path(file_path)
        suffix = path.suffix.lower()
        
        try:
            if suffix in ['.yaml', '.yml']:
                return self._parse_yaml(content)
            elif suffix == '.json':
                return self._parse_json(content)
            else:
                raise ValueError(f"Unsupported file format: {suffix}")
        except Exception as e:
            logger.error(f"Error parsing config from {file_path}: {str(e)}")
            raise
    
    def _parse_yaml(self, content: str) -> Dict[str, Any]:
        """Parse YAML content."""
        try:
            data = yaml.safe_load(content)
            if not isinstance(data, dict):
                raise ValueError("YAML content must be a dictionary/object")
            return data
        except yaml.YAMLError as e:
            raise ValueError(f"Invalid YAML: {str(e)}")
    
    def _parse_json(self, content: str) -> Dict[str, Any]:
        """Parse JSON content."""
        try:
            data = json.loads(content)
            if not isinstance(data, dict):
                raise ValueError("JSON content must be an object")
            return data
        except json.JSONDecodeError as e:
            raise ValueError(f"Invalid JSON: {str(e)}")
    
    def validate_config(self, config: Dict[str, Any]) -> None:
        """
        Validate that configuration has required fields.
        
        Args:
            config: Configuration dictionary
            
        Raises:
            ValueError: If required fields are missing
        """
        required_fields = ['name']
        missing = [field for field in required_fields if field not in config]
        
        if missing:
            raise ValueError(f"Missing required fields: {', '.join(missing)}")
        
        # Validate name is non-empty
        if not config['name'] or not config['name'].strip():
            raise ValueError("API name cannot be empty")


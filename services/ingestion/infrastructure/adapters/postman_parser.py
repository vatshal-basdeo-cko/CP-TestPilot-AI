"""Postman collection parser adapter."""
import logging
from typing import Dict, Any, List

logger = logging.getLogger(__name__)


class PostmanParser:
    """Adapter for parsing Postman collections into our API configuration format."""
    
    async def parse(self, collection_data: Dict[str, Any]) -> Dict[str, Any]:
        """
        Parse Postman collection to our API configuration format.
        
        Args:
            collection_data: Postman collection JSON
            
        Returns:
            API configuration dictionary
            
        Raises:
            ValueError: If collection is invalid
        """
        try:
            # Extract collection info
            info = collection_data.get('info', {})
            name = info.get('name', 'Unnamed API')
            description = info.get('description', '')
            version = info.get('version', '1.0.0')
            
            # Parse items (requests)
            items = collection_data.get('item', [])
            endpoints = []
            
            for item in items:
                endpoints.extend(self._parse_item(item))
            
            # Build our configuration format
            config = {
                'name': name,
                'version': version,
                'description': description,
                'endpoints': endpoints,
                'source': 'postman',
                'postman_id': info.get('_postman_id'),
            }
            
            logger.info(f"Parsed Postman collection: {name} with {len(endpoints)} endpoints")
            return config
            
        except Exception as e:
            logger.error(f"Error parsing Postman collection: {str(e)}")
            raise ValueError(f"Invalid Postman collection: {str(e)}")
    
    def _parse_item(self, item: Dict[str, Any], prefix: str = '') -> List[Dict[str, Any]]:
        """
        Recursively parse Postman items (can be requests or folders).
        
        Args:
            item: Postman item
            prefix: Path prefix for nested items
            
        Returns:
            List of endpoint dictionaries
        """
        endpoints = []
        
        # Check if this is a folder (has nested items)
        if 'item' in item:
            folder_name = item.get('name', '')
            new_prefix = f"{prefix}/{folder_name}" if prefix else folder_name
            
            for nested_item in item['item']:
                endpoints.extend(self._parse_item(nested_item, new_prefix))
        
        # Check if this is a request
        elif 'request' in item:
            endpoint = self._parse_request(item, prefix)
            if endpoint:
                endpoints.append(endpoint)
        
        return endpoints
    
    def _parse_request(self, item: Dict[str, Any], prefix: str) -> Dict[str, Any]:
        """
        Parse a Postman request item.
        
        Args:
            item: Postman request item
            prefix: Path prefix
            
        Returns:
            Endpoint dictionary
        """
        request = item['request']
        
        # Extract method
        method = request.get('method', 'GET').upper()
        
        # Extract URL
        url = request.get('url', {})
        if isinstance(url, str):
            path = url
        elif isinstance(url, dict):
            raw_url = url.get('raw', '')
            # Try to extract path from raw URL
            path = self._extract_path_from_url(raw_url)
        else:
            path = ''
        
        # Add prefix
        if prefix:
            path = f"/{prefix}{path}" if not path.startswith('/') else f"/{prefix}{path}"
        
        # Extract description
        description = item.get('name', '')
        if 'description' in item:
            desc_obj = item['description']
            if isinstance(desc_obj, str):
                description = desc_obj
            elif isinstance(desc_obj, dict):
                description = desc_obj.get('content', description)
        
        # Extract headers
        headers = {}
        if 'header' in request:
            for header in request['header']:
                if not header.get('disabled', False):
                    key = header.get('key', '')
                    value = header.get('value', '')
                    if key:
                        headers[key] = value
        
        # Extract body
        body_schema = None
        if 'body' in request:
            body = request['body']
            mode = body.get('mode', '')
            
            if mode == 'raw':
                # Try to infer schema from raw body
                raw = body.get('raw', '')
                body_schema = self._infer_schema_from_example(raw)
        
        # Extract parameters from URL
        parameters = []
        if isinstance(url, dict) and 'variable' in url:
            for var in url['variable']:
                parameters.append({
                    'name': var.get('key', ''),
                    'type': 'string',
                    'description': var.get('description', ''),
                    'required': not var.get('disabled', False),
                })
        
        endpoint = {
            'name': item.get('name', ''),
            'path': path,
            'method': method,
            'description': description,
            'headers': headers,
            'parameters': parameters,
        }
        
        if body_schema:
            endpoint['request_schema'] = body_schema
        
        return endpoint
    
    def _extract_path_from_url(self, url: str) -> str:
        """Extract path from full URL."""
        # Remove protocol and domain
        if '://' in url:
            url = url.split('://', 1)[1]
        if '/' in url:
            path = '/' + url.split('/', 1)[1]
        else:
            path = '/'
        
        # Remove query string
        if '?' in path:
            path = path.split('?', 1)[0]
        
        return path
    
    def _infer_schema_from_example(self, example: str) -> Dict[str, Any]:
        """Try to infer JSON schema from example data."""
        try:
            import json
            data = json.loads(example)
            
            # Simple schema inference
            schema = {
                'type': 'object',
                'properties': {}
            }
            
            if isinstance(data, dict):
                for key, value in data.items():
                    schema['properties'][key] = {
                        'type': self._get_json_type(value)
                    }
            
            return schema
        except:
            return {'type': 'string'}
    
    def _get_json_type(self, value: Any) -> str:
        """Get JSON schema type from Python value."""
        if isinstance(value, bool):
            return 'boolean'
        elif isinstance(value, int):
            return 'integer'
        elif isinstance(value, float):
            return 'number'
        elif isinstance(value, str):
            return 'string'
        elif isinstance(value, list):
            return 'array'
        elif isinstance(value, dict):
            return 'object'
        else:
            return 'string'


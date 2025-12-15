package prompts

import (
	"fmt"
	"strings"
)

// SystemPrompt is the base system prompt for API construction
const SystemPrompt = `You are an AI assistant that helps construct API requests from natural language descriptions.
Your role is to:
1. Understand the user's intent
2. Map it to the appropriate API endpoint
3. Extract parameters from the request
4. Generate any missing test data using realistic values

Always respond in valid JSON format.`

// ParseRequestPrompt generates a prompt for parsing natural language requests
func ParseRequestPrompt(naturalLanguage string, apiContext string) string {
	return fmt.Sprintf(`%s

## Available APIs
%s

## User Request
"%s"

## Task
Parse the user's request and extract:
1. The intent (what action they want to perform)
2. Which API/endpoint matches their request
3. Any parameters mentioned in their request
4. Any missing required parameters that need clarification

Respond in this JSON format:
{
    "intent": "description of what the user wants",
    "api_name": "name of the matching API",
    "endpoint": "the specific endpoint path",
    "method": "HTTP method",
    "parameters": {
        "param_name": "value or null if not specified"
    },
    "missing_required": ["list of required params not provided"],
    "confidence": 0.0 to 1.0
}`, SystemPrompt, apiContext, naturalLanguage)
}

// ConstructRequestPrompt generates a prompt for constructing the API call
func ConstructRequestPrompt(parseResult string, apiConfig string, generatedData map[string]interface{}) string {
	dataStr := ""
	if len(generatedData) > 0 {
		var pairs []string
		for k, v := range generatedData {
			pairs = append(pairs, fmt.Sprintf("  %s: %v", k, v))
		}
		dataStr = "\n## Generated Test Data\n" + strings.Join(pairs, "\n")
	}

	return fmt.Sprintf(`%s

## Parse Result
%s

## API Configuration
%s
%s

## Task
Construct the complete API request with:
1. Full URL with path parameters replaced
2. All required headers
3. Query parameters if applicable
4. Request body with all required fields

Respond in this JSON format:
{
    "method": "GET/POST/PUT/DELETE",
    "url": "full URL with path params replaced",
    "path": "the endpoint path",
    "headers": {"header_name": "value"},
    "query_params": {"param": "value"},
    "body": {"field": "value"},
    "confidence": 0.0 to 1.0
}`, SystemPrompt, parseResult, apiConfig, dataStr)
}

// ClarificationPrompt generates a prompt for requesting clarification
func ClarificationPrompt(missingParams []string, apiContext string) string {
	return fmt.Sprintf(`The user's request is missing some required information.

## Missing Parameters
%s

## API Context
%s

Generate a friendly clarification request. Respond in JSON:
{
    "message": "friendly message asking for the missing information",
    "type": "multiple_choice or free_text",
    "options": [
        {"value": "option1", "description": "description"},
        {"value": "option2", "description": "description"}
    ],
    "field_name": "the parameter being clarified"
}`, strings.Join(missingParams, ", "), apiContext)
}

// BuildAPIContext builds context string from API configurations
func BuildAPIContext(contexts []map[string]interface{}) string {
	if len(contexts) == 0 {
		return "No API configurations available."
	}

	var parts []string
	for _, ctx := range contexts {
		name, _ := ctx["api_name"].(string)
		version, _ := ctx["version"].(string)
		desc, _ := ctx["description"].(string)

		part := fmt.Sprintf("### %s (v%s)\n%s", name, version, desc)
		
		if endpoints, ok := ctx["endpoints"].([]interface{}); ok {
			part += "\n\nEndpoints:"
			for _, ep := range endpoints {
				if epMap, ok := ep.(map[string]interface{}); ok {
					method, _ := epMap["method"].(string)
					path, _ := epMap["path"].(string)
					epDesc, _ := epMap["description"].(string)
					part += fmt.Sprintf("\n- %s %s: %s", method, path, epDesc)
				}
			}
		}

		parts = append(parts, part)
	}

	return strings.Join(parts, "\n\n---\n\n")
}


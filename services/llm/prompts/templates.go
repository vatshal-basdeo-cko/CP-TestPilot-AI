package prompts

import (
	"encoding/json"
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

// ConstructRequestPrompt generates a prompt for constructing the API call (legacy, without API context)
func ConstructRequestPrompt(parseResult string, apiConfig string, generatedData map[string]interface{}) string {
	return ConstructRequestPromptWithContext(parseResult, apiConfig, "", generatedData)
}

// ConstructRequestPromptWithContext generates a prompt for constructing the API call with API context
func ConstructRequestPromptWithContext(parseResult string, apiConfig string, apiContext string, generatedData map[string]interface{}) string {
	dataStr := ""
	if len(generatedData) > 0 {
		var pairs []string
		for k, v := range generatedData {
			pairs = append(pairs, fmt.Sprintf("  %s: %v", k, v))
		}
		dataStr = "\n## Generated Test Data\n" + strings.Join(pairs, "\n")
	}

	contextStr := ""
	if apiContext != "" {
		contextStr = "\n## Available API Context (use base_url from here)\n" + apiContext
	}

	return fmt.Sprintf(`%s

## Parse Result
%s

## API Configuration
%s
%s
%s

## Task
Construct the complete API request. **CRITICAL REQUIREMENTS**:

1. **USE THE EXAMPLE AS A TEMPLATE**: Start with the complete example request from the API context and modify only the fields the user specified. Keep ALL other fields from the example.

2. **INCLUDE ALL REQUIRED NESTED OBJECTS**: The request body MUST include:
   - All top-level required fields
   - Complete "sender" object with: first_name, last_name, address, city, country, account_number
   - Complete "recipient" object with: first_name, last_name, address, city, country, account_number  
   - Complete "card" object with: card_number, expiry_month, expiry_year
   - Complete "merchant_acquirer_configuration" with both:
     - "processing_profile": processing_profile_name, entity_id, business_model, acquirer_key, status, schemes, processing_type, business_settings
     - "acquirer": acquirer_key, acquirer_name, acquirer_country_code, forwarding_institution_id, processor_key, custom_settings

3. **BUILD THE FULL URL**: Combine base_url + endpoint path with path parameters replaced
   - Example: base_url "http://cp-ptc.qa.internal/mastercard" + path "/authorizations/{payment_id}/paytocard" 
   - Result: "http://cp-ptc.qa.internal/mastercard/authorizations/pay_xyz123/paytocard"

4. **USE EXAMPLE VALUES FOR UNSPECIFIED FIELDS**: For any field the user didn't mention, use the value from the example in the API context.

5. All required headers including Content-Type: application/json

Respond in this JSON format (body must be a complete, valid payload):
{
    "method": "GET/POST/PUT/DELETE",
    "url": "FULL URL starting with http:// or https:// with all path params replaced",
    "path": "the endpoint path",
    "headers": {"Content-Type": "application/json"},
    "query_params": {},
    "body": {COMPLETE request body with ALL required nested objects},
    "confidence": 0.0 to 1.0
}`, SystemPrompt, parseResult, apiConfig, contextStr, dataStr)
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

		// Try to get config - it can be a string (raw) or map (already parsed)
		var config map[string]interface{}
		if configStr, ok := ctx["config"].(string); ok && configStr != "" {
			// Config is a JSON string - parse it
			_ = json.Unmarshal([]byte(configStr), &config)
		} else if configMap, ok := ctx["config"].(map[string]interface{}); ok {
			// Config is already a map
			config = configMap
		}

		if config != nil {
			// Include base_url - CRITICAL for constructing full URLs
			if baseURL, ok := config["base_url"].(string); ok && baseURL != "" {
				part += fmt.Sprintf("\n\n**Base URL**: %s", baseURL)
			}

			// Include detailed endpoint information
			if endpoints, ok := config["endpoints"].([]interface{}); ok {
				part += "\n\n**Endpoints:**"
				for _, ep := range endpoints {
					if epMap, ok := ep.(map[string]interface{}); ok {
						method, _ := epMap["method"].(string)
						path, _ := epMap["path"].(string)
						epDesc, _ := epMap["description"].(string)
						part += fmt.Sprintf("\n- %s %s: %s", method, path, epDesc)

						// Include parameters info
						if params, ok := epMap["parameters"].([]interface{}); ok && len(params) > 0 {
							part += "\n  Parameters:"
							for _, p := range params {
								if pMap, ok := p.(map[string]interface{}); ok {
									pName, _ := pMap["name"].(string)
									pIn, _ := pMap["in"].(string)
									pRequired, _ := pMap["required"].(bool)
									pExample, _ := pMap["example"].(string)
									if pIn == "" {
										pIn = "body"
									}
									reqStr := ""
									if pRequired {
										reqStr = " (required)"
									}
									if pExample != "" {
										part += fmt.Sprintf("\n    - %s [%s]%s example: %s", pName, pIn, reqStr, pExample)
									} else {
										part += fmt.Sprintf("\n    - %s [%s]%s", pName, pIn, reqStr)
									}
								}
							}
						}

						// Include example request if available
						if examples, ok := epMap["examples"].([]interface{}); ok && len(examples) > 0 {
							if ex, ok := examples[0].(map[string]interface{}); ok {
								if exReq, ok := ex["request"].(map[string]interface{}); ok {
									exJSON, _ := json.MarshalIndent(exReq, "    ", "  ")
									part += fmt.Sprintf("\n  Example request body:\n    %s", string(exJSON))
								}
							}
						}
					}
				}
			}
		} else if endpoints, ok := ctx["endpoints"].([]interface{}); ok {
			// Fallback to old format if no config available
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

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
4. For parameters NOT mentioned by the user, use the string "[AUTO]" to indicate they should be auto-generated with test data

IMPORTANT: For payment/test fields like card_number, expiry_month, expiry_year, cvv, account_number, first_name, last_name, address, city, country, phone, email - use "[AUTO]" instead of null. These will be automatically filled with realistic test data.

Only use null for truly business-critical fields that the user MUST specify (like payment_id, transaction type, or specific amounts the user wants).

Respond in this JSON format:
{
    "intent": "description of what the user wants",
    "api_name": "name of the matching API",
    "endpoint": "the specific endpoint path",
    "method": "HTTP method",
    "parameters": {
        "param_name": "value, '[AUTO]' for auto-generate, or null if user must specify"
    },
    "missing_required": ["list of required params user MUST provide - not auto-generatable ones"],
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

1. **COPY THE ENTIRE EXAMPLE REQUEST**: You MUST start by COPYING the complete example request body from the API context EXACTLY as-is. Then ONLY modify the specific values the user mentioned (like amount, currency). DO NOT remove ANY fields from the example.

2. **THE BODY MUST INCLUDE ALL OF THESE (copy from example)**:
   - service_type, action_id, transaction_id, transaction_purpose, entity_id, processing_key
   - amount, currency, merchant_category_code, fund_transfer_type
   - All card_acceptor_* fields
   - cko_vault_id, sequence_number (use example value: 1000000000)
   - Complete "card" object with card_number, expiry_month, expiry_year (INTEGERS, not strings!)
   - Complete "sender" object with ALL fields from the example
   - Complete "recipient" object with ALL fields from the example
   - Complete "merchant_acquirer_configuration" with BOTH nested objects from the example

3. **BUILD THE FULL URL - CRITICAL**:
   - The HOST is ALWAYS: http://cp-ptc.qa.internal (do NOT use base_url if it contains "{{" variables)
   - The PATH comes from the matched API's endpoint path (e.g., /visa/authorizations/{payment_id}/paytocard or /mastercard/authorizations/{payment_id}/paytocard)
   - Replace {payment_id} or {{payment_id}} with a generated payment ID like pay_kgaohkfk72ketfpxkf55gpytwu
   - EXAMPLE for Visa: http://cp-ptc.qa.internal/visa/authorizations/pay_xyz123/paytocard
   - EXAMPLE for Mastercard: http://cp-ptc.qa.internal/mastercard/authorizations/pay_xyz123/paytocard
   - The path ALREADY includes /visa/ or /mastercard/ - just append it to the host!

4. **PATH PARAMETERS**: Generate a valid CKO payment_id like pay_kgaohkfk72ketfpxkf55gpytwu

5. **ONLY MODIFY WHAT USER SPECIFIED**: If user says "100 USD", only change amount to 100 and currency to USD. Keep EVERYTHING else from the example.

6. Headers: {"Content-Type": "application/json"}

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
						hasExample := false
						if examples, ok := epMap["examples"].([]interface{}); ok && len(examples) > 0 {
							if ex, ok := examples[0].(map[string]interface{}); ok {
								if exReq, ok := ex["request"].(map[string]interface{}); ok {
									exJSON, _ := json.MarshalIndent(exReq, "    ", "  ")
									part += fmt.Sprintf("\n  Example request body:\n    %s", string(exJSON))
									hasExample = true
								}
							}
						}
						// If no examples, use request_schema as the template (for Postman imports)
						if !hasExample {
							if reqSchema, ok := epMap["request_schema"].(map[string]interface{}); ok && len(reqSchema) > 0 {
								schemaJSON, _ := json.MarshalIndent(reqSchema, "    ", "  ")
								part += fmt.Sprintf("\n  Request body template (USE THIS EXACTLY, only modify values user specified):\n    %s", string(schemaJSON))
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

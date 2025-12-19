package adapters

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/testpilot-ai/ingestion/domain/entities"
)

// PostmanParser handles parsing of Postman collections
type PostmanParser struct{}

// NewPostmanParser creates a new Postman parser
func NewPostmanParser() *PostmanParser {
	return &PostmanParser{}
}

// ParseCollection parses a Postman collection file
func (p *PostmanParser) ParseCollection(filePath string) (*entities.APIConfig, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	return p.ParseCollectionData(data)
}

// ParseCollectionData parses Postman collection from bytes
func (p *PostmanParser) ParseCollectionData(data []byte) (*entities.APIConfig, string, error) {
	var collection entities.PostmanCollection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, "", fmt.Errorf("failed to parse Postman collection: %w", err)
	}

	// Calculate hash
	fileParser := NewFileParser()
	contentHash := fileParser.CalculateHash(data)

	// Extract base_url from the first request's host
	baseURL := p.extractBaseURL(collection.Item)

	// Convert to APIConfig
	config := &entities.APIConfig{
		Name:        collection.Info.Name,
		Version:     "1.0.0",
		Description: collection.Info.Description,
		BaseURL:     baseURL,
		Endpoints:   p.extractEndpoints(collection.Item),
	}

	return config, contentHash, nil
}

// extractBaseURL extracts the base URL from the first request in the collection
func (p *PostmanParser) extractBaseURL(items []entities.PostmanItem) string {
	for _, item := range items {
		// If it's a folder, recurse
		if len(item.Item) > 0 {
			if baseURL := p.extractBaseURL(item.Item); baseURL != "" {
				return baseURL
			}
			continue
		}

		// If it has a request with a host, extract base URL
		if item.Request != nil && len(item.Request.URL.Host) > 0 {
			host := strings.Join(item.Request.URL.Host, ".")
			// Check if it looks like a variable (e.g., {{base_url}})
			if strings.Contains(host, "{{") {
				// Try to get from raw URL
				if item.Request.URL.Raw != "" {
					// Parse raw URL to extract host
					raw := item.Request.URL.Raw
					// Remove variable syntax and try to find a real host
					if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
						// Extract host from raw URL
						parts := strings.SplitN(raw, "/", 4)
						if len(parts) >= 3 {
							return parts[0] + "//" + parts[2]
						}
					}
				}
				continue
			}
			// Build proper base URL
			protocol := "http://"
			if strings.Contains(item.Request.URL.Raw, "https://") {
				protocol = "https://"
			}
			return protocol + host
		}
	}
	return ""
}

// extractEndpoints recursively extracts endpoints from Postman items
func (p *PostmanParser) extractEndpoints(items []entities.PostmanItem) []entities.APIEndpoint {
	var endpoints []entities.APIEndpoint

	for _, item := range items {
		// If it's a folder, recurse (do NOT prepend folder names to paths)
		if len(item.Item) > 0 {
			endpoints = append(endpoints, p.extractEndpoints(item.Item)...)
			continue
		}

		// If it has a request, convert it
		if item.Request != nil {
			endpoint := p.convertRequest(item)
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

// convertRequest converts a Postman request to an APIEndpoint
func (p *PostmanParser) convertRequest(item entities.PostmanItem) entities.APIEndpoint {
	req := item.Request

	// Build path from URL (just the path, no folder prefixes)
	path := "/" + strings.Join(req.URL.Path, "/")

	endpoint := entities.APIEndpoint{
		Name:        item.Name,
		Path:        path,
		Method:      req.Method,
		Description: item.Name,
		Parameters:  p.extractParameters(req),
	}

	// Extract auth from headers
	for _, header := range req.Header {
		if strings.ToLower(header.Key) == "authorization" ||
			strings.ToLower(header.Key) == "x-api-key" {
			endpoint.Authentication = &entities.AuthConfig{
				Type:   "api_key",
				Header: header.Key,
			}
			break
		}
	}

	// Parse request body if present - preserve proper types
	if req.Body != nil && req.Body.Mode == "raw" && req.Body.Raw != "" {
		var bodySchema map[string]interface{}
		if err := json.Unmarshal([]byte(req.Body.Raw), &bodySchema); err == nil {
			endpoint.RequestSchema = bodySchema
		}
	}

	return endpoint
}

// extractParameters extracts parameters from a Postman request
func (p *PostmanParser) extractParameters(req *entities.PostmanRequest) []entities.Parameter {
	var params []entities.Parameter

	// Extract query parameters
	for _, q := range req.URL.Query {
		params = append(params, entities.Parameter{
			Name:     q.Key,
			Type:     "string",
			Required: false,
		})
	}

	return params
}

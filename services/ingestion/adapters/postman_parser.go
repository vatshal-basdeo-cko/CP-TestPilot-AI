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

	// Convert to APIConfig
	config := &entities.APIConfig{
		Name:        collection.Info.Name,
		Version:     "1.0.0",
		Description: collection.Info.Description,
		Endpoints:   p.extractEndpoints(collection.Item, ""),
	}

	return config, contentHash, nil
}

// extractEndpoints recursively extracts endpoints from Postman items
func (p *PostmanParser) extractEndpoints(items []entities.PostmanItem, prefix string) []entities.APIEndpoint {
	var endpoints []entities.APIEndpoint

	for _, item := range items {
		// If it's a folder, recurse
		if len(item.Item) > 0 {
			folderPrefix := prefix
			if item.Name != "" {
				folderPrefix = prefix + "/" + item.Name
			}
			endpoints = append(endpoints, p.extractEndpoints(item.Item, folderPrefix)...)
			continue
		}

		// If it has a request, convert it
		if item.Request != nil {
			endpoint := p.convertRequest(item, prefix)
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

// convertRequest converts a Postman request to an APIEndpoint
func (p *PostmanParser) convertRequest(item entities.PostmanItem, prefix string) entities.APIEndpoint {
	req := item.Request

	// Build path from URL
	path := "/" + strings.Join(req.URL.Path, "/")
	if prefix != "" {
		path = prefix + path
	}

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

	// Parse request body if present
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


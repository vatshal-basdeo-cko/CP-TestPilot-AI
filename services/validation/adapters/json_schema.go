package adapters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/testpilot-ai/validation/domain/entities"
	"github.com/xeipuuv/gojsonschema"
)

// JSONSchemaValidator handles JSON schema validation
type JSONSchemaValidator struct{}

// NewJSONSchemaValidator creates a new JSON schema validator
func NewJSONSchemaValidator() *JSONSchemaValidator {
	return &JSONSchemaValidator{}
}

// ValidateSchema validates a response against a JSON schema
func (v *JSONSchemaValidator) ValidateSchema(response map[string]interface{}, schema map[string]interface{}) *entities.SchemaCheckResult {
	result := &entities.SchemaCheckResult{IsValid: true}

	if schema == nil {
		return result // No schema to validate against
	}

	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewGoLoader(response)

	validationResult, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		result.IsValid = false
		result.Errors = []string{fmt.Sprintf("Schema validation error: %s", err)}
		return result
	}

	if !validationResult.Valid() {
		result.IsValid = false
		for _, desc := range validationResult.Errors() {
			result.Errors = append(result.Errors, desc.String())
		}
	}

	return result
}

// ValidateStatus validates the response status code
func (v *JSONSchemaValidator) ValidateStatus(actual, expected int) *entities.StatusCheckResult {
	return &entities.StatusCheckResult{
		Expected: expected,
		Actual:   actual,
		IsValid:  actual == expected,
	}
}

// CompareResponses compares two responses and returns differences
func (v *JSONSchemaValidator) CompareResponses(current, previous map[string]interface{}) *entities.DiffResult {
	result := &entities.DiffResult{HasDifferences: false}

	if previous == nil {
		return result // Nothing to compare with
	}

	v.compareObjects("", current, previous, result)
	result.HasDifferences = len(result.Additions) > 0 || len(result.Deletions) > 0 || len(result.Modifications) > 0

	return result
}

// compareObjects recursively compares two objects
func (v *JSONSchemaValidator) compareObjects(path string, current, previous map[string]interface{}, result *entities.DiffResult) {
	// Check for additions and modifications in current
	for key, currVal := range current {
		keyPath := v.buildPath(path, key)
		prevVal, exists := previous[key]

		if !exists {
			result.Additions = append(result.Additions, entities.DiffEntry{
				Path:     keyPath,
				NewValue: currVal,
			})
			continue
		}

		// Both exist, check if they differ
		if !v.deepEqual(currVal, prevVal) {
			// If both are maps, recurse
			currMap, currIsMap := currVal.(map[string]interface{})
			prevMap, prevIsMap := prevVal.(map[string]interface{})

			if currIsMap && prevIsMap {
				v.compareObjects(keyPath, currMap, prevMap, result)
			} else {
				result.Modifications = append(result.Modifications, entities.DiffEntry{
					Path:     keyPath,
					OldValue: prevVal,
					NewValue: currVal,
				})
			}
		}
	}

	// Check for deletions in current
	for key, prevVal := range previous {
		keyPath := v.buildPath(path, key)
		if _, exists := current[key]; !exists {
			result.Deletions = append(result.Deletions, entities.DiffEntry{
				Path:     keyPath,
				OldValue: prevVal,
			})
		}
	}
}

// buildPath builds a JSON path string
func (v *JSONSchemaValidator) buildPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}

// deepEqual compares two values for deep equality
func (v *JSONSchemaValidator) deepEqual(a, b interface{}) bool {
	// Convert both to JSON and compare strings for reliability
	aJSON, err1 := json.Marshal(a)
	bJSON, err2 := json.Marshal(b)

	if err1 != nil || err2 != nil {
		return reflect.DeepEqual(a, b)
	}

	return strings.TrimSpace(string(aJSON)) == strings.TrimSpace(string(bJSON))
}

// ApplyCustomRules applies custom validation rules
func (v *JSONSchemaValidator) ApplyCustomRules(response map[string]interface{}, rules []entities.ValidationRule) []entities.CustomCheckResult {
	var results []entities.CustomCheckResult

	for _, rule := range rules {
		if rule.RuleType != "custom" {
			continue
		}

		result := v.applyCustomRule(response, rule)
		results = append(results, result)
	}

	return results
}

// applyCustomRule applies a single custom rule
func (v *JSONSchemaValidator) applyCustomRule(response map[string]interface{}, rule entities.ValidationRule) entities.CustomCheckResult {
	result := entities.CustomCheckResult{
		RuleName: rule.RuleDefinition["name"].(string),
		IsValid:  true,
	}

	// Get rule type
	ruleType, ok := rule.RuleDefinition["type"].(string)
	if !ok {
		result.IsValid = false
		result.Message = "Invalid rule definition: missing type"
		return result
	}

	switch ruleType {
	case "field_required":
		field := rule.RuleDefinition["field"].(string)
		if _, exists := response[field]; !exists {
			result.IsValid = false
			result.Message = fmt.Sprintf("Required field '%s' is missing", field)
		}

	case "field_not_empty":
		field := rule.RuleDefinition["field"].(string)
		val, exists := response[field]
		if !exists || val == nil || val == "" {
			result.IsValid = false
			result.Message = fmt.Sprintf("Field '%s' must not be empty", field)
		}

	case "field_matches":
		field := rule.RuleDefinition["field"].(string)
		expected := rule.RuleDefinition["value"]
		actual, exists := response[field]
		if !exists || !v.deepEqual(actual, expected) {
			result.IsValid = false
			result.Message = fmt.Sprintf("Field '%s' does not match expected value", field)
		}

	default:
		result.IsValid = true
		result.Message = fmt.Sprintf("Unknown rule type: %s", ruleType)
	}

	return result
}


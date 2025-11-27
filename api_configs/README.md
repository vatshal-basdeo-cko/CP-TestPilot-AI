# API Configuration Schema

This document describes the YAML schema for API configuration files used by TestPilot AI.

## Overview

API configuration files define the structure, endpoints, and behavior of APIs that TestPilot AI can test. These configurations are used by the LLM service to understand API specifications and construct valid requests from natural language.

## Schema Structure

### Root Level

```yaml
name: string (required)
version: string (required)
description: string (optional)
base_url: string (required)
endpoints: array (required)
test_scenarios: array (optional)
metadata: object (optional)
```

### Root Fields

#### `name` (required)
- **Type:** string
- **Description:** Unique name for the API
- **Example:** `"Mastercard PTC Authorisation"`

#### `version` (required)
- **Type:** string
- **Format:** Semantic versioning (major.minor.patch)
- **Description:** API version
- **Example:** `"1.0.0"`

#### `description` (optional)
- **Type:** string
- **Description:** Human-readable description of the API
- **Example:** `"Process Mastercard payment authorization transactions"`

#### `base_url` (required)
- **Type:** string
- **Description:** Base URL for the API. Can use environment variable substitution.
- **Example:** `"${ENV_BASE_URL}"` or `"https://api.example.com"`
- **Note:** `${ENV_BASE_URL}` will be replaced at runtime from environment configuration

---

## Endpoints

Each endpoint in the `endpoints` array defines a single API operation.

### Endpoint Structure

```yaml
- name: string (required)
  path: string (required)
  method: string (required)
  description: string (optional)
  authentication: object (optional)
  parameters: array (optional)
  request_schema: object (optional)
  response_schema: object (optional)
  expected_status_codes: array (optional)
  examples: array (optional)
```

### Endpoint Fields

#### `name` (required)
- **Type:** string
- **Description:** Unique name for the endpoint
- **Example:** `"authorize"`, `"create_payment"`

#### `path` (required)
- **Type:** string
- **Description:** API path (can include path parameters in curly braces)
- **Example:** `"/api/v1/authorize"`, `"/api/v2/payments/{payment_id}"`

#### `method` (required)
- **Type:** string
- **Values:** `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
- **Description:** HTTP method

#### `description` (optional)
- **Type:** string
- **Description:** Human-readable description of what the endpoint does

#### `authentication` (optional)
- **Type:** object
- **Fields:**
  - `type`: `"api_key"`, `"bearer"`, `"basic"`, `"oauth2"`, `"none"`
  - `header`: Header name (for api_key type)
  - `description`: Authentication description

**Example:**
```yaml
authentication:
  type: "api_key"
  header: "X-API-Key"
  description: "API key authentication required"
```

---

## Parameters

Parameters define the inputs to an endpoint.

### Parameter Structure

```yaml
- name: string (required)
  type: string (required)
  required: boolean (optional, default: false)
  in: string (optional)
  format: string (optional)
  description: string (optional)
  example: any (optional)
  default: any (optional)
  enum: array (optional)
```

### Parameter Fields

#### `name` (required)
- **Type:** string
- **Description:** Parameter name

#### `type` (required)
- **Type:** string
- **Values:** `string`, `number`, `integer`, `boolean`, `array`, `object`

#### `required` (optional)
- **Type:** boolean
- **Default:** `false`
- **Description:** Whether the parameter is required

#### `in` (optional)
- **Type:** string
- **Values:** `query`, `header`, `path`, `body`
- **Default:** `body` for POST/PUT, `query` for GET
- **Description:** Where the parameter should be sent

#### `format` (optional)
- **Type:** string
- **Values:** `email`, `date`, `date-time`, `uuid`, `card`, custom formats
- **Description:** Additional format specification

#### `example` (optional)
- **Type:** any
- **Description:** Example value for the parameter

#### `default` (optional)
- **Type:** any
- **Description:** Default value if not provided

#### `enum` (optional)
- **Type:** array
- **Description:** List of allowed values

**Example:**
```yaml
parameters:
  - name: "amount"
    type: "number"
    required: true
    description: "Transaction amount"
    example: 200.00
  
  - name: "currency"
    type: "string"
    required: false
    default: "USD"
    enum: ["USD", "EUR", "GBP"]
```

---

## Schemas

Schemas define the structure of request and response bodies using JSON Schema.

### Request Schema

```yaml
request_schema:
  type: "object"
  required: ["field1", "field2"]
  properties:
    field1:
      type: "string"
      minLength: 3
    field2:
      type: "number"
      minimum: 0
```

### Response Schema

```yaml
response_schema:
  type: "object"
  properties:
    id:
      type: "string"
    status:
      type: "string"
      enum: ["success", "failed"]
```

**Supported JSON Schema Keywords:**
- `type`, `properties`, `required`
- `minimum`, `maximum`, `minLength`, `maxLength`
- `pattern`, `format`, `enum`
- `description`, `default`

---

## Examples

Examples help the LLM understand how to construct valid requests.

### Example Structure

```yaml
examples:
  - name: string (required)
    request: object (required)
    response: object (optional)
```

**Example:**
```yaml
examples:
  - name: "Successful authorization"
    request:
      amount: 200.00
      currency: "USD"
      card_number: "5555555555554444"
    response:
      transaction_id: "txn_123"
      status: "approved"
```

---

## Test Scenarios

Test scenarios define how natural language requests map to API calls.

### Test Scenario Structure

```yaml
test_scenarios:
  - scenario: string (required)
    natural_language: string (required)
    expected_endpoint: string (optional)
    expected_amount: number (optional)
    # ... other expected fields
```

**Example:**
```yaml
test_scenarios:
  - scenario: "Small transaction"
    natural_language: "Test authorization with amount 50"
    expected_endpoint: "authorize"
    expected_amount: 50.00
  
  - scenario: "Large transaction"
    natural_language: "Authorize payment of 5000 dollars"
    expected_endpoint: "authorize"
    expected_amount: 5000.00
```

---

## Metadata

Metadata provides additional information about the API.

### Metadata Structure

```yaml
metadata:
  api_provider: string
  environment: string
  rate_limit: string
  documentation_url: string
  support_contact: string
  version_history: array
```

**Example:**
```yaml
metadata:
  api_provider: "Mastercard"
  environment: "test"
  rate_limit: "100 requests per minute"
  documentation_url: "https://developer.mastercard.com/docs"
  support_contact: "api-support@mastercard.com"
  version_history:
    - version: "1.0.0"
      date: "2024-01-15"
      changes: "Initial release"
```

---

## Complete Example

See the following example configurations:
- `mastercard_ptc.yaml` - Payment authorization API
- `payment_api.yaml` - Generic payment processing
- `user_management_api.yaml` - User account management

---

## Best Practices

1. **Be Descriptive:** Use clear names and descriptions
2. **Include Examples:** Provide multiple examples for complex endpoints
3. **Document Parameters:** Explain what each parameter does
4. **Use Schemas:** Define request/response schemas for validation
5. **Test Scenarios:** Include natural language test scenarios
6. **Version Control:** Update version number when making changes
7. **Environment Variables:** Use `${ENV_BASE_URL}` for base URLs

---

## Environment Variable Substitution

The following environment variables can be used in configuration files:

- `${ENV_BASE_URL}` - Base URL for the API (replaced from environment config)
- `${ENV_API_KEY}` - API key (if using api_key authentication)
- `${ENV_*}` - Any environment variable matching this pattern

These are replaced at runtime based on the selected environment configuration in the system.

---

## Validation

API configurations are validated when ingested:

1. **Required Fields:** All required fields must be present
2. **Schema Validation:** JSON schemas must be valid
3. **URL Format:** URLs must be properly formatted
4. **Version Format:** Versions must follow semantic versioning
5. **Unique Names:** API names and versions must be unique

---

## Ingestion

To ingest these configurations:

```bash
# Via API
curl -X POST http://localhost:8001/api/v1/ingest/folder \
  -H "Content-Type: application/json" \
  -d '{"folder_path": "/app/api_configs", "recursive": false}'

# Via Makefile
make test-ingestion
```

---

## Extending the Schema

This schema can be extended with additional fields as needed. The system will store all fields in the metadata and make them available to the LLM service.

Custom fields are preserved and can be used for:
- Custom validation rules
- Additional metadata
- Tool-specific configuration
- Documentation links
- Examples and test data


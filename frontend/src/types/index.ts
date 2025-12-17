// Auth types
export interface User {
  id: string;
  username: string;
  role: 'admin' | 'user';
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

// API types
export interface APISpecification {
  id: string;
  name: string;
  version: string;
  source_type: string;
  source_path: string;
  content_hash: string;
  metadata: {
    base_url?: string;
    description?: string;
    endpoints?: number;
  };
  created_at: string;
  updated_at: string;
}

// LLM types
export interface ParseResult {
  intent: string;
  api_name?: string;
  endpoint?: string;
  method?: string;
  parameters?: Record<string, unknown>;
  missing_required?: string[];
  confidence: number;
  needs_clarification?: boolean;
  clarification?: Clarification;
}

export interface Clarification {
  id: string;
  message: string;
  type: 'multiple_choice' | 'free_text';
  options?: { value: string; description: string }[];
  field_name: string;
}

export interface ConstructedRequest {
  method: string;
  url: string;
  path: string;
  headers: Record<string, string>;
  query_params?: Record<string, string>;
  body?: Record<string, unknown>;
  confidence: number;
}

// Execution types
export interface ExecuteRequest {
  method: string;
  url: string;
  headers?: Record<string, string>;
  body?: unknown;
  environment_id?: string;
}

export interface ExecuteResponse {
  status_code: number;
  headers: Record<string, string>;
  body: unknown;
  execution_time_ms: number;
}

export interface Environment {
  id: string;
  name: string;
  base_url: string;
  auth_config: Record<string, unknown>;
  active: boolean;
  created_at: string;
  updated_at: string;
}

// Validation types
export interface ValidationResult {
  is_valid: boolean;
  status_check?: {
    expected: number;
    actual: number;
    is_valid: boolean;
  };
  schema_check?: {
    is_valid: boolean;
    errors?: string[];
  };
  errors?: string[];
  validated_at: string;
}

// History types
export interface TestExecution {
  id: string;
  user_id: string;
  api_spec_id?: string;
  natural_language_request: string;
  constructed_request: ConstructedRequest;
  response: ExecuteResponse;
  validation_result: ValidationResult;
  status: 'success' | 'failed' | 'error';
  execution_time_ms: number;
  created_at: string;
}

export interface HistoryResponse {
  executions: TestExecution[];
  total: number;
}


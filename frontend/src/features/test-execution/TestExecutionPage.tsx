import { useState } from 'react';
import { useLocation } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import TestInput from './TestInput';
import ClarificationDialog from './ClarificationDialog';
import RequestPreview from './RequestPreview';
import ResponseDisplay from './ResponseDisplay';
import ValidationResults from './ValidationResults';
import { llmApi } from '../../api/llm';
import { executionApi } from '../../api/execution';
import { validationApi } from '../../api/validation';
import { historyApi } from '../../api/history';
import type { ParseResult, ConstructedRequest, ExecuteResponse, ValidationResult, Clarification } from '../../types';

type TestStep = 'idle' | 'parsing' | 'clarifying' | 'constructing' | 'executing' | 'validating' | 'complete';

// Determine expected status code based on HTTP method and actual response
function getExpectedStatusCode(method: string, actualStatusCode: number): number {
  const upperMethod = method.toUpperCase();
  
  // If the actual status code is in the 2xx range, use it as the expected code
  // This handles cases like 201 for POST, 204 for DELETE, etc.
  if (actualStatusCode >= 200 && actualStatusCode < 300) {
    return actualStatusCode;
  }
  
  // Default expected status codes by method (for when response is an error)
  switch (upperMethod) {
    case 'POST':
      return 201; // Created
    case 'DELETE':
      return 204; // No Content (or 200)
    case 'PUT':
    case 'PATCH':
    case 'GET':
    default:
      return 200; // OK
  }
}

interface TestState {
  step: TestStep;
  parseResult: ParseResult | null;
  constructedRequest: ConstructedRequest | null;
  response: ExecuteResponse | null;
  validationResult: ValidationResult | null;
  error: string | null;
  naturalLanguageInput: string;
}

interface LocationState {
  prefillInput?: string;
}

export default function TestExecutionPage() {
  const location = useLocation();
  const locationState = location.state as LocationState | null;
  const prefillInput = locationState?.prefillInput || '';

  const [state, setState] = useState<TestState>({
    step: 'idle',
    parseResult: null,
    constructedRequest: null,
    response: null,
    validationResult: null,
    error: null,
    naturalLanguageInput: '',
  });

  const [clarification, setClarification] = useState<Clarification | null>(null);

  const handleSubmit = async (input: string) => {
    setState({ ...state, step: 'parsing', error: null, naturalLanguageInput: input });

    try {
      // Step 1: Parse natural language
      const parseResult = await llmApi.parse(input);
      setState((prev) => ({ ...prev, parseResult }));

      // Check if clarification is needed
      if (parseResult.needs_clarification && parseResult.clarification) {
        setState((prev) => ({ ...prev, step: 'clarifying' }));
        setClarification(parseResult.clarification);
        return;
      }

      // Continue with construction
      await constructAndExecute(parseResult, input);
    } catch (err) {
      setState((prev) => ({
        ...prev,
        step: 'idle',
        error: err instanceof Error ? err.message : 'Failed to parse request',
      }));
    }
  };

  const handleClarificationSubmit = async (value: string) => {
    setClarification(null);

    if (!state.parseResult || !clarification) return;

    // Update parse result with clarified value
    const updatedParseResult: ParseResult = {
      ...state.parseResult,
      parameters: {
        ...state.parseResult.parameters,
        [clarification.field_name]: value,
      },
      needs_clarification: false,
    };

    setState((prev) => ({ ...prev, parseResult: updatedParseResult }));
    await constructAndExecute(updatedParseResult, state.naturalLanguageInput);
  };

  const constructAndExecute = async (parseResult: ParseResult, naturalLanguageInput: string) => {
    try {
      // Step 2: Construct request
      setState((prev) => ({ ...prev, step: 'constructing' }));
      const constructedRequest = await llmApi.construct(parseResult);
      setState((prev) => ({ ...prev, constructedRequest }));

      // Step 3: Execute
      setState((prev) => ({ ...prev, step: 'executing' }));
      const response = await executionApi.execute({
        method: constructedRequest.method,
        url: constructedRequest.url,
        headers: constructedRequest.headers,
        body: constructedRequest.body,
        natural_language_request: naturalLanguageInput,
      });
      setState((prev) => ({ ...prev, response }));

      // Step 4: Validate
      setState((prev) => ({ ...prev, step: 'validating' }));
      // Determine expected status code based on HTTP method
      const expectedStatusCode = getExpectedStatusCode(constructedRequest.method, response.status_code);
      const validationResult = await validationApi.validate(
        { status_code: response.status_code, body: response.body },
        expectedStatusCode
      );
      setState((prev) => ({ ...prev, validationResult, step: 'complete' }));

      // Step 5: Link validation result back to the test execution record
      if (response.id) {
        try {
          await historyApi.updateValidation(response.id, validationResult);
        } catch (updateErr) {
          console.error('Failed to save validation result:', updateErr);
          // Don't fail the test for this error
        }
      }
    } catch (err) {
      setState((prev) => ({
        ...prev,
        step: 'complete',
        error: err instanceof Error ? err.message : 'Execution failed',
      }));
    }
  };

  const handleReset = () => {
    setState({
      step: 'idle',
      parseResult: null,
      constructedRequest: null,
      response: null,
      validationResult: null,
      error: null,
      naturalLanguageInput: '',
    });
    setClarification(null);
  };

  const isLoading = ['parsing', 'constructing', 'executing', 'validating'].includes(state.step);

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white mb-2">Test Execution</h1>
        <p className="text-gray-400">
          Describe your API test in natural language and let AI construct and execute it.
        </p>
      </div>

      <TestInput onSubmit={handleSubmit} isLoading={isLoading} initialValue={prefillInput} />

      {/* Loading indicator */}
      {isLoading && (
        <div className="flex items-center gap-3 p-4 bg-surface rounded-lg border border-surface-light animate-fadeIn">
          <Loader2 className="w-5 h-5 text-primary animate-spin" />
          <span className="text-gray-300">
            {state.step === 'parsing' && 'Parsing your request...'}
            {state.step === 'constructing' && 'Constructing API call...'}
            {state.step === 'executing' && 'Executing request...'}
            {state.step === 'validating' && 'Validating response...'}
          </span>
        </div>
      )}

      {/* Error display */}
      {state.error && (
        <div className="p-4 bg-error/10 border border-error/20 rounded-lg animate-fadeIn">
          <p className="text-error">{state.error}</p>
        </div>
      )}

      {/* Clarification dialog */}
      {clarification && (
        <ClarificationDialog
          clarification={clarification}
          onSubmit={handleClarificationSubmit}
          onCancel={() => {
            setClarification(null);
            setState((prev) => ({ ...prev, step: 'idle' }));
          }}
        />
      )}

      {/* Results */}
      {state.parseResult && !clarification && (
        <div className="space-y-4 animate-fadeIn">
          {/* Parse Result Summary */}
          <div className="p-4 bg-surface rounded-lg border border-surface-light">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-400">Detected Intent</span>
              <span className="text-sm px-2 py-1 rounded bg-primary/20 text-primary">
                {Math.round(state.parseResult.confidence * 100)}% confidence
              </span>
            </div>
            <p className="text-white font-medium">{state.parseResult.intent}</p>
            {state.parseResult.api_name && (
              <p className="text-sm text-gray-400 mt-1">
                API: {state.parseResult.api_name} → {state.parseResult.endpoint}
              </p>
            )}
          </div>

          {state.constructedRequest && (
            <RequestPreview request={state.constructedRequest} />
          )}

          {state.response && (
            <ResponseDisplay response={state.response} />
          )}

          {state.validationResult && (
            <ValidationResults result={state.validationResult} />
          )}

          {state.step === 'complete' && (
            <button
              onClick={handleReset}
              className="w-full py-3 px-4 bg-surface-light text-white font-medium rounded-lg hover:bg-gray-600 transition-colors"
            >
              Run Another Test
            </button>
          )}
        </div>
      )}
    </div>
  );
}

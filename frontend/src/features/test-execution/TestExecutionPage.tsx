import { useState } from 'react';
import { Loader2 } from 'lucide-react';
import TestInput from './TestInput';
import ClarificationDialog from './ClarificationDialog';
import RequestPreview from './RequestPreview';
import ResponseDisplay from './ResponseDisplay';
import ValidationResults from './ValidationResults';
import { llmApi } from '../../api/llm';
import { executionApi } from '../../api/execution';
import { validationApi } from '../../api/validation';
import type { ParseResult, ConstructedRequest, ExecuteResponse, ValidationResult, Clarification } from '../../types';

type TestStep = 'idle' | 'parsing' | 'clarifying' | 'constructing' | 'executing' | 'validating' | 'complete';

interface TestState {
  step: TestStep;
  parseResult: ParseResult | null;
  constructedRequest: ConstructedRequest | null;
  response: ExecuteResponse | null;
  validationResult: ValidationResult | null;
  error: string | null;
}

export default function TestExecutionPage() {
  const [state, setState] = useState<TestState>({
    step: 'idle',
    parseResult: null,
    constructedRequest: null,
    response: null,
    validationResult: null,
    error: null,
  });

  const [clarification, setClarification] = useState<Clarification | null>(null);

  const handleSubmit = async (input: string) => {
    setState({ ...state, step: 'parsing', error: null });

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
      await constructAndExecute(parseResult);
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
    await constructAndExecute(updatedParseResult);
  };

  const constructAndExecute = async (parseResult: ParseResult) => {
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
      });
      setState((prev) => ({ ...prev, response }));

      // Step 4: Validate
      setState((prev) => ({ ...prev, step: 'validating' }));
      const validationResult = await validationApi.validate(
        { status_code: response.status_code, body: response.body },
        200
      );
      setState((prev) => ({ ...prev, validationResult, step: 'complete' }));
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

      <TestInput onSubmit={handleSubmit} isLoading={isLoading} />

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


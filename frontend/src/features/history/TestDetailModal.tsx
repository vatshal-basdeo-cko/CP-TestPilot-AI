import { X, Play, CheckCircle, XCircle, Clock } from 'lucide-react';
import JsonViewer from '../../components/JsonViewer';
import type { TestExecution } from '../../types';

interface TestDetailModalProps {
  execution: TestExecution;
  onClose: () => void;
}

export default function TestDetailModal({ execution, onClose }: TestDetailModalProps) {
  const isSuccess = execution.status === 'success';

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-surface rounded-xl border border-surface-light max-w-4xl w-full max-h-[90vh] overflow-hidden animate-slideIn">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-surface-light">
          <div className="flex items-center gap-3">
            <div className={`p-2 rounded-lg ${isSuccess ? 'bg-success/20' : 'bg-error/20'}`}>
              {isSuccess ? (
                <CheckCircle className="w-5 h-5 text-success" />
              ) : (
                <XCircle className="w-5 h-5 text-error" />
              )}
            </div>
            <div>
              <h2 className="text-lg font-semibold text-white">Test Execution Details</h2>
              <p className="text-sm text-gray-400">
                {new Date(execution.created_at).toLocaleString()}
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-white transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-4 space-y-4 overflow-y-auto max-h-[calc(90vh-130px)]">
          {/* Natural Language Request */}
          <div>
            <h3 className="text-sm font-medium text-gray-400 mb-2">Natural Language Request</h3>
            <p className="p-3 bg-background rounded-lg text-white">
              {execution.natural_language_request}
            </p>
          </div>

          {/* Metrics */}
          <div className="flex gap-4">
            <div className="flex items-center gap-2 px-4 py-2 bg-background rounded-lg">
              <Clock className="w-4 h-4 text-primary" />
              <span className="text-gray-400">Duration:</span>
              <span className="text-white font-medium">{execution.execution_time_ms}ms</span>
            </div>
            <div className="flex items-center gap-2 px-4 py-2 bg-background rounded-lg">
              <span className="text-gray-400">Status Code:</span>
              <span className="text-white font-medium">{execution.response?.status_code || '-'}</span>
            </div>
          </div>

          {/* Constructed Request */}
          {execution.constructed_request && (
            <div>
              <h3 className="text-sm font-medium text-gray-400 mb-2">Constructed Request</h3>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <span className="px-2 py-1 text-xs font-medium bg-primary/20 text-primary rounded">
                    {execution.constructed_request.method}
                  </span>
                  <code className="text-sm text-gray-300 break-all">
                    {execution.constructed_request.url}
                  </code>
                </div>
                {execution.constructed_request.body && (
                  <JsonViewer data={execution.constructed_request.body} maxHeight="200px" />
                )}
              </div>
            </div>
          )}

          {/* Response */}
          {execution.response && (
            <div>
              <h3 className="text-sm font-medium text-gray-400 mb-2">Response</h3>
              <JsonViewer data={execution.response.body} maxHeight="300px" />
            </div>
          )}

          {/* Validation Result */}
          {execution.validation_result && (
            <div>
              <h3 className="text-sm font-medium text-gray-400 mb-2">Validation Result</h3>
              <div className={`p-3 rounded-lg border ${
                execution.validation_result.is_valid
                  ? 'bg-success/10 border-success/20'
                  : 'bg-error/10 border-error/20'
              }`}>
                <div className="flex items-center gap-2 mb-2">
                  {execution.validation_result.is_valid ? (
                    <CheckCircle className="w-4 h-4 text-success" />
                  ) : (
                    <XCircle className="w-4 h-4 text-error" />
                  )}
                  <span className={execution.validation_result.is_valid ? 'text-success' : 'text-error'}>
                    {execution.validation_result.is_valid ? 'Validation Passed' : 'Validation Failed'}
                  </span>
                </div>
                {execution.validation_result.errors && execution.validation_result.errors.length > 0 && (
                  <ul className="text-sm text-gray-300 space-y-1">
                    {execution.validation_result.errors.map((error, index) => (
                      <li key={index}>• {error}</li>
                    ))}
                  </ul>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-3 p-4 border-t border-surface-light">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
          >
            Close
          </button>
          <button
            className="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
          >
            <Play className="w-4 h-4" />
            Re-run Test
          </button>
        </div>
      </div>
    </div>
  );
}


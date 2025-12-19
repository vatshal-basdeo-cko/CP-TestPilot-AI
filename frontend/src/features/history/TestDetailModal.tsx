import { useState } from 'react';
import { X, Play, CheckCircle, XCircle, Clock, Trash2, AlertTriangle } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import JsonViewer from '../../components/JsonViewer';
import type { TestExecution } from '../../types';

interface TestDetailModalProps {
  execution: TestExecution;
  onClose: () => void;
  onDelete?: (id: string) => Promise<void>;
}

export default function TestDetailModal({ execution, onClose, onDelete }: TestDetailModalProps) {
  const navigate = useNavigate();
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const isSuccess = execution.status === 'success';

  const handleRerunTest = () => {
    onClose();
    // Navigate to test execution page with prefilled input
    navigate('/test', { 
      state: { 
        prefillInput: execution.natural_language_request 
      } 
    });
  };

  const handleDelete = async () => {
    if (!onDelete) return;
    setIsDeleting(true);
    try {
      await onDelete(execution.id);
      onClose();
    } catch (error) {
      console.error('Failed to delete execution:', error);
    } finally {
      setIsDeleting(false);
    }
  };

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
        <div className="flex items-center justify-between p-4 border-t border-surface-light">
          <div>
            {onDelete && (
              <button
                onClick={() => setShowDeleteConfirm(true)}
                className="flex items-center gap-2 px-4 py-2 bg-error/20 text-error rounded-lg hover:bg-error/30 transition-colors"
              >
                <Trash2 className="w-4 h-4" />
                Delete
              </button>
            )}
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={onClose}
              className="px-4 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
            >
              Close
            </button>
            <button
              onClick={handleRerunTest}
              className="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
            >
              <Play className="w-4 h-4" />
              Re-run Test
            </button>
          </div>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      {showDeleteConfirm && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-[60]">
          <div className="bg-surface rounded-xl border border-surface-light p-6 max-w-md w-full mx-4 animate-slideIn">
            <div className="flex items-center gap-3 mb-4">
              <div className="p-3 bg-error/20 rounded-lg">
                <AlertTriangle className="w-6 h-6 text-error" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white">Delete Test Execution</h3>
                <p className="text-sm text-gray-400">This action cannot be undone</p>
              </div>
            </div>
            <p className="text-gray-300 mb-6">
              Are you sure you want to delete this test execution record? This will permanently remove it from your history.
            </p>
            <div className="flex items-center justify-end gap-3">
              <button
                onClick={() => setShowDeleteConfirm(false)}
                disabled={isDeleting}
                className="px-4 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                onClick={handleDelete}
                disabled={isDeleting}
                className="flex items-center gap-2 px-4 py-2 bg-error text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50"
              >
                {isDeleting ? (
                  <>
                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Deleting...
                  </>
                ) : (
                  <>
                    <Trash2 className="w-4 h-4" />
                    Delete
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}


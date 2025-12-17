import { CheckCircle, XCircle, AlertCircle, Shield } from 'lucide-react';
import CollapsibleSection from '../../components/CollapsibleSection';
import type { ValidationResult } from '../../types';

interface ValidationResultsProps {
  result: ValidationResult;
}

export default function ValidationResults({ result }: ValidationResultsProps) {
  const StatusIcon = result.is_valid ? CheckCircle : XCircle;
  const statusColor = result.is_valid ? 'text-success' : 'text-error';
  const statusBg = result.is_valid ? 'bg-success/20' : 'bg-error/20';

  return (
    <CollapsibleSection
      title="Validation Results"
      icon={<Shield className="w-5 h-5" />}
      defaultOpen={true}
      badge={
        <span className={`flex items-center gap-1 px-2 py-1 text-xs font-medium rounded ${statusBg} ${statusColor}`}>
          <StatusIcon className="w-3 h-3" />
          {result.is_valid ? 'Passed' : 'Failed'}
        </span>
      }
    >
      <div className="space-y-3">
        {/* Status Code Check */}
        {result.status_check && (
          <div className="flex items-center gap-3 p-3 bg-background rounded-lg">
            {result.status_check.is_valid ? (
              <CheckCircle className="w-5 h-5 text-success" />
            ) : (
              <XCircle className="w-5 h-5 text-error" />
            )}
            <div>
              <span className="text-white">Status Code</span>
              <p className="text-sm text-gray-400">
                Expected: {result.status_check.expected}, Actual: {result.status_check.actual}
              </p>
            </div>
          </div>
        )}

        {/* Schema Check */}
        {result.schema_check && (
          <div className="flex items-center gap-3 p-3 bg-background rounded-lg">
            {result.schema_check.is_valid ? (
              <CheckCircle className="w-5 h-5 text-success" />
            ) : (
              <XCircle className="w-5 h-5 text-error" />
            )}
            <div>
              <span className="text-white">Schema Validation</span>
              {result.schema_check.errors && result.schema_check.errors.length > 0 && (
                <ul className="text-sm text-gray-400 mt-1">
                  {result.schema_check.errors.map((error, index) => (
                    <li key={index}>• {error}</li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        )}

        {/* General Errors */}
        {result.errors && result.errors.length > 0 && (
          <div className="p-3 bg-error/10 rounded-lg border border-error/20">
            <div className="flex items-center gap-2 text-error mb-2">
              <AlertCircle className="w-4 h-4" />
              <span className="font-medium">Errors</span>
            </div>
            <ul className="text-sm text-gray-300 space-y-1">
              {result.errors.map((error, index) => (
                <li key={index}>• {error}</li>
              ))}
            </ul>
          </div>
        )}

        {/* Timestamp */}
        <p className="text-xs text-gray-500">
          Validated at: {new Date(result.validated_at).toLocaleString()}
        </p>
      </div>
    </CollapsibleSection>
  );
}


import { FileJson, Clock } from 'lucide-react';
import CollapsibleSection from '../../components/CollapsibleSection';
import JsonViewer from '../../components/JsonViewer';
import type { ExecuteResponse } from '../../types';

interface ResponseDisplayProps {
  response: ExecuteResponse;
}

export default function ResponseDisplay({ response }: ResponseDisplayProps) {
  const getStatusColor = (status: number) => {
    if (status >= 200 && status < 300) return 'bg-success/20 text-success';
    if (status >= 400 && status < 500) return 'bg-warning/20 text-warning';
    if (status >= 500) return 'bg-error/20 text-error';
    return 'bg-gray-600 text-gray-300';
  };

  return (
    <CollapsibleSection
      title="API Response"
      icon={<FileJson className="w-5 h-5" />}
      defaultOpen={true}
      badge={
        <div className="flex items-center gap-2">
          <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusColor(response.status_code)}`}>
            {response.status_code}
          </span>
          <span className="flex items-center gap-1 text-xs text-gray-400">
            <Clock className="w-3 h-3" />
            {response.execution_time_ms}ms
          </span>
        </div>
      }
    >
      <div className="space-y-4">
        {/* Response Headers */}
        {response.headers && Object.keys(response.headers).length > 0 && (
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">Response Headers</label>
            <div className="bg-background rounded-lg p-3 space-y-1 max-h-32 overflow-auto">
              {Object.entries(response.headers).map(([key, value]) => (
                <div key={key} className="font-mono text-sm">
                  <span className="text-cyan-400">{key}</span>
                  <span className="text-gray-500">: </span>
                  <span className="text-gray-300">{value}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Response Body */}
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">Response Body</label>
          <JsonViewer data={response.body} />
        </div>
      </div>
    </CollapsibleSection>
  );
}


import { Code } from 'lucide-react';
import CollapsibleSection from '../../components/CollapsibleSection';
import JsonViewer from '../../components/JsonViewer';
import type { ConstructedRequest } from '../../types';

interface RequestPreviewProps {
  request: ConstructedRequest;
}

export default function RequestPreview({ request }: RequestPreviewProps) {
  const methodColors: Record<string, string> = {
    GET: 'bg-success/20 text-success',
    POST: 'bg-primary/20 text-primary',
    PUT: 'bg-warning/20 text-warning',
    PATCH: 'bg-warning/20 text-warning',
    DELETE: 'bg-error/20 text-error',
  };

  return (
    <CollapsibleSection
      title="Constructed Request"
      icon={<Code className="w-5 h-5" />}
      defaultOpen={true}
      badge={
        <span className={`px-2 py-1 text-xs font-medium rounded ${methodColors[request.method] || 'bg-gray-600 text-gray-300'}`}>
          {request.method}
        </span>
      }
    >
      <div className="space-y-4">
        {/* URL */}
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">URL</label>
          <code className="block p-3 bg-background rounded-lg text-primary font-mono text-sm break-all">
            {request.url}
          </code>
        </div>

        {/* Headers */}
        {request.headers && Object.keys(request.headers).length > 0 && (
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">Headers</label>
            <div className="bg-background rounded-lg p-3 space-y-1">
              {Object.entries(request.headers).map(([key, value]) => (
                <div key={key} className="font-mono text-sm">
                  <span className="text-cyan-400">{key}</span>
                  <span className="text-gray-500">: </span>
                  <span className="text-gray-300">{value}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Query Parameters */}
        {request.query_params && Object.keys(request.query_params).length > 0 && (
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">Query Parameters</label>
            <JsonViewer data={request.query_params} maxHeight="200px" />
          </div>
        )}

        {/* Body */}
        {request.body && (
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">Request Body</label>
            <JsonViewer data={request.body} />
          </div>
        )}
      </div>
    </CollapsibleSection>
  );
}


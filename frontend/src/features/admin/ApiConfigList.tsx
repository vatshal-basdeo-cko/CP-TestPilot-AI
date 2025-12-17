import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Trash2, FileJson, Loader2, FolderSync } from 'lucide-react';
import { ingestionApi } from '../../api/ingestion';
import type { APISpecification } from '../../types';

export default function ApiConfigList() {
  const queryClient = useQueryClient();
  const [isIngesting, setIsIngesting] = useState(false);

  const { data, isLoading, error } = useQuery({
    queryKey: ['apis'],
    queryFn: ingestionApi.listAPIs,
  });

  const ingestMutation = useMutation({
    mutationFn: (folderPath: string) => ingestionApi.ingestFolder(folderPath),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apis'] });
      setIsIngesting(false);
    },
    onError: () => {
      setIsIngesting(false);
    },
  });

  const handleIngest = () => {
    setIsIngesting(true);
    ingestMutation.mutate('/app/api_configs');
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="w-6 h-6 text-primary animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 bg-error/10 border border-error/20 rounded-lg">
        <p className="text-error">Failed to load API configurations</p>
      </div>
    );
  }

  const apis = data?.apis || [];

  return (
    <div className="space-y-4">
      {/* Actions */}
      <div className="flex items-center justify-between">
        <p className="text-gray-400">
          {apis.length} API configuration{apis.length !== 1 ? 's' : ''} ingested
        </p>
        <button
          onClick={handleIngest}
          disabled={isIngesting}
          className="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark disabled:opacity-50 transition-colors"
        >
          {isIngesting ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            <FolderSync className="w-4 h-4" />
          )}
          Re-ingest from Folder
        </button>
      </div>

      {/* Ingest Result */}
      {ingestMutation.isSuccess && ingestMutation.data && (
        <div className="p-3 bg-success/10 border border-success/20 rounded-lg text-success text-sm animate-fadeIn">
          Ingested: {ingestMutation.data.ingested}, Skipped: {ingestMutation.data.skipped}, Failed: {ingestMutation.data.failed}
        </div>
      )}

      {/* API List */}
      {apis.length === 0 ? (
        <div className="text-center py-12 bg-surface rounded-lg border border-surface-light">
          <FileJson className="w-12 h-12 text-gray-600 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-white mb-2">No APIs Ingested</h3>
          <p className="text-gray-400 mb-4">
            Click "Re-ingest from Folder" to load API configurations.
          </p>
        </div>
      ) : (
        <div className="bg-surface rounded-lg border border-surface-light overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-surface-light">
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Name</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Version</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Source</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Endpoints</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Updated</th>
                <th className="px-4 py-3 text-right text-sm font-medium text-gray-400">Actions</th>
              </tr>
            </thead>
            <tbody>
              {apis.map((api: APISpecification) => (
                <tr key={api.id} className="border-b border-surface-light last:border-0 hover:bg-surface-light/50">
                  <td className="px-4 py-3">
                    <div>
                      <p className="text-white font-medium">{api.name}</p>
                      <p className="text-xs text-gray-500">{api.metadata?.description?.slice(0, 50)}...</p>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <span className="px-2 py-1 text-xs bg-primary/20 text-primary rounded">
                      v{api.version}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-400">
                    {api.source_type}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-400">
                    {api.metadata?.endpoints || '-'}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-400">
                    {new Date(api.updated_at).toLocaleDateString()}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      className="p-2 text-gray-400 hover:text-error transition-colors"
                      title="Delete"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}


import { CheckCircle, XCircle, AlertCircle, Clock, Eye } from 'lucide-react';
import type { TestExecution } from '../../types';

interface HistoryTableProps {
  executions: TestExecution[];
  onSelect: (execution: TestExecution) => void;
}

export default function HistoryTable({ executions, onSelect }: HistoryTableProps) {
  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="w-4 h-4 text-success" />;
      case 'failed':
        return <XCircle className="w-4 h-4 text-error" />;
      default:
        return <AlertCircle className="w-4 h-4 text-warning" />;
    }
  };

  const getStatusBadge = (status: string) => {
    const styles = {
      success: 'bg-success/20 text-success',
      failed: 'bg-error/20 text-error',
      error: 'bg-warning/20 text-warning',
    };
    return styles[status as keyof typeof styles] || 'bg-gray-600 text-gray-300';
  };

  return (
    <div className="bg-surface rounded-lg border border-surface-light overflow-hidden">
      <table className="w-full">
        <thead>
          <tr className="border-b border-surface-light">
            <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Status</th>
            <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Request</th>
            <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Method</th>
            <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Duration</th>
            <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Date</th>
            <th className="px-4 py-3 text-right text-sm font-medium text-gray-400">Actions</th>
          </tr>
        </thead>
        <tbody>
          {executions.map((execution) => (
            <tr
              key={execution.id}
              className="border-b border-surface-light last:border-0 hover:bg-surface-light/50 cursor-pointer"
              onClick={() => onSelect(execution)}
            >
              <td className="px-4 py-3">
                <span className={`flex items-center gap-2 px-2 py-1 text-xs font-medium rounded w-fit ${getStatusBadge(execution.status)}`}>
                  {getStatusIcon(execution.status)}
                  {execution.status}
                </span>
              </td>
              <td className="px-4 py-3">
                <p className="text-white text-sm truncate max-w-xs" title={execution.natural_language_request}>
                  {execution.natural_language_request}
                </p>
              </td>
              <td className="px-4 py-3">
                <span className="text-sm text-gray-400">
                  {execution.constructed_request?.method || '-'}
                </span>
              </td>
              <td className="px-4 py-3">
                <span className="flex items-center gap-1 text-sm text-gray-400">
                  <Clock className="w-3 h-3" />
                  {execution.execution_time_ms}ms
                </span>
              </td>
              <td className="px-4 py-3 text-sm text-gray-400">
                {new Date(execution.created_at).toLocaleString()}
              </td>
              <td className="px-4 py-3 text-right">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onSelect(execution);
                  }}
                  className="p-2 text-gray-400 hover:text-primary transition-colors"
                  title="View Details"
                >
                  <Eye className="w-4 h-4" />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}


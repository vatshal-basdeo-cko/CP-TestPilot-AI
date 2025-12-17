import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Calendar, Filter, Loader2, Search, RefreshCw } from 'lucide-react';
import { historyApi, HistoryFilters } from '../../api/history';
import HistoryTable from './HistoryTable';
import TestDetailModal from './TestDetailModal';
import type { TestExecution } from '../../types';

export default function HistoryPage() {
  const [filters, setFilters] = useState<HistoryFilters>({
    limit: 20,
    offset: 0,
  });
  const [selectedExecution, setSelectedExecution] = useState<TestExecution | null>(null);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['history', filters],
    queryFn: () => historyApi.list(filters),
  });

  const handleStatusFilter = (status: string) => {
    if (status === 'all') {
      const { status: _, ...rest } = filters;
      setFilters({ ...rest, offset: 0 });
    } else {
      setFilters({ ...filters, status: status as 'success' | 'failed' | 'error', offset: 0 });
    }
  };

  const handlePageChange = (newOffset: number) => {
    setFilters({ ...filters, offset: newOffset });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white mb-2">Test History</h1>
          <p className="text-gray-400">View and analyze past test executions.</p>
        </div>
        <button
          onClick={() => refetch()}
          className="flex items-center gap-2 px-4 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
        >
          <RefreshCw className="w-4 h-4" />
          Refresh
        </button>
      </div>

      {/* Filters */}
      <div className="bg-surface rounded-lg border border-surface-light p-4">
        <div className="flex flex-wrap items-center gap-4">
          {/* Status Filter */}
          <div className="flex items-center gap-2">
            <Filter className="w-4 h-4 text-gray-400" />
            <select
              value={filters.status || 'all'}
              onChange={(e) => handleStatusFilter(e.target.value)}
              className="px-3 py-2 bg-background border border-surface-light rounded-lg text-white text-sm"
            >
              <option value="all">All Status</option>
              <option value="success">Success</option>
              <option value="failed">Failed</option>
              <option value="error">Error</option>
            </select>
          </div>

          {/* Date Filters */}
          <div className="flex items-center gap-2">
            <Calendar className="w-4 h-4 text-gray-400" />
            <input
              type="date"
              value={filters.from_date || ''}
              onChange={(e) => setFilters({ ...filters, from_date: e.target.value || undefined, offset: 0 })}
              className="px-3 py-2 bg-background border border-surface-light rounded-lg text-white text-sm"
              placeholder="From"
            />
            <span className="text-gray-500">to</span>
            <input
              type="date"
              value={filters.to_date || ''}
              onChange={(e) => setFilters({ ...filters, to_date: e.target.value || undefined, offset: 0 })}
              className="px-3 py-2 bg-background border border-surface-light rounded-lg text-white text-sm"
              placeholder="To"
            />
          </div>

          {/* Result count */}
          {data && (
            <div className="ml-auto text-sm text-gray-400">
              Showing {data.executions?.length || 0} of {data.total || 0} results
            </div>
          )}
        </div>
      </div>

      {/* Content */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="w-6 h-6 text-primary animate-spin" />
        </div>
      ) : error ? (
        <div className="p-4 bg-error/10 border border-error/20 rounded-lg">
          <p className="text-error">Failed to load test history</p>
        </div>
      ) : data?.executions && data.executions.length > 0 ? (
        <>
          <HistoryTable
            executions={data.executions}
            onSelect={setSelectedExecution}
          />

          {/* Pagination */}
          <div className="flex items-center justify-between">
            <button
              onClick={() => handlePageChange(Math.max(0, (filters.offset || 0) - (filters.limit || 20)))}
              disabled={(filters.offset || 0) === 0}
              className="px-4 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Previous
            </button>
            <span className="text-gray-400">
              Page {Math.floor((filters.offset || 0) / (filters.limit || 20)) + 1}
            </span>
            <button
              onClick={() => handlePageChange((filters.offset || 0) + (filters.limit || 20))}
              disabled={data.executions.length < (filters.limit || 20)}
              className="px-4 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Next
            </button>
          </div>
        </>
      ) : (
        <div className="text-center py-12 bg-surface rounded-lg border border-surface-light">
          <Search className="w-12 h-12 text-gray-600 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-white mb-2">No Test History</h3>
          <p className="text-gray-400">Run some tests to see your history here.</p>
        </div>
      )}

      {/* Detail Modal */}
      {selectedExecution && (
        <TestDetailModal
          execution={selectedExecution}
          onClose={() => setSelectedExecution(null)}
        />
      )}
    </div>
  );
}


import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { 
  BarChart3, 
  TrendingUp, 
  TrendingDown, 
  Clock, 
  CheckCircle, 
  XCircle, 
  Loader2,
  Calendar,
  RefreshCw,
  Activity
} from 'lucide-react';
import { analyticsApi, AnalyticsFilters } from '../../api/analytics';

export default function AnalyticsPage() {
  const [filters, setFilters] = useState<AnalyticsFilters>({});

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['analytics', filters],
    queryFn: () => analyticsApi.getOverview(filters),
  });

  const successRate = data?.success_rate ?? 0;
  const isHighSuccessRate = successRate >= 80;
  const isMediumSuccessRate = successRate >= 50 && successRate < 80;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white mb-2">Analytics Dashboard</h1>
          <p className="text-gray-400">Monitor your API testing performance and trends.</p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => refetch()}
            className="flex items-center gap-2 px-4 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </button>
        </div>
      </div>

      {/* Date Filters */}
      <div className="bg-surface rounded-lg border border-surface-light p-4">
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex items-center gap-2">
            <Calendar className="w-4 h-4 text-gray-400" />
            <span className="text-sm text-gray-400">Date Range:</span>
          </div>
          <input
            type="date"
            value={filters.start_date || ''}
            onChange={(e) => setFilters({ ...filters, start_date: e.target.value || undefined })}
            className="px-3 py-2 bg-background border border-surface-light rounded-lg text-white text-sm"
          />
          <span className="text-gray-500">to</span>
          <input
            type="date"
            value={filters.end_date || ''}
            onChange={(e) => setFilters({ ...filters, end_date: e.target.value || undefined })}
            className="px-3 py-2 bg-background border border-surface-light rounded-lg text-white text-sm"
          />
          {(filters.start_date || filters.end_date) && (
            <button
              onClick={() => setFilters({})}
              className="px-3 py-2 text-sm text-gray-400 hover:text-white transition-colors"
            >
              Clear
            </button>
          )}
        </div>
      </div>

      {/* Content */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="w-8 h-8 text-primary animate-spin" />
        </div>
      ) : error ? (
        <div className="p-4 bg-error/10 border border-error/20 rounded-lg">
          <p className="text-error">Failed to load analytics data</p>
        </div>
      ) : data ? (
        <>
          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Total Tests */}
            <div className="bg-surface rounded-xl border border-surface-light p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="p-3 bg-primary/20 rounded-lg">
                  <BarChart3 className="w-6 h-6 text-primary" />
                </div>
                <Activity className="w-5 h-5 text-gray-500" />
              </div>
              <div className="text-3xl font-bold text-white mb-1">
                {data.total_tests.toLocaleString()}
              </div>
              <div className="text-sm text-gray-400">Total Tests</div>
            </div>

            {/* Success Rate */}
            <div className="bg-surface rounded-xl border border-surface-light p-6">
              <div className="flex items-center justify-between mb-4">
                <div className={`p-3 rounded-lg ${
                  isHighSuccessRate ? 'bg-success/20' : isMediumSuccessRate ? 'bg-warning/20' : 'bg-error/20'
                }`}>
                  {isHighSuccessRate ? (
                    <TrendingUp className="w-6 h-6 text-success" />
                  ) : (
                    <TrendingDown className={`w-6 h-6 ${isMediumSuccessRate ? 'text-warning' : 'text-error'}`} />
                  )}
                </div>
                <span className={`text-xs font-medium px-2 py-1 rounded ${
                  isHighSuccessRate ? 'bg-success/20 text-success' : 
                  isMediumSuccessRate ? 'bg-warning/20 text-warning' : 'bg-error/20 text-error'
                }`}>
                  {isHighSuccessRate ? 'Good' : isMediumSuccessRate ? 'Fair' : 'Poor'}
                </span>
              </div>
              <div className="text-3xl font-bold text-white mb-1">
                {successRate.toFixed(1)}%
              </div>
              <div className="text-sm text-gray-400">Success Rate</div>
              {/* Progress bar */}
              <div className="mt-3 h-2 bg-surface-light rounded-full overflow-hidden">
                <div 
                  className={`h-full transition-all duration-500 ${
                    isHighSuccessRate ? 'bg-success' : isMediumSuccessRate ? 'bg-warning' : 'bg-error'
                  }`}
                  style={{ width: `${successRate}%` }}
                />
              </div>
            </div>

            {/* Successful Tests */}
            <div className="bg-surface rounded-xl border border-surface-light p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="p-3 bg-success/20 rounded-lg">
                  <CheckCircle className="w-6 h-6 text-success" />
                </div>
                <span className="text-xs font-medium px-2 py-1 rounded bg-success/20 text-success">
                  Passed
                </span>
              </div>
              <div className="text-3xl font-bold text-white mb-1">
                {data.successful_tests.toLocaleString()}
              </div>
              <div className="text-sm text-gray-400">Successful Tests</div>
            </div>

            {/* Failed Tests */}
            <div className="bg-surface rounded-xl border border-surface-light p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="p-3 bg-error/20 rounded-lg">
                  <XCircle className="w-6 h-6 text-error" />
                </div>
                <span className="text-xs font-medium px-2 py-1 rounded bg-error/20 text-error">
                  Failed
                </span>
              </div>
              <div className="text-3xl font-bold text-white mb-1">
                {data.failed_tests.toLocaleString()}
              </div>
              <div className="text-sm text-gray-400">Failed Tests</div>
            </div>
          </div>

          {/* Secondary Stats */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Average Execution Time */}
            <div className="bg-surface rounded-xl border border-surface-light p-6">
              <div className="flex items-center gap-3 mb-4">
                <div className="p-3 bg-primary/20 rounded-lg">
                  <Clock className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-white">Average Execution Time</h3>
                  <p className="text-sm text-gray-400">Mean response time across all tests</p>
                </div>
              </div>
              <div className="flex items-baseline gap-2">
                <span className="text-4xl font-bold text-white">
                  {data.avg_execution_time_ms?.toFixed(0) || 0}
                </span>
                <span className="text-xl text-gray-400">ms</span>
              </div>
              {/* Performance indicator */}
              <div className="mt-4 flex items-center gap-2">
                {(data.avg_execution_time_ms || 0) < 500 ? (
                  <>
                    <div className="w-2 h-2 rounded-full bg-success" />
                    <span className="text-sm text-success">Excellent performance</span>
                  </>
                ) : (data.avg_execution_time_ms || 0) < 2000 ? (
                  <>
                    <div className="w-2 h-2 rounded-full bg-warning" />
                    <span className="text-sm text-warning">Moderate performance</span>
                  </>
                ) : (
                  <>
                    <div className="w-2 h-2 rounded-full bg-error" />
                    <span className="text-sm text-error">Slow performance</span>
                  </>
                )}
              </div>
            </div>

            {/* Test Distribution */}
            <div className="bg-surface rounded-xl border border-surface-light p-6">
              <div className="flex items-center gap-3 mb-4">
                <div className="p-3 bg-primary/20 rounded-lg">
                  <BarChart3 className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-white">Test Distribution</h3>
                  <p className="text-sm text-gray-400">Breakdown of test outcomes</p>
                </div>
              </div>
              {data.total_tests > 0 ? (
                <div className="space-y-4">
                  {/* Visual bar */}
                  <div className="h-8 rounded-lg overflow-hidden flex">
                    <div 
                      className="bg-success transition-all duration-500 flex items-center justify-center"
                      style={{ width: `${(data.successful_tests / data.total_tests) * 100}%` }}
                    >
                      {data.successful_tests > 0 && (
                        <span className="text-xs font-medium text-white">
                          {((data.successful_tests / data.total_tests) * 100).toFixed(0)}%
                        </span>
                      )}
                    </div>
                    <div 
                      className="bg-error transition-all duration-500 flex items-center justify-center"
                      style={{ width: `${(data.failed_tests / data.total_tests) * 100}%` }}
                    >
                      {data.failed_tests > 0 && (
                        <span className="text-xs font-medium text-white">
                          {((data.failed_tests / data.total_tests) * 100).toFixed(0)}%
                        </span>
                      )}
                    </div>
                  </div>
                  {/* Legend */}
                  <div className="flex items-center gap-6">
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 rounded bg-success" />
                      <span className="text-sm text-gray-400">Success ({data.successful_tests})</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 rounded bg-error" />
                      <span className="text-sm text-gray-400">Failed ({data.failed_tests})</span>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  No test data available
                </div>
              )}
            </div>
          </div>

          {/* Top APIs Section (if available) */}
          {data.top_apis && data.top_apis.length > 0 && (
            <div className="bg-surface rounded-xl border border-surface-light p-6">
              <div className="flex items-center gap-3 mb-6">
                <div className="p-3 bg-primary/20 rounded-lg">
                  <TrendingUp className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-white">Top APIs by Usage</h3>
                  <p className="text-sm text-gray-400">Most frequently tested APIs</p>
                </div>
              </div>
              <div className="space-y-4">
                {data.top_apis.map((api, index) => (
                  <div key={index} className="flex items-center gap-4">
                    <div className="w-8 h-8 rounded-lg bg-surface-light flex items-center justify-center text-sm font-medium text-gray-400">
                      {index + 1}
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between mb-1">
                        <span className="font-medium text-white">{api.api_name || 'Unknown API'}</span>
                        <span className="text-sm text-gray-400">{api.test_count} tests</span>
                      </div>
                      <div className="h-2 bg-surface-light rounded-full overflow-hidden">
                        <div 
                          className="h-full bg-primary transition-all duration-500"
                          style={{ width: `${api.success_rate}%` }}
                        />
                      </div>
                    </div>
                    <div className={`text-sm font-medium ${
                      api.success_rate >= 80 ? 'text-success' : 
                      api.success_rate >= 50 ? 'text-warning' : 'text-error'
                    }`}>
                      {api.success_rate.toFixed(0)}%
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </>
      ) : null}
    </div>
  );
}



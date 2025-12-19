import type { TestExecution } from '../types';

/**
 * Export data as JSON file
 */
export function exportToJSON(data: TestExecution[], filename = 'test-history') {
  const jsonStr = JSON.stringify(data, null, 2);
  const blob = new Blob([jsonStr], { type: 'application/json' });
  downloadBlob(blob, `${filename}-${getTimestamp()}.json`);
}

/**
 * Export data as CSV file
 */
export function exportToCSV(data: TestExecution[], filename = 'test-history') {
  if (data.length === 0) {
    return;
  }

  // Define CSV headers
  const headers = [
    'ID',
    'Status',
    'Natural Language Request',
    'Method',
    'URL',
    'Status Code',
    'Execution Time (ms)',
    'Created At',
    'Validation Valid',
    'Errors',
  ];

  // Build CSV rows
  const rows = data.map((execution) => [
    execution.id,
    execution.status,
    escapeCSV(execution.natural_language_request || ''),
    execution.constructed_request?.method || '',
    escapeCSV(execution.constructed_request?.url || ''),
    execution.response?.status_code || '',
    execution.execution_time_ms,
    execution.created_at,
    execution.validation_result?.is_valid ? 'Yes' : 'No',
    escapeCSV((execution.validation_result?.errors || []).join('; ')),
  ]);

  // Combine headers and rows
  const csvContent = [
    headers.join(','),
    ...rows.map((row) => row.join(',')),
  ].join('\n');

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  downloadBlob(blob, `${filename}-${getTimestamp()}.csv`);
}

/**
 * Escape CSV values that contain commas, quotes, or newlines
 */
function escapeCSV(value: string): string {
  if (value.includes(',') || value.includes('"') || value.includes('\n')) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}

/**
 * Get current timestamp for filename
 */
function getTimestamp(): string {
  const now = new Date();
  return now.toISOString().split('T')[0];
}

/**
 * Trigger browser download of a blob
 */
function downloadBlob(blob: Blob, filename: string): void {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', filename);
  document.body.appendChild(link);
  link.click();
  link.parentNode?.removeChild(link);
  window.URL.revokeObjectURL(url);
}



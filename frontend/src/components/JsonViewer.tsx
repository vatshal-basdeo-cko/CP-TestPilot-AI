import { Copy, Check } from 'lucide-react';
import { useState } from 'react';

interface JsonViewerProps {
  data: unknown;
  maxHeight?: string;
}

export default function JsonViewer({ data, maxHeight = '400px' }: JsonViewerProps) {
  const [copied, setCopied] = useState(false);
  const jsonString = JSON.stringify(data, null, 2);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(jsonString);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  // Syntax highlighting
  const highlightJson = (json: string) => {
    return json
      .replace(/"([^"]+)":/g, '<span class="text-cyan-400">"$1"</span>:')
      .replace(/: "([^"]*)"/g, ': <span class="text-success">"$1"</span>')
      .replace(/: (\d+)/g, ': <span class="text-warning">$1</span>')
      .replace(/: (true|false)/g, ': <span class="text-primary">$1</span>')
      .replace(/: (null)/g, ': <span class="text-gray-500">$1</span>');
  };

  return (
    <div className="relative">
      <button
        onClick={handleCopy}
        className="absolute top-2 right-2 p-2 rounded-lg bg-surface-light hover:bg-gray-600 transition-colors"
        title="Copy to clipboard"
      >
        {copied ? (
          <Check className="w-4 h-4 text-success" />
        ) : (
          <Copy className="w-4 h-4 text-gray-400" />
        )}
      </button>
      <pre
        className="font-mono text-sm bg-background rounded-lg p-4 overflow-auto"
        style={{ maxHeight }}
        dangerouslySetInnerHTML={{ __html: highlightJson(jsonString) }}
      />
    </div>
  );
}


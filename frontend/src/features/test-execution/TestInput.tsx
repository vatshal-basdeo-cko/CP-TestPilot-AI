import { useState, useEffect, FormEvent } from 'react';
import { Send, Sparkles } from 'lucide-react';

interface TestInputProps {
  onSubmit: (input: string) => void;
  isLoading: boolean;
  initialValue?: string;
}

const examplePrompts = [
  'Authorize a Mastercard payment for 150 dollars',
  'Register a new user with email test@example.com',
  'Process a refund for transaction abc123',
  'Get user profile for user ID 12345',
  'Create a new payment of 99.99 USD',
];

export default function TestInput({ onSubmit, isLoading, initialValue = '' }: TestInputProps) {
  const [input, setInput] = useState(initialValue);

  // Update input when initialValue changes (e.g., from re-run navigation)
  useEffect(() => {
    if (initialValue) {
      setInput(initialValue);
    }
  }, [initialValue]);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (input.trim() && !isLoading) {
      onSubmit(input.trim());
    }
  };

  const handleExample = (example: string) => {
    setInput(example);
  };

  return (
    <div className="bg-surface rounded-xl p-6 border border-surface-light">
      <form onSubmit={handleSubmit}>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          Describe your API test
        </label>
        <div className="relative">
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="e.g., Authorize a Mastercard payment for 200 dollars with currency EUR"
            className="w-full h-32 px-4 py-3 bg-background border border-surface-light rounded-lg text-white placeholder-gray-500 resize-none focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
            disabled={isLoading}
          />
          <button
            type="submit"
            disabled={!input.trim() || isLoading}
            className="absolute bottom-3 right-3 p-2 bg-primary hover:bg-primary-dark text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Send className="w-5 h-5" />
          </button>
        </div>
      </form>

      {/* Example prompts */}
      <div className="mt-4">
        <div className="flex items-center gap-2 text-sm text-gray-400 mb-2">
          <Sparkles className="w-4 h-4" />
          <span>Try an example:</span>
        </div>
        <div className="flex flex-wrap gap-2">
          {examplePrompts.map((example, index) => (
            <button
              key={index}
              onClick={() => handleExample(example)}
              disabled={isLoading}
              className="px-3 py-1.5 text-sm bg-surface-light text-gray-300 rounded-full hover:bg-gray-600 hover:text-white transition-colors disabled:opacity-50"
            >
              {example}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}


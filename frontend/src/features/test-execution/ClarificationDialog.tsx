import { useState, FormEvent } from 'react';
import { HelpCircle, X } from 'lucide-react';
import type { Clarification } from '../../types';

interface ClarificationDialogProps {
  clarification: Clarification;
  onSubmit: (value: string) => void;
  onCancel: () => void;
}

export default function ClarificationDialog({
  clarification,
  onSubmit,
  onCancel,
}: ClarificationDialogProps) {
  const [value, setValue] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (value.trim()) {
      onSubmit(value.trim());
    }
  };

  const handleOptionSelect = (optionValue: string) => {
    onSubmit(optionValue);
  };

  return (
    <div className="bg-surface rounded-xl p-6 border border-primary/30 animate-slideIn">
      <div className="flex items-start gap-3 mb-4">
        <div className="p-2 bg-primary/20 rounded-lg">
          <HelpCircle className="w-5 h-5 text-primary" />
        </div>
        <div className="flex-1">
          <h3 className="text-lg font-semibold text-white mb-1">
            Need More Information
          </h3>
          <p className="text-gray-400">{clarification.message}</p>
        </div>
        <button
          onClick={onCancel}
          className="p-1 text-gray-400 hover:text-white transition-colors"
        >
          <X className="w-5 h-5" />
        </button>
      </div>

      {clarification.type === 'multiple_choice' && clarification.options ? (
        <div className="space-y-2">
          {clarification.options.map((option, index) => (
            <button
              key={index}
              onClick={() => handleOptionSelect(option.value)}
              className="w-full p-3 text-left bg-surface-light rounded-lg hover:bg-gray-600 transition-colors group"
            >
              <span className="text-white group-hover:text-primary transition-colors">
                {option.value}
              </span>
              {option.description && (
                <p className="text-sm text-gray-400 mt-1">{option.description}</p>
              )}
            </button>
          ))}
        </div>
      ) : (
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder={`Enter ${clarification.field_name}...`}
            className="w-full px-4 py-3 bg-background border border-surface-light rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
            autoFocus
          />
          <div className="flex gap-3 mt-4">
            <button
              type="button"
              onClick={onCancel}
              className="flex-1 py-2 px-4 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={!value.trim()}
              className="flex-1 py-2 px-4 bg-primary text-white rounded-lg hover:bg-primary-dark disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Continue
            </button>
          </div>
        </form>
      )}
    </div>
  );
}


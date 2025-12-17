import { useState, ReactNode } from 'react';
import { ChevronDown, ChevronRight } from 'lucide-react';

interface CollapsibleSectionProps {
  title: string;
  icon?: ReactNode;
  defaultOpen?: boolean;
  badge?: ReactNode;
  children: ReactNode;
}

export default function CollapsibleSection({
  title,
  icon,
  defaultOpen = false,
  badge,
  children,
}: CollapsibleSectionProps) {
  const [isOpen, setIsOpen] = useState(defaultOpen);

  return (
    <div className="bg-surface rounded-lg border border-surface-light overflow-hidden">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full px-4 py-3 flex items-center gap-3 hover:bg-surface-light transition-colors"
      >
        {isOpen ? (
          <ChevronDown className="w-5 h-5 text-gray-400" />
        ) : (
          <ChevronRight className="w-5 h-5 text-gray-400" />
        )}
        {icon && <span className="text-primary">{icon}</span>}
        <span className="font-medium text-white">{title}</span>
        {badge && <span className="ml-auto">{badge}</span>}
      </button>
      {isOpen && (
        <div className="px-4 pb-4 animate-slideIn">
          {children}
        </div>
      )}
    </div>
  );
}


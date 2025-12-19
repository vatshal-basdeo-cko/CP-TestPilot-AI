interface LogoProps {
  size?: number;
  className?: string;
  variant?: 'default' | 'hero';
}

export default function Logo({ size = 36, className = '', variant = 'default' }: LogoProps) {
  if (variant === 'hero') {
    // Larger, more prominent logo for login/landing pages
    return (
      <div className={`relative ${className}`}>
        {/* Glow effect behind logo */}
        <div 
          className="absolute inset-0 bg-primary/20 rounded-2xl blur-xl scale-150"
          style={{ width: size, height: size }}
        />
        <svg 
          viewBox="0 0 32 32" 
          fill="none" 
          xmlns="http://www.w3.org/2000/svg"
          width={size}
          height={size}
          className="relative z-10"
        >
          {/* Background with gradient */}
          <defs>
            <linearGradient id="bgGradient" x1="0" y1="0" x2="32" y2="32">
              <stop offset="0%" stopColor="#1a1a24"/>
              <stop offset="100%" stopColor="#0d0d0f"/>
            </linearGradient>
            <linearGradient id="ringGradient" x1="6" y1="6" x2="26" y2="26">
              <stop offset="0%" stopColor="#00f5c4"/>
              <stop offset="100%" stopColor="#00d4aa"/>
            </linearGradient>
          </defs>
          <rect width="32" height="32" rx="8" fill="url(#bgGradient)" />
          {/* Outer ring with gradient */}
          <circle cx="16" cy="16" r="10" stroke="url(#ringGradient)" strokeWidth="2" fill="none"/>
          {/* Inner circle */}
          <circle cx="16" cy="16" r="5" stroke="#00d4aa" strokeWidth="1.5" fill="none" opacity="0.5"/>
          {/* Checkmark */}
          <path d="M11 16L14.5 19.5L21 12" stroke="#ffffff" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"/>
          {/* Crosshair lines */}
          <line x1="16" y1="3" x2="16" y2="6" stroke="#00d4aa" strokeWidth="1.5" opacity="0.5"/>
          <line x1="16" y1="26" x2="16" y2="29" stroke="#00d4aa" strokeWidth="1.5" opacity="0.5"/>
          <line x1="3" y1="16" x2="6" y2="16" stroke="#00d4aa" strokeWidth="1.5" opacity="0.5"/>
          <line x1="26" y1="16" x2="29" y2="16" stroke="#00d4aa" strokeWidth="1.5" opacity="0.5"/>
        </svg>
      </div>
    );
  }

  // Default compact logo for navbar
  return (
    <svg 
      viewBox="0 0 32 32" 
      fill="none" 
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      className={className}
    >
      {/* Background */}
      <rect width="32" height="32" rx="8" fill="#0d0d0f" />
      {/* Outer circle */}
      <circle cx="16" cy="16" r="10" stroke="#00d4aa" strokeWidth="2" fill="none"/>
      {/* Inner circle */}
      <circle cx="16" cy="16" r="5" stroke="#00d4aa" strokeWidth="1.5" fill="none" opacity="0.6"/>
      {/* Checkmark */}
      <path d="M11 16L14.5 19.5L21 12" stroke="#ffffff" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"/>
      {/* Crosshair lines */}
      <line x1="16" y1="3" x2="16" y2="6" stroke="#00d4aa" strokeWidth="1.5" opacity="0.4"/>
      <line x1="16" y1="26" x2="16" y2="29" stroke="#00d4aa" strokeWidth="1.5" opacity="0.4"/>
      <line x1="3" y1="16" x2="6" y2="16" stroke="#00d4aa" strokeWidth="1.5" opacity="0.4"/>
      <line x1="26" y1="16" x2="29" y2="16" stroke="#00d4aa" strokeWidth="1.5" opacity="0.4"/>
    </svg>
  );
}

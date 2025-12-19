/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Checkout.com-inspired dark palette
        background: '#0d0d0f',
        surface: '#141417',
        'surface-light': '#1c1c21',
        'surface-hover': '#252529',
        
        // Accent - Checkout.com teal/green
        primary: '#00d4aa',
        'primary-dark': '#00b894',
        'primary-light': '#00f5c4',
        
        // Secondary accent (subtle blue for links/info)
        accent: '#5c6bc0',
        'accent-light': '#7986cb',
        
        // Status colors
        success: '#00d4aa',
        'success-dark': '#00b894',
        error: '#ff5252',
        'error-dark': '#ff1744',
        warning: '#ffab00',
        'warning-dark': '#ff8f00',
        info: '#40c4ff',
        
        // Text colors - high contrast whites and grays
        'text-primary': '#ffffff',
        'text-secondary': '#a1a1aa',
        'text-muted': '#71717a',
        
        // Border colors - subtle
        'border-default': '#27272a',
        'border-light': '#3f3f46',
        'border-subtle': '#1f1f23',
      },
      fontFamily: {
        // Clean geometric sans-serif like Checkout.com
        sans: ['DM Sans', 'Helvetica Neue', 'Arial', 'sans-serif'],
        display: ['DM Sans', 'Helvetica Neue', 'Arial', 'sans-serif'],
        mono: ['Space Mono', 'JetBrains Mono', 'SF Mono', 'monospace'],
      },
      boxShadow: {
        'glow': '0 0 40px rgba(0, 212, 170, 0.12)',
        'glow-lg': '0 0 60px rgba(0, 212, 170, 0.18)',
        'card': '0 4px 24px rgba(0, 0, 0, 0.4)',
        'card-hover': '0 8px 40px rgba(0, 0, 0, 0.5)',
        'inner-glow': 'inset 0 1px 0 rgba(255, 255, 255, 0.03)',
      },
      backgroundImage: {
        // Subtle mesh gradient like Checkout.com
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'mesh-gradient': 'radial-gradient(at 20% 30%, rgba(0, 212, 170, 0.04) 0px, transparent 50%), radial-gradient(at 80% 20%, rgba(92, 107, 192, 0.03) 0px, transparent 50%)',
        'hero-gradient': 'linear-gradient(180deg, rgba(0, 212, 170, 0.03) 0%, transparent 40%)',
      },
      animation: {
        'fade-in': 'fadeIn 0.4s ease-out',
        'slide-up': 'slideUp 0.5s cubic-bezier(0.22, 1, 0.36, 1)',
        'slide-down': 'slideDown 0.4s ease-out',
        'pulse-glow': 'pulseGlow 3s ease-in-out infinite',
        'stagger': 'slideUp 0.5s cubic-bezier(0.22, 1, 0.36, 1) forwards',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        pulseGlow: {
          '0%, 100%': { boxShadow: '0 0 30px rgba(0, 212, 170, 0.1)' },
          '50%': { boxShadow: '0 0 50px rgba(0, 212, 170, 0.2)' },
        },
      },
      borderRadius: {
        'xl': '12px',
        '2xl': '16px',
        '3xl': '24px',
      },
    },
  },
  plugins: [],
}

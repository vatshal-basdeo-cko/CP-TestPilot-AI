import { LogOut, User, Shield, ChevronDown } from 'lucide-react';
import { useAuthStore } from '../store/auth';
import { useNavigate } from 'react-router-dom';
import Logo from './Logo';

export default function Navbar() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="fixed top-0 left-0 right-0 h-16 bg-background/95 backdrop-blur-xl border-b border-border-subtle z-50">
      <div className="h-full px-6 flex items-center justify-between max-w-screen-2xl mx-auto">
        {/* Logo */}
        <div className="flex items-center gap-3">
          <Logo size={36} />
          <div>
            <h1 className="text-lg font-bold text-text-primary tracking-tight">
              TestPilot<span className="text-primary">.AI</span>
            </h1>
          </div>
        </div>

        {/* User menu */}
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-3 px-4 py-2 rounded-xl bg-surface border border-border-subtle hover:border-border-default transition-colors cursor-pointer">
            <div className="w-7 h-7 rounded-full bg-surface-light flex items-center justify-center">
              <User className="w-3.5 h-3.5 text-text-muted" />
            </div>
            <div className="flex flex-col">
              <span className="text-sm text-text-primary font-medium leading-tight">{user?.username}</span>
              <span className="text-[10px] text-text-muted uppercase tracking-wider flex items-center gap-1">
                {user?.role === 'admin' && <Shield className="w-2.5 h-2.5" />}
                {user?.role}
              </span>
            </div>
            <ChevronDown className="w-4 h-4 text-text-muted ml-1" />
          </div>
          
          <button
            onClick={handleLogout}
            className="p-2.5 rounded-xl border border-transparent hover:bg-surface hover:border-border-subtle transition-all text-text-muted hover:text-text-primary"
            title="Sign out"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </div>
    </nav>
  );
}

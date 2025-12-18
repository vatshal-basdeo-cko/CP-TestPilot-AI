import { Zap, LogOut, User, Shield } from 'lucide-react';
import { useAuthStore } from '../store/auth';
import { useNavigate } from 'react-router-dom';

export default function Navbar() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="fixed top-0 left-0 right-0 h-14 bg-surface/95 backdrop-blur-md border-b border-border-default z-50">
      <div className="h-full px-5 flex items-center justify-between">
        {/* Logo */}
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-gradient-to-br from-primary to-primary-light rounded-lg flex items-center justify-center shadow-glow">
            <Zap className="w-4 h-4 text-white" />
          </div>
          <div>
            <h1 className="text-base font-semibold text-text-primary tracking-tight">
              TestPilot<span className="text-primary">AI</span>
            </h1>
          </div>
        </div>

        {/* User menu */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2.5 px-3 py-1.5 rounded-lg bg-surface-light border border-border-default">
            <User className="w-4 h-4 text-text-muted" />
            <span className="text-sm text-text-secondary font-medium">{user?.username}</span>
            <span className={`flex items-center gap-1 text-xs px-2 py-0.5 rounded-full font-medium ${
              user?.role === 'admin' 
                ? 'bg-primary/15 text-primary-light' 
                : 'bg-surface-hover text-text-muted'
            }`}>
              {user?.role === 'admin' && <Shield className="w-3 h-3" />}
              {user?.role}
            </span>
          </div>
          <button
            onClick={handleLogout}
            className="p-2 rounded-lg hover:bg-surface-light border border-transparent hover:border-border-default transition-all text-text-muted hover:text-text-primary"
            title="Sign out"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </div>
    </nav>
  );
}

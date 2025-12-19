import { NavLink } from 'react-router-dom';
import { FlaskConical, Settings, History, BarChart3, Target } from 'lucide-react';
import { useAuthStore } from '../store/auth';

const navItems = [
  { to: '/test', icon: FlaskConical, label: 'Test Execution', category: 'EXECUTE' },
  { to: '/history', icon: History, label: 'History', category: 'MANAGE' },
  { to: '/analytics', icon: BarChart3, label: 'Analytics', category: 'ANALYZE' },
  { to: '/admin', icon: Settings, label: 'Admin', category: 'CONFIGURE', adminOnly: true },
];

export default function Sidebar() {
  const { user } = useAuthStore();
  const isAdmin = user?.role === 'admin';

  return (
    <aside className="fixed left-0 top-16 bottom-0 w-64 bg-surface border-r border-border-subtle">
      <div className="flex flex-col h-full">
        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-1">
          {navItems.map((item) => {
            if (item.adminOnly && !isAdmin) return null;
            
            return (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  `group flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 ${
                    isActive
                      ? 'bg-primary/10 text-primary'
                      : 'text-text-secondary hover:bg-surface-light hover:text-text-primary'
                  }`
                }
              >
                <item.icon className="w-4 h-4 flex-shrink-0" />
                <div className="flex flex-col">
                  <span className="text-sm font-medium">{item.label}</span>
                  <span className="text-[9px] tracking-[0.1em] uppercase text-text-muted group-hover:text-primary/60 transition-colors">
                    {item.category}
                  </span>
                </div>
              </NavLink>
            );
          })}
        </nav>

        {/* Status indicator */}
        <div className="p-4 border-t border-border-subtle">
          <div className="p-4 rounded-xl bg-surface-light border border-border-subtle">
            <div className="flex items-center gap-2 mb-3">
              <Target className="w-3.5 h-3.5 text-primary" />
              <span className="text-[10px] font-semibold tracking-[0.1em] uppercase text-text-muted">
                System Status
              </span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-2 h-2 rounded-full bg-success animate-pulse"></div>
              <span className="text-xs text-text-secondary">All services operational</span>
            </div>
          </div>
        </div>
      </div>
    </aside>
  );
}

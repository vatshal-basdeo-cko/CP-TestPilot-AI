import { NavLink } from 'react-router-dom';
import { FlaskConical, Settings, History, Activity, BarChart3 } from 'lucide-react';
import { useAuthStore } from '../store/auth';

const navItems = [
  { to: '/test', icon: FlaskConical, label: 'Test Execution', description: 'Run API tests' },
  { to: '/history', icon: History, label: 'History', description: 'View past results' },
  { to: '/analytics', icon: BarChart3, label: 'Analytics', description: 'View statistics' },
  { to: '/admin', icon: Settings, label: 'Admin', description: 'System settings', adminOnly: true },
];

export default function Sidebar() {
  const { user } = useAuthStore();
  const isAdmin = user?.role === 'admin';

  return (
    <aside className="fixed left-0 top-14 bottom-0 w-60 bg-surface border-r border-border-default">
      <div className="flex flex-col h-full">
        {/* Navigation */}
        <nav className="flex-1 p-3 space-y-1">
          {navItems.map((item) => {
            if (item.adminOnly && !isAdmin) return null;
            
            return (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all ${
                    isActive
                      ? 'bg-primary/10 text-primary border border-primary/20'
                      : 'text-text-secondary hover:bg-surface-light hover:text-text-primary border border-transparent'
                  }`
                }
              >
                <item.icon className="w-4 h-4 flex-shrink-0" />
                <span className="text-sm font-medium">{item.label}</span>
              </NavLink>
            );
          })}
        </nav>

        {/* Status indicator */}
        <div className="p-3 border-t border-border-default">
          <div className="p-3 rounded-lg bg-surface-light">
            <div className="flex items-center gap-2 text-xs text-text-muted mb-2">
              <Activity className="w-3.5 h-3.5" />
              <span>System Status</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-1.5 h-1.5 rounded-full bg-success"></div>
              <span className="text-xs text-text-secondary">All services operational</span>
            </div>
          </div>
        </div>
      </div>
    </aside>
  );
}

import { NavLink } from 'react-router-dom';
import { FlaskConical, Settings, History, Database } from 'lucide-react';
import { useAuthStore } from '../store/auth';

const navItems = [
  { to: '/test', icon: FlaskConical, label: 'Test Execution' },
  { to: '/history', icon: History, label: 'History' },
  { to: '/admin', icon: Settings, label: 'Admin Panel', adminOnly: true },
];

export default function Sidebar() {
  const { user } = useAuthStore();
  const isAdmin = user?.role === 'admin';

  return (
    <aside className="fixed left-0 top-16 bottom-0 w-64 bg-surface border-r border-surface-light p-4">
      <nav className="space-y-2">
        {navItems.map((item) => {
          if (item.adminOnly && !isAdmin) return null;
          
          return (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                  isActive
                    ? 'bg-primary/20 text-primary'
                    : 'text-gray-400 hover:bg-surface-light hover:text-white'
                }`
              }
            >
              <item.icon className="w-5 h-5" />
              <span className="font-medium">{item.label}</span>
            </NavLink>
          );
        })}
      </nav>

      <div className="absolute bottom-4 left-4 right-4">
        <div className="p-4 rounded-lg bg-surface-light">
          <div className="flex items-center gap-2 text-sm text-gray-400 mb-2">
            <Database className="w-4 h-4" />
            <span>System Status</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-success animate-pulse"></div>
            <span className="text-xs text-gray-300">All services healthy</span>
          </div>
        </div>
      </div>
    </aside>
  );
}


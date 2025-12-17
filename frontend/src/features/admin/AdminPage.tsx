import { useState } from 'react';
import { Database, Users, Settings } from 'lucide-react';
import ApiConfigList from './ApiConfigList';
import UserManagement from './UserManagement';
import SystemSettings from './SystemSettings';
import { useAuthStore } from '../../store/auth';

type TabId = 'apis' | 'users' | 'settings';

const tabs: { id: TabId; label: string; icon: React.ElementType }[] = [
  { id: 'apis', label: 'API Configurations', icon: Database },
  { id: 'users', label: 'User Management', icon: Users },
  { id: 'settings', label: 'System Settings', icon: Settings },
];

export default function AdminPage() {
  const [activeTab, setActiveTab] = useState<TabId>('apis');
  const { user } = useAuthStore();

  if (user?.role !== 'admin') {
    return (
      <div className="text-center py-12">
        <h2 className="text-xl font-semibold text-white mb-2">Access Denied</h2>
        <p className="text-gray-400">You need admin privileges to access this page.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white mb-2">Admin Panel</h1>
        <p className="text-gray-400">Manage API configurations, users, and system settings.</p>
      </div>

      {/* Tabs */}
      <div className="border-b border-surface-light">
        <div className="flex gap-1">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex items-center gap-2 px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-primary text-primary'
                  : 'border-transparent text-gray-400 hover:text-white hover:border-gray-600'
              }`}
            >
              <tab.icon className="w-4 h-4" />
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Tab Content */}
      <div className="animate-fadeIn">
        {activeTab === 'apis' && <ApiConfigList />}
        {activeTab === 'users' && <UserManagement />}
        {activeTab === 'settings' && <SystemSettings />}
      </div>
    </div>
  );
}


import { useState } from 'react';
import { UserPlus, User, Shield } from 'lucide-react';

interface UserData {
  id: string;
  username: string;
  role: 'admin' | 'user';
  created_at: string;
}

// Mock data - replace with actual API call when backend supports it
const mockUsers: UserData[] = [
  { id: '1', username: 'admin', role: 'admin', created_at: '2024-01-01T00:00:00Z' },
];

export default function UserManagement() {
  const [users] = useState<UserData[]>(mockUsers);
  const [isAddingUser, setIsAddingUser] = useState(false);
  const [newUser, setNewUser] = useState({ username: '', password: '', role: 'user' as 'admin' | 'user' });

  const handleAddUser = async (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement user creation API
    console.log('Create user:', newUser);
    setIsAddingUser(false);
    setNewUser({ username: '', password: '', role: 'user' });
  };

  return (
    <div className="space-y-4">
      {/* Actions */}
      <div className="flex items-center justify-between">
        <p className="text-gray-400">{users.length} user{users.length !== 1 ? 's' : ''}</p>
        <button
          onClick={() => setIsAddingUser(true)}
          className="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors"
        >
          <UserPlus className="w-4 h-4" />
          Add User
        </button>
      </div>

      {/* Add User Form */}
      {isAddingUser && (
        <div className="bg-surface rounded-lg border border-surface-light p-4 animate-slideIn">
          <h3 className="text-lg font-medium text-white mb-4">Add New User</h3>
          <form onSubmit={handleAddUser} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Username</label>
              <input
                type="text"
                value={newUser.username}
                onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
                className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Password</label>
              <input
                type="password"
                value={newUser.password}
                onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
                className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Role</label>
              <select
                value={newUser.role}
                onChange={(e) => setNewUser({ ...newUser, role: e.target.value as 'admin' | 'user' })}
                className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white"
              >
                <option value="user">User</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => setIsAddingUser(false)}
                className="flex-1 py-2 px-4 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="flex-1 py-2 px-4 bg-primary text-white rounded-lg hover:bg-primary-dark"
              >
                Create User
              </button>
            </div>
          </form>
        </div>
      )}

      {/* User List */}
      <div className="bg-surface rounded-lg border border-surface-light overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-surface-light">
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">User</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Role</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Created</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id} className="border-b border-surface-light last:border-0 hover:bg-surface-light/50">
                <td className="px-4 py-3">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-surface-light rounded-full flex items-center justify-center">
                      <User className="w-4 h-4 text-gray-400" />
                    </div>
                    <span className="text-white font-medium">{user.username}</span>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <span className={`flex items-center gap-1 px-2 py-1 text-xs rounded ${
                    user.role === 'admin' ? 'bg-primary/20 text-primary' : 'bg-gray-600 text-gray-300'
                  }`}>
                    {user.role === 'admin' && <Shield className="w-3 h-3" />}
                    {user.role}
                  </span>
                </td>
                <td className="px-4 py-3 text-sm text-gray-400">
                  {new Date(user.created_at).toLocaleDateString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}


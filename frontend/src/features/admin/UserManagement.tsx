import { useState, useEffect } from 'react';
import { UserPlus, User, Shield, Trash2, Loader2, AlertCircle } from 'lucide-react';
import { usersApi, User as UserData } from '../../api/users';

export default function UserManagement() {
  const [users, setUsers] = useState<UserData[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isAddingUser, setIsAddingUser] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [newUser, setNewUser] = useState({ username: '', password: '', role: 'user' as 'admin' | 'user' });

  // Fetch users on mount
  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const usersList = await usersApi.list();
      setUsers(usersList);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);
    
    try {
      await usersApi.create({
        username: newUser.username,
        password: newUser.password,
        role: newUser.role,
      });
      setIsAddingUser(false);
      setNewUser({ username: '', password: '', role: 'user' });
      await fetchUsers(); // Refresh the list
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create user';
      // Try to get error from axios response
      const axiosError = err as { response?: { data?: { error?: string } } };
      setError(axiosError.response?.data?.error || errorMessage);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteUser = async (id: string, username: string) => {
    if (!confirm(`Are you sure you want to delete user "${username}"?`)) {
      return;
    }
    
    setDeletingId(id);
    setError(null);
    
    try {
      await usersApi.delete(id);
      await fetchUsers(); // Refresh the list
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete user';
      const axiosError = err as { response?: { data?: { error?: string } } };
      setError(axiosError.response?.data?.error || errorMessage);
    } finally {
      setDeletingId(null);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Error Message */}
      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400">
          <AlertCircle className="w-5 h-5 flex-shrink-0" />
          <span>{error}</span>
        </div>
      )}

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
                className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white focus:border-primary focus:outline-none"
                placeholder="Enter username"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Password</label>
              <input
                type="password"
                value={newUser.password}
                onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
                className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white focus:border-primary focus:outline-none"
                placeholder="Minimum 8 characters"
                minLength={8}
                required
                disabled={isSubmitting}
              />
              <p className="mt-1 text-xs text-gray-500">Password must be at least 8 characters</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Role</label>
              <select
                value={newUser.role}
                onChange={(e) => setNewUser({ ...newUser, role: e.target.value as 'admin' | 'user' })}
                className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white focus:border-primary focus:outline-none"
                disabled={isSubmitting}
              >
                <option value="user">User</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => {
                  setIsAddingUser(false);
                  setNewUser({ username: '', password: '', role: 'user' });
                  setError(null);
                }}
                className="flex-1 py-2 px-4 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
                disabled={isSubmitting}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="flex-1 py-2 px-4 bg-primary text-white rounded-lg hover:bg-primary-dark transition-colors flex items-center justify-center gap-2"
                disabled={isSubmitting}
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    Creating...
                  </>
                ) : (
                  'Create User'
                )}
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
              <th className="px-4 py-3 text-right text-sm font-medium text-gray-400">Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-4 py-8 text-center text-gray-500">
                  No users found
                </td>
              </tr>
            ) : (
              users.map((user) => (
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
                    <span className={`flex items-center gap-1 px-2 py-1 text-xs rounded w-fit ${
                      user.role === 'admin' ? 'bg-primary/20 text-primary' : 'bg-gray-600 text-gray-300'
                    }`}>
                      {user.role === 'admin' && <Shield className="w-3 h-3" />}
                      {user.role}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-400">
                    {new Date(user.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      onClick={() => handleDeleteUser(user.id, user.username)}
                      disabled={deletingId === user.id}
                      className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-colors disabled:opacity-50"
                      title="Delete user"
                    >
                      {deletingId === user.id ? (
                        <Loader2 className="w-4 h-4 animate-spin" />
                      ) : (
                        <Trash2 className="w-4 h-4" />
                      )}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

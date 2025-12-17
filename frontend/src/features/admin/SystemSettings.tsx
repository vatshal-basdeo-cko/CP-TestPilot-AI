import { useState } from 'react';
import { Save, RotateCcw } from 'lucide-react';

interface Settings {
  learningThreshold: number;
  historyRetentionDays: number;
  defaultEnvironment: string;
  rateLimitPerMinute: number;
}

const defaultSettings: Settings = {
  learningThreshold: 5,
  historyRetentionDays: 90,
  defaultEnvironment: 'QA1',
  rateLimitPerMinute: 60,
};

export default function SystemSettings() {
  const [settings, setSettings] = useState<Settings>(defaultSettings);
  const [isSaving, setIsSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  const handleSave = async () => {
    setIsSaving(true);
    // TODO: Implement settings save API
    await new Promise((resolve) => setTimeout(resolve, 1000));
    setIsSaving(false);
    setSaved(true);
    setTimeout(() => setSaved(false), 3000);
  };

  const handleReset = () => {
    setSettings(defaultSettings);
  };

  return (
    <div className="max-w-2xl space-y-6">
      <div className="bg-surface rounded-lg border border-surface-light p-6">
        <h3 className="text-lg font-semibold text-white mb-6">Learning Settings</h3>
        
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Learning Threshold
            </label>
            <p className="text-xs text-gray-500 mb-2">
              Number of successful tests before the system learns from patterns
            </p>
            <input
              type="number"
              min="1"
              max="100"
              value={settings.learningThreshold}
              onChange={(e) => setSettings({ ...settings, learningThreshold: parseInt(e.target.value) || 5 })}
              className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              History Retention (days)
            </label>
            <p className="text-xs text-gray-500 mb-2">
              How long to keep test execution history
            </p>
            <input
              type="number"
              min="7"
              max="365"
              value={settings.historyRetentionDays}
              onChange={(e) => setSettings({ ...settings, historyRetentionDays: parseInt(e.target.value) || 90 })}
              className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white"
            />
          </div>
        </div>
      </div>

      <div className="bg-surface rounded-lg border border-surface-light p-6">
        <h3 className="text-lg font-semibold text-white mb-6">Environment Settings</h3>
        
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Default Environment
            </label>
            <select
              value={settings.defaultEnvironment}
              onChange={(e) => setSettings({ ...settings, defaultEnvironment: e.target.value })}
              className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white"
            >
              <option value="QA1">QA1</option>
              <option value="QA2">QA2</option>
              <option value="Staging">Staging</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Rate Limit (requests/minute)
            </label>
            <input
              type="number"
              min="1"
              max="1000"
              value={settings.rateLimitPerMinute}
              onChange={(e) => setSettings({ ...settings, rateLimitPerMinute: parseInt(e.target.value) || 60 })}
              className="w-full px-4 py-2 bg-background border border-surface-light rounded-lg text-white"
            />
          </div>
        </div>
      </div>

      {/* Actions */}
      <div className="flex items-center gap-3">
        <button
          onClick={handleSave}
          disabled={isSaving}
          className="flex items-center gap-2 px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark disabled:opacity-50 transition-colors"
        >
          <Save className="w-4 h-4" />
          {isSaving ? 'Saving...' : 'Save Changes'}
        </button>
        <button
          onClick={handleReset}
          className="flex items-center gap-2 px-6 py-2 bg-surface-light text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
        >
          <RotateCcw className="w-4 h-4" />
          Reset to Defaults
        </button>
        {saved && (
          <span className="text-success text-sm animate-fadeIn">Settings saved successfully!</span>
        )}
      </div>
    </div>
  );
}


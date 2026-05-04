import { useState, useEffect } from 'react';
import { useScanner } from '../context/ScannerContext';
import { Save, SlidersHorizontal, HardDrive, Check, Globe } from 'lucide-react';
import { toast } from 'sonner';
import { motion } from 'motion/react';
import { translations } from '../i18n';

export function Settings() {
  const { ember2018, ember2024, updateConfig, language, setLanguage } = useScanner();
  const t = translations[language];

  // Local state for forms
  const [e2018Path, setE2018Path] = useState(ember2018.path);
  const [e2018Threshold, setE2018Threshold] = useState(ember2018.threshold);
  
  const [e2024Path, setE2024Path] = useState(ember2024.path);
  const [e2024Threshold, setE2024Threshold] = useState(ember2024.threshold);

  useEffect(() => {
    setE2018Path(ember2018.path);
    setE2018Threshold(ember2018.threshold);
    setE2024Path(ember2024.path);
    setE2024Threshold(ember2024.threshold);
  }, [ember2018, ember2024]);

  const handleSave = () => {
    updateConfig('ember2018', { path: e2018Path, threshold: e2018Threshold });
    updateConfig('ember2024', { path: e2024Path, threshold: e2024Threshold });
    toast.success(t.settings_saved);
  };

  return (
    <div className="p-8 max-w-4xl mx-auto space-y-8 pb-20">
      <div>
        <h2 className="text-3xl font-bold text-slate-100 mb-2">{t.setting_title}</h2>
        <p className="text-slate-400">{t.setting_desc}</p>
      </div>

      <div className="grid gap-6">
        {/* Language Selection */}
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
          <div className="flex items-center gap-3 mb-4">
            <Globe className="w-5 h-5 text-slate-400" />
            <h3 className="text-lg font-semibold text-slate-200">{t.language}</h3>
          </div>
          <div className="flex gap-4">
            <button
              onClick={() => setLanguage('zh')}
              className={`flex-1 flex items-center justify-between p-4 rounded-lg border-2 transition-all ${
                language === 'zh'
                  ? 'border-cyan-500 bg-cyan-500/10'
                  : 'border-slate-800 bg-slate-950 hover:border-slate-700'
              }`}
            >
              <span className={`font-semibold ${language === 'zh' ? 'text-cyan-400' : 'text-slate-300'}`}>
                {t.lang_zh}
              </span>
              {language === 'zh' && <Check className="text-cyan-500 w-5 h-5" />}
            </button>
            <button
              onClick={() => setLanguage('en')}
              className={`flex-1 flex items-center justify-between p-4 rounded-lg border-2 transition-all ${
                language === 'en'
                  ? 'border-cyan-500 bg-cyan-500/10'
                  : 'border-slate-800 bg-slate-950 hover:border-slate-700'
              }`}
            >
              <span className={`font-semibold ${language === 'en' ? 'text-cyan-400' : 'text-slate-300'}`}>
                {t.lang_en}
              </span>
              {language === 'en' && <Check className="text-cyan-500 w-5 h-5" />}
            </button>
          </div>
        </div>

        {/* Ember 2018 Config */}
        <div className={`bg-slate-900 border rounded-xl p-6 transition-colors duration-300 border-slate-700`}>
          <div className="flex items-center gap-3 mb-6">
            <div className={`p-2 rounded-lg bg-cyan-500/20 text-cyan-400`}>
              <HardDrive className="w-5 h-5" />
            </div>
            <h3 className="text-lg font-semibold text-slate-200">{t.ember2018}</h3>
          </div>

          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">{t.model_path}</label>
              <input
                type="text"
                value={e2018Path}
                onChange={(e) => setE2018Path(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2.5 text-slate-300 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                placeholder="C:\path\to\ember2018.bin"
              />
            </div>

            <div>
              <div className="flex justify-between items-center mb-2">
                <label className="flex items-center gap-2 text-sm font-medium text-slate-400">
                  <SlidersHorizontal className="w-4 h-4" />
                  {t.malware_threshold}
                </label>
                <span className="text-sm font-mono text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded">
                  {e2018Threshold.toFixed(2)}
                </span>
              </div>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                value={e2018Threshold}
                onChange={(e) => setE2018Threshold(parseFloat(e.target.value))}
                className="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
              />
              <p className="text-xs text-slate-500 mt-2">
                {t.threshold_help_2018}
              </p>
            </div>
          </div>
        </div>

        {/* Ember 2024 Config */}
        <div className={`bg-slate-900 border rounded-xl p-6 transition-colors duration-300 border-slate-700`}>
          <div className="flex items-center gap-3 mb-6">
            <div className={`p-2 rounded-lg bg-cyan-500/20 text-cyan-400`}>
              <HardDrive className="w-5 h-5" />
            </div>
            <h3 className="text-lg font-semibold text-slate-200">{t.ember2024}</h3>
          </div>

          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">{t.model_path}</label>
              <input
                type="text"
                value={e2024Path}
                onChange={(e) => setE2024Path(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2.5 text-slate-300 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                placeholder="C:\path\to\ember2024.bin"
              />
            </div>

            <div>
              <div className="flex justify-between items-center mb-2">
                <label className="flex items-center gap-2 text-sm font-medium text-slate-400">
                  <SlidersHorizontal className="w-4 h-4" />
                  {t.malware_threshold}
                </label>
                <span className="text-sm font-mono text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded">
                  {e2024Threshold.toFixed(2)}
                </span>
              </div>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                value={e2024Threshold}
                onChange={(e) => setE2024Threshold(parseFloat(e.target.value))}
                className="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
              />
              <p className="text-xs text-slate-500 mt-2">
                {t.threshold_help_2024}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div className="fixed bottom-0 left-64 right-0 p-4 bg-slate-900/80 backdrop-blur-md border-t border-slate-800 flex justify-end gap-4">
        <motion.button
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={handleSave}
          className="flex items-center gap-2 bg-cyan-500 hover:bg-cyan-400 text-slate-950 px-6 py-2.5 rounded-lg font-medium transition-colors shadow-lg shadow-cyan-500/20"
        >
          <Save className="w-4 h-4" />
          {t.apply_settings}
        </motion.button>
      </div>
    </div>
  );
}

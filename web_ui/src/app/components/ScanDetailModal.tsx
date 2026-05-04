import { motion, AnimatePresence } from 'motion/react';
import { X, ShieldAlert, ShieldCheck, Database, HardDrive, FileText } from 'lucide-react';
import { ScanResult } from '../context/ScannerContext';
import { translations } from '../i18n';

interface ScanDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  result: ScanResult | null;
  language: 'zh' | 'en';
}

export function ScanDetailModal({ isOpen, onClose, result, language }: ScanDetailModalProps) {
  const t = translations[language];

  if (!result) return null;

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm"
          />
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-full max-w-md bg-slate-900 border border-slate-700 rounded-2xl shadow-2xl shadow-cyan-900/20 overflow-hidden"
          >
            <div className="flex items-center justify-between p-4 border-b border-slate-800">
              <h3 className="text-lg font-semibold text-slate-200">{t.detail_title}</h3>
              <button
                onClick={onClose}
                className="text-slate-400 hover:text-slate-200 transition-colors p-1"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="p-6 space-y-6">
              <div className="flex flex-col items-center justify-center p-6 bg-slate-950/50 rounded-xl border border-slate-800">
                {result.isMalicious ? (
                  <div className="w-16 h-16 bg-red-500/10 text-red-400 rounded-full flex items-center justify-center mb-4 border border-red-500/20">
                    <ShieldAlert className="w-8 h-8" />
                  </div>
                ) : (
                  <div className="w-16 h-16 bg-emerald-500/10 text-emerald-400 rounded-full flex items-center justify-center mb-4 border border-emerald-500/20">
                    <ShieldCheck className="w-8 h-8" />
                  </div>
                )}
                
                <h4 className={`text-xl font-bold mb-1 ${result.isMalicious ? 'text-red-400' : 'text-emerald-400'}`}>
                  {result.isMalicious ? t.malicious : t.safe}
                </h4>
                <p className="text-sm text-slate-400 font-mono text-center break-all px-4">
                  {result.filename}
                </p>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-between py-2 border-b border-slate-800/50">
                  <div className="flex items-center gap-2 text-slate-400">
                    <Database className="w-4 h-4" />
                    <span className="text-sm">{t.lightgbm_score}</span>
                  </div>
                  <span className={`text-lg font-mono font-bold ${result.isMalicious ? 'text-red-400' : 'text-cyan-400'}`}>
                    {result.score.toFixed(2)}
                  </span>
                </div>

                <div className="flex items-center justify-between py-2 border-b border-slate-800/50">
                  <div className="flex items-center gap-2 text-slate-400">
                    <HardDrive className="w-4 h-4" />
                    <span className="text-sm">{t.engine}</span>
                  </div>
                  <span className="text-sm font-medium text-slate-200">{result.modelUsed}</span>
                </div>

                <div className="flex items-center justify-between py-2 border-b border-slate-800/50">
                  <div className="flex items-center gap-2 text-slate-400">
                    <FileText className="w-4 h-4" />
                    <span className="text-sm">{t.file_size}</span>
                  </div>
                  <span className="text-sm font-medium text-slate-200">
                    {(result.size / 1024).toFixed(1)} KB
                  </span>
                </div>
              </div>
            </div>
            
            <div className="p-4 border-t border-slate-800 bg-slate-900/50 flex justify-end">
              <button
                onClick={onClose}
                className="px-6 py-2 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-lg transition-colors font-medium text-sm"
              >
                Close
              </button>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}

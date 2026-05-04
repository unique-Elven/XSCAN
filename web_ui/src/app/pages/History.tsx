import { useState } from 'react';
import { useScanner, ScanResult } from '../context/ScannerContext';
import { AlertTriangle, CheckCircle, Trash2, Calendar, HardDrive } from 'lucide-react';
import { translations } from '../i18n';
import { ScanDetailModal } from '../components/ScanDetailModal';

export function History() {
  const { scanHistory, clearHistory, language } = useScanner();
  const t = translations[language];

  const [selectedResult, setSelectedResult] = useState<ScanResult | null>(null);

  return (
    <div className="p-8 max-w-5xl mx-auto space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold text-slate-100 mb-2">{t.history_title}</h2>
          <p className="text-slate-400">{t.history_desc}</p>
        </div>
        {scanHistory.length > 0 && (
          <button
            onClick={clearHistory}
            className="flex items-center gap-2 bg-slate-900 border border-slate-800 hover:border-red-500/50 hover:text-red-400 text-slate-400 px-4 py-2 rounded-lg transition-colors"
          >
            <Trash2 className="w-4 h-4" />
            {t.clear_log}
          </button>
        )}
      </div>

      {scanHistory.length === 0 ? (
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-16 flex flex-col items-center justify-center text-center">
          <div className="w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center mb-4 text-slate-500">
            <Calendar className="w-8 h-8" />
          </div>
          <h3 className="text-xl font-medium text-slate-300">{t.no_history}</h3>
          <p className="text-slate-500 mt-2">{t.no_history_desc}</p>
        </div>
      ) : (
        <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-slate-800 bg-slate-900/50 text-slate-400 text-xs uppercase tracking-wider">
                  <th className="px-6 py-4 font-medium">{t.status}</th>
                  <th className="px-6 py-4 font-medium">{t.filename}</th>
                  <th className="px-6 py-4 font-medium">{t.engine}</th>
                  <th className="px-6 py-4 font-medium">{t.score}</th>
                  <th className="px-6 py-4 font-medium">{t.date}</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800">
                {scanHistory.map((result) => (
                  <tr 
                    key={result.id} 
                    onClick={() => setSelectedResult(result)}
                    className="hover:bg-slate-800/50 transition-colors cursor-pointer"
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      {result.isMalicious ? (
                        <div className="flex items-center gap-2 text-red-400 bg-red-400/10 w-fit px-2.5 py-1 rounded-full text-xs font-medium border border-red-400/20">
                          <AlertTriangle className="w-3.5 h-3.5" />
                          {t.malicious}
                        </div>
                      ) : (
                        <div className="flex items-center gap-2 text-emerald-400 bg-emerald-400/10 w-fit px-2.5 py-1 rounded-full text-xs font-medium border border-emerald-400/20">
                          <CheckCircle className="w-3.5 h-3.5" />
                          {t.safe}
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col">
                        <span className="text-slate-200 text-sm font-medium truncate max-w-[250px]" title={result.filename}>
                          {result.filename}
                        </span>
                        <span className="text-slate-500 text-xs">
                          {(result.size / 1024).toFixed(1)} KB
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center gap-1.5 text-xs text-slate-300">
                        <HardDrive className="w-3.5 h-3.5 text-slate-500" />
                        {result.modelUsed}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`text-sm font-mono ${result.isMalicious ? 'text-red-400' : 'text-slate-300'}`}>
                        {result.score.toFixed(2)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-500">
                      {new Date(result.timestamp).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
      
      <ScanDetailModal 
        isOpen={selectedResult !== null} 
        onClose={() => setSelectedResult(null)} 
        result={selectedResult}
        language={language}
      />
    </div>
  );
}

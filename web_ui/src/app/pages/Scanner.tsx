import { useState, useRef, ChangeEvent } from 'react';
import { useScanner, ScanResult } from '../context/ScannerContext';
import { motion, AnimatePresence } from 'motion/react';
import { Upload, FolderOpen, FileText, AlertTriangle, CheckCircle, ShieldCheck, Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { translations } from '../i18n';
import { ScanDetailModal } from '../components/ScanDetailModal';

export function Scanner() {
  const { ember2018, ember2024, addScanResults, language } = useScanner();
  const t = translations[language];

  const [isScanning, setIsScanning] = useState(false);
  const [progress, setProgress] = useState(0);
  const [currentFile, setCurrentFile] = useState('');
  const [recentScans, setRecentScans] = useState<ScanResult[]>([]);
  
  const [selectedResult, setSelectedResult] = useState<ScanResult | null>(null);

  const fileInputRef = useRef<HTMLInputElement>(null);
  const dirInputRef = useRef<HTMLInputElement>(null);

  const triggerFileSelect = () => {
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
      fileInputRef.current.click();
    }
  };

  const triggerDirSelect = () => {
    if (dirInputRef.current) {
      dirInputRef.current.value = '';
      dirInputRef.current.click();
    }
  };

  const handleFiles = async (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    setIsScanning(true);
    setProgress(0);
    setRecentScans([]);
    
    const fileList = Array.from(files);
    const totalFiles = fileList.length;
    const newResults: ScanResult[] = [];

    // Mock scanning process
    for (let i = 0; i < totalFiles; i++) {
      const file = fileList[i];
      setCurrentFile(file.name);
      
      // Artificial delay to simulate ML processing
      const delay = Math.max(100, Math.min(800, file.size / 1024)); 
      await new Promise(resolve => setTimeout(resolve, delay));

      // Generate mock score based on some heuristics (random here for demo)
      // Make smaller files slightly more likely to be benign for demo realism
      const randomFactor = Math.random();
      const baseScore = randomFactor * 1.2; // Can go up to 1.2, clamped later
      const score = Math.min(1.0, Math.max(0.0, baseScore));
      
      // Simulate backend auto-optimization: Randomly select which model 'caught' or 'processed' it best
      const used2018 = Math.random() > 0.5;
      const currentModelName = used2018 ? 'Ember2018' : 'Ember2024';
      const currentThreshold = used2018 ? ember2018.threshold : ember2024.threshold;
      
      const isMalicious = score >= currentThreshold;

      const result: ScanResult = {
        id: crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).substring(7),
        filename: file.webkitRelativePath || file.name,
        size: file.size,
        modelUsed: currentModelName,
        score: Number(score.toFixed(4)),
        isMalicious,
        timestamp: Date.now()
      };

      newResults.push(result);
      setRecentScans(prev => [result, ...prev]);
      setProgress(Math.round(((i + 1) / totalFiles) * 100));
    }

    addScanResults(newResults);
    
    setIsScanning(false);
    setCurrentFile('');
    
    const maliciousCount = newResults.filter(r => r.isMalicious).length;
    if (maliciousCount > 0) {
      toast.error(t.scan_complete_malicious);
    } else {
      toast.success(t.scan_complete_safe);
    }
  };

  return (
    <div className="relative min-h-full">
      {/* Cool CPU Background Effect in the bottom area */}
      <div className="absolute inset-x-0 bottom-0 h-96 pointer-events-none overflow-hidden opacity-30">
        <div className="absolute inset-0 flex items-end justify-center">
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 2 }}
            className="w-[800px] h-[400px] relative"
          >
            {/* Grid Lines */}
            <div className="absolute inset-0 bg-[linear-gradient(to_right,#0891b222_1px,transparent_1px),linear-gradient(to_bottom,#0891b222_1px,transparent_1px)] bg-[size:40px_40px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_100%,#000_70%,transparent_100%)]" />
            
            {/* CPU Node */}
            <motion.div 
              animate={{ 
                boxShadow: ['0 0 20px 0px rgba(6,182,212,0.3)', '0 0 60px 10px rgba(6,182,212,0.6)', '0 0 20px 0px rgba(6,182,212,0.3)']
              }}
              transition={{ repeat: Infinity, duration: 3, ease: "easeInOut" }}
              className="absolute bottom-10 left-1/2 -translate-x-1/2 w-32 h-32 bg-slate-900 border border-cyan-500/50 rounded-xl flex items-center justify-center z-10"
            >
              <div className="w-24 h-24 border border-cyan-400/30 rounded-lg flex items-center justify-center">
                <div className="w-16 h-16 bg-cyan-950 rounded flex items-center justify-center relative overflow-hidden">
                  <div className="absolute inset-0 bg-cyan-500/20" />
                  <motion.div
                    animate={{ y: ['-100%', '100%'] }}
                    transition={{ repeat: Infinity, duration: 2, ease: "linear" }}
                    className="absolute inset-0 h-1/2 bg-gradient-to-b from-transparent via-cyan-400/50 to-transparent"
                  />
                  <ShieldCheck className="w-8 h-8 text-cyan-400 z-10" />
                </div>
              </div>
            </motion.div>

            {/* Circuit Traces */}
            <svg className="absolute inset-0 w-full h-full" style={{ filter: 'drop-shadow(0 0 8px rgba(6,182,212,0.5))' }}>
              <motion.path
                d="M400,320 L250,320 L250,200 L150,200"
                fill="none"
                stroke="rgba(6,182,212,0.5)"
                strokeWidth="2"
                initial={{ pathLength: 0 }}
                animate={{ pathLength: 1 }}
                transition={{ duration: 2, repeat: Infinity, repeatType: "reverse", ease: "easeInOut" }}
              />
              <motion.path
                d="M400,320 L550,320 L550,150 L650,150"
                fill="none"
                stroke="rgba(6,182,212,0.5)"
                strokeWidth="2"
                initial={{ pathLength: 0 }}
                animate={{ pathLength: 1 }}
                transition={{ duration: 2.5, repeat: Infinity, repeatType: "reverse", ease: "easeInOut" }}
              />
              <motion.path
                d="M400,280 L400,100 L300,50"
                fill="none"
                stroke="rgba(6,182,212,0.5)"
                strokeWidth="2"
                initial={{ pathLength: 0 }}
                animate={{ pathLength: 1 }}
                transition={{ duration: 1.8, repeat: Infinity, repeatType: "reverse", ease: "easeInOut" }}
              />
            </svg>
          </motion.div>
        </div>
      </div>

      <div className="p-8 max-w-5xl mx-auto space-y-8 relative z-10">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold text-slate-100 mb-2">{t.scan_title}</h2>
            <p className="text-slate-400">{t.scan_desc}</p>
          </div>
          <div className="bg-slate-900 px-4 py-2 rounded-lg border border-slate-800 flex items-center gap-3 shadow-lg shadow-black/50">
            <ShieldCheck className="w-5 h-5 text-cyan-500" />
            <div className="flex flex-col">
              <span className="text-xs text-slate-500 font-medium uppercase">{t.dual_engine}</span>
              <span className="text-sm text-slate-200 font-medium">{t.backend_auto}</span>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <motion.button
            whileHover={{ scale: 1.02, y: -2 }}
            whileTap={{ scale: 0.98 }}
            onClick={triggerFileSelect}
            disabled={isScanning}
            className="relative group overflow-hidden bg-slate-900/80 backdrop-blur-sm border-2 border-dashed border-slate-700 hover:border-cyan-500/50 rounded-2xl p-10 flex flex-col items-center justify-center gap-4 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-xl"
          >
            <div className="p-4 bg-slate-800 rounded-full group-hover:bg-cyan-500/10 transition-colors text-slate-400 group-hover:text-cyan-400">
              <FileText className="w-10 h-10" />
            </div>
            <div className="text-center">
              <h3 className="text-lg font-semibold text-slate-200">{t.scan_files}</h3>
              <p className="text-sm text-slate-500 mt-1">{t.scan_files_desc}</p>
            </div>
          </motion.button>

          <motion.button
            whileHover={{ scale: 1.02, y: -2 }}
            whileTap={{ scale: 0.98 }}
            onClick={triggerDirSelect}
            disabled={isScanning}
            className="relative group overflow-hidden bg-slate-900/80 backdrop-blur-sm border-2 border-dashed border-slate-700 hover:border-cyan-500/50 rounded-2xl p-10 flex flex-col items-center justify-center gap-4 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-xl"
          >
            <div className="p-4 bg-slate-800 rounded-full group-hover:bg-cyan-500/10 transition-colors text-slate-400 group-hover:text-cyan-400">
              <FolderOpen className="w-10 h-10" />
            </div>
            <div className="text-center">
              <h3 className="text-lg font-semibold text-slate-200">{t.scan_dir}</h3>
              <p className="text-sm text-slate-500 mt-1">{t.scan_dir_desc}</p>
            </div>
          </motion.button>
        </div>

        {/* Hidden inputs */}
        <input
          type="file"
          multiple
          className="hidden"
          ref={fileInputRef}
          onChange={handleFiles}
        />
        <input
          type="file"
          // @ts-ignore
          webkitdirectory="true"
          directory="true"
          multiple
          className="hidden"
          ref={dirInputRef}
          onChange={handleFiles}
        />

        <AnimatePresence>
          {isScanning && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, height: 0 }}
              className="bg-slate-900/90 backdrop-blur-sm border border-slate-800 rounded-xl p-6 space-y-4 shadow-xl"
            >
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-3 text-cyan-400">
                  <Loader2 className="w-5 h-5 animate-spin" />
                  <span className="font-medium">{t.analyzing} {currentFile}</span>
                </div>
                <span className="text-slate-400">{progress}%</span>
              </div>
              <div className="h-2 w-full bg-slate-800 rounded-full overflow-hidden">
                <motion.div
                  className="h-full bg-cyan-500"
                  initial={{ width: 0 }}
                  animate={{ width: `${progress}%` }}
                  transition={{ ease: "linear" }}
                />
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Recent Scan Results view inside scanner */}
        {recentScans.length > 0 && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="space-y-4"
          >
            <h3 className="text-lg font-semibold text-slate-200">{t.current_results}</h3>
            <div className="bg-slate-900/90 backdrop-blur-sm border border-slate-800 rounded-xl overflow-hidden shadow-xl">
              <div className="overflow-x-auto">
                <table className="w-full text-left border-collapse">
                  <thead>
                    <tr className="border-b border-slate-800 bg-slate-900/50 text-slate-400 text-xs uppercase tracking-wider">
                      <th className="px-6 py-4 font-medium">{t.status}</th>
                      <th className="px-6 py-4 font-medium">{t.filename}</th>
                      <th className="px-6 py-4 font-medium">{t.score}</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800">
                    {recentScans.slice(0, 10).map((result) => (
                      <tr 
                        key={result.id} 
                        onClick={() => setSelectedResult(result)}
                        className="hover:bg-slate-800/50 transition-colors cursor-pointer"
                      >
                        <td className="px-6 py-4">
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
                            <span className="text-slate-200 text-sm font-medium truncate max-w-[300px]" title={result.filename}>
                              {result.filename}
                            </span>
                            <span className="text-slate-500 text-xs">
                              {(result.size / 1024).toFixed(1)} KB
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div className="text-sm font-mono text-slate-300 w-12">
                              {result.score.toFixed(2)}
                            </div>
                            <div className="w-24 h-1.5 bg-slate-800 rounded-full overflow-hidden">
                              <div 
                                className={`h-full ${result.isMalicious ? 'bg-red-500' : 'bg-emerald-500'}`}
                                style={{ width: `${Math.min(100, result.score * 100)}%` }}
                              />
                            </div>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {recentScans.length > 10 && (
                  <div className="px-6 py-4 border-t border-slate-800 text-center">
                    <span className="text-sm text-slate-500">
                      {t.showing_partial}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </motion.div>
        )}
      </div>

      <ScanDetailModal 
        isOpen={selectedResult !== null} 
        onClose={() => setSelectedResult(null)} 
        result={selectedResult}
        language={language}
      />
    </div>
  );
}

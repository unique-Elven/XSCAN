import { createContext, useContext, useState, ReactNode } from 'react';

export interface ModelConfig {
  path: string;
  threshold: number; // 0 to 1
}

export interface ScanResult {
  id: string;
  filename: string;
  size: number;
  modelUsed: 'Ember2018' | 'Ember2024';
  score: number;
  isMalicious: boolean;
  timestamp: number;
}

interface ScannerState {
  language: 'zh' | 'en';
  setLanguage: (lang: 'zh' | 'en') => void;
  ember2018: ModelConfig;
  ember2024: ModelConfig;
  updateConfig: (model: 'ember2018' | 'ember2024', config: Partial<ModelConfig>) => void;
  scanHistory: ScanResult[];
  addScanResults: (results: ScanResult[]) => void;
  clearHistory: () => void;
}

const ScannerContext = createContext<ScannerState | undefined>(undefined);

export function ScannerProvider({ children }: { children: ReactNode }) {
  const [language, setLanguage] = useState<'zh' | 'en'>('zh');

  const [ember2018, setEmber2018] = useState<ModelConfig>({
    path: 'C:\\Models\\ember2018\\model.bin',
    threshold: 0.65,
  });

  const [ember2024, setEmber2024] = useState<ModelConfig>({
    path: 'C:\\Models\\ember2024\\model.bin',
    threshold: 0.85,
  });

  const [scanHistory, setScanHistory] = useState<ScanResult[]>([]);

  const updateConfig = (model: 'ember2018' | 'ember2024', config: Partial<ModelConfig>) => {
    if (model === 'ember2018') {
      setEmber2018(prev => ({ ...prev, ...config }));
    } else {
      setEmber2024(prev => ({ ...prev, ...config }));
    }
  };

  const addScanResults = (results: ScanResult[]) => {
    setScanHistory(prev => [...results, ...prev]);
  };

  const clearHistory = () => {
    setScanHistory([]);
  };

  return (
    <ScannerContext.Provider
      value={{
        language,
        setLanguage,
        ember2018,
        ember2024,
        updateConfig,
        scanHistory,
        addScanResults,
        clearHistory
      }}
    >
      {children}
    </ScannerContext.Provider>
  );
}

export function useScanner() {
  const context = useContext(ScannerContext);
  if (context === undefined) {
    throw new Error('useScanner must be used within a ScannerProvider');
  }
  return context;
}

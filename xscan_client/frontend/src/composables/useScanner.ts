import { inject, provide, ref, type InjectionKey, type Ref } from 'vue'

export interface ModelConfig {
  path: string
  threshold: number
}

export interface ScanResult {
  id: string
  filename: string
  size: number
  /** Present when the row comes from a live scan */
  modelUsed?: 'Ember2018' | 'Ember2024'
  score?: number
  isMalicious: boolean
  timestamp: number
  /** SHA-256 hex when loaded from DB */
  fileHash?: string
  /** SQLite row id when loaded from disk (for bulk delete) */
  historyDbId?: number
  /** Persisted engine tag: 2018 | 2024 */
  engine?: string
}

export type ScannerLang = 'zh' | 'en'

export interface ScannerContext {
  language: Ref<ScannerLang>
  setLanguage: (lang: ScannerLang) => void
  soundEnabled: Ref<boolean>
  modelStrategy: Ref<string>
  ember2018: Ref<ModelConfig>
  ember2024: Ref<ModelConfig>
  updateConfig: (model: 'ember2018' | 'ember2024', config: Partial<ModelConfig>) => void
  scanHistory: Ref<ScanResult[]>
  addScanResults: (results: ScanResult[]) => void
  clearHistory: () => Promise<void>
  refreshHistoryFromDisk: () => Promise<void>
}

const ScannerKey: InjectionKey<ScannerContext> = Symbol('scanner')

function isWailsRuntime(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { go?: unknown }).go
}

export async function hydrateScannerFromDisk(ctx: ScannerContext): Promise<void> {
  if (!isWailsRuntime()) return
  try {
    const { GetAppConfig } = await import('../../wailsjs/go/main/App')
    const cfg = await GetAppConfig()
    ctx.ember2018.value = {
      path: cfg.path2018 ?? '',
      threshold: Number(cfg.threshold2018) || 0.65,
    }
    ctx.ember2024.value = {
      path: cfg.path2024 ?? '',
      threshold: Number(cfg.threshold2024) || 0.85,
    }
    if (cfg.language === 'en' || cfg.language === 'zh') {
      ctx.setLanguage(cfg.language)
    }
    ctx.soundEnabled.value = !!cfg.soundEnabled
    ctx.modelStrategy.value = (cfg.modelStrategy || 'auto').trim() || 'auto'
  } catch (e) {
    console.error(e)
  }
}

function rowToScanResult(r: {
  id: number
  scannedAt: number
  filePath: string
  fileHash: string
  verdict: string
  fileSize: number
  engine?: string
  score?: number
}): ScanResult {
  const eng = (r.engine || '').trim()
  let modelUsed: 'Ember2018' | 'Ember2024' | undefined
  if (eng === '2024') {
    modelUsed = 'Ember2024'
  } else if (eng === '2018') {
    modelUsed = 'Ember2018'
  }
  const sc = r.score
  return {
    id: String(r.id),
    historyDbId: r.id,
    filename: r.filePath,
    size: Number(r.fileSize) || 0,
    isMalicious: r.verdict === 'malicious',
    timestamp: r.scannedAt,
    fileHash: r.fileHash || undefined,
    engine: eng || undefined,
    score: sc !== undefined && sc !== null ? Number(sc) : undefined,
    modelUsed,
  }
}

export function provideScanner(): ScannerContext {
  const language = ref<ScannerLang>('zh')
  const soundEnabled = ref(false)
  const modelStrategy = ref('auto')

  const ember2018 = ref<ModelConfig>({
    path: '',
    threshold: 0.65,
  })

  const ember2024 = ref<ModelConfig>({
    path: '',
    threshold: 0.85,
  })

  const scanHistory = ref<ScanResult[]>([])

  const setLanguage = (lang: ScannerLang) => {
    language.value = lang
  }

  const updateConfig = (model: 'ember2018' | 'ember2024', config: Partial<ModelConfig>) => {
    if (model === 'ember2018') {
      ember2018.value = { ...ember2018.value, ...config }
    } else {
      ember2024.value = { ...ember2024.value, ...config }
    }
  }

  const addScanResults = (results: ScanResult[]) => {
    scanHistory.value = [...results, ...scanHistory.value]
  }

  const refreshHistoryFromDisk = async () => {
    if (!isWailsRuntime()) return
    try {
      const { GetHistory } = await import('../../wailsjs/go/main/App')
      const rows = await GetHistory(500)
      scanHistory.value = rows.map(rowToScanResult)
    } catch (e) {
      console.error(e)
    }
  }

  const clearHistory = async () => {
    if (isWailsRuntime()) {
      try {
        const { ClearScanHistory } = await import('../../wailsjs/go/main/App')
        await ClearScanHistory()
      } catch (e) {
        console.error(e)
      }
    }
    scanHistory.value = []
  }

  const ctx: ScannerContext = {
    language,
    setLanguage,
    soundEnabled,
    modelStrategy,
    ember2018,
    ember2024,
    updateConfig,
    scanHistory,
    addScanResults,
    clearHistory,
    refreshHistoryFromDisk,
  }

  provide(ScannerKey, ctx)
  return ctx
}

export function useScanner(): ScannerContext {
  const ctx = inject(ScannerKey)
  if (!ctx) {
    throw new Error('useScanner must be used within a ScannerProvider')
  }
  return ctx
}

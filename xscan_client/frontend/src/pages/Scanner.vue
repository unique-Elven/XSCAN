<script lang="ts">
export default {
  name: 'ScannerPage',
}
</script>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  FolderOpen,
  FileText,
  AlertTriangle,
  CheckCircle,
  ShieldCheck,
  Loader2,
} from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { useScanner, type ScanResult } from '../composables/useScanner'
import { translations } from '../i18n'
import ScanDetailModal from '../components/ScanDetailModal.vue'
import {
  ConfigureScanner,
  ScanFiles,
  PickScanFiles,
  PickScanDirectory,
} from '../../wailsjs/go/main/App'

const { ember2018, ember2024, addScanResults, language } = useScanner()

const t = computed(() => translations[language.value])

const isScanning = ref(false)
const progress = ref(0)
const currentFile = ref('')
const recentScans = ref<ScanResult[]>([])
const selectedResult = ref<ScanResult | null>(null)

const fileInputRef = ref<HTMLInputElement | null>(null)
const dirInputRef = ref<HTMLInputElement | null>(null)

function isWails(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { go?: unknown }).go
}

function basename(p: string): string {
  const norm = p.replace(/\\/g, '/')
  const i = norm.lastIndexOf('/')
  return i >= 0 ? norm.slice(i + 1) : p
}

async function startScan(_options?: { paths?: string[] }): Promise<void> {
  const paths = _options?.paths
  if (paths?.length) {
    await runWailsScan(paths)
  }
}

async function triggerFileSelect() {
  if (isWails()) {
    try {
      const paths = await PickScanFiles()
      if (paths?.length) await runWailsScan(paths)
    } catch (e) {
      toast.error(String(e))
    }
    return
  }
  const el = fileInputRef.value
  if (el) {
    el.value = ''
    el.click()
  }
}

async function triggerDirSelect() {
  if (isWails()) {
    try {
      const paths = await PickScanDirectory()
      if (paths?.length) await runWailsScan(paths)
    } catch (e) {
      toast.error(String(e))
    }
    return
  }
  const el = dirInputRef.value
  if (el) {
    el.value = ''
    el.click()
  }
}

async function runWailsScan(absPaths: string[]) {
  const p18 = ember2018.value.path?.trim()
  const p24 = ember2024.value.path?.trim()
  if (!p18 || !p24) {
    toast.error(t.value.scan_models_both_required)
    return
  }

  isScanning.value = true
  progress.value = 0
  recentScans.value = []

  try {
    await ConfigureScanner({ path2018: p18, path2024: p24 })
    const batch = await ScanFiles(absPaths)
    const newResults: ScanResult[] = []
    let i = 0
    for (const row of batch) {
      i++
      currentFile.value = basename(row.path)
      progress.value = Math.round((i / batch.length) * 100)

      if (row.errMsg) {
        toast.error(`${basename(row.path)}: ${row.errMsg}`)
        continue
      }

      const score = row.score
      const modelUsed =
        row.modelUsed === 'Ember2024' ? 'Ember2024' : 'Ember2018'
      const threshold =
        modelUsed === 'Ember2024' ? ember2024.value.threshold : ember2018.value.threshold

      const result: ScanResult = {
        id: crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).substring(7),
        filename: row.path,
        size: Number(row.size),
        modelUsed,
        score: Number(score.toFixed(4)),
        isMalicious: score >= threshold,
        timestamp: Date.now(),
      }
      newResults.push(result)
      recentScans.value = [result, ...recentScans.value]
    }

    addScanResults(newResults)

    const maliciousCount = newResults.filter((r) => r.isMalicious).length
    if (newResults.length === 0) {
      /* only errors */
    } else if (maliciousCount > 0) {
      /* 风险隔离弹窗由后端 scan_completed 事件在 App.vue 全局挂载 */
    } else {
      toast.success(t.value.scan_complete_safe)
    }
  } catch (e) {
    toast.error(String(e))
  } finally {
    isScanning.value = false
    currentFile.value = ''
  }
}

/** Vite 纯前端开发：无 Wails 时使用模拟扫描 */
async function handleFiles(e: Event) {
  const target = e.target as HTMLInputElement
  const files = target.files
  if (!files || files.length === 0) return

  if (isWails()) {
    target.value = ''
    return
  }

  await startScan({
    paths: Array.from(files).map((f) => f.webkitRelativePath || f.name),
  })

  isScanning.value = true
  progress.value = 0
  recentScans.value = []

  const fileList = Array.from(files)
  const totalFiles = fileList.length
  const newResults: ScanResult[] = []

  for (let i = 0; i < totalFiles; i++) {
    const file = fileList[i]
    currentFile.value = file.name

    const delay = Math.max(100, Math.min(800, file.size / 1024))
    await new Promise((resolve) => setTimeout(resolve, delay))

    const randomFactor = Math.random()
    const baseScore = randomFactor * 1.2
    const score = Math.min(1.0, Math.max(0.0, baseScore))

    const used2018 = Math.random() > 0.5
    const currentModelName = used2018 ? 'Ember2018' : 'Ember2024'
    const currentThreshold = used2018 ? ember2018.value.threshold : ember2024.value.threshold

    const isMalicious = score >= currentThreshold

    const result: ScanResult = {
      id: crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).substring(7),
      filename: file.webkitRelativePath || file.name,
      size: file.size,
      modelUsed: currentModelName,
      score: Number(score.toFixed(4)),
      isMalicious,
      timestamp: Date.now(),
    }

    newResults.push(result)
    recentScans.value = [result, ...recentScans.value]
    progress.value = Math.round(((i + 1) / totalFiles) * 100)
  }

  addScanResults(newResults)

  isScanning.value = false
  currentFile.value = ''

  const maliciousCount = newResults.filter((r) => r.isMalicious).length
  if (maliciousCount > 0) {
    toast.error(t.value.scan_complete_malicious)
  } else {
    toast.success(t.value.scan_complete_safe)
  }
}

function openDetail(result: ScanResult) {
  selectedResult.value = result
}

function closeDetail() {
  selectedResult.value = null
}
</script>

<template>
  <div class="relative min-h-full min-w-0 max-w-full">
    <div class="absolute inset-x-0 bottom-0 h-96 pointer-events-none overflow-hidden opacity-30">
      <div class="absolute inset-0 flex items-end justify-center">
        <div class="w-[800px] h-[400px] relative animate-fade-in-cpu">
          <div
            class="absolute inset-0 bg-[linear-gradient(to_right,#0891b222_1px,transparent_1px),linear-gradient(to_bottom,#0891b222_1px,transparent_1px)] bg-[size:40px_40px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_100%,#000_70%,transparent_100%)]"
          />

          <div
            class="absolute bottom-10 left-1/2 -translate-x-1/2 w-32 h-32 bg-slate-900 border border-cyan-500/50 rounded-xl flex items-center justify-center z-10 animate-cpu-glow"
          >
            <div class="w-24 h-24 border border-cyan-400/30 rounded-lg flex items-center justify-center">
              <div
                class="w-16 h-16 bg-cyan-950 rounded flex items-center justify-center relative overflow-hidden"
              >
                <div class="absolute inset-0 bg-cyan-500/20" />
                <div
                  class="absolute inset-0 h-1/2 bg-gradient-to-b from-transparent via-cyan-400/50 to-transparent animate-cpu-scan"
                />
                <ShieldCheck class="w-8 h-8 text-cyan-400 z-10" />
              </div>
            </div>
          </div>

          <svg
            class="absolute inset-0 w-full h-full"
            style="filter: drop-shadow(0 0 8px rgba(6, 182, 212, 0.5))"
          >
            <path
              d="M400,320 L250,320 L250,200 L150,200"
              fill="none"
              stroke="rgba(6,182,212,0.5)"
              stroke-width="2"
              pathLength="1"
              class="[stroke-dasharray:1] animate-path-draw"
            />
            <path
              d="M400,320 L550,320 L550,150 L650,150"
              fill="none"
              stroke="rgba(6,182,212,0.5)"
              stroke-width="2"
              pathLength="1"
              class="[stroke-dasharray:1] animate-path-draw-slow"
            />
            <path
              d="M400,280 L400,100 L300,50"
              fill="none"
              stroke="rgba(6,182,212,0.5)"
              stroke-width="2"
              pathLength="1"
              class="[stroke-dasharray:1] animate-path-draw-fast"
            />
          </svg>
        </div>
      </div>
    </div>

    <div class="p-8 max-w-5xl mx-auto space-y-8 relative z-10">
      <div class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
        <div class="flex-1 space-y-3">
          <div>
            <h2 class="text-3xl font-bold text-slate-100 mb-2">{{ t.scan_title }}</h2>
            <p class="text-slate-400">{{ t.scan_desc }}</p>
          </div>
          <p class="text-sm text-slate-500">{{ t.scan_backend_route_hint }}</p>
        </div>
        <div
          class="bg-slate-900 px-4 py-2 rounded-lg border border-slate-800 flex items-center gap-3 shadow-lg shadow-black/50 shrink-0"
        >
          <ShieldCheck class="w-5 h-5 text-cyan-500" />
          <div class="flex flex-col">
            <span class="text-xs text-slate-500 font-medium uppercase">{{ t.dual_engine }}</span>
            <span class="text-sm text-slate-200 font-medium">{{ t.backend_auto }}</span>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <button
          type="button"
          class="relative group overflow-hidden bg-slate-900/80 backdrop-blur-sm border-2 border-dashed border-slate-700 hover:border-cyan-500/50 rounded-2xl p-10 flex flex-col items-center justify-center gap-4 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-xl hover:scale-[1.02] hover:-translate-y-0.5 active:scale-[0.98]"
          :disabled="isScanning"
          @click="triggerFileSelect"
        >
          <div
            class="p-4 bg-slate-800 rounded-full group-hover:bg-cyan-500/10 transition-colors text-slate-400 group-hover:text-cyan-400"
          >
            <FileText class="w-10 h-10" />
          </div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-slate-200">{{ t.scan_files }}</h3>
            <p class="text-sm text-slate-500 mt-1">{{ t.scan_files_desc }}</p>
          </div>
        </button>

        <button
          type="button"
          class="relative group overflow-hidden bg-slate-900/80 backdrop-blur-sm border-2 border-dashed border-slate-700 hover:border-cyan-500/50 rounded-2xl p-10 flex flex-col items-center justify-center gap-4 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-xl hover:scale-[1.02] hover:-translate-y-0.5 active:scale-[0.98]"
          :disabled="isScanning"
          @click="triggerDirSelect"
        >
          <div
            class="p-4 bg-slate-800 rounded-full group-hover:bg-cyan-500/10 transition-colors text-slate-400 group-hover:text-cyan-400"
          >
            <FolderOpen class="w-10 h-10" />
          </div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-slate-200">{{ t.scan_dir }}</h3>
            <p class="text-sm text-slate-500 mt-1">{{ t.scan_dir_desc }}</p>
          </div>
        </button>
      </div>

      <input
        ref="fileInputRef"
        type="file"
        multiple
        class="hidden"
        @change="handleFiles"
      />
      <input
        ref="dirInputRef"
        type="file"
        webkitdirectory
        directory
        multiple
        class="hidden"
        @change="handleFiles"
      />

      <Transition name="scan-panel">
        <div
          v-if="isScanning"
          class="bg-slate-900/90 backdrop-blur-sm border border-slate-800 rounded-xl p-6 space-y-4 shadow-xl"
        >
          <div class="flex items-center justify-between text-sm">
            <div class="flex items-center gap-3 text-cyan-400">
              <Loader2 class="w-5 h-5 animate-spin" />
              <span class="font-medium">{{ t.analyzing }} {{ currentFile }}</span>
            </div>
            <span class="text-slate-400">{{ progress }}%</span>
          </div>
          <div class="h-2 w-full bg-slate-800 rounded-full overflow-hidden">
            <div
              class="h-full bg-cyan-500 transition-[width] duration-150 ease-linear"
              :style="{ width: `${progress}%` }"
            />
          </div>
        </div>
      </Transition>

      <div v-if="recentScans.length > 0" class="space-y-4 animate-fade-in-short">
        <h3 class="text-lg font-semibold text-slate-200">{{ t.current_results }}</h3>
        <div class="bg-slate-900/90 backdrop-blur-sm border border-slate-800 rounded-xl overflow-hidden shadow-xl">
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr
                  class="border-b border-slate-800 bg-slate-900/50 text-slate-400 text-xs uppercase tracking-wider"
                >
                  <th class="px-6 py-4 font-medium">{{ t.status }}</th>
                  <th class="px-6 py-4 font-medium">{{ t.filename }}</th>
                  <th class="px-6 py-4 font-medium">{{ t.score }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-800">
                <tr
                  v-for="result in recentScans.slice(0, 10)"
                  :key="result.id"
                  class="hover:bg-slate-800/50 transition-colors cursor-pointer"
                  @click="openDetail(result)"
                >
                  <td class="px-6 py-4">
                    <div
                      v-if="result.isMalicious"
                      class="flex items-center gap-2 text-red-400 bg-red-400/10 w-fit px-2.5 py-1 rounded-full text-xs font-medium border border-red-400/20"
                    >
                      <AlertTriangle class="w-3.5 h-3.5" />
                      {{ t.malicious }}
                    </div>
                    <div
                      v-else
                      class="flex items-center gap-2 text-emerald-400 bg-emerald-400/10 w-fit px-2.5 py-1 rounded-full text-xs font-medium border border-emerald-400/20"
                    >
                      <CheckCircle class="w-3.5 h-3.5" />
                      {{ t.safe }}
                    </div>
                  </td>
                  <td class="px-6 py-4">
                    <div class="flex flex-col">
                      <span
                        class="text-slate-200 text-sm font-medium truncate max-w-[300px]"
                        :title="result.filename"
                      >
                        {{ result.filename }}
                      </span>
                      <span class="text-slate-500 text-xs">
                        {{ (result.size / 1024).toFixed(1) }} KB
                      </span>
                    </div>
                  </td>
                  <td class="px-6 py-4">
                    <div class="flex items-center gap-3">
                      <div class="text-sm font-mono text-slate-300 w-12">
                        {{ result.score !== undefined ? result.score.toFixed(2) : '—' }}
                      </div>
                      <div class="w-24 h-1.5 bg-slate-800 rounded-full overflow-hidden">
                        <div
                          class="h-full"
                          :class="result.isMalicious ? 'bg-red-500' : 'bg-emerald-500'"
                          :style="{
                            width: `${Math.min(100, (result.score ?? 0) * 100)}%`,
                          }"
                        />
                      </div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            <div
              v-if="recentScans.length > 10"
              class="px-6 py-4 border-t border-slate-800 text-center"
            >
              <span class="text-sm text-slate-500">
                {{ t.showing_partial }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ScanDetailModal
      :is-open="selectedResult !== null"
      :result="selectedResult"
      :language="language"
      @close="closeDetail"
    />
  </div>
</template>

<style scoped>
.scan-panel-enter-active,
.scan-panel-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.scan-panel-enter-from,
.scan-panel-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>

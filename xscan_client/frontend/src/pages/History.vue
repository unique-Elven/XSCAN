<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { AlertTriangle, CheckCircle, Trash2, Calendar } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { useScanner, type ScanResult } from '../composables/useScanner'
import { translations } from '../i18n'
import ScanDetailModal from '../components/ScanDetailModal.vue'
import { DeleteHistory } from '../../wailsjs/go/main/App'

const { scanHistory, clearHistory, refreshHistoryFromDisk, language } = useScanner()

const t = computed(() => translations[language.value])

const selectedIds = ref<number[]>([])
const selectedResult = ref<ScanResult | null>(null)

const selectableIds = computed(() =>
  scanHistory.value.map((r) => r.historyDbId).filter((id): id is number => typeof id === 'number')
)

const allSelected = computed(() => {
  const ids = selectableIds.value
  return ids.length > 0 && ids.every((id) => selectedIds.value.includes(id))
})

onMounted(() => {
  refreshHistoryFromDisk()
})

function formatEngine(r: ScanResult): string {
  if (r.engine === '2018' || r.engine === '2024') {
    return r.engine === '2024' ? 'Ember 2024' : 'Ember 2018'
  }
  return r.modelUsed ?? '—'
}

function toggleRow(id: number | undefined, checked: boolean) {
  if (id === undefined) return
  if (checked) {
    if (!selectedIds.value.includes(id)) {
      selectedIds.value = [...selectedIds.value, id]
    }
  } else {
    selectedIds.value = selectedIds.value.filter((x) => x !== id)
  }
}

function toggleSelectAll(checked: boolean) {
  if (!checked) {
    selectedIds.value = []
    return
  }
  selectedIds.value = [...selectableIds.value]
}

function isRowSelected(id: number | undefined): boolean {
  return id !== undefined && selectedIds.value.includes(id)
}

function openDetail(result: ScanResult) {
  selectedResult.value = result
}

function closeDetail() {
  selectedResult.value = null
}

function isWails(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { go?: unknown }).go
}

async function deleteSelected() {
  const ids = selectedIds.value
  if (!ids.length) {
    toast.info(t.value.history_none_selected)
    return
  }
  if (!isWails()) {
    toast.error('Wails only')
    return
  }
  try {
    await DeleteHistory(ids)
    selectedIds.value = []
    toast.success(t.value.history_delete_success)
    await refreshHistoryFromDisk()
  } catch (e) {
    toast.error(String(e))
  }
}
</script>

<template>
  <div class="max-w-full min-w-0 p-6 md:p-8 mx-auto space-y-8">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div class="min-w-0">
        <h2 class="text-3xl font-bold text-slate-100 mb-2">{{ t.history_title }}</h2>
        <p class="text-slate-400">{{ t.history_desc }}</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 shrink-0">
        <button
          v-if="scanHistory.length > 0 && selectedIds.length > 0"
          type="button"
          class="flex items-center gap-2 bg-red-600/90 hover:bg-red-500 text-white px-4 py-2 rounded-lg transition-colors text-sm font-medium shadow-lg shadow-red-900/20"
          @click="deleteSelected"
        >
          <Trash2 class="w-4 h-4" />
          {{ t.history_delete_selected }}
        </button>
        <button
          v-if="scanHistory.length > 0"
          type="button"
          class="flex items-center gap-2 bg-slate-900 border border-slate-800 hover:border-red-500/50 hover:text-red-400 text-slate-400 px-4 py-2 rounded-lg transition-colors text-sm"
          @click="() => clearHistory()"
        >
          <Trash2 class="w-4 h-4" />
          {{ t.clear_log }}
        </button>
      </div>
    </div>

    <div
      v-if="scanHistory.length === 0"
      class="bg-slate-900 border border-slate-800 rounded-xl p-16 flex flex-col items-center justify-center text-center"
    >
      <div class="w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center mb-4 text-slate-500">
        <Calendar class="w-8 h-8" />
      </div>
      <h3 class="text-xl font-medium text-slate-300">{{ t.no_history }}</h3>
      <p class="text-slate-500 mt-2">{{ t.no_history_desc }}</p>
    </div>

    <div v-else class="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden min-w-0">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[640px]">
          <thead>
            <tr class="border-b border-slate-800 bg-slate-900/50 text-slate-400 text-xs uppercase tracking-wider">
              <th class="w-10 px-3 py-3 font-medium">
                <input
                  type="checkbox"
                  class="rounded border-slate-600 bg-slate-900 text-cyan-500 focus:ring-cyan-500/40"
                  :checked="allSelected"
                  @change="toggleSelectAll(($event.target as HTMLInputElement).checked)"
                />
              </th>
              <th class="px-4 py-3 font-medium">{{ t.status }}</th>
              <th class="px-4 py-3 font-medium">{{ t.filename }}</th>
              <th class="px-4 py-3 font-medium whitespace-nowrap">{{ t.history_col_engine }}</th>
              <th class="px-4 py-3 font-medium whitespace-nowrap">{{ t.history_col_score }}</th>
              <th class="px-4 py-3 font-medium whitespace-nowrap">{{ t.date }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800">
            <tr
              v-for="result in scanHistory"
              :key="result.id"
              class="hover:bg-slate-800/50 transition-colors cursor-pointer"
              @click="openDetail(result)"
            >
              <td class="px-3 py-3 align-middle" @click.stop>
                <input
                  type="checkbox"
                  class="rounded border-slate-600 bg-slate-900 text-cyan-500 focus:ring-cyan-500/40"
                  :checked="isRowSelected(result.historyDbId)"
                  @change="
                    toggleRow(result.historyDbId, ($event.target as HTMLInputElement).checked)
                  "
                />
              </td>
              <td class="px-4 py-3 whitespace-nowrap">
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
              <td class="px-4 py-3 min-w-0">
                <div class="flex flex-col min-w-0">
                  <span
                    class="text-slate-200 text-sm font-medium truncate max-w-[280px] md:max-w-md"
                    :title="result.filename"
                  >
                    {{ result.filename }}
                  </span>
                  <span class="text-slate-500 text-xs">
                    {{ (result.size / 1024).toFixed(1) }} KB
                  </span>
                  <span
                    v-if="result.fileHash"
                    class="text-slate-600 text-[10px] font-mono truncate max-w-[280px] md:max-w-md block"
                    :title="result.fileHash"
                  >
                    {{ result.fileHash.slice(0, 12) }}…{{ result.fileHash.slice(-8) }}
                  </span>
                </div>
              </td>
              <td class="px-4 py-3 whitespace-nowrap text-xs text-slate-300 font-medium">
                {{ formatEngine(result) }}
              </td>
              <td class="px-4 py-3 whitespace-nowrap">
                <span
                  class="text-sm font-mono"
                  :class="result.isMalicious ? 'text-red-400' : 'text-slate-300'"
                >
                  {{ result.score !== undefined ? result.score.toFixed(4) : '—' }}
                </span>
              </td>
              <td class="px-4 py-3 whitespace-nowrap text-sm text-slate-500">
                {{ new Date(result.timestamp).toLocaleString() }}
              </td>
            </tr>
          </tbody>
        </table>
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

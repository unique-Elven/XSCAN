<script setup lang="ts">
import { computed } from 'vue'
import { X, ShieldAlert, ShieldCheck, Database, HardDrive, FileText } from 'lucide-vue-next'
import type { ScanResult } from '../composables/useScanner'
import { translations } from '../i18n'

const props = defineProps<{
  isOpen: boolean
  result: ScanResult | null
  language: 'zh' | 'en'
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const t = computed(() => translations[props.language])

function onClose() {
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-overlay">
      <div
        v-if="isOpen && result"
        class="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm"
        @click="onClose"
      />
    </Transition>
    <Transition name="modal-panel">
      <div
        v-if="isOpen && result"
        class="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-full max-w-md bg-slate-900 border border-slate-700 rounded-2xl shadow-2xl shadow-cyan-900/20 overflow-hidden"
      >
        <div class="flex items-center justify-between p-4 border-b border-slate-800">
          <h3 class="text-lg font-semibold text-slate-200">{{ t.detail_title }}</h3>
          <button
            type="button"
            class="text-slate-400 hover:text-slate-200 transition-colors p-1"
            @click="onClose"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-6 space-y-6">
          <div
            class="flex flex-col items-center justify-center p-6 bg-slate-950/50 rounded-xl border border-slate-800"
          >
            <div
              v-if="result.isMalicious"
              class="w-16 h-16 bg-red-500/10 text-red-400 rounded-full flex items-center justify-center mb-4 border border-red-500/20"
            >
              <ShieldAlert class="w-8 h-8" />
            </div>
            <div
              v-else
              class="w-16 h-16 bg-emerald-500/10 text-emerald-400 rounded-full flex items-center justify-center mb-4 border border-emerald-500/20"
            >
              <ShieldCheck class="w-8 h-8" />
            </div>

            <h4
              class="text-xl font-bold mb-1"
              :class="result.isMalicious ? 'text-red-400' : 'text-emerald-400'"
            >
              {{ result.isMalicious ? t.malicious : t.safe }}
            </h4>
            <p class="text-sm text-slate-400 font-mono text-center break-all px-4">
              {{ result.filename }}
            </p>
          </div>

          <div class="space-y-4">
            <div class="flex items-center justify-between py-2 border-b border-slate-800/50">
              <div class="flex items-center gap-2 text-slate-400">
                <Database class="w-4 h-4" />
                <span class="text-sm">{{ t.lightgbm_score }}</span>
              </div>
              <span
                class="text-lg font-mono font-bold"
                :class="result.isMalicious ? 'text-red-400' : 'text-cyan-400'"
              >
                {{ result.score !== undefined ? result.score.toFixed(2) : '—' }}
              </span>
            </div>

            <div class="flex items-center justify-between py-2 border-b border-slate-800/50">
              <div class="flex items-center gap-2 text-slate-400">
                <HardDrive class="w-4 h-4" />
                <span class="text-sm">{{ t.engine }}</span>
              </div>
              <span class="text-sm font-medium text-slate-200">{{ result.modelUsed ?? '—' }}</span>
            </div>

            <div
              v-if="result.fileHash"
              class="flex items-start justify-between gap-2 py-2 border-b border-slate-800/50"
            >
              <span class="text-sm text-slate-400 shrink-0">SHA256</span>
              <span class="text-xs font-mono text-slate-300 text-right break-all">{{ result.fileHash }}</span>
            </div>

            <div class="flex items-center justify-between py-2 border-b border-slate-800/50">
              <div class="flex items-center gap-2 text-slate-400">
                <FileText class="w-4 h-4" />
                <span class="text-sm">{{ t.file_size }}</span>
              </div>
              <span class="text-sm font-medium text-slate-200">
                {{ (result.size / 1024).toFixed(1) }} KB
              </span>
            </div>
          </div>
        </div>

        <div class="p-4 border-t border-slate-800 bg-slate-900/50 flex justify-end">
          <button
            type="button"
            class="px-6 py-2 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-lg transition-colors font-medium text-sm"
            @click="onClose"
          >
            Close
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-overlay-enter-active,
.modal-overlay-leave-active {
  transition: opacity 0.2s ease;
}
.modal-overlay-enter-from,
.modal-overlay-leave-to {
  opacity: 0;
}

.modal-panel-enter-active,
.modal-panel-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.modal-panel-enter-from,
.modal-panel-leave-to {
  opacity: 0;
  transform: translate(-50%, -50%) scale(0.95) translateY(20px);
}
</style>

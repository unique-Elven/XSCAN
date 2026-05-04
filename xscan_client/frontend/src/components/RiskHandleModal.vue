<script setup lang="ts">
import { ref, watch } from 'vue'
import { X } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  title: string
  paths: string[]
  ignoreLabel: string
  confirmLabel: string
}>()

const emit = defineEmits<{
  close: []
  quarantine: [paths: string[]]
}>()

/** Paths selected for quarantine (default: all) */
const checkedSet = ref<Set<string>>(new Set())

watch(
  () => [props.open, props.paths] as const,
  ([isOpen, paths]) => {
    if (!isOpen) return
    checkedSet.value = new Set(paths)
  },
  { immediate: true }
)

function toggle(path: string, on: boolean) {
  const next = new Set(checkedSet.value)
  if (on) {
    next.add(path)
  } else {
    next.delete(path)
  }
  checkedSet.value = next
}

function onIgnore() {
  emit('close')
}

function onConfirm() {
  const selected = props.paths.filter((p) => checkedSet.value.has(p))
  emit('quarantine', selected)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="risk-overlay">
      <div
        v-if="open"
        class="fixed inset-0 z-[350] flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
        role="dialog"
        aria-modal="true"
        @click.self="onIgnore"
      >
        <Transition name="risk-panel">
          <div
            v-if="open"
            class="w-full max-w-lg max-h-[85vh] flex flex-col bg-slate-900 border border-slate-700 rounded-xl shadow-2xl shadow-black/50 overflow-hidden"
            @click.stop
          >
            <div class="flex items-start justify-between gap-3 px-5 pt-5 pb-3 border-b border-slate-800">
              <h3 class="text-lg font-semibold text-red-400 leading-snug">
                {{ title }}
              </h3>
              <button
                type="button"
                class="p-1 rounded-lg text-slate-500 hover:text-slate-300 hover:bg-slate-800 shrink-0"
                aria-label="close"
                @click="onIgnore"
              >
                <X class="w-5 h-5" />
              </button>
            </div>

            <div class="px-5 py-3 flex-1 min-h-0 flex flex-col gap-2">
              <div
                class="flex-1 min-h-[160px] max-h-[45vh] overflow-y-auto rounded-lg border border-slate-800 bg-slate-950/80 py-2"
              >
                <ul class="divide-y divide-slate-800/80">
                  <li
                    v-for="p in paths"
                    :key="p"
                    class="flex items-start gap-3 px-3 py-2.5 hover:bg-slate-800/40"
                  >
                    <input
                      type="checkbox"
                      class="mt-1 rounded border-slate-600 bg-slate-900 text-cyan-500 focus:ring-cyan-500/40 shrink-0"
                      :checked="checkedSet.has(p)"
                      @change="toggle(p, ($event.target as HTMLInputElement).checked)"
                    />
                    <span class="text-sm text-slate-300 break-all font-mono leading-relaxed" :title="p">
                      {{ p }}
                    </span>
                  </li>
                </ul>
              </div>
            </div>

            <div class="flex justify-end gap-3 px-5 py-4 border-t border-slate-800 bg-slate-900/95">
              <button
                type="button"
                class="px-4 py-2 rounded-lg border border-slate-600 text-slate-300 text-sm font-medium hover:bg-slate-800 transition-colors"
                @click="onIgnore"
              >
                {{ ignoreLabel }}
              </button>
              <button
                type="button"
                class="px-4 py-2 rounded-lg bg-red-600 hover:bg-red-500 text-white text-sm font-semibold shadow-lg shadow-red-900/30 transition-colors"
                @click="onConfirm"
              >
                {{ confirmLabel }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.risk-overlay-enter-active,
.risk-overlay-leave-active {
  transition: opacity 0.2s ease;
}
.risk-overlay-enter-from,
.risk-overlay-leave-to {
  opacity: 0;
}
.risk-panel-enter-active,
.risk-panel-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.risk-panel-enter-from,
.risk-panel-leave-to {
  opacity: 0;
  transform: scale(0.96) translateY(8px);
}
</style>

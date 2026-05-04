<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Save, SlidersHorizontal, HardDrive, Check, Globe, FolderOpen } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { useScanner } from '../composables/useScanner'
import { translations } from '../i18n'
import { ConfigureScanner, PickLightGBMModelFile, SaveAppConfig } from '../../wailsjs/go/main/App'

const {
  ember2018,
  ember2024,
  updateConfig,
  language,
  setLanguage,
  soundEnabled,
  modelStrategy,
} = useScanner()

const t = computed(() => translations[language.value])

const e2018Path = ref(ember2018.value.path)
const e2018Threshold = ref(ember2018.value.threshold)
const e2024Path = ref(ember2024.value.path)
const e2024Threshold = ref(ember2024.value.threshold)

watch(
  [ember2018, ember2024],
  () => {
    e2018Path.value = ember2018.value.path
    e2018Threshold.value = ember2018.value.threshold
    e2024Path.value = ember2024.value.path
    e2024Threshold.value = ember2024.value.threshold
  },
  { deep: true }
)

async function handleSave() {
  updateConfig('ember2018', { path: e2018Path.value, threshold: e2018Threshold.value })
  updateConfig('ember2024', { path: e2024Path.value, threshold: e2024Threshold.value })
  if (isWails()) {
    try {
      await SaveAppConfig({
        modelStrategy: modelStrategy.value || 'auto',
        soundEnabled: soundEnabled.value,
        language: language.value,
        path2018: e2018Path.value.trim(),
        path2024: e2024Path.value.trim(),
        threshold2018: e2018Threshold.value,
        threshold2024: e2024Threshold.value,
      })
      await ConfigureScanner({
        path2018: e2018Path.value.trim(),
        path2024: e2024Path.value.trim(),
      })
    } catch (e) {
      toast.error(String(e))
      return
    }
  }
  toast.success(t.value.settings_saved)
}

function isWails(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { go?: unknown }).go
}

async function browseModel(which: '2018' | '2024') {
  if (!isWails()) {
    toast.info(t.value.pick_model_desktop_only)
    return
  }
  try {
    const p = await PickLightGBMModelFile()
    if (p) {
      if (which === '2018') {
        e2018Path.value = p
      } else {
        e2024Path.value = p
      }
    }
  } catch (e) {
    toast.error(String(e))
  }
}
</script>

<template>
  <div class="p-8 max-w-4xl mx-auto space-y-8 pb-20">
    <div>
      <h2 class="text-3xl font-bold text-slate-100 mb-2">{{ t.setting_title }}</h2>
      <p class="text-slate-400">{{ t.setting_desc }}</p>
    </div>

    <div class="grid gap-6">
      <div class="bg-slate-900 border border-slate-800 rounded-xl p-6">
        <div class="flex items-center gap-3 mb-4">
          <Globe class="w-5 h-5 text-slate-400" />
          <h3 class="text-lg font-semibold text-slate-200">{{ t.language }}</h3>
        </div>
        <div class="flex gap-4">
          <button
            type="button"
            class="flex-1 flex items-center justify-between p-4 rounded-lg border-2 transition-all"
            :class="
              language === 'zh'
                ? 'border-cyan-500 bg-cyan-500/10'
                : 'border-slate-800 bg-slate-950 hover:border-slate-700'
            "
            @click="setLanguage('zh')"
          >
            <span class="font-semibold" :class="language === 'zh' ? 'text-cyan-400' : 'text-slate-300'">
              {{ t.lang_zh }}
            </span>
            <Check v-if="language === 'zh'" class="text-cyan-500 w-5 h-5" />
          </button>
          <button
            type="button"
            class="flex-1 flex items-center justify-between p-4 rounded-lg border-2 transition-all"
            :class="
              language === 'en'
                ? 'border-cyan-500 bg-cyan-500/10'
                : 'border-slate-800 bg-slate-950 hover:border-slate-700'
            "
            @click="setLanguage('en')"
          >
            <span class="font-semibold" :class="language === 'en' ? 'text-cyan-400' : 'text-slate-300'">
              {{ t.lang_en }}
            </span>
            <Check v-if="language === 'en'" class="text-cyan-500 w-5 h-5" />
          </button>
        </div>
      </div>

      <div class="bg-slate-900 border border-slate-800 rounded-xl p-6">
        <h3 class="text-lg font-semibold text-slate-200 mb-4">{{ t.behavior_prefs }}</h3>
        <label class="flex items-center gap-3 text-slate-300 cursor-pointer select-none">
          <input v-model="soundEnabled" type="checkbox" class="rounded border-slate-600 bg-slate-950" />
          <span>{{ t.sound_alerts }}</span>
        </label>
        <p class="text-xs text-slate-500 mt-3">
          {{ t.model_strategy_label }}: <span class="text-slate-400 font-mono">{{ modelStrategy }}</span>
          （{{ t.model_strategy_auto_hint }}）
        </p>
      </div>

      <div class="bg-slate-900 border rounded-xl p-6 transition-colors duration-300 border-slate-700">
        <div class="flex items-center gap-3 mb-6">
          <div class="p-2 rounded-lg bg-cyan-500/20 text-cyan-400">
            <HardDrive class="w-5 h-5" />
          </div>
          <h3 class="text-lg font-semibold text-slate-200">{{ t.ember2018 }}</h3>
        </div>

        <div class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">{{ t.model_path }}</label>
            <div class="flex gap-2">
              <input
                v-model="e2018Path"
                type="text"
                class="flex-1 min-w-0 bg-slate-950 border border-slate-800 rounded-lg px-4 py-2.5 text-slate-300 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                placeholder="C:\path\to\ember2018_model.txt"
              />
              <button
                type="button"
                class="shrink-0 px-3 py-2 rounded-lg border border-slate-700 bg-slate-800 text-slate-300 hover:border-cyan-500/50 hover:text-cyan-400 flex items-center gap-2 text-sm"
                @click="browseModel('2018')"
              >
                <FolderOpen class="w-4 h-4" />
                {{ t.browse_model_file }}
              </button>
            </div>
          </div>

          <div>
            <div class="flex justify-between items-center mb-2">
              <label class="flex items-center gap-2 text-sm font-medium text-slate-400">
                <SlidersHorizontal class="w-4 h-4" />
                {{ t.malware_threshold }}
              </label>
              <span class="text-sm font-mono text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded">
                {{ e2018Threshold.toFixed(2) }}
              </span>
            </div>
            <input
              v-model.number="e2018Threshold"
              type="range"
              min="0"
              max="1"
              step="0.01"
              class="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
            />
            <p class="text-xs text-slate-500 mt-2">
              {{ t.threshold_help_2018 }}
            </p>
          </div>
        </div>
      </div>

      <div class="bg-slate-900 border rounded-xl p-6 transition-colors duration-300 border-slate-700">
        <div class="flex items-center gap-3 mb-6">
          <div class="p-2 rounded-lg bg-cyan-500/20 text-cyan-400">
            <HardDrive class="w-5 h-5" />
          </div>
          <h3 class="text-lg font-semibold text-slate-200">{{ t.ember2024 }}</h3>
        </div>

        <div class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">{{ t.model_path }}</label>
            <div class="flex gap-2">
              <input
                v-model="e2024Path"
                type="text"
                class="flex-1 min-w-0 bg-slate-950 border border-slate-800 rounded-lg px-4 py-2.5 text-slate-300 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                placeholder="C:\path\to\ember2024_model.txt"
              />
              <button
                type="button"
                class="shrink-0 px-3 py-2 rounded-lg border border-slate-700 bg-slate-800 text-slate-300 hover:border-cyan-500/50 hover:text-cyan-400 flex items-center gap-2 text-sm"
                @click="browseModel('2024')"
              >
                <FolderOpen class="w-4 h-4" />
                {{ t.browse_model_file }}
              </button>
            </div>
          </div>

          <div>
            <div class="flex justify-between items-center mb-2">
              <label class="flex items-center gap-2 text-sm font-medium text-slate-400">
                <SlidersHorizontal class="w-4 h-4" />
                {{ t.malware_threshold }}
              </label>
              <span class="text-sm font-mono text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded">
                {{ e2024Threshold.toFixed(2) }}
              </span>
            </div>
            <input
              v-model.number="e2024Threshold"
              type="range"
              min="0"
              max="1"
              step="0.01"
              class="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
            />
            <p class="text-xs text-slate-500 mt-2">
              {{ t.threshold_help_2024 }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <div
      class="fixed bottom-0 left-24 right-0 p-4 bg-slate-900/80 backdrop-blur-md border-t border-slate-800 flex justify-end gap-4"
    >
      <button
        type="button"
        class="flex items-center gap-2 bg-cyan-500 hover:bg-cyan-400 text-slate-950 px-6 py-2.5 rounded-lg font-medium transition-colors shadow-lg shadow-cyan-500/20 hover:scale-[1.02] active:scale-[0.98]"
        @click="handleSave"
      >
        <Save class="w-4 h-4" />
        {{ t.apply_settings }}
      </button>
    </div>
  </div>
</template>

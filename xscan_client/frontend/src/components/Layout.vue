<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import {
  Shield,
  Settings as SettingsIcon,
  Activity,
  ScanSearch,
  X,
  Minus,
  Maximize2,
  Copy,
} from 'lucide-vue-next'
import { Toaster } from 'vue-sonner'
import { useScanner } from '../composables/useScanner'
import { translations } from '../i18n'
import {
  WindowMinimise,
  WindowToggleMaximise,
  WindowIsMaximised,
  Quit,
} from '../../wailsjs/runtime/runtime'

const route = useRoute()
const { language } = useScanner()

const t = computed(() => translations[language.value])

const isMaximized = ref(false)

function navItemClass(targetPath: string) {
  const p = route.path
  const normalized = p.length > 1 && p.endsWith('/') ? p.slice(0, -1) : p
  const active =
    targetPath === '/'
      ? normalized === '/' || normalized === ''
      : normalized === targetPath
  return [
    'flex flex-col items-center justify-center gap-0.5 rounded-lg px-1 py-2 text-center text-[11px] font-medium leading-tight transition-colors min-h-[3.25rem]',
    active
      ? 'bg-cyan-500/10 text-cyan-400 border border-cyan-500/20'
      : 'text-slate-400 hover:bg-slate-800/60 hover:text-slate-200 border border-transparent',
  ]
}

function isWailsRuntime(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { runtime?: unknown }).runtime
}

onMounted(async () => {
  if (!isWailsRuntime()) return
  try {
    isMaximized.value = await WindowIsMaximised()
  } catch {
    /* dev / browser */
  }
})

function minimizeWin() {
  if (!isWailsRuntime()) return
  WindowMinimise()
}

async function toggleMaximize() {
  if (!isWailsRuntime()) return
  await WindowToggleMaximise()
  try {
    isMaximized.value = await WindowIsMaximised()
  } catch {
    isMaximized.value = !isMaximized.value
  }
}

function closeWin() {
  if (!isWailsRuntime()) return
  Quit()
}
</script>

<template>
  <div
    class="flex h-screen w-full flex-col bg-slate-950 text-slate-200 font-sans overflow-hidden selection:bg-cyan-500/30"
  >
    <!-- 无边框窗口：可拖拽标题栏（--wails-draggable） -->
    <header
      class="titlebar-drag flex h-10 shrink-0 items-center justify-between border-b border-slate-800 bg-slate-900 pl-3 pr-1 select-none"
      style="--wails-draggable: drag"
    >
      <div class="flex min-w-0 flex-1 items-center gap-2 pr-2">
        <Shield class="h-4 w-4 shrink-0 text-cyan-400" />
        <span class="truncate text-xs font-medium tracking-wide text-slate-300">
          {{ t.window_title }}
        </span>
      </div>
      <div
        class="flex shrink-0 items-center gap-0.5 [&_button]:h-8 [&_button]:w-10 [&_button]:rounded-md [&_button]:flex [&_button]:items-center [&_button]:justify-center [&_button]:text-slate-400 [&_button]:transition-colors hover:[&_button]:bg-slate-800 hover:[&_button]:text-slate-100"
        style="--wails-draggable: none"
      >
        <button type="button" title="最小化" @click="minimizeWin">
          <Minus class="h-4 w-4" />
        </button>
        <button type="button" :title="isMaximized ? '还原' : '最大化'" @click="toggleMaximize">
          <Copy v-if="isMaximized" class="h-3.5 w-3.5" />
          <Maximize2 v-else class="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          title="关闭"
          class="hover:!text-red-400"
          @click="closeWin"
        >
          <X class="h-4 w-4" />
        </button>
      </div>
    </header>

    <div class="flex min-h-0 flex-1 overflow-hidden">
      <aside
        class="z-10 flex w-24 shrink-0 flex-col border-r border-slate-800 bg-slate-900 overflow-hidden"
      >
        <nav class="flex-1 space-y-1 px-1.5 py-4">
          <RouterLink to="/" :class="navItemClass('/')">
            <ScanSearch class="w-4 h-4 shrink-0" />
            <span class="break-words text-[10px] leading-snug px-0.5">{{ t.scanner }}</span>
          </RouterLink>

          <RouterLink to="/history" :class="navItemClass('/history')">
            <Activity class="w-4 h-4 shrink-0" />
            <span class="break-words text-[10px] leading-snug px-0.5">{{ t.history }}</span>
          </RouterLink>
        </nav>

        <div class="border-t border-slate-800 px-1.5 py-4">
          <RouterLink to="/settings" :class="navItemClass('/settings')">
            <SettingsIcon class="w-4 h-4 shrink-0" />
            <span class="break-words text-[10px] leading-snug px-0.5">{{ t.settings }}</span>
          </RouterLink>
        </div>
      </aside>

      <main class="relative flex min-w-0 flex-1 flex-col overflow-hidden">
        <div class="flex-1 min-h-0 overflow-auto">
          <RouterView v-slot="{ Component }">
            <keep-alive include="ScannerPage">
              <component :is="Component" />
            </keep-alive>
          </RouterView>
        </div>
      </main>
    </div>

    <Toaster
      theme="dark"
      position="bottom-right"
      class="toaster-override"
      :toast-options="{
        classes: {
          toast: '!bg-slate-900 !border-slate-800 !text-slate-200',
        },
      }"
    />
  </div>
</template>

<style scoped>
:deep(.toaster-override [data-sonner-toast]) {
  background-color: rgb(15 23 42);
  border-color: rgb(30 41 59);
  color: rgb(226 232 240);
}
</style>

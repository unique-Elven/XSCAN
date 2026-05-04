<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterView } from 'vue-router'
import { toast } from 'vue-sonner'
import RiskHandleModal from './components/RiskHandleModal.vue'
import { provideScanner, hydrateScannerFromDisk } from './composables/useScanner'
import { translations } from './i18n'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { QuarantineFiles } from '../wailsjs/go/main/App'

const scanner = provideScanner()

const t = computed(() => translations[scanner.language.value])

const riskModalOpen = ref(false)
const riskModalPaths = ref<string[]>([])
const riskModalTitle = ref('')

let unsubscribeScanCompleted: (() => void) | undefined

function parseMaliciousPathsFromPayload(...args: unknown[]): string[] {
  const payload = args[0]
  if (payload == null) return []
  if (Array.isArray(payload)) {
    return payload.filter((x): x is string => typeof x === 'string')
  }
  if (typeof payload === 'object' && payload !== null && 'maliciousPaths' in payload) {
    const mp = (payload as { maliciousPaths?: unknown }).maliciousPaths
    if (Array.isArray(mp)) {
      return mp.filter((x): x is string => typeof x === 'string')
    }
  }
  if (typeof payload === 'string') {
    try {
      const o = JSON.parse(payload) as unknown
      return parseMaliciousPathsFromPayload(o)
    } catch {
      return []
    }
  }
  return []
}

function isWailsRuntime(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { runtime?: unknown }).runtime
}

onMounted(() => {
  hydrateScannerFromDisk(scanner)
  if (!isWailsRuntime()) return
  unsubscribeScanCompleted = EventsOn('scan_completed', (...data: unknown[]) => {
    const paths = parseMaliciousPathsFromPayload(...data)
    if (paths.length === 0) return
    riskModalPaths.value = paths
    const lang = scanner.language.value
    riskModalTitle.value = translations[lang].risk_modal_title.replace('{n}', String(paths.length))
    riskModalOpen.value = true
  })
})

onUnmounted(() => {
  unsubscribeScanCompleted?.()
})

function closeRiskModal() {
  riskModalOpen.value = false
}

async function onGlobalQuarantine(paths: string[]) {
  if (!paths.length) {
    closeRiskModal()
    return
  }
  try {
    const n = await QuarantineFiles(paths)
    closeRiskModal()
    const lang = scanner.language.value
    toast.success(translations[lang].quarantine_success.replace('{n}', String(n)))
  } catch (e) {
    toast.error(String(e))
  }
}
</script>

<template>
  <RouterView />
  <RiskHandleModal
    :open="riskModalOpen"
    :title="riskModalTitle"
    :paths="riskModalPaths"
    :ignore-label="t.risk_ignore_all"
    :confirm-label="t.risk_action_now"
    @close="closeRiskModal"
    @quarantine="onGlobalQuarantine"
  />
</template>

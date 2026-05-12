<script>
  import { onMount } from 'svelte'
  import { getSettings, updateSettings } from '../lib/api.js'

  let patientName = $state('')
  let defaultArm = $state('left')
  let defaultPosition = $state('sitting')
  let deviceName = $state('')

  let loading = $state(true)
  let saving = $state(false)
  let saved = $state(false)
  let error = $state('')

  onMount(async () => {
    try {
      const s = await getSettings()
      patientName = s.patient_name || ''
      defaultArm = s.default_arm || 'left'
      defaultPosition = s.default_position || 'sitting'
      deviceName = s.device_name || ''
    } catch {
      // No settings yet
    } finally {
      loading = false
    }
  })

  async function handleSave() {
    saving = true
    saved = false
    error = ''

    try {
      await updateSettings({
        patient_name: patientName,
        default_arm: defaultArm,
        default_position: defaultPosition,
        device_name: deviceName,
      })
      saved = true
      setTimeout(() => saved = false, 3000)
    } catch (err) {
      error = err.message || 'Failed to save settings'
    } finally {
      saving = false
    }
  }
</script>

{#if loading}
  <div class="flex justify-center py-12">
    <p class="text-gray-500 dark:text-gray-400">Loading...</p>
  </div>
{:else}
  <div class="space-y-6">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Settings</h1>

    <form onsubmit={e => { e.preventDefault(); handleSave() }} class="space-y-6 max-w-lg">
      <div>
        <label for="patient-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Patient Name</label>
        <input
          id="patient-name"
          type="text"
          bind:value={patientName}
          placeholder="Used in report headers"
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100 dark:placeholder-gray-400"
        />
        <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">Shown on printed reports</p>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label for="default-arm" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Default Arm</label>
          <select
            id="default-arm"
            bind:value={defaultArm}
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
          >
            <option value="left">Left</option>
            <option value="right">Right</option>
          </select>
        </div>
        <div>
          <label for="default-position" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Default Position</label>
          <select
            id="default-position"
            bind:value={defaultPosition}
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
          >
            <option value="sitting">Sitting</option>
            <option value="lying">Lying</option>
            <option value="standing">Standing</option>
          </select>
        </div>
      </div>

      <div>
        <label for="device-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Device Name</label>
        <input
          id="device-name"
          type="text"
          bind:value={deviceName}
          placeholder="e.g. Omron M3"
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100 dark:placeholder-gray-400"
        />
      </div>

      {#if error}
        <p class="text-sm text-red-600 dark:text-red-400">{error}</p>
      {/if}

      <div class="flex items-center gap-3">
        <button
          type="submit"
          disabled={saving}
          class="px-6 py-2 bg-teal-600 text-white rounded-md text-sm font-medium
            hover:bg-teal-700 transition-colors disabled:opacity-50"
        >
          {saving ? 'Saving...' : 'Save Settings'}
        </button>
        {#if saved}
          <span class="text-sm text-teal-600 dark:text-teal-400 font-medium">Settings saved</span>
        {/if}
      </div>
    </form>
  </div>
{/if}

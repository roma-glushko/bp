<script>
  import { onMount, tick } from 'svelte'
  import ReadingRow from './ReadingRow.svelte'
  import BPIndicator from './BPIndicator.svelte'
  import { classifyBP } from './bp.js'
  import { getSettings } from './api.js'

  let { initialData = null, onSubmit, submitLabel = 'Save Session' } = $props()

  function roundToFiveMin(d) {
    const m = Math.round(d.getMinutes() / 5) * 5
    return `${String(d.getHours()).padStart(2, '0')}:${String(m % 60).padStart(2, '0')}`
  }

  function todayStr() {
    const d = new Date()
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  }

  function emptyReading() {
    return { systolic: '', diastolic: '', pulse: '' }
  }

  function parseInitialData(data) {
    if (!data) return null
    const dt = new Date(data.measured_at)
    return {
      date: `${dt.getFullYear()}-${String(dt.getMonth() + 1).padStart(2, '0')}-${String(dt.getDate()).padStart(2, '0')}`,
      time: `${String(dt.getHours()).padStart(2, '0')}:${String(dt.getMinutes()).padStart(2, '0')}`,
      period: data.period || 'morning',
      arm: data.arm || 'left',
      position: data.position || 'sitting',
      readings: data.readings.map(r => ({
        systolic: r.systolic,
        diastolic: r.diastolic,
        pulse: r.pulse ?? '',
      })),
      notes: data.notes || '',
    }
  }

  // Capture initial prop value once for form defaults — intentionally non-reactive
  const init = initialData ? parseInitialData(initialData) : null

  let date = $state(init?.date ?? todayStr())
  let time = $state(init?.time ?? roundToFiveMin(new Date()))
  let period = $state(init?.period ?? 'morning')
  let arm = $state(init?.arm ?? 'left')
  let position = $state(init?.position ?? 'sitting')
  let readings = $state(init?.readings ?? [emptyReading()])
  let notes = $state(init?.notes ?? '')

  let submitting = $state(false)
  let errors = $state([])

  let sessionAverage = $derived.by(() => {
    const filled = readings.filter(r => r.systolic > 0 && r.diastolic > 0)
    if (filled.length === 0) return null

    const avgSys = Math.round(filled.reduce((s, r) => s + Number(r.systolic), 0) / filled.length)
    const avgDia = Math.round(filled.reduce((s, r) => s + Number(r.diastolic), 0) / filled.length)

    const pulsed = filled.filter(r => r.pulse !== '' && r.pulse > 0)
    const avgPulse = pulsed.length > 0
      ? Math.round(pulsed.reduce((s, r) => s + Number(r.pulse), 0) / pulsed.length)
      : null

    return {
      avgSys,
      avgDia,
      avgPulse,
      indicator: classifyBP(avgSys, avgDia),
    }
  })

  let readingsContainer

  async function addReading() {
    if (readings.length < 10) {
      readings = [...readings, emptyReading()]
      await tick()
      const inputs = readingsContainer?.querySelectorAll('input[type="number"]')
      if (inputs?.length) {
        inputs[inputs.length - 3]?.focus()
      }
    }
  }

  function removeReading(index) {
    readings = readings.filter((_, i) => i !== index)
  }

  async function handleSubmit() {
    errors = []
    submitting = true

    try {
      const measuredAt = new Date(`${date}T${time}:00`).toISOString()

      const sessionData = {
        measured_at: measuredAt,
        period,
        arm,
        position,
        readings: readings.map(r => ({
          systolic: Number(r.systolic),
          diastolic: Number(r.diastolic),
          pulse: r.pulse !== '' && r.pulse > 0 ? Number(r.pulse) : undefined,
        })),
        notes: notes || undefined,
      }

      await onSubmit(sessionData)
    } catch (err) {
      if (err.body?.errors) {
        errors = err.body.errors
      } else {
        errors = [{ field: '', message: err.message || 'Something went wrong' }]
      }
    } finally {
      submitting = false
    }
  }

  onMount(async () => {
    if (initialData) return
    try {
      const settings = await getSettings()
      if (settings.default_arm) arm = settings.default_arm
      if (settings.default_position) position = settings.default_position
    } catch {
      // Settings not configured yet — use defaults
    }
  })
</script>

<form onsubmit={e => { e.preventDefault(); handleSubmit() }} class="space-y-6">
  {#if errors.length > 0}
    <div class="bg-red-50 border border-red-200 rounded-md p-4">
      <p class="text-sm font-medium text-red-800 mb-1">Please fix the following errors:</p>
      <ul class="text-sm text-red-600 list-disc list-inside">
        {#each errors as err}
          <li>{err.field ? `${err.field}: ` : ''}{err.message}</li>
        {/each}
      </ul>
    </div>
  {/if}

  <div class="grid grid-cols-2 gap-4">
    <div>
      <label for="date" class="block text-sm font-medium text-gray-700 mb-1">Date</label>
      <input
        id="date"
        type="date"
        bind:value={date}
        class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
      />
    </div>
    <div>
      <label for="time" class="block text-sm font-medium text-gray-700 mb-1">Time</label>
      <input
        id="time"
        type="time"
        bind:value={time}
        class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
      />
    </div>
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
    <div>
      <label for="period" class="block text-sm font-medium text-gray-700 mb-1">Period</label>
      <select
        id="period"
        bind:value={period}
        class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
      >
        <option value="morning">Morning</option>
        <option value="evening">Evening</option>
        <option value="custom">Custom</option>
      </select>
    </div>
    <div>
      <label for="arm" class="block text-sm font-medium text-gray-700 mb-1">Arm</label>
      <select
        id="arm"
        bind:value={arm}
        class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
      >
        <option value="left">Left</option>
        <option value="right">Right</option>
      </select>
    </div>
    <div>
      <label for="position" class="block text-sm font-medium text-gray-700 mb-1">Position</label>
      <select
        id="position"
        bind:value={position}
        class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
      >
        <option value="sitting">Sitting</option>
        <option value="lying">Lying</option>
        <option value="standing">Standing</option>
      </select>
    </div>
  </div>

  <div>
    <div class="flex items-center justify-between mb-2">
      <span class="block text-sm font-medium text-gray-700">Readings</span>
      <div class="flex items-center gap-3">
        {#if sessionAverage}
          <span class="text-sm text-gray-500">
            Avg: {sessionAverage.avgSys}/{sessionAverage.avgDia}{sessionAverage.avgPulse != null ? `, pulse ${sessionAverage.avgPulse}` : ''}
          </span>
          <BPIndicator {...sessionAverage.indicator} />
        {/if}
      </div>
    </div>

    <div class="hidden sm:flex items-center gap-2 mb-2 text-xs text-gray-400 pl-8">
      <div class="flex-1 grid grid-cols-3 gap-2">
        <span>SYS</span>
        <span>DIA</span>
        <span>Pulse</span>
      </div>
      <div class="w-24"></div>
      <div class="w-8"></div>
    </div>

    <div class="space-y-2" bind:this={readingsContainer}>
      {#each readings as reading, i}
        <ReadingRow
          {reading}
          index={i}
          canRemove={readings.length > 1}
          onRemove={() => removeReading(i)}
        />
      {/each}
    </div>

    {#if readings.length < 10}
      <button
        type="button"
        onclick={addReading}
        class="mt-2 text-sm text-teal-600 hover:text-teal-800 font-medium"
      >
        + Add another reading
      </button>
    {/if}
  </div>

  <div>
    <label for="notes" class="block text-sm font-medium text-gray-700 mb-1">Notes</label>
    <textarea
      id="notes"
      bind:value={notes}
      rows="3"
      placeholder="Optional notes about this session..."
      class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 resize-none"
    ></textarea>
  </div>

  <button
    type="submit"
    disabled={submitting}
    class="w-full py-2.5 px-4 bg-teal-600 text-white rounded-md text-sm font-medium
      hover:bg-teal-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
  >
    {submitting ? 'Saving...' : submitLabel}
  </button>
</form>

<script>
  import BPIndicator from './BPIndicator.svelte'
  import { classifyBP, isValidSystolic, isValidDiastolic, isValidPulse } from './bp.js'

  let { reading, index, canRemove, onRemove } = $props()

  let indicator = $derived(
    reading.systolic > 0 && reading.diastolic > 0
      ? classifyBP(reading.systolic, reading.diastolic)
      : null
  )

  let sysInvalid = $derived(reading.systolic !== '' && reading.systolic > 0 && !isValidSystolic(reading.systolic))
  let diaInvalid = $derived(reading.diastolic !== '' && reading.diastolic > 0 && !isValidDiastolic(reading.diastolic))
  let pulseInvalid = $derived(reading.pulse !== '' && reading.pulse > 0 && !isValidPulse(reading.pulse))
  let sysLessDia = $derived(
    reading.systolic > 0 && reading.diastolic > 0 && reading.systolic <= reading.diastolic
  )
</script>

<div class="flex items-center gap-2">
  <span class="text-sm text-gray-500 w-6 text-right">{index + 1}</span>

  <div class="flex-1 grid grid-cols-3 gap-2">
    <div>
      <input
        type="number"
        bind:value={reading.systolic}
        placeholder="SYS"
        min="70"
        max="250"
        class="w-full px-3 py-2 border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500
          {sysInvalid || sysLessDia ? 'border-red-400 bg-red-50' : 'border-gray-300'}"
      />
    </div>
    <div>
      <input
        type="number"
        bind:value={reading.diastolic}
        placeholder="DIA"
        min="40"
        max="150"
        class="w-full px-3 py-2 border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500
          {diaInvalid || sysLessDia ? 'border-red-400 bg-red-50' : 'border-gray-300'}"
      />
    </div>
    <div>
      <input
        type="number"
        bind:value={reading.pulse}
        placeholder="Pulse"
        min="30"
        max="220"
        class="w-full px-3 py-2 border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500
          {pulseInvalid ? 'border-red-400 bg-red-50' : 'border-gray-300'}"
      />
    </div>
  </div>

  <div class="hidden sm:flex w-24 items-center">
    {#if indicator}
      <BPIndicator {...indicator} />
    {/if}
  </div>

  <div class="w-8">
    {#if canRemove}
      <button
        type="button"
        onclick={onRemove}
        class="text-gray-400 hover:text-red-500 transition-colors p-1"
        title="Remove reading"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    {/if}
  </div>
</div>

{#if sysInvalid}
  <p class="text-xs text-red-500 ml-8 mt-0.5">Systolic must be between 70 and 250</p>
{/if}
{#if diaInvalid}
  <p class="text-xs text-red-500 ml-8 mt-0.5">Diastolic must be between 40 and 150</p>
{/if}
{#if pulseInvalid}
  <p class="text-xs text-red-500 ml-8 mt-0.5">Pulse must be between 30 and 220</p>
{/if}
{#if sysLessDia}
  <p class="text-xs text-red-500 ml-8 mt-0.5">Systolic must be greater than diastolic</p>
{/if}

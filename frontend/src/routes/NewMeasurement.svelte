<script>
  import MeasurementForm from '../lib/MeasurementForm.svelte'
  import BPIndicator from '../lib/BPIndicator.svelte'
  import { createSession } from '../lib/api.js'

  let result = $state(null)

  async function handleSubmit(data) {
    result = await createSession(data)
  }

  function reset() {
    result = null
  }
</script>

{#if result}
  <div class="max-w-md mx-auto text-center space-y-6">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Saved.</h1>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-6 space-y-3">
      <p class="text-sm text-gray-500 uppercase tracking-wide">Session Average</p>
      <p class="text-3xl font-semibold text-gray-900">
        {result.average.avg_systolic} / {result.average.avg_diastolic}
        {#if result.average.avg_pulse}
          <span class="text-xl text-gray-500">, pulse {result.average.avg_pulse}</span>
        {/if}
      </p>
      <BPIndicator {...result.average.indicator} />
      <p class="text-sm text-gray-500 mt-2">
        This session has {result.readings.length} reading{result.readings.length !== 1 ? 's' : ''}.
      </p>
    </div>

    <div class="flex gap-3 justify-center">
      <button
        onclick={reset}
        class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-colors"
      >
        Add Another
      </button>
      <a
        href="#/history"
        class="px-4 py-2 border border-gray-300 text-gray-700 rounded-md text-sm font-medium hover:bg-gray-50 transition-colors"
      >
        View History
      </a>
    </div>
  </div>
{:else}
  <div>
    <h1 class="text-2xl font-bold text-gray-900 mb-6">New Measurement</h1>
    <MeasurementForm onSubmit={handleSubmit} />
  </div>
{/if}

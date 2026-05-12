<script>
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import MeasurementForm from '../lib/MeasurementForm.svelte'
  import { getSession, updateSession } from '../lib/api.js'

  let { params = {} } = $props()

  let session = $state(null)
  let loading = $state(true)
  let notFound = $state(false)

  onMount(async () => {
    try {
      session = await getSession(params.id)
    } catch (err) {
      if (err.status === 404) {
        notFound = true
      }
    } finally {
      loading = false
    }
  })

  async function handleSubmit(data) {
    await updateSession(params.id, data)
    push('/history')
  }
</script>

{#if loading}
  <div class="flex justify-center py-12">
    <p class="text-gray-500 dark:text-gray-400">Loading...</p>
  </div>
{:else if notFound}
  <div class="text-center py-12">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">Session Not Found</h1>
    <p class="text-gray-500 dark:text-gray-400 mb-4">This measurement session doesn't exist.</p>
    <a
      href="#/history"
      class="text-teal-600 dark:text-teal-400 hover:text-teal-800 dark:hover:text-teal-300 font-medium text-sm"
    >
      Back to History
    </a>
  </div>
{:else}
  <div>
    <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">Edit Measurement</h1>
    <MeasurementForm
      initialData={session}
      onSubmit={handleSubmit}
      submitLabel="Update Session"
    />
  </div>
{/if}

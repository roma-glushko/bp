<script>
  import { onMount } from 'svelte'
  import BPIndicator from '../lib/BPIndicator.svelte'
  import { listSessions, deleteSession } from '../lib/api.js'
  import { classifyBP } from '../lib/bp.js'

  let loading = $state(true)
  let sessions = $state([])
  let expandedId = $state(null)
  let deleteConfirmId = $state(null)
  let deleting = $state(false)

  let fromDate = $state('')
  let toDate = $state('')

  function defaultRange() {
    const now = new Date()
    const from = new Date(now)
    from.setDate(from.getDate() - 30)
    return {
      from: from.toISOString().slice(0, 10),
      to: now.toISOString().slice(0, 10),
    }
  }

  const defaults = defaultRange()
  fromDate = defaults.from
  toDate = defaults.to

  function formatDate(iso) {
    return new Date(iso).toLocaleDateString([], { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
  }

  function formatTime(iso) {
    return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  function dayKey(iso) {
    return new Date(iso).toISOString().slice(0, 10)
  }

  function groupByDay(sessions) {
    const groups = new Map()
    for (const s of sessions) {
      const key = dayKey(s.measured_at)
      if (!groups.has(key)) groups.set(key, [])
      groups.get(key).push(s)
    }

    return Array.from(groups.entries()).map(([date, daySessions]) => {
      const avgSys = Math.round(daySessions.reduce((sum, s) => sum + Number(s.average.avg_systolic), 0) / daySessions.length)
      const avgDia = Math.round(daySessions.reduce((sum, s) => sum + Number(s.average.avg_diastolic), 0) / daySessions.length)
      return {
        date,
        label: formatDate(date + 'T00:00:00'),
        sessions: daySessions,
        avgSys,
        avgDia,
        indicator: classifyBP(avgSys, avgDia),
      }
    })
  }

  let dayGroups = $derived(groupByDay(sessions))

  async function loadSessions() {
    loading = true
    try {
      const from = fromDate ? new Date(fromDate + 'T00:00:00').toISOString() : undefined
      const to = toDate ? new Date(toDate + 'T23:59:59').toISOString() : undefined
      const resp = await listSessions(from, to)
      sessions = resp.sessions || []
    } catch {
      sessions = []
    } finally {
      loading = false
    }
  }

  function toggleExpand(id) {
    expandedId = expandedId === id ? null : id
  }

  async function confirmDelete(id) {
    deleting = true
    try {
      await deleteSession(id)
      sessions = sessions.filter(s => s.id !== id)
      deleteConfirmId = null
    } catch {
      // Failed to delete
    } finally {
      deleting = false
    }
  }

  onMount(loadSessions)
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">History</h1>
    <a
      href="#/measurements/new"
      class="px-4 py-2 bg-teal-600 text-white rounded-md text-sm font-medium hover:bg-teal-700 transition-colors"
    >
      Add Reading
    </a>
  </div>

  <div class="flex items-end gap-3">
    <div>
      <label for="from" class="block text-xs text-gray-500 dark:text-gray-400 mb-1">From</label>
      <input
        id="from"
        type="date"
        bind:value={fromDate}
        class="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
      />
    </div>
    <div>
      <label for="to" class="block text-xs text-gray-500 dark:text-gray-400 mb-1">To</label>
      <input
        id="to"
        type="date"
        bind:value={toDate}
        class="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
      />
    </div>
    <button
      onclick={loadSessions}
      class="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
    >
      Filter
    </button>
  </div>

  {#if loading}
    <div class="flex justify-center py-12">
      <p class="text-gray-500 dark:text-gray-400">Loading...</p>
    </div>
  {:else if sessions.length === 0}
    <div class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">No sessions found in this date range.</p>
    </div>
  {:else}
    <div class="space-y-6">
      {#each dayGroups as day}
        <div>
          <div class="flex items-center justify-between mb-2 pb-1 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-300">{day.label}</h2>
            <div class="flex items-center gap-2">
              <span class="text-sm text-gray-500 dark:text-gray-400">{day.avgSys}/{day.avgDia}</span>
              <BPIndicator {...day.indicator} />
            </div>
          </div>

          <div class="space-y-2">
            {#each day.sessions as session}
              <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                <button
                  onclick={() => toggleExpand(session.id)}
                  class="w-full px-4 py-3 flex items-center justify-between text-left hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                >
                  <div class="flex items-center gap-2 sm:gap-4 min-w-0">
                    <span class="text-sm text-gray-500 dark:text-gray-400 shrink-0">{formatTime(session.measured_at)}</span>
                    <span class="text-xs text-gray-400 dark:text-gray-500 capitalize hidden sm:inline">{session.period}</span>
                    <span class="text-base font-semibold text-gray-900 dark:text-gray-100 shrink-0">
                      {session.average.avg_systolic}/{session.average.avg_diastolic}
                    </span>
                    {#if session.average.avg_pulse}
                      <span class="text-sm text-gray-500 dark:text-gray-400 hidden sm:inline">pulse {session.average.avg_pulse}</span>
                    {/if}
                  </div>
                  <div class="flex items-center gap-3">
                    <BPIndicator {...session.average.indicator} />
                    <span class="text-xs text-gray-400 dark:text-gray-500">{session.readings.length}r</span>
                    <svg
                      class="w-4 h-4 text-gray-400 transition-transform {expandedId === session.id ? 'rotate-180' : ''}"
                      fill="none" stroke="currentColor" viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                    </svg>
                  </div>
                </button>

                {#if expandedId === session.id}
                  <div class="px-4 pb-4 border-t border-gray-100 dark:border-gray-700">
                    <table class="w-full mt-3 text-sm">
                      <thead>
                        <tr class="text-xs text-gray-400 dark:text-gray-500 uppercase">
                          <th class="text-left pb-1 font-medium">#</th>
                          <th class="text-left pb-1 font-medium">SYS</th>
                          <th class="text-left pb-1 font-medium">DIA</th>
                          <th class="text-left pb-1 font-medium">Pulse</th>
                          <th class="text-left pb-1 font-medium"></th>
                        </tr>
                      </thead>
                      <tbody>
                        {#each session.readings as reading, i}
                          {@const ind = classifyBP(reading.systolic, reading.diastolic)}
                          <tr class="text-gray-700 dark:text-gray-300">
                            <td class="py-0.5 text-gray-400 dark:text-gray-500">{i + 1}</td>
                            <td class="py-0.5">{reading.systolic}</td>
                            <td class="py-0.5">{reading.diastolic}</td>
                            <td class="py-0.5">{reading.pulse ?? '—'}</td>
                            <td class="py-0.5"><BPIndicator {...ind} /></td>
                          </tr>
                        {/each}
                      </tbody>
                    </table>

                    {#if session.notes}
                      <p class="mt-2 text-sm text-gray-500 dark:text-gray-400 italic">{session.notes}</p>
                    {/if}

                    {#if session.arm || session.position}
                      <p class="mt-1 text-xs text-gray-400 dark:text-gray-500">
                        {#if session.arm}<span class="capitalize">{session.arm} arm</span>{/if}
                        {#if session.arm && session.position} · {/if}
                        {#if session.position}<span class="capitalize">{session.position}</span>{/if}
                      </p>
                    {/if}

                    <div class="mt-3 flex gap-2">
                      <a
                        href="#/measurements/{session.id}/edit"
                        class="px-3 py-1 border border-gray-300 dark:border-gray-600 rounded text-xs font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                      >
                        Edit
                      </a>
                      {#if deleteConfirmId === session.id}
                        <span class="text-xs text-red-600 dark:text-red-400 flex items-center gap-2">
                          Delete this session?
                          <button
                            onclick={() => confirmDelete(session.id)}
                            disabled={deleting}
                            class="px-2 py-1 bg-red-600 text-white rounded text-xs font-medium hover:bg-red-700 disabled:opacity-50"
                          >
                            {deleting ? 'Deleting...' : 'Yes, delete'}
                          </button>
                          <button
                            onclick={() => deleteConfirmId = null}
                            class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-xs text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
                          >
                            Cancel
                          </button>
                        </span>
                      {:else}
                        <button
                          onclick={() => deleteConfirmId = session.id}
                          class="px-3 py-1 border border-red-200 dark:border-red-800 rounded text-xs font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                        >
                          Delete
                        </button>
                      {/if}
                    </div>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

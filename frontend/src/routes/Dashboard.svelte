<script>
  import { onMount } from 'svelte'
  import BPIndicator from '../lib/BPIndicator.svelte'
  import BPTrendChart from '../lib/BPTrendChart.svelte'
  import { listSessions } from '../lib/api.js'
  import { classifyBP } from '../lib/bp.js'

  let loading = $state(true)
  let todaySessions = $state([])
  let lastSession = $state(null)
  let weekAvg = $state(null)
  let monthAvg = $state(null)
  let allSessions = $state([])

  function startOfDay(d) {
    return new Date(d.getFullYear(), d.getMonth(), d.getDate())
  }

  function formatTime(iso) {
    const d = new Date(iso)
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  function formatDate(iso) {
    const d = new Date(iso)
    return d.toLocaleDateString([], { month: 'short', day: 'numeric' })
  }

  function computeAvg(sessions) {
    if (sessions.length === 0) return null
    const avgSys = Math.round(sessions.reduce((s, x) => s + Number(x.average.avg_systolic), 0) / sessions.length)
    const avgDia = Math.round(sessions.reduce((s, x) => s + Number(x.average.avg_diastolic), 0) / sessions.length)
    const withPulse = sessions.filter(x => x.average.avg_pulse)
    const avgPulse = withPulse.length > 0
      ? Math.round(withPulse.reduce((s, x) => s + Number(x.average.avg_pulse), 0) / withPulse.length)
      : null
    return { avgSys, avgDia, avgPulse, count: sessions.length, indicator: classifyBP(avgSys, avgDia) }
  }

  onMount(async () => {
    try {
      const now = new Date()
      const monthAgo = new Date(now)
      monthAgo.setDate(monthAgo.getDate() - 30)

      const resp = await listSessions(monthAgo.toISOString(), now.toISOString())
      allSessions = resp.sessions || []

      if (allSessions.length > 0) {
        lastSession = allSessions[0]
      }

      const todayStart = startOfDay(now)
      todaySessions = allSessions.filter(s => new Date(s.measured_at) >= todayStart)

      const weekAgo = new Date(now)
      weekAgo.setDate(weekAgo.getDate() - 7)
      const weekSessions = allSessions.filter(s => new Date(s.measured_at) >= weekAgo)

      weekAvg = computeAvg(weekSessions)
      monthAvg = computeAvg(allSessions)
    } catch {
      // API not available — show empty state
    } finally {
      loading = false
    }
  })
</script>

{#if loading}
  <div class="flex justify-center py-12">
    <p class="text-gray-500 dark:text-gray-400">Loading...</p>
  </div>
{:else if !lastSession}
  <div class="text-center py-12">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">Welcome to BP Journal</h1>
    <p class="text-gray-600 dark:text-gray-400 mb-6">Start tracking your blood pressure by adding your first reading.</p>
    <a
      href="#/measurements/new"
      class="inline-block px-6 py-3 bg-teal-600 text-white rounded-md text-sm font-medium hover:bg-teal-700 transition-colors"
    >
      Add First Reading
    </a>
  </div>
{:else}
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Dashboard</h1>
      <a
        href="#/measurements/new"
        class="px-4 py-2 bg-teal-600 text-white rounded-md text-sm font-medium hover:bg-teal-700 transition-colors"
      >
        Add Reading
      </a>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Last Reading</p>
        <p class="text-2xl font-semibold text-gray-900 dark:text-gray-100">
          {lastSession.average.avg_systolic} / {lastSession.average.avg_diastolic}
        </p>
        {#if lastSession.average.avg_pulse}
          <p class="text-sm text-gray-500 dark:text-gray-400">Pulse {lastSession.average.avg_pulse}</p>
        {/if}
        <div class="mt-1">
          <BPIndicator {...lastSession.average.indicator} />
        </div>
        <p class="text-xs text-gray-400 dark:text-gray-500 mt-2">
          {formatDate(lastSession.measured_at)} at {formatTime(lastSession.measured_at)}
        </p>
      </div>

      {#if todaySessions.length > 0}
        {@const todaySys = Math.round(todaySessions.reduce((s, t) => s + Number(t.average.avg_systolic), 0) / todaySessions.length)}
        {@const todayDia = Math.round(todaySessions.reduce((s, t) => s + Number(t.average.avg_diastolic), 0) / todaySessions.length)}
        {@const todayIndicator = classifyBP(todaySys, todayDia)}
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Today's Average</p>
          <p class="text-2xl font-semibold text-gray-900 dark:text-gray-100">{todaySys} / {todayDia}</p>
          <div class="mt-1">
            <BPIndicator {...todayIndicator} />
          </div>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-2">{todaySessions.length} session{todaySessions.length !== 1 ? 's' : ''}</p>
        </div>
      {/if}

      {#if weekAvg}
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">7-Day Average</p>
          <p class="text-2xl font-semibold text-gray-900 dark:text-gray-100">{weekAvg.avgSys} / {weekAvg.avgDia}</p>
          {#if weekAvg.avgPulse}
            <p class="text-sm text-gray-500 dark:text-gray-400">Pulse {weekAvg.avgPulse}</p>
          {/if}
          <div class="mt-1">
            <BPIndicator {...weekAvg.indicator} />
          </div>
        </div>
      {/if}

      {#if monthAvg}
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">30-Day Average</p>
          <p class="text-2xl font-semibold text-gray-900 dark:text-gray-100">{monthAvg.avgSys} / {monthAvg.avgDia}</p>
          {#if monthAvg.avgPulse}
            <p class="text-sm text-gray-500 dark:text-gray-400">Pulse {monthAvg.avgPulse}</p>
          {/if}
          <div class="mt-1">
            <BPIndicator {...monthAvg.indicator} />
          </div>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-2">{monthAvg.count} session{monthAvg.count !== 1 ? 's' : ''}</p>
        </div>
      {/if}
    </div>

    {#if allSessions.length > 1}
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">30-Day Trend</h2>
        <BPTrendChart sessions={allSessions} />
      </div>
    {/if}

    {#if todaySessions.length > 0}
      <div>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-3">Today's Sessions</h2>
        <div class="space-y-2">
          {#each todaySessions as session}
            <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-gray-100">
                    {formatTime(session.measured_at)}
                    <span class="text-gray-500 dark:text-gray-400 capitalize">· {session.period}</span>
                  </p>
                  <p class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    {session.average.avg_systolic} / {session.average.avg_diastolic}
                    {#if session.average.avg_pulse}
                      <span class="text-sm text-gray-500 dark:text-gray-400 font-normal">pulse {session.average.avg_pulse}</span>
                    {/if}
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <BPIndicator {...session.average.indicator} />
                <span class="text-xs text-gray-400 dark:text-gray-500">{session.readings.length} reading{session.readings.length !== 1 ? 's' : ''}</span>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
{/if}

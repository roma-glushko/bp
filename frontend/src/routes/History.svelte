<script>
  import { onMount } from 'svelte'
  import BPIndicator from '../lib/BPIndicator.svelte'
  import CalendarView from '../lib/CalendarView.svelte'
  import { listSessions, deleteSession, listAnnotations, upsertDailyNote, upsertWeeklyNote } from '../lib/api.js'
  import { classifyBP } from '../lib/bp.js'

  let loading = $state(true)
  let sessions = $state([])
  let expandedId = $state(null)
  let deleteConfirmId = $state(null)
  let deleting = $state(false)

  let fromDate = $state('')
  let toDate = $state('')

  let viewMode = $state('list')
  let dailyNotes = $state({})
  let weeklyNotes = $state({})
  let editingDayNote = $state(null)
  let editingWeekNote = $state(null)
  let editNoteText = $state('')

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

  function isoWeekKey(dateStr) {
    const d = new Date(dateStr + 'T00:00:00')
    const dayOfWeek = d.getDay() || 7
    d.setDate(d.getDate() + 4 - dayOfWeek)
    const yearStart = new Date(d.getFullYear(), 0, 1)
    const weekNo = Math.ceil((((d - yearStart) / 86400000) + 1) / 7)
    return `${d.getFullYear()}-W${String(weekNo).padStart(2, '0')}`
  }

  function groupByWeekAndDay(sessions) {
    const dayMap = new Map()
    for (const s of sessions) {
      const key = dayKey(s.measured_at)
      if (!dayMap.has(key)) dayMap.set(key, [])
      dayMap.get(key).push(s)
    }

    const weekMap = new Map()
    for (const [date, daySessions] of dayMap) {
      const week = isoWeekKey(date)
      if (!weekMap.has(week)) weekMap.set(week, [])

      const avgSys = Math.round(daySessions.reduce((sum, s) => sum + Number(s.average.avg_systolic), 0) / daySessions.length)
      const avgDia = Math.round(daySessions.reduce((sum, s) => sum + Number(s.average.avg_diastolic), 0) / daySessions.length)
      weekMap.get(week).push({
        date,
        label: formatDate(date + 'T00:00:00'),
        sessions: daySessions,
        avgSys,
        avgDia,
        indicator: classifyBP(avgSys, avgDia),
      })
    }

    return Array.from(weekMap.entries()).map(([week, days]) => ({
      week,
      days: days.sort((a, b) => b.date.localeCompare(a.date)),
    }))
  }

  let weekGroups = $derived(groupByWeekAndDay(sessions))

  let calendarDayData = $derived.by(() => {
    const map = {}
    for (const s of sessions) {
      const key = dayKey(s.measured_at)
      if (!map[key]) map[key] = { sysList: [], diaList: [], sessions: 0 }
      map[key].sysList.push(Number(s.average.avg_systolic))
      map[key].diaList.push(Number(s.average.avg_diastolic))
      map[key].sessions++
    }
    const result = {}
    for (const [key, d] of Object.entries(map)) {
      const avgSys = Math.round(d.sysList.reduce((a, b) => a + b, 0) / d.sysList.length)
      const avgDia = Math.round(d.diaList.reduce((a, b) => a + b, 0) / d.diaList.length)
      result[key] = { avgSys, avgDia, sessions: d.sessions }
    }
    return result
  })

  async function loadSessions() {
    loading = true
    try {
      const from = fromDate ? new Date(fromDate + 'T00:00:00').toISOString() : undefined
      const to = toDate ? new Date(toDate + 'T23:59:59').toISOString() : undefined

      const [sessResp, annResp] = await Promise.all([
        listSessions(from, to),
        listAnnotations(from, to),
      ])

      sessions = sessResp.sessions || []

      dailyNotes = {}
      if (annResp.daily_notes) {
        for (const dn of annResp.daily_notes) {
          dailyNotes[dn.date] = dn.notes
        }
      }

      weeklyNotes = {}
      if (annResp.weekly_notes) {
        for (const wn of annResp.weekly_notes) {
          weeklyNotes[wn.week] = wn.notes
        }
      }
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

  function startEditDayNote(date) {
    editingDayNote = date
    editingWeekNote = null
    editNoteText = dailyNotes[date] || ''
  }

  function startEditWeekNote(week) {
    editingWeekNote = week
    editingDayNote = null
    editNoteText = weeklyNotes[week] || ''
  }

  function cancelEditNote() {
    editingDayNote = null
    editingWeekNote = null
    editNoteText = ''
  }

  async function saveDayNote(date) {
    try {
      await upsertDailyNote(date, editNoteText.trim())
      if (editNoteText.trim()) {
        dailyNotes[date] = editNoteText.trim()
      } else {
        delete dailyNotes[date]
      }
      dailyNotes = dailyNotes
    } catch {
      // Failed to save
    }
    editingDayNote = null
    editNoteText = ''
  }

  async function saveWeekNote(week) {
    try {
      await upsertWeeklyNote(week, editNoteText.trim())
      if (editNoteText.trim()) {
        weeklyNotes[week] = editNoteText.trim()
      } else {
        delete weeklyNotes[week]
      }
      weeklyNotes = weeklyNotes
    } catch {
      // Failed to save
    }
    editingWeekNote = null
    editNoteText = ''
  }

  onMount(loadSessions)
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">History</h1>
    <div class="flex items-center gap-2">
      <div class="flex rounded-md border border-gray-300 dark:border-gray-600 overflow-hidden">
        <button
          onclick={() => viewMode = 'list'}
          class="px-2.5 py-1.5 text-sm transition-colors {viewMode === 'list' ? 'bg-teal-600 text-white' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700'}"
          title="List view"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <button
          onclick={() => viewMode = 'calendar'}
          class="px-2.5 py-1.5 text-sm transition-colors {viewMode === 'calendar' ? 'bg-teal-600 text-white' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700'}"
          title="Calendar view"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </button>
      </div>
      <a
        href="#/measurements/new"
        class="px-4 py-2 bg-teal-600 text-white rounded-md text-sm font-medium hover:bg-teal-700 transition-colors"
      >
        Add Reading
      </a>
    </div>
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
  {:else if viewMode === 'calendar'}
    <CalendarView
      dayData={calendarDayData}
      {dailyNotes}
      {weeklyNotes}
      onEditDayNote={startEditDayNote}
      onEditWeekNote={startEditWeekNote}
    />
    {#if editingDayNote}
      <div class="mt-3 p-3 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Note for {editingDayNote}</p>
        <textarea
          bind:value={editNoteText}
          rows="2"
          class="w-full px-3 py-1.5 border border-teal-300 dark:border-teal-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
          placeholder="Day note..."
        ></textarea>
        <div class="flex gap-2 mt-1">
          <button
            onclick={() => saveDayNote(editingDayNote)}
            class="px-2 py-1 bg-teal-600 text-white rounded text-xs font-medium hover:bg-teal-700"
          >
            Save
          </button>
          <button
            onclick={cancelEditNote}
            class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-xs text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
          >
            Cancel
          </button>
        </div>
      </div>
    {/if}
    {#if editingWeekNote}
      <div class="mt-3 p-3 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Note for {editingWeekNote}</p>
        <textarea
          bind:value={editNoteText}
          rows="2"
          class="w-full px-3 py-1.5 border border-teal-300 dark:border-teal-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
          placeholder="Week note..."
        ></textarea>
        <div class="flex gap-2 mt-1">
          <button
            onclick={() => saveWeekNote(editingWeekNote)}
            class="px-2 py-1 bg-teal-600 text-white rounded text-xs font-medium hover:bg-teal-700"
          >
            Save
          </button>
          <button
            onclick={cancelEditNote}
            class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-xs text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
          >
            Cancel
          </button>
        </div>
      </div>
    {/if}
  {:else}
    <div class="space-y-8">
      {#each weekGroups as weekGroup}
        <div>
          <!-- Week header -->
          <div class="mb-3">
            <div class="flex items-center gap-2 mb-1">
              <h2 class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase tracking-wide">{weekGroup.week}</h2>
              {#if !weeklyNotes[weekGroup.week] && editingWeekNote !== weekGroup.week}
                <button
                  onclick={() => startEditWeekNote(weekGroup.week)}
                  class="text-xs text-gray-400 dark:text-gray-500 hover:text-teal-600 dark:hover:text-teal-400 transition-colors"
                >
                  + note
                </button>
              {/if}
            </div>
            {#if editingWeekNote === weekGroup.week}
              <div class="mt-1">
                <textarea
                  bind:value={editNoteText}
                  rows="2"
                  class="w-full px-3 py-1.5 border border-teal-300 dark:border-teal-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
                  placeholder="Week note..."
                ></textarea>
                <div class="flex gap-2 mt-1">
                  <button
                    onclick={() => saveWeekNote(weekGroup.week)}
                    class="px-2 py-1 bg-teal-600 text-white rounded text-xs font-medium hover:bg-teal-700"
                  >
                    Save
                  </button>
                  <button
                    onclick={cancelEditNote}
                    class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-xs text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            {:else if weeklyNotes[weekGroup.week]}
              <button
                onclick={() => startEditWeekNote(weekGroup.week)}
                class="text-sm text-gray-500 dark:text-gray-400 italic hover:text-gray-700 dark:hover:text-gray-300 transition-colors text-left"
              >
                {weeklyNotes[weekGroup.week]}
              </button>
            {/if}
          </div>

          <div class="space-y-6">
            {#each weekGroup.days as day}
              <div>
                <div class="flex items-center justify-between mb-1 pb-1 border-b border-gray-200 dark:border-gray-700">
                  <div class="flex items-center gap-2">
                    <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300">{day.label}</h3>
                    {#if !dailyNotes[day.date] && editingDayNote !== day.date}
                      <button
                        onclick={() => startEditDayNote(day.date)}
                        class="text-xs text-gray-400 dark:text-gray-500 hover:text-teal-600 dark:hover:text-teal-400 transition-colors"
                      >
                        + note
                      </button>
                    {/if}
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-sm text-gray-500 dark:text-gray-400">{day.avgSys}/{day.avgDia}</span>
                    <BPIndicator {...day.indicator} />
                  </div>
                </div>

                {#if editingDayNote === day.date}
                  <div class="mb-2">
                    <textarea
                      bind:value={editNoteText}
                      rows="2"
                      class="w-full px-3 py-1.5 border border-teal-300 dark:border-teal-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 dark:bg-gray-700 dark:text-gray-100"
                      placeholder="Day note..."
                    ></textarea>
                    <div class="flex gap-2 mt-1">
                      <button
                        onclick={() => saveDayNote(day.date)}
                        class="px-2 py-1 bg-teal-600 text-white rounded text-xs font-medium hover:bg-teal-700"
                      >
                        Save
                      </button>
                      <button
                        onclick={cancelEditNote}
                        class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-xs text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                {:else if dailyNotes[day.date]}
                  <button
                    onclick={() => startEditDayNote(day.date)}
                    class="mb-2 text-sm text-gray-500 dark:text-gray-400 italic hover:text-gray-700 dark:hover:text-gray-300 transition-colors text-left block"
                  >
                    {dailyNotes[day.date]}
                  </button>
                {/if}

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
        </div>
      {/each}
    </div>
  {/if}
</div>

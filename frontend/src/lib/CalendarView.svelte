<script>
  import { classifyBP } from './bp.js'

  let { dayData = {}, dailyNotes = {}, weeklyNotes = {}, onEditDayNote, onEditWeekNote } = $props()

  let currentMonth = $state(new Date())

  const WEEKDAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

  function toDateKey(d) {
    return d.toISOString().slice(0, 10)
  }

  function isoWeekKey(dateStr) {
    const d = new Date(dateStr + 'T00:00:00')
    const dayOfWeek = d.getDay() || 7
    d.setDate(d.getDate() + 4 - dayOfWeek)
    const yearStart = new Date(d.getFullYear(), 0, 1)
    const weekNo = Math.ceil((((d - yearStart) / 86400000) + 1) / 7)
    return `${d.getFullYear()}-W${String(weekNo).padStart(2, '0')}`
  }

  function getMonthGrid(year, month) {
    const firstDay = new Date(year, month, 1)
    const lastDay = new Date(year, month + 1, 0)

    let startDow = firstDay.getDay() || 7
    const startDate = new Date(firstDay)
    startDate.setDate(startDate.getDate() - (startDow - 1))

    const weeks = []
    const cursor = new Date(startDate)

    while (cursor <= lastDay || weeks.length < 6) {
      const week = []
      for (let i = 0; i < 7; i++) {
        week.push(new Date(cursor))
        cursor.setDate(cursor.getDate() + 1)
      }
      weeks.push(week)
      if (cursor > lastDay && weeks.length >= 4) break
    }

    return weeks
  }

  function prevMonth() {
    const d = new Date(currentMonth)
    d.setMonth(d.getMonth() - 1)
    currentMonth = d
  }

  function nextMonth() {
    const d = new Date(currentMonth)
    d.setMonth(d.getMonth() + 1)
    currentMonth = d
  }

  function goToToday() {
    currentMonth = new Date()
  }

  let monthLabel = $derived(
    currentMonth.toLocaleDateString([], { month: 'long', year: 'numeric' })
  )

  let weeks = $derived(
    getMonthGrid(currentMonth.getFullYear(), currentMonth.getMonth())
  )

  let isCurrentMonth = $derived(
    currentMonth.getFullYear() === new Date().getFullYear() &&
    currentMonth.getMonth() === new Date().getMonth()
  )

  let visibleWeeks = $derived.by(() => {
    return [...new Set(weeks.flat().filter(d => d.getMonth() === currentMonth.getMonth()).map(d => isoWeekKey(toDateKey(d))))]
  })

  const colorClasses = {
    green: 'bg-green-500',
    yellow: 'bg-yellow-500',
    red: 'bg-red-500',
    darkred: 'bg-red-800',
  }
</script>

<div>
  <!-- Month navigation -->
  <div class="flex items-center justify-between mb-4">
    <button
      onclick={prevMonth}
      aria-label="Previous month"
      class="p-1.5 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
    >
      <svg class="w-5 h-5 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>
    <div class="flex items-center gap-2">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{monthLabel}</h2>
      {#if !isCurrentMonth}
        <button
          onclick={goToToday}
          class="text-xs px-2 py-0.5 rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
        >
          Today
        </button>
      {/if}
    </div>
    <button
      onclick={nextMonth}
      aria-label="Next month"
      class="p-1.5 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
    >
      <svg class="w-5 h-5 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>
  </div>

  <!-- Weekday headers -->
  <div class="grid grid-cols-7 mb-1">
    {#each WEEKDAYS as wd}
      <div class="text-center text-xs font-medium text-gray-500 dark:text-gray-400 py-1">{wd}</div>
    {/each}
  </div>

  <!-- Calendar grid -->
  <div class="grid grid-cols-7 border-t border-l border-gray-200 dark:border-gray-700">
    {#each weeks as week}
      {#each week as day}
        {@const key = toDateKey(day)}
        {@const inMonth = day.getMonth() === currentMonth.getMonth()}
        {@const isToday = key === toDateKey(new Date())}
        {@const data = dayData[key]}
        {@const note = dailyNotes[key]}
        {@const weekKey = isoWeekKey(key)}
        <div
          class="min-h-[5rem] p-2 border-r border-b border-gray-200 dark:border-gray-700 {inMonth ? '' : 'bg-gray-50 dark:bg-gray-800/50'}"
        >
          <div class="flex items-start justify-between">
            <span
              class="text-xs font-medium leading-tight {isToday ? 'bg-teal-600 text-white w-5 h-5 rounded-full flex items-center justify-center' : inMonth ? 'text-gray-700 dark:text-gray-300' : 'text-gray-400 dark:text-gray-600'}"
            >
              {day.getDate()}
            </span>
            {#if note}
              <button
                onclick={() => onEditDayNote?.(key)}
                class="text-teal-500 dark:text-teal-400 hover:text-teal-700 dark:hover:text-teal-300"
                title={note}
              >
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z"/>
                </svg>
              </button>
            {/if}
          </div>
          {#if data}
            {@const indicator = classifyBP(data.avgSys, data.avgDia)}
            <div class="mt-0.5">
              <div class="flex items-center gap-1">
                <span class="inline-block w-2 h-2 rounded-full {colorClasses[indicator.color] || 'bg-gray-400'}"></span>
                <span class="text-xs font-semibold text-gray-800 dark:text-gray-200">{data.avgSys}/{data.avgDia}</span>
              </div>
              {#if data.sessions > 1}
                <p class="text-[10px] text-gray-400 dark:text-gray-500 mt-0.5">{data.sessions}s</p>
              {/if}
            </div>
          {/if}
        </div>
      {/each}
    {/each}
  </div>

  <!-- Week notes legend -->
  {#if visibleWeeks.some(w => weeklyNotes[w])}
    <div class="mt-3 space-y-1">
      {#each visibleWeeks as wk}
        {#if weeklyNotes[wk]}
          <button
            onclick={() => onEditWeekNote?.(wk)}
            class="block text-xs text-gray-500 dark:text-gray-400 italic hover:text-gray-700 dark:hover:text-gray-300 text-left"
          >
            <span class="font-medium text-gray-600 dark:text-gray-400 not-italic">{wk}:</span> {weeklyNotes[wk]}
          </button>
        {/if}
      {/each}
    </div>
  {/if}
</div>

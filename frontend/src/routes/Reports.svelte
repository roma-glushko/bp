<script>
  import BPIndicator from '../lib/BPIndicator.svelte'
  import { getReportPreview } from '../lib/api.js'
  import { classifyBP } from '../lib/bp.js'

  function defaultFrom() {
    const d = new Date()
    d.setMonth(d.getMonth() - 1)
    return d.toISOString().slice(0, 10)
  }

  let fromDate = $state(defaultFrom())
  let toDate = $state(new Date().toISOString().slice(0, 10))

  let sections = $state({
    summary: true,
    monthly: true,
    weekly: true,
    daily: true,
    sessions: true,
    readings: true,
    notes: true,
  })

  let loading = $state(false)
  let report = $state(null)
  let error = $state('')

  async function generateReport() {
    loading = true
    error = ''
    report = null

    try {
      const selected = Object.entries(sections).filter(([, v]) => v).map(([k]) => k)
      const from = new Date(fromDate + 'T00:00:00').toISOString()
      const to = new Date(toDate + 'T23:59:59').toISOString()
      report = await getReportPreview(from, to, selected)
    } catch (err) {
      error = err.message || 'Failed to generate report'
    } finally {
      loading = false
    }
  }

  function formatPulse(avg) {
    if (!avg) return '—'
    return String(Math.round(Number(avg)))
  }

  function formatAvg(v) {
    return String(Math.round(Number(v)))
  }

  function printReport() {
    window.print()
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between print:hidden">
    <h1 class="text-2xl font-bold text-gray-900">Reports</h1>
    {#if report}
      <button
        onclick={printReport}
        class="px-4 py-2 bg-teal-600 text-white rounded-md text-sm font-medium hover:bg-teal-700 transition-colors"
      >
        Print / Save as PDF
      </button>
    {/if}
  </div>

  <div class="bg-white rounded-lg border border-gray-200 p-4 print:hidden">
    <div class="grid grid-cols-2 gap-4 mb-4">
      <div>
        <label for="report-from" class="block text-xs text-gray-500 mb-1">From</label>
        <input
          id="report-from"
          type="date"
          bind:value={fromDate}
          class="w-full px-3 py-1.5 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
        />
      </div>
      <div>
        <label for="report-to" class="block text-xs text-gray-500 mb-1">To</label>
        <input
          id="report-to"
          type="date"
          bind:value={toDate}
          class="w-full px-3 py-1.5 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
        />
      </div>
    </div>

    <fieldset>
      <legend class="text-sm font-medium text-gray-700 mb-2">Include sections</legend>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
        {#each [
          ['summary', 'Summary'],
          ['monthly', 'Monthly Averages'],
          ['weekly', 'Weekly Averages'],
          ['daily', 'Daily Averages'],
          ['sessions', 'Session Averages'],
          ['readings', 'Raw Measurements'],
          ['notes', 'Notes'],
        ] as [key, label]}
          <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
            <input type="checkbox" bind:checked={sections[key]} class="rounded border-gray-300" />
            {label}
          </label>
        {/each}
      </div>
    </fieldset>

    <button
      onclick={generateReport}
      disabled={loading}
      class="mt-4 w-full py-2 px-4 bg-teal-600 text-white rounded-md text-sm font-medium
        hover:bg-teal-700 transition-colors disabled:opacity-50"
    >
      {loading ? 'Generating...' : 'Preview Report'}
    </button>

    {#if error}
      <p class="mt-2 text-sm text-red-600">{error}</p>
    {/if}
  </div>

  {#if report}
    <div class="report-preview bg-white rounded-lg border border-gray-200 p-6 print:border-0 print:p-0 print:rounded-none">
      <div class="text-center mb-8">
        <h2 class="text-xl font-bold text-gray-900">Blood Pressure Report</h2>
        {#if report.patient}
          <p class="text-gray-700 mt-1">Patient: {report.patient}</p>
        {/if}
        {#if report.device}
          <p class="text-sm text-gray-500 mt-1">Measured By: {report.device}</p>
        {/if}
        <p class="text-sm text-gray-500 mt-1">Period: {report.from} — {report.to}</p>
        <p class="text-sm text-gray-500">Generated: {report.generated}</p>
      </div>

      {#if report.summary}
        <section class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-3 border-b border-gray-200 pb-1">1. Summary</h3>
          <div class="grid grid-cols-3 gap-4 mb-4">
            <div>
              <p class="text-xs text-gray-500 uppercase">Days Measured</p>
              <p class="text-xl font-semibold">{report.summary.days_measured}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase">Sessions</p>
              <p class="text-xl font-semibold">{report.summary.session_count}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase">Readings</p>
              <p class="text-xl font-semibold">{report.summary.reading_count}</p>
            </div>
          </div>

          <table class="w-full text-sm mb-4">
            <thead>
              <tr class="text-xs text-gray-400 uppercase border-b">
                <th class="text-left py-1 font-medium"></th>
                <th class="text-left py-1 font-medium">Avg BP</th>
                <th class="text-left py-1 font-medium">Pulse</th>
                <th class="text-left py-1 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {#if report.summary.overall_avg.count > 0}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5 font-medium">Overall</td>
                  <td>{formatAvg(report.summary.overall_avg.avg_systolic)}/{formatAvg(report.summary.overall_avg.avg_diastolic)}</td>
                  <td>{formatPulse(report.summary.overall_avg.avg_pulse)}</td>
                  <td><BPIndicator {...report.summary.overall_avg.indicator} /><span class="print-indicator hidden print:inline"> ({report.summary.overall_avg.indicator.label})</span></td>
                </tr>
              {/if}
              {#if report.summary.morning_avg.count > 0}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5 font-medium">Morning</td>
                  <td>{formatAvg(report.summary.morning_avg.avg_systolic)}/{formatAvg(report.summary.morning_avg.avg_diastolic)}</td>
                  <td>{formatPulse(report.summary.morning_avg.avg_pulse)}</td>
                  <td><BPIndicator {...report.summary.morning_avg.indicator} /><span class="print-indicator hidden print:inline"> ({report.summary.morning_avg.indicator.label})</span></td>
                </tr>
              {/if}
              {#if report.summary.evening_avg.count > 0}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5 font-medium">Evening</td>
                  <td>{formatAvg(report.summary.evening_avg.avg_systolic)}/{formatAvg(report.summary.evening_avg.avg_diastolic)}</td>
                  <td>{formatPulse(report.summary.evening_avg.avg_pulse)}</td>
                  <td><BPIndicator {...report.summary.evening_avg.indicator} /><span class="print-indicator hidden print:inline"> ({report.summary.evening_avg.indicator.label})</span></td>
                </tr>
              {/if}
            </tbody>
          </table>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p class="text-xs text-gray-500 uppercase mb-1">Blood Pressure Range</p>
              <p>Systolic: {report.summary.lowest_systolic} — {report.summary.highest_systolic}</p>
              <p>Diastolic: {report.summary.lowest_diastolic} — {report.summary.highest_diastolic}</p>
            </div>
            {#if report.summary.highest_pulse}
              <div>
                <p class="text-xs text-gray-500 uppercase mb-1">Pulse Range</p>
                <p>{report.summary.lowest_pulse} — {report.summary.highest_pulse}</p>
              </div>
            {/if}
          </div>
        </section>
      {/if}

      {#if report.month_averages?.length > 0}
        <section class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-3 border-b border-gray-200 pb-1">2. Monthly Averages</h3>
          <table class="w-full text-sm">
            <thead>
              <tr class="text-xs text-gray-400 uppercase border-b">
                <th class="text-left py-1 font-medium">Month</th>
                <th class="text-left py-1 font-medium">Avg BP</th>
                <th class="text-left py-1 font-medium">Pulse</th>
                <th class="text-left py-1 font-medium">Days</th>
                <th class="text-left py-1 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {#each report.month_averages as m}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5">{m.month}</td>
                  <td>{formatAvg(m.average.avg_systolic)}/{formatAvg(m.average.avg_diastolic)}</td>
                  <td>{formatPulse(m.average.avg_pulse)}</td>
                  <td>{m.days}</td>
                  <td><BPIndicator {...m.average.indicator} /><span class="hidden print:inline"> ({m.average.indicator.label})</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </section>
      {/if}

      {#if report.week_averages?.length > 0}
        <section class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-3 border-b border-gray-200 pb-1">3. Weekly Averages</h3>
          <table class="w-full text-sm">
            <thead>
              <tr class="text-xs text-gray-400 uppercase border-b">
                <th class="text-left py-1 font-medium">Week</th>
                <th class="text-left py-1 font-medium">Avg BP</th>
                <th class="text-left py-1 font-medium">Pulse</th>
                <th class="text-left py-1 font-medium">Days</th>
                <th class="text-left py-1 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {#each report.week_averages as w}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5">{w.week}</td>
                  <td>{formatAvg(w.average.avg_systolic)}/{formatAvg(w.average.avg_diastolic)}</td>
                  <td>{formatPulse(w.average.avg_pulse)}</td>
                  <td>{w.days}</td>
                  <td><BPIndicator {...w.average.indicator} /><span class="hidden print:inline"> ({w.average.indicator.label})</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </section>
      {/if}

      {#if report.day_averages?.length > 0}
        <section class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-3 border-b border-gray-200 pb-1">4. Daily Averages</h3>
          <table class="w-full text-sm">
            <thead>
              <tr class="text-xs text-gray-400 uppercase border-b">
                <th class="text-left py-1 font-medium">Date</th>
                <th class="text-left py-1 font-medium">Avg BP</th>
                <th class="text-left py-1 font-medium">Pulse</th>
                <th class="text-left py-1 font-medium">Sessions</th>
                <th class="text-left py-1 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {#each report.day_averages as d}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5">{d.date}</td>
                  <td>{formatAvg(d.average.avg_systolic)}/{formatAvg(d.average.avg_diastolic)}</td>
                  <td>{formatPulse(d.average.avg_pulse)}</td>
                  <td>{d.sessions}</td>
                  <td><BPIndicator {...d.average.indicator} /><span class="hidden print:inline"> ({d.average.indicator.label})</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </section>
      {/if}

      {#if report.sessions?.length > 0}
        <section class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-3 border-b border-gray-200 pb-1">5. Session Averages</h3>
          <table class="w-full text-sm">
            <thead>
              <tr class="text-xs text-gray-400 uppercase border-b">
                <th class="text-left py-1 font-medium">Date</th>
                <th class="text-left py-1 font-medium">Time</th>
                <th class="text-left py-1 font-medium">Period</th>
                <th class="text-left py-1 font-medium">Avg BP</th>
                <th class="text-left py-1 font-medium">Pulse</th>
                <th class="text-left py-1 font-medium">Readings</th>
                <th class="text-left py-1 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {#each report.sessions as s}
                {@const dt = new Date(s.measured_at)}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5">{dt.toISOString().slice(0, 10)}</td>
                  <td>{dt.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</td>
                  <td class="capitalize">{s.period}</td>
                  <td>{formatAvg(s.average.avg_systolic)}/{formatAvg(s.average.avg_diastolic)}</td>
                  <td>{formatPulse(s.average.avg_pulse)}</td>
                  <td>{s.readings.length}</td>
                  <td><BPIndicator {...s.average.indicator} /><span class="hidden print:inline"> ({s.average.indicator.label})</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </section>
      {/if}

      {#if report.readings?.length > 0}
        <section class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-3 border-b border-gray-200 pb-1">6. Raw Measurements</h3>
          <table class="w-full text-sm">
            <thead>
              <tr class="text-xs text-gray-400 uppercase border-b">
                <th class="text-left py-1 font-medium">Date</th>
                <th class="text-left py-1 font-medium">Time</th>
                <th class="text-left py-1 font-medium">Period</th>
                <th class="text-left py-1 font-medium">#</th>
                <th class="text-left py-1 font-medium">SYS</th>
                <th class="text-left py-1 font-medium">DIA</th>
                <th class="text-left py-1 font-medium">Pulse</th>
                <th class="text-left py-1 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {#each report.readings as r}
                {@const ind = classifyBP(r.systolic, r.diastolic)}
                <tr class="border-b border-gray-100">
                  <td class="py-1.5">{r.date}</td>
                  <td>{r.time}</td>
                  <td class="capitalize">{r.period}</td>
                  <td>{r.attempt}</td>
                  <td>{r.systolic}</td>
                  <td>{r.diastolic}</td>
                  <td>{r.pulse ?? '—'}</td>
                  <td><BPIndicator {...ind} /><span class="hidden print:inline"> ({ind.label})</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </section>
      {/if}

      <footer class="text-xs text-gray-400 text-center border-t border-gray-200 pt-4 mt-8">
        This report summarizes self-measured blood pressure records and is not a medical diagnosis.
      </footer>
    </div>
  {/if}
</div>

<style>
  @media print {
    :global(nav) { display: none !important; }
    :global(main) { padding: 0 !important; max-width: none !important; }
    :global(body) { background: white !important; }
  }
</style>

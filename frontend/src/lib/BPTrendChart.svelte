<script>
  import { onMount, onDestroy } from 'svelte'
  import { Chart, LineController, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, Filler } from 'chart.js'

  Chart.register(LineController, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, Filler)

  let { sessions = [] } = $props()

  let canvas
  let chart

  function dayKey(iso) {
    return new Date(iso).toISOString().slice(0, 10)
  }

  function buildDailyPoints(sessions) {
    const map = new Map()
    for (const s of sessions) {
      const key = dayKey(s.measured_at)
      if (!map.has(key)) map.set(key, { sysList: [], diaList: [], pulseList: [] })
      const d = map.get(key)
      d.sysList.push(Number(s.average.avg_systolic))
      d.diaList.push(Number(s.average.avg_diastolic))
      if (s.average.avg_pulse) d.pulseList.push(Number(s.average.avg_pulse))
    }

    return Array.from(map.entries())
      .map(([date, d]) => ({
        date,
        label: new Date(date + 'T00:00:00').toLocaleDateString([], { month: 'short', day: 'numeric' }),
        sys: Math.round(d.sysList.reduce((a, b) => a + b, 0) / d.sysList.length),
        dia: Math.round(d.diaList.reduce((a, b) => a + b, 0) / d.diaList.length),
        pulse: d.pulseList.length > 0 ? Math.round(d.pulseList.reduce((a, b) => a + b, 0) / d.pulseList.length) : null,
      }))
      .sort((a, b) => a.date.localeCompare(b.date))
  }

  function isDark() {
    return document.documentElement.classList.contains('dark')
  }

  function createChart() {
    const points = buildDailyPoints(sessions)
    if (points.length === 0) return

    const hasPulse = points.some(p => p.pulse !== null)
    const dark = isDark()
    const gridColor = dark ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)'
    const textColor = dark ? '#9ca3af' : '#6b7280'

    chart = new Chart(canvas, {
      type: 'line',
      data: {
        labels: points.map(p => p.label),
        datasets: [
          {
            label: 'Systolic',
            data: points.map(p => p.sys),
            borderColor: '#0d9488',
            backgroundColor: 'rgba(13, 148, 136, 0.08)',
            borderWidth: 2,
            pointRadius: 4,
            pointBackgroundColor: '#0d9488',
            pointHoverRadius: 6,
            tension: 0.3,
            fill: true,
          },
          {
            label: 'Diastolic',
            data: points.map(p => p.dia),
            borderColor: 'rgba(13, 148, 136, 0.5)',
            backgroundColor: 'rgba(13, 148, 136, 0.03)',
            borderWidth: 2,
            borderDash: [6, 3],
            pointRadius: 3,
            pointBackgroundColor: 'rgba(13, 148, 136, 0.5)',
            pointHoverRadius: 5,
            tension: 0.3,
            fill: true,
          },
          ...(hasPulse ? [{
            label: 'Pulse',
            data: points.map(p => p.pulse),
            borderColor: '#e879797a',
            backgroundColor: 'transparent',
            borderWidth: 1.5,
            pointRadius: 2.5,
            pointBackgroundColor: '#e87979',
            pointHoverRadius: 4,
            tension: 0.3,
            fill: false,
            yAxisID: 'pulse',
          }] : []),
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          mode: 'index',
          intersect: false,
        },
        plugins: {
          legend: {
            position: 'bottom',
            labels: {
              color: textColor,
              usePointStyle: true,
              pointStyleWidth: 20,
              padding: 16,
              font: { size: 12 },
            },
          },
          tooltip: {
            backgroundColor: dark ? '#374151' : '#111827',
            titleFont: { size: 12 },
            bodyFont: { size: 12 },
            padding: 8,
            cornerRadius: 6,
            callbacks: {
              label: (ctx) => {
                const unit = ctx.dataset.yAxisID === 'pulse' ? 'bpm' : 'mmHg'
                return `${ctx.dataset.label}: ${ctx.parsed.y} ${unit}`
              },
            },
          },
        },
        scales: {
          x: {
            grid: { color: gridColor },
            ticks: {
              color: textColor,
              font: { size: 11 },
              maxRotation: 0,
              autoSkip: true,
              maxTicksLimit: 8,
            },
          },
          y: {
            grid: { color: gridColor },
            ticks: {
              color: textColor,
              font: { size: 11 },
              stepSize: 20,
            },
            suggestedMin: 60,
            suggestedMax: 160,
          },
          ...(hasPulse ? {
            pulse: {
              position: 'right',
              grid: { drawOnChartArea: false },
              ticks: {
                color: '#e8797980',
                font: { size: 11 },
                stepSize: 10,
              },
              suggestedMin: 40,
              suggestedMax: 120,
              title: {
                display: true,
                text: 'bpm',
                color: '#e8797980',
                font: { size: 10 },
              },
            },
          } : {}),
        },
      },
    })
  }

  onMount(() => {
    createChart()
  })

  onDestroy(() => {
    chart?.destroy()
  })
</script>

<div class="h-56">
  <canvas bind:this={canvas}></canvas>
</div>

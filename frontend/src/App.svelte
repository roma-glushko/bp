<script>
  import { onMount } from 'svelte'
  import Router from 'svelte-spa-router'
  import Nav from './lib/Nav.svelte'
  import { initTheme } from './lib/theme.js'

  import Dashboard from './routes/Dashboard.svelte'
  import NewMeasurement from './routes/NewMeasurement.svelte'
  import EditMeasurement from './routes/EditMeasurement.svelte'
  import History from './routes/History.svelte'
  import Reports from './routes/Reports.svelte'
  import Settings from './routes/Settings.svelte'

  const routes = {
    '/': Dashboard,
    '/measurements/new': NewMeasurement,
    '/measurements/:id/edit': EditMeasurement,
    '/history': History,
    '/reports': Reports,
    '/settings': Settings,
  }

  onMount(() => initTheme())
</script>

<div class="min-h-screen bg-gray-50 dark:bg-gray-950 transition-colors">
  <Nav />
  <main class="max-w-4xl mx-auto px-4 py-8">
    <Router {routes} />
  </main>
</div>

<script>
  import Router from 'svelte-spa-router'
  import Nav from './lib/Nav.svelte'

  import Dashboard from './routes/Dashboard.svelte'
  import NewMeasurement from './routes/NewMeasurement.svelte'
  import History from './routes/History.svelte'
  import Reports from './routes/Reports.svelte'
  import Settings from './routes/Settings.svelte'

  const routes = {
    '/': Dashboard,
    '/measurements/new': NewMeasurement,
    '/history': History,
    '/reports': Reports,
    '/settings': Settings,
  }
</script>

<div class="min-h-screen bg-gray-50">
  <Nav />
  <main class="max-w-4xl mx-auto px-4 py-8">
    <Router {routes} />
  </main>
</div>

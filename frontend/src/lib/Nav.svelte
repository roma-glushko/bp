<script>
  import { location } from 'svelte-spa-router'
  import { getTheme, toggleTheme } from './theme.js'

  const links = [
    { href: '/', label: 'Dashboard' },
    { href: '/measurements/new', label: 'New Reading' },
    { href: '/history', label: 'History' },
    { href: '/reports', label: 'Reports' },
    { href: '/settings', label: 'Settings' },
  ]

  let menuOpen = $state(false)
  let dark = $state(false)

  $effect(() => {
    dark = getTheme() === 'dark'
  })

  function handleToggle() {
    const next = toggleTheme()
    dark = next === 'dark'
  }
</script>

<nav class="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 print:hidden transition-colors">
  <div class="max-w-4xl mx-auto px-4">
    <div class="flex items-center justify-between h-14">
      <a href="#/" class="text-lg font-semibold text-gray-900 dark:text-gray-100">🩺 BP Journal</a>

      <div class="hidden sm:flex gap-1">
        {#each links as link}
          <a
            href="#{link.href}"
            class="px-3 py-2 rounded-md text-sm font-medium transition-colors
              {$location === link.href
                ? 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800'}"
          >
            {link.label}
          </a>
        {/each}
      </div>

      <div class="flex items-center gap-1">
        <button
          onclick={handleToggle}
          class="p-2 text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
          aria-label="Toggle theme"
        >
          {#if dark}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          {:else}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          {/if}
        </button>

        <button
          onclick={() => menuOpen = !menuOpen}
          class="sm:hidden p-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100"
          aria-label="Toggle menu"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            {#if menuOpen}
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            {:else}
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            {/if}
          </svg>
        </button>
      </div>
    </div>
  </div>

  {#if menuOpen}
    <div class="sm:hidden border-t border-gray-100 dark:border-gray-800 px-4 py-2 space-y-1">
      {#each links as link}
        <a
          href="#{link.href}"
          onclick={() => menuOpen = false}
          class="block px-3 py-2 rounded-md text-sm font-medium transition-colors
            {$location === link.href
              ? 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800'}"
        >
          {link.label}
        </a>
      {/each}
    </div>
  {/if}
</nav>

<script>
  import { location } from 'svelte-spa-router'

  const links = [
    { href: '/', label: 'Dashboard' },
    { href: '/measurements/new', label: 'New Reading' },
    { href: '/history', label: 'History' },
    { href: '/reports', label: 'Reports' },
    { href: '/settings', label: 'Settings' },
  ]

  let menuOpen = $state(false)
</script>

<nav class="bg-white border-b border-gray-200 print:hidden">
  <div class="max-w-4xl mx-auto px-4">
    <div class="flex items-center justify-between h-14">
      <a href="#/" class="text-lg font-semibold text-gray-900">🩺 BP Journal</a>

      <button
        onclick={() => menuOpen = !menuOpen}
        class="sm:hidden p-2 text-gray-600 hover:text-gray-900"
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

      <div class="hidden sm:flex gap-1">
        {#each links as link}
          <a
            href="#{link.href}"
            class="px-3 py-2 rounded-md text-sm font-medium transition-colors
              {$location === link.href
                ? 'bg-gray-100 text-gray-900'
                : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'}"
          >
            {link.label}
          </a>
        {/each}
      </div>
    </div>
  </div>

  {#if menuOpen}
    <div class="sm:hidden border-t border-gray-100 px-4 py-2 space-y-1">
      {#each links as link}
        <a
          href="#{link.href}"
          onclick={() => menuOpen = false}
          class="block px-3 py-2 rounded-md text-sm font-medium transition-colors
            {$location === link.href
              ? 'bg-gray-100 text-gray-900'
              : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'}"
        >
          {link.label}
        </a>
      {/each}
    </div>
  {/if}
</nav>

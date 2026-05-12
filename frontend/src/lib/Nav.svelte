<script>
  import { location } from 'svelte-spa-router'

  const links = [
    { href: '/', label: 'Dashboard' },
    { href: '/measurements/new', label: 'New Reading' },
    { href: '/history', label: 'History' },
    { href: '/reports', label: 'Reports' },
    { href: '/settings', label: 'Settings' },
  ]
</script>

<nav class="bg-white border-b border-gray-200">
  <div class="max-w-4xl mx-auto px-4">
    <div class="flex items-center justify-between h-14">
      <a href="#/" class="text-lg font-semibold text-gray-900">BP Journal</a>
      <div class="flex gap-1">
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
</nav>

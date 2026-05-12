<script>
  let { status, color, label } = $props()

  const colorClasses = {
    green: 'bg-green-500',
    yellow: 'bg-yellow-500',
    red: 'bg-red-500',
    darkred: 'bg-red-800',
  }
</script>

<span class="inline-flex items-center gap-1.5">
  <span class="inline-block w-2.5 h-2.5 rounded-full {colorClasses[color] || 'bg-gray-400'}"></span>
  <span class="text-sm font-medium">{label}</span>
</span>

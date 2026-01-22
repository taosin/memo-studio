<script>
  import { cn } from '$lib/utils.js';

  export let className = '';
</script>

<h3 class={cn('text-2xl font-semibold leading-none tracking-tight', className)}>
  <slot />
</h3>

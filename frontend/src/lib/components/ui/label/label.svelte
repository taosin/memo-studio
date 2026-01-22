<script>
  import { cn } from '$lib/utils.js';

  export let className = '';
</script>

<label
  class={cn(
    'text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70',
    className
  )}
>
  <slot />
</label>

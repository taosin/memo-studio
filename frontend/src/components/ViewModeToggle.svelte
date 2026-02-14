<script>
  import { createEventDispatcher } from 'svelte';
  import Button from '$lib/components/ui/button/button.svelte';

  export let mode = 'waterfall';
  const dispatch = createEventDispatcher();

  function setMode(newMode) {
    mode = newMode;
    dispatch('change', newMode);
  }
</script>

<div class="flex gap-1 bg-secondary/50 rounded-lg p-1">
  <button
    class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium transition-all duration-200"
    class:bg-primary={mode === 'waterfall'}
    class:text-primary-foreground={mode === 'waterfall'}
    class:text-muted-foreground={mode !== 'waterfall'}
    class:hover:bg-accent={mode !== 'waterfall'}
    on:click={() => setMode('waterfall')}
  >
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <rect x="3" y="3" width="7" height="7"></rect>
      <rect x="14" y="3" width="7" height="7"></rect>
      <rect x="14" y="14" width="7" height="7"></rect>
      <rect x="3" y="14" width="7" height="7"></rect>
    </svg>
    <span class="hidden sm:inline">瀑布流</span>
  </button>
  <button
    class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium transition-all duration-200"
    class:bg-primary={mode === 'timeline'}
    class:text-primary-foreground={mode === 'timeline'}
    class:text-muted-foreground={mode !== 'timeline'}
    class:hover:bg-accent={mode !== 'timeline'}
    on:click={() => setMode('timeline')}
  >
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <line x1="8" y1="6" x2="21" y2="6"></line>
      <line x1="8" y1="12" x2="21" y2="12"></line>
      <line x1="8" y1="18" x2="21" y2="18"></line>
      <line x1="3" y1="6" x2="3.01" y2="6"></line>
      <line x1="3" y1="12" x2="3.01" y2="12"></line>
      <line x1="3" y1="18" x2="3.01" y2="18"></line>
    </svg>
    <span class="hidden sm:inline">时间线</span>
  </button>
</div>

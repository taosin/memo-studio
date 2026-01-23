<script>
  import { createEventDispatcher } from 'svelte';
  import Button from '$lib/components/ui/button/button.svelte';

  export let mode = 'waterfall'; // 'waterfall' | 'timeline'
  const dispatch = createEventDispatcher();

  function setMode(newMode) {
    mode = newMode;
    dispatch('change', newMode);
  }
</script>

<div class="flex gap-2">
  <Button
    variant={mode === 'waterfall' ? 'default' : 'outline'}
    size="sm"
    on:click={() => setMode('waterfall')}
  >
    瀑布流
  </Button>
  <Button
    variant={mode === 'timeline' ? 'default' : 'outline'}
    size="sm"
    on:click={() => setMode('timeline')}
  >
    时间线
  </Button>
</div>

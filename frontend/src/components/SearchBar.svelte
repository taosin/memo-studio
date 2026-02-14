<script>
  import { createEventDispatcher } from 'svelte';
  import Input from '$lib/components/ui/input/input.svelte';

  export let value = '';
  const dispatch = createEventDispatcher();

  let inputValue = value;
  let isFocused = false;

  $: {
    if (inputValue !== value) {
      inputValue = value;
    }
  }

  function handleInput() {
    dispatch('search', inputValue);
  }

  function handleClear() {
    inputValue = '';
    dispatch('search', '');
  }
</script>

<div class="w-full relative">
  <div class="relative {isFocused ? 'animate-scale-in' : ''}">
    <div class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"></circle>
        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
      </svg>
    </div>
    <input
      type="text"
      placeholder="搜索笔记..."
      bind:value={inputValue}
      on:input={handleInput}
      on:focus={() => isFocused = true}
      on:blur={() => isFocused = false}
      class="w-full h-10 pl-10 pr-10 rounded-lg border border-border bg-background text-sm transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
    />
    {#if inputValue}
      <button
        class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
        on:click={handleClear}
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    {/if}
  </div>
</div>

<script>
  import { createEventDispatcher } from 'svelte';
  import Input from '$lib/components/ui/input/input.svelte';

  export let value = '';
  const dispatch = createEventDispatcher();

  let inputValue = value;

  $: {
    if (inputValue !== value) {
      inputValue = value;
    }
  }

  function handleInput() {
    dispatch('search', inputValue);
  }
</script>

<div class="w-full">
  <Input
    placeholder="搜索笔记..."
    bind:value={inputValue}
    on:input={handleInput}
    className="w-full"
  />
</div>

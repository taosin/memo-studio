<script>
  import { cn } from '$lib/utils.js';
  import { createEventDispatcher } from 'svelte';

  export let className = '';
  export let value = '';
  export let placeholder = '';
  export let disabled = false;
  export let type = 'text';
  export let id = '';
  export let autocomplete = '';

  const dispatch = createEventDispatcher();

  function handleInput(event) {
    value = event.target.value;
    dispatch('input', event);
  }

  // Svelte 不允许在 bind:value 时使用动态 type
  // 所以对于非 text 类型，使用 value 和 on:input
  // 对于 text 类型，使用 bind:value 以保持向后兼容
</script>

{#if type === 'text'}
  <input
    type="text"
    {id}
    {placeholder}
    {disabled}
    {autocomplete}
    class={cn(
      'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
      className
    )}
    bind:value
  />
{:else}
  <input
    {type}
    {id}
    {placeholder}
    {disabled}
    {autocomplete}
    {value}
    on:input={handleInput}
    class={cn(
      'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
      className
    )}
  />
{/if}

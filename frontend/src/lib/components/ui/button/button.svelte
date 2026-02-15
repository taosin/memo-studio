<script>
  import { cn } from '$lib/utils.js';
  import { createEventDispatcher } from 'svelte';

  export let variant = 'default';
  export let size = 'default';
  export let className = '';
  export let disabled = false;
  export let type = 'button';
  export let loading = false;

  const dispatch = createEventDispatcher();

  const variants = {
    default: 'bg-primary text-primary-foreground hover:bg-primary/90 hover:shadow-md hover:-translate-y-0.5 active:translate-y-0',
    destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/90 hover:shadow-md',
    outline: 'border border-input bg-background hover:bg-accent hover:text-accent-foreground hover:border-primary/50',
    secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/80 hover:shadow-sm',
    ghost: 'hover:bg-accent hover:text-accent-foreground',
    link: 'text-primary underline-offset-4 hover:underline',
    gradient: 'bg-gradient-to-r from-primary to-primary-light text-primary-foreground hover:shadow-lg hover:shadow-primary/25 hover:-translate-y-0.5',
  };

  const sizes = {
    default: 'h-10 px-4 py-2',
    sm: 'h-8 px-3 text-sm',
    lg: 'h-11 rounded-lg px-8',
    icon: 'h-10 w-10',
  };

  function handleClick(event) {
    if (!disabled && !loading) {
      dispatch('click', event);
    }
  }
</script>

<button
  type={type}
  class={cn(
    'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-all duration-200 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
    variants[variant],
    sizes[size],
    className
  )}
  disabled={disabled || loading}
  on:click={handleClick}
>
  {#if loading}
    <svg class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  {/if}
  <slot />
</button>

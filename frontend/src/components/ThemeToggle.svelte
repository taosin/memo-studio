<script>
  import { themeStore } from '../stores/theme.js';

  let isAnimating = false;

  function toggleTheme() {
    if (isAnimating) return;
    isAnimating = true;
    $themeStore = $themeStore === 'light' ? 'dark' : 'light';
    setTimeout(() => isAnimating = false, 300);
  }
</script>

<button 
  class="theme-toggle" 
  on:click={toggleTheme} 
  aria-label="切换主题"
  class:light={$themeStore === 'light'}
>
  <div class="icon-wrapper" class:rotate={isAnimating}>
    {#if $themeStore === 'light'}
      🌙
    {:else}
      ☀️
    {/if}
  </div>
</button>

<style>
  .theme-toggle {
    background: none;
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 50%;
    width: 2.5rem;
    height: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    background-color: hsl(var(--accent) / 0.5);
  }

  .theme-toggle:hover {
    background-color: hsl(var(--accent));
    transform: scale(1.05);
  }

  .theme-toggle:active {
    transform: scale(0.95);
  }

  .icon-wrapper {
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .icon-wrapper.rotate {
    transform: rotate(360deg);
  }
</style>

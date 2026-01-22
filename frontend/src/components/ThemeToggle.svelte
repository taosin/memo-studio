<script>
  import { themeStore } from '../stores/theme.js';

  function toggleTheme() {
    $themeStore = $themeStore === 'light' ? 'dark' : 'light';
  }
</script>

<button class="theme-toggle" on:click={toggleTheme} aria-label="切换主题">
  {#if $themeStore === 'light'}
    🌙
  {:else}
    ☀️
  {/if}
</button>

<style>
  .theme-toggle {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 50%;
    width: 2.5rem;
    height: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s ease;
  }

  .theme-toggle:hover {
    background-color: var(--bg-secondary);
  }
</style>

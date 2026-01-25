<script>
  import { onMount } from 'svelte';
  import { theme, applyTheme, toggleTheme } from '$lib/theme.js';

  let current = 'dark';
  const unsub = theme.subscribe((t) => {
    current = t;
    applyTheme(t);
  });

  onMount(() => {
    applyTheme(current);
    // 生产环境注册 Service Worker
    if (typeof navigator !== 'undefined' && 'serviceWorker' in navigator && import.meta.env.PROD) {
      navigator.serviceWorker.register('/service-worker.js').catch(() => {});
    }
    return () => unsub();
  });
</script>

<svelte:head>
  <meta name="theme-color" content="#0f172a" />
</svelte:head>

<div class="app">
  <header class="topbar">
    <div class="brand">Memo Studio</div>
    <div class="hint">极简记录 · Ctrl/Cmd + Enter 保存</div>
    <div class="spacer" />
    <button class="iconBtn" on:click={toggleTheme} aria-label="切换主题" title="切换主题">
      {#if current === 'dark'}
        <span class="icon">🌙</span>
      {:else}
        <span class="icon">☀️</span>
      {/if}
    </button>
  </header>
  <main class="main">
    <slot />
  </main>
</div>

<style>
  :global(:root) {
    --bg: #0b1220;
    --panel: rgba(2, 6, 23, 0.35);
    --panel-2: rgba(2, 6, 23, 0.30);
    --text: #e5e7eb;
    --muted: rgba(148, 163, 184, 0.9);
    --border: rgba(148, 163, 184, 0.16);
    --border-2: rgba(148, 163, 184, 0.18);
    --topbar: rgba(11, 18, 32, 0.75);
    --accent: rgba(34, 197, 94, 1);
    --accent-soft: rgba(34, 197, 94, 0.16);
    --danger: rgba(248, 113, 113, 1);
  }
  :global(:root[data-theme='light']) {
    --bg: #f8fafc;
    --panel: rgba(255, 255, 255, 0.85);
    --panel-2: rgba(255, 255, 255, 0.75);
    --text: #0f172a;
    --muted: rgba(71, 85, 105, 0.9);
    --border: rgba(15, 23, 42, 0.12);
    --border-2: rgba(15, 23, 42, 0.16);
    --topbar: rgba(248, 250, 252, 0.75);
    --accent: rgba(16, 185, 129, 1);
    --accent-soft: rgba(16, 185, 129, 0.14);
  }

  :global(html, body) {
    height: 100%;
  }
  :global(body) {
    margin: 0;
    background: var(--bg);
    color: var(--text);
    font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica,
      Arial, "Apple Color Emoji", "Segoe UI Emoji";
  }
  :global(a) {
    color: inherit;
  }
  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }
  .topbar {
    position: sticky;
    top: 0;
    z-index: 10;
    backdrop-filter: blur(10px);
    background: var(--topbar);
    border-bottom: 1px solid var(--border);
    padding: 14px 16px;
    display: flex;
    align-items: baseline;
    gap: 12px;
  }
  .brand {
    font-weight: 700;
    letter-spacing: 0.2px;
  }
  .hint {
    font-size: 12px;
    color: var(--muted);
  }
  .spacer {
    flex: 1;
  }
  .iconBtn {
    border-radius: 10px;
    border: 1px solid var(--border);
    background: var(--panel);
    color: inherit;
    padding: 6px 10px;
    cursor: pointer;
  }
  .iconBtn:hover {
    filter: brightness(1.02);
  }
  .icon {
    font-size: 14px;
  }
  .main {
    padding: 16px;
    max-width: 1100px;
    width: 100%;
    margin: 0 auto;
    box-sizing: border-box;
  }
</style>


<script>
  import { onMount } from 'svelte';
  import { afterNavigate, goto } from '$app/navigation';
  import { theme, applyTheme, toggleTheme } from '$lib/theme.js';
  import { api } from '$lib/api.js';

  let current = 'dark';
  let authed = false;
  let isAdmin = false;
  const unsub = theme.subscribe((t) => {
    current = t;
    applyTheme(t);
  });

  function syncAuth() {
    try {
      if (typeof localStorage === 'undefined') {
        authed = false;
        isAdmin = false;
        return;
      }
      const t = localStorage.getItem('token');
      authed = !!t;
      const u = JSON.parse(localStorage.getItem('user') || '{}');
      isAdmin = !!u?.is_admin;
    } catch {
      authed = false;
      isAdmin = false;
    }
  }

  // 需要在组件初始化阶段注册（不能放到 onMount 里）
  afterNavigate(() => syncAuth());

  async function logout() {
    await api.logout();
    syncAuth();
  }

  onMount(() => {
    applyTheme(current);
    // 生产环境注册 Service Worker
    if (typeof navigator !== 'undefined' && 'serviceWorker' in navigator && import.meta.env.PROD) {
      navigator.serviceWorker.register('/service-worker.js').catch(() => {});
    }

    syncAuth();
    window.addEventListener('storage', syncAuth);

    return () => {
      window.removeEventListener('storage', syncAuth);
      unsub();
    };
  });
</script>

<svelte:head>
  <meta name="theme-color" content="#0f172a" />
</svelte:head>

<div class="app">
  <header class="topbar">
    <a href="/" class="brand" on:click|preventDefault={() => goto('/')}>Memo Studio</a>
    <div class="hint">极简记录 · Ctrl/Cmd + Enter 保存</div>
    <div class="spacer" />
    {#if authed}
      <a class="nav" href="/profile">个人信息</a>
      {#if isAdmin}
        <a class="nav" href="/admin/users">用户管理</a>
      {/if}
      <button class="navBtn" on:click={logout}>登出</button>
    {:else}
      <a class="nav" href="/login">登录</a>
    {/if}
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
    color: inherit;
    text-decoration: none;
    cursor: pointer;
  }
  .brand:hover {
    filter: brightness(1.1);
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
  .nav {
    font-size: 12px;
    color: inherit;
    text-decoration: none;
    padding: 6px 10px;
    border-radius: 10px;
    border: 1px solid var(--border);
    background: var(--panel);
  }
  .nav:hover {
    filter: brightness(1.02);
  }
  .navBtn {
    font-size: 12px;
    color: inherit;
    padding: 6px 10px;
    border-radius: 10px;
    border: 1px solid var(--border);
    background: var(--panel);
    cursor: pointer;
  }
  .navBtn:hover {
    filter: brightness(1.02);
  }
  .main {
    padding: 16px;
    max-width: 1100px;
    width: 100%;
    margin: 0 auto;
    box-sizing: border-box;
  }
</style>


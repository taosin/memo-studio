<script>
  // 极简全局样式（不依赖额外 CSS 框架）
</script>

<svelte:head>
  <meta name="theme-color" content="#0f172a" />
</svelte:head>

<div class="app">
  <header class="topbar">
    <div class="brand">Memo Studio</div>
    <div class="hint">输入一条想法，回车保存</div>
  </header>
  <main class="main">
    <slot />
  </main>
</div>

<style>
  :global(html, body) {
    height: 100%;
  }
  :global(body) {
    margin: 0;
    background: #0b1220;
    color: #e5e7eb;
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
    background: rgba(11, 18, 32, 0.75);
    border-bottom: 1px solid rgba(148, 163, 184, 0.15);
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
    color: rgba(148, 163, 184, 0.9);
  }
  .main {
    padding: 16px;
    max-width: 1100px;
    width: 100%;
    margin: 0 auto;
    box-sizing: border-box;
  }
</style>


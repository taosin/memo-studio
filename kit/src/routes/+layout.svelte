<script>
  import { onMount } from "svelte";
  import { afterNavigate, goto } from "$app/navigation";
  import { theme, applyTheme, toggleTheme } from "$lib/theme.js";
  import { api } from "$lib/api.js";

  let current = "dark";
  let authed = false;
  let isAdmin = false;
  let navOpen = false;
  const unsub = theme.subscribe((t) => {
    current = t;
    applyTheme(t);
  });

  function closeNav() {
    navOpen = false;
  }

  function syncAuth() {
    try {
      if (typeof localStorage === "undefined") {
        authed = false;
        isAdmin = false;
        return;
      }
      const t = localStorage.getItem("token");
      authed = !!t;
      const u = JSON.parse(localStorage.getItem("user") || "{}");
      isAdmin = !!u?.is_admin;
    } catch {
      authed = false;
      isAdmin = false;
    }
  }

  afterNavigate(() => {
    syncAuth();
    closeNav();
  });

  async function logout() {
    await api.logout();
    syncAuth();
  }

  onMount(() => {
    applyTheme(current);
    // Service Worker 已禁用
    // 如需开启，取消下方注释
    /*
    if (
      typeof navigator !== "undefined" &&
      "serviceWorker" in navigator &&
      import.meta.env.PROD
    ) {
      navigator.serviceWorker.register("/service-worker.js").catch(() => {});
    }
    */

    syncAuth();
    window.addEventListener("storage", syncAuth);

    return () => {
      window.removeEventListener("storage", syncAuth);
      unsub();
    };
  });
</script>

<svelte:head>
  <meta name="theme-color" content="#0f172a" />
</svelte:head>

<div class="app">
  <header class="topbar">
    <a href="/" class="brand" on:click|preventDefault={() => goto("/")}>
      <img src="/favicon.svg" alt="" class="brandIcon" width="24" height="24" />
      Memo Studio
    </a>
    <div class="spacer"></div>
    {#if authed}
      <nav class="navWrap">
        <a class="nav" href="/notebooks" on:click={closeNav}>笔记本</a>
        <a class="nav" href="/tags" on:click={closeNav}>标签</a>
        <a class="nav" href="/stats" on:click={closeNav}>统计</a>
      </nav>
      <div class="navRight">
        <a class="nav" href="/profile" on:click={closeNav} title="个人信息">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
        </a>
        <button
          class="iconBtn menuBtn"
          on:click={() => (navOpen = !navOpen)}
          aria-label="更多菜单"
          aria-expanded={navOpen}
          title="更多"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <circle cx="12" cy="12" r="1"></circle>
            <circle cx="19" cy="12" r="1"></circle>
            <circle cx="5" cy="12" r="1"></circle>
          </svg>
        </button>
      </div>
      {#if navOpen}
        <div
          class="navOverlay"
          role="button"
          tabindex="-1"
          on:click={closeNav}
          on:keydown={(e) => e.key === "Escape" && closeNav()}
        ></div>
        <div class="navDropdown" role="menu">
          <div class="navSection">
            <span class="navSectionTitle">功能</span>
            <a href="/insights" on:click={closeNav}>🧠 洞察</a>
            <a href="/locations" on:click={closeNav}>📍 位置</a>
            <a href="/resources" on:click={closeNav}>🖼️ 资源</a>
            <a href="/stocks" on:click={closeNav}>📈 股票</a>
            <a href="/voice" on:click={closeNav}>🎤 语音</a>
          </div>
          <div class="navSection">
            <span class="navSectionTitle">数据</span>
            <a href="/export" on:click={closeNav}>导出</a>
            <a href="/import" on:click={closeNav}>导入</a>
          </div>
          <div class="navSection">
            <span class="navSectionTitle">设置</span>
            <a href="/settings" on:click={closeNav}>应用设置</a>
            <a href="/help" on:click={closeNav}>帮助</a>
            {#if isAdmin}
              <a href="/admin/users" on:click={closeNav}>用户管理</a>
            {/if}
            <a
              href="#"
              on:click|preventDefault={() => {
                logout();
                closeNav();
              }}
              class="logout">登出</a
            >
          </div>
        </div>
      {/if}
    {:else}
      <a class="nav" href="/login">登录</a>
    {/if}
    <button
      class="iconBtn"
      on:click={toggleTheme}
      aria-label="切换主题"
      title="切换主题"
    >
      {#if current === "dark"}
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
    --panel-2: rgba(2, 6, 23, 0.3);
    --text: #e5e7eb;
    --muted: rgba(148, 163, 184, 0.9);
    --border: rgba(148, 163, 184, 0.16);
    --border-2: rgba(148, 163, 184, 0.18);
    --topbar: rgba(11, 18, 32, 0.75);
    --accent: rgba(34, 197, 94, 1);
    --accent-soft: rgba(34, 197, 94, 0.16);
    --danger: rgba(248, 113, 113, 1);
  }
  :global(:root[data-theme="light"]) {
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
    font-family:
      ui-sans-serif,
      system-ui,
      -apple-system,
      Segoe UI,
      Roboto,
      Helvetica,
      Arial,
      "Apple Color Emoji",
      "Segoe UI Emoji";
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
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-weight: 700;
    letter-spacing: 0.2px;
    color: inherit;
    text-decoration: none;
    cursor: pointer;
  }
  .brand:hover {
    filter: brightness(1.1);
  }
  .brandIcon {
    flex-shrink: 0;
    border-radius: 6px;
  }
  .hint {
    font-size: 12px;
    color: var(--muted);
  }
  .spacer {
    flex: 1;
  }
  .navWrap {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }
  .navRight {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .menuBtn {
    display: none;
  }
  .navOverlay {
    display: none;
  }
  .navDropdown {
    display: none;
  }
  @media (max-width: 768px) {
    .hint {
      display: none;
    }
    .navWrap {
      display: none;
    }
    .menuBtn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }
    .navOverlay {
      display: block;
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.35);
      z-index: 11;
      animation: fadeIn 0.2s ease;
    }
    .navDropdown {
      display: flex;
      flex-direction: column;
      position: fixed;
      top: 56px;
      right: 16px;
      min-width: 200px;
      max-width: 90vw;
      max-height: 70vh;
      overflow: auto;
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: 12px;
      box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
      z-index: 12;
      padding: 8px 0;
      animation: slideDown 0.2s ease;
    }
    .navSection {
      padding: 4px 0;
      border-bottom: 1px solid var(--border);
    }
    .navSection:last-child {
      border-bottom: none;
    }
    .navSectionTitle {
      display: block;
      padding: 6px 12px;
      font-size: 11px;
      font-weight: 600;
      color: var(--muted);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
    .navDropdown a {
      display: block;
      padding: 8px 12px;
      margin: 0 8px;
      border-radius: 6px;
      text-decoration: none;
      color: inherit;
      font-size: 14px;
      transition: background 0.15s ease;
    }
    .navDropdown a:hover {
      background: rgba(148, 163, 184, 0.12);
    }
    .navDropdown a.logout {
      color: rgba(248, 113, 113, 0.9);
    }
    .navDropdown a.logout:hover {
      background: rgba(248, 113, 113, 0.12);
    }
  }
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
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
    transition:
      background 0.15s ease,
      border-color 0.15s ease;
  }
  .nav:hover {
    background: rgba(148, 163, 184, 0.1);
    border-color: rgba(148, 163, 184, 0.22);
  }
  .navBtn {
    font-size: 12px;
    color: inherit;
    padding: 6px 10px;
    border-radius: 10px;
    border: 1px solid var(--border);
    background: var(--panel);
    cursor: pointer;
    transition:
      background 0.15s ease,
      border-color 0.15s ease;
  }
  .navBtn:hover {
    background: rgba(148, 163, 184, 0.1);
    border-color: rgba(148, 163, 184, 0.22);
  }
  .main {
    padding: 16px;
    max-width: 1100px;
    width: 100%;
    margin: 0 auto;
    box-sizing: border-box;
  }
  @media (max-width: 600px) {
    .topbar {
      padding: 10px 12px;
    }
    .main {
      padding: 12px;
    }
  }
</style>

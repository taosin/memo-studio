<script>
  import { theme, applyTheme } from "$lib/theme.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  const PAGE_SIZE_KEY = "memo-page-size";
  const DEFAULT_VIEW_KEY = "memo-default-view";

  let currentTheme = "dark";
  let pageSize = 20;
  let defaultView = "list";
  let toast = "";

  const unsub = theme.subscribe((t) => {
    currentTheme = t;
  });

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ""), 1500);
  }

  function getStoredPageSize() {
    if (typeof localStorage === "undefined") return 20;
    const v = parseInt(localStorage.getItem(PAGE_SIZE_KEY), 10);
    return Number.isFinite(v) && v >= 5 && v <= 100 ? v : 20;
  }

  function getStoredDefaultView() {
    if (typeof localStorage === "undefined") return "list";
    const v = localStorage.getItem(DEFAULT_VIEW_KEY);
    return v === "card" || v === "list" ? v : "list";
  }

  function savePageSize() {
    const n = Math.min(100, Math.max(5, pageSize));
    pageSize = n;
    try {
      localStorage.setItem(PAGE_SIZE_KEY, String(n));
      setToast("已保存每页条数");
    } catch (e) {
      setToast("保存失败");
    }
  }

  function saveDefaultView() {
    try {
      localStorage.setItem(DEFAULT_VIEW_KEY, defaultView);
      setToast("已保存默认视图");
    } catch (e) {
      setToast("保存失败");
    }
  }

  function setTheme(t) {
    currentTheme = t;
    applyTheme(t);
    setToast("主题已切换");
  }

  onMount(() => {
    pageSize = getStoredPageSize();
    defaultView = getStoredDefaultView();
  });
</script>

<svelte:head>
  <title>设置 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">设置</div>
      <a class="link" href="/">返回首页</a>
    </div>
    <p class="hint">应用偏好仅保存在本机，不会同步到服务器。</p>

    {#if toast}
      <div class="toast">{toast}</div>
    {/if}

    <section class="section">
      <h3 class="sectionTitle">外观</h3>
      <div class="option">
        <span class="optionLabel">主题</span>
        <div class="optionValue">
          <button
            class="chip"
            class:active={currentTheme === "dark"}
            on:click={() => setTheme("dark")}
          >
            深色
          </button>
          <button
            class="chip"
            class:active={currentTheme === "light"}
            on:click={() => setTheme("light")}
          >
            浅色
          </button>
        </div>
      </div>
    </section>

    <section class="section">
      <h3 class="sectionTitle">列表与分页</h3>
      <div class="option">
        <span class="optionLabel">每页显示条数（资源库等列表）</span>
        <div class="optionValue">
          <input
            type="number"
            class="input"
            bind:value={pageSize}
            min="5"
            max="100"
          />
          <button class="btn small" on:click={savePageSize}>保存</button>
        </div>
      </div>
      <div class="option">
        <span class="optionLabel">默认列表视图</span>
        <div class="optionValue">
          <button
            class="chip"
            class:active={defaultView === "list"}
            on:click={() => {
              defaultView = "list";
              saveDefaultView();
            }}
          >
            列表
          </button>
          <button
            class="chip"
            class:active={defaultView === "card"}
            on:click={() => {
              defaultView = "card";
              saveDefaultView();
            }}
          >
            卡片
          </button>
        </div>
      </div>
    </section>

    <section class="section">
      <h3 class="sectionTitle">关于</h3>
      <p class="muted">
        设置项存储在浏览器本地，清除站点数据会恢复默认。主题与顶栏切换按钮同步。
      </p>
    </section>
  </div>
</div>

<style>
  .wrap {
    padding: 24px 16px;
    max-width: 560px;
    margin: 0 auto;
  }
  .card {
    width: 100%;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 16px;
  }
  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    margin-bottom: 12px;
  }
  .title {
    font-weight: 800;
  }
  .link {
    font-size: 12px;
    text-decoration: underline;
    color: inherit;
  }
  .hint {
    font-size: 13px;
    color: var(--muted);
    margin-bottom: 20px;
  }
  .toast {
    border: 1px solid rgba(34, 197, 94, 0.35);
    background: var(--accent-soft);
    border-radius: 12px;
    padding: 10px 12px;
    margin-bottom: 16px;
  }
  .section {
    margin-bottom: 24px;
  }
  .section:last-child {
    margin-bottom: 0;
  }
  .sectionTitle {
    font-size: 14px;
    font-weight: 700;
    margin: 0 0 12px;
    color: var(--text);
  }
  .option {
    margin-bottom: 14px;
  }
  .optionLabel {
    display: block;
    font-size: 13px;
    color: var(--muted);
    margin-bottom: 8px;
  }
  .optionValue {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
  }
  .chip {
    border-radius: 10px;
    border: 1px solid var(--border);
    background: var(--panel-2);
    color: inherit;
    padding: 8px 14px;
    cursor: pointer;
    font-size: 13px;
  }
  .chip.active {
    border-color: var(--accent);
    background: var(--accent-soft);
  }
  .input {
    width: 80px;
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 8px 12px;
    font-size: 14px;
  }
  .btn.small {
    border-radius: 8px;
    border: 1px solid var(--border);
    background: var(--panel);
    color: inherit;
    padding: 6px 12px;
    cursor: pointer;
    font-size: 12px;
  }
  .muted {
    font-size: 13px;
    color: var(--muted);
    margin: 0;
  }
</style>

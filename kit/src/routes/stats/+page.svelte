<script>
  import { api } from "$lib/api.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let stats = null;
  let loading = false;
  let error = "";

  async function load() {
    loading = true;
    error = "";
    try {
      stats = await api.getStats();
    } catch (e) {
      error = e?.message || "加载失败";
      if (String(error).includes("401") || String(error).includes("未登录")) {
        await goto("/login");
      }
    } finally {
      loading = false;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>统计 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">数据统计</div>
      <a class="link" href="/">返回首页</a>
    </div>
    <p class="hint">当前账号下的笔记、标签、资源与笔记本数量概览。</p>

    {#if error}
      <div class="error">{error}</div>
    {:else if loading}
      <div class="muted">加载中…</div>
    {:else if stats}
      <div class="grid">
        <div class="statCard">
          <div class="statValue">{stats.notes_count ?? 0}</div>
          <div class="statLabel">笔记</div>
        </div>
        <div class="statCard">
          <div class="statValue">{stats.tags_count ?? 0}</div>
          <div class="statLabel">标签</div>
        </div>
        <div class="statCard">
          <div class="statValue">{stats.resources_count ?? 0}</div>
          <div class="statLabel">资源</div>
        </div>
        <div class="statCard">
          <div class="statValue">{stats.notebooks_count ?? 0}</div>
          <div class="statLabel">笔记本</div>
        </div>
        <div class="statCard accent">
          <div class="statValue">{stats.pinned_count ?? 0}</div>
          <div class="statLabel">置顶笔记</div>
        </div>
        <div class="statCard">
          <div class="statValue">{stats.notes_created_7d ?? 0}</div>
          <div class="statLabel">近 7 天新建</div>
        </div>
        <div class="statCard">
          <div class="statValue">{stats.notes_updated_7d ?? 0}</div>
          <div class="statLabel">近 7 天更新</div>
        </div>
      </div>
    {:else}
      <div class="empty">暂无数据</div>
    {/if}
  </div>
</div>

<style>
  .wrap {
    padding: 24px 16px;
    max-width: 700px;
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
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 12px;
  }
  .muted {
    color: var(--muted);
  }
  .empty {
    color: var(--muted);
    padding: 24px;
    text-align: center;
  }
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 16px;
  }
  .statCard {
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 20px;
    background: var(--panel-2);
    text-align: center;
  }
  .statCard.accent {
    border-color: rgba(34, 197, 94, 0.4);
    background: var(--accent-soft);
  }
  .statValue {
    font-size: 28px;
    font-weight: 800;
    line-height: 1.2;
    margin-bottom: 6px;
  }
  .statLabel {
    font-size: 13px;
    color: var(--muted);
  }
</style>

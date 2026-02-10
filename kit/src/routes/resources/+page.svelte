<script>
  import { api } from "$lib/api.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let items = [];
  let total = 0;
  let loading = false;
  let error = "";
  let toast = "";
  let uploading = false;
  let uploadError = "";
  const pageSize = 20;
  let page = 0;
  let fileInput;

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ""), 1500);
  }

  async function load() {
    loading = true;
    error = "";
    try {
      const res = await api.listResources(pageSize, page * pageSize);
      items = res?.items ?? [];
      total = res?.total ?? 0;
    } catch (e) {
      error = e?.message || "加载失败";
      if (String(error).includes("401") || String(error).includes("未登录")) {
        await goto("/login");
      }
    } finally {
      loading = false;
    }
  }

  async function deleteOne(id) {
    if (!confirm("确定删除该资源吗？")) return;
    try {
      await api.deleteResource(id);
      setToast("已删除");
      await load();
    } catch (e) {
      setToast(e?.message || "删除失败");
    }
  }

  function triggerUpload() {
    fileInput?.click();
  }

  async function onFileChange(ev) {
    const file = ev.target.files?.[0];
    if (!file) return;
    uploading = true;
    uploadError = "";
    try {
      await api.uploadResource(file);
      setToast("上传成功");
      await load();
    } catch (e) {
      uploadError = e?.message || "上传失败";
    } finally {
      uploading = false;
      ev.target.value = "";
    }
  }

  function formatSize(bytes) {
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    return (bytes / (1024 * 1024)).toFixed(1) + " MB";
  }

  function formatDate(d) {
    try {
      return new Date(d).toLocaleString("zh-CN");
    } catch {
      return d;
    }
  }

  $: totalPages = Math.max(1, Math.ceil(total / pageSize));
  $: hasPrev = page > 0;
  $: hasNext = page < totalPages - 1;

  onMount(load);
</script>

<svelte:head>
  <title>资源库 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">资源库</div>
      <a class="link" href="/">返回首页</a>
    </div>

    <div class="toolbar">
      <input
        type="file"
        bind:this={fileInput}
        on:change={onFileChange}
        accept="image/*,.pdf,.doc,.docx,.txt,.md"
        style="display: none;"
      />
      <button class="btn" on:click={triggerUpload} disabled={uploading}>
        {uploading ? "上传中…" : "上传文件"}
      </button>
      <span class="muted">共 {total} 个资源</span>
    </div>

    {#if uploadError}
      <div class="error">{uploadError}</div>
    {/if}
    {#if error}
      <div class="error">{error}</div>
    {/if}
    {#if toast}
      <div class="toast">{toast}</div>
    {/if}

    {#if loading && items.length === 0}
      <div class="muted">加载中…</div>
    {:else if items.length === 0}
      <div class="empty">暂无资源，点击「上传文件」添加。</div>
    {:else}
      <div class="list">
        {#each items as r (r.id)}
          <div class="item">
            <div class="preview">
              {#if r.mime_type?.startsWith("image/")}
                <a
                  href={r.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="thumb"
                >
                  <img src={r.url} alt={r.filename} />
                </a>
              {:else}
                <div class="thumb placehold">📄</div>
              {/if}
            </div>
            <div class="meta">
              <a
                href={r.url}
                target="_blank"
                rel="noopener noreferrer"
                class="name">{r.filename}</a
              >
              <span class="size">{formatSize(r.size)}</span>
              <span class="date">{formatDate(r.created_at)}</span>
            </div>
            <div class="actions">
              <a
                href={r.url}
                target="_blank"
                rel="noopener noreferrer"
                class="mini">打开</a
              >
              <button class="mini danger" on:click={() => deleteOne(r.id)}
                >删除</button
              >
            </div>
          </div>
        {/each}
      </div>

      {#if totalPages > 1}
        <div class="pagination">
          <button
            class="mini"
            disabled={!hasPrev}
            on:click={() => (page = page - 1) && load()}
          >
            上一页
          </button>
          <span class="pageInfo">第 {page + 1} / {totalPages} 页</span>
          <button
            class="mini"
            disabled={!hasNext}
            on:click={() => (page = page + 1) && load()}
          >
            下一页
          </button>
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .wrap {
    display: flex;
    justify-content: center;
    padding: 24px 16px;
  }
  .card {
    width: 100%;
    max-width: 900px;
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
  .toolbar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
  }
  .btn {
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 8px 14px;
    cursor: pointer;
    font-weight: 600;
  }
  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .muted {
    font-size: 12px;
    color: var(--muted);
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 10px 12px;
    margin-bottom: 10px;
  }
  .toast {
    border: 1px solid rgba(34, 197, 94, 0.35);
    background: var(--accent-soft);
    border-radius: 12px;
    padding: 10px 12px;
    margin-bottom: 10px;
  }
  .empty {
    color: var(--muted);
    padding: 24px;
    text-align: center;
  }
  .list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .item {
    display: grid;
    grid-template-columns: 72px 1fr auto;
    gap: 12px;
    align-items: center;
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 10px 12px;
    background: var(--panel-2);
  }
  .preview {
    flex-shrink: 0;
  }
  .thumb {
    display: block;
    width: 64px;
    height: 64px;
    border-radius: 8px;
    overflow: hidden;
    background: rgba(15, 23, 42, 0.06);
  }
  .thumb img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .thumb.placehold {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
  }
  .meta {
    min-width: 0;
  }
  .meta .name {
    display: block;
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: inherit;
    text-decoration: none;
  }
  .meta .name:hover {
    text-decoration: underline;
  }
  .meta .size,
  .meta .date {
    font-size: 12px;
    color: var(--muted);
    margin-top: 4px;
  }
  .actions {
    display: flex;
    gap: 8px;
  }
  .mini {
    border-radius: 8px;
    border: 1px solid var(--border);
    background: var(--panel);
    color: inherit;
    padding: 6px 10px;
    cursor: pointer;
    font-size: 12px;
    text-decoration: none;
  }
  .mini.danger {
    border-color: rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
  }
  .pagination {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 16px;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid var(--border);
  }
  .pageInfo {
    font-size: 12px;
    color: var(--muted);
  }
</style>

<script>
  import { api } from "$lib/api.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let format = "json";
  let limit = 500;
  let loading = false;
  let error = "";
  let toast = "";

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ""), 2500);
  }

  async function doExport() {
    loading = true;
    error = "";
    try {
      if (format === "markdown") {
        const { blob, filename } = await api.exportNotes("markdown", limit);
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = filename;
        a.click();
        URL.revokeObjectURL(url);
        setToast("已下载 Markdown 文件");
      } else {
        const data = await api.exportNotes("json", limit);
        const blob = new Blob([JSON.stringify(data, null, 2)], {
          type: "application/json; charset=utf-8",
        });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = data.notes
          ? `memo-export-${new Date().toISOString().slice(0, 19).replace(/:/g, "-")}.json`
          : "memo-export.json";
        a.click();
        URL.revokeObjectURL(url);
        setToast("已下载 JSON 文件");
      }
    } catch (e) {
      error = e?.message || "导出失败";
      if (String(error).includes("401") || String(error).includes("未登录")) {
        await goto("/login");
      }
    } finally {
      loading = false;
    }
  }

  onMount(() => {});
</script>

<svelte:head>
  <title>导出 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">导出笔记</div>
      <a class="link" href="/">返回首页</a>
    </div>
    <p class="hint">
      将当前账号下的笔记导出为 JSON 或 Markdown
      文件，便于备份或迁移。单次最多导出 2000 条。
    </p>

    <div class="form">
      <label class="label" for="export-format">导出格式</label>
      <select id="export-format" class="input" bind:value={format}>
        <option value="json">JSON（含标签、资源等完整数据）</option>
        <option value="markdown">Markdown（纯文本，便于阅读）</option>
      </select>

      <label class="label" for="export-limit">最多导出条数</label>
      <input
        id="export-limit"
        type="number"
        class="input"
        bind:value={limit}
        min="1"
        max="2000"
        step="1"
      />

      {#if error}
        <div class="error">{error}</div>
      {/if}
      {#if toast}
        <div class="toast">{toast}</div>
      {/if}

      <button class="btn" disabled={loading} on:click={doExport}>
        {loading ? "导出中…" : "下载导出文件"}
      </button>
    </div>
  </div>
</div>

<style>
  .wrap {
    padding: 24px 16px;
    max-width: 500px;
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
  .form {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .label {
    font-size: 12px;
    color: var(--muted);
  }
  .input {
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 10px 12px;
    font-size: 14px;
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 10px 12px;
  }
  .toast {
    border: 1px solid rgba(34, 197, 94, 0.35);
    background: var(--accent-soft);
    border-radius: 12px;
    padding: 10px 12px;
  }
  .btn {
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 10px 16px;
    cursor: pointer;
    font-weight: 600;
    margin-top: 8px;
  }
  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>

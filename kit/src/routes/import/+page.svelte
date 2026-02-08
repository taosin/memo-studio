<script>
  import { api } from "$lib/api.js";
  import { goto } from "$app/navigation";

  let jsonText = "";
  let loading = false;
  let error = "";
  let toast = "";
  let result = null;

  const SAMPLE = `[
  { "title": "示例笔记", "content": "这是导入的示例内容。", "tags": ["示例"] }
]`;

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ""), 2500);
  }

  function pasteSample() {
    jsonText = SAMPLE;
  }

  async function doImport() {
    const raw = String(jsonText || "").trim();
    if (!raw) {
      error = "请粘贴或输入 JSON 数组（格式见下方说明）";
      return;
    }
    loading = true;
    error = "";
    result = null;
    try {
      let notes = [];
      try {
        notes = JSON.parse(raw);
      } catch (e) {
        error = "JSON 格式错误：" + (e?.message || "请检查括号与逗号");
        loading = false;
        return;
      }
      if (!Array.isArray(notes)) {
        error = "根节点必须是数组，例如 [{ \"title\": \"...\", \"content\": \"...\", \"tags\": [] }]";
        loading = false;
        return;
      }
      const items = notes.map((n) => ({
        title: typeof n.title === "string" ? n.title : "",
        content: typeof n.content === "string" ? n.content : "",
        tags: Array.isArray(n.tags) ? n.tags.map((t) => String(t)) : [],
      }));
      result = await api.importNotes(items);
      setToast(`导入完成：成功 ${result.created ?? 0} 条，失败 ${result.failed ?? 0} 条`);
      if ((result.created ?? 0) > 0) jsonText = "";
    } catch (e) {
      error = e?.message || "导入失败";
      if (String(error).includes("401") || String(error).includes("未登录")) {
        await goto("/login");
      }
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>导入 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">导入笔记</div>
      <a class="link" href="/">返回首页</a>
    </div>
    <p class="hint">
      将 JSON 格式的笔记数组粘贴到下方，每条需包含 <code>title</code>、<code>content</code>、<code>tags</code>（可选）。单次最多 500 条。
    </p>

    <div class="form">
      <label class="label">JSON 数据</label>
      <textarea
        class="textarea"
        bind:value={jsonText}
        placeholder='[{"title":"标题","content":"内容","tags":["标签1"]}]'
        rows="12"
      ></textarea>
      <button type="button" class="btn ghost" on:click={pasteSample}>填充示例</button>

      {#if error}
        <div class="error">{error}</div>
      {/if}
      {#if toast}
        <div class="toast">{toast}</div>
      {/if}
      {#if result}
        <div class="result">
          成功：{result.created ?? 0}，失败：{result.failed ?? 0}，共：{result.total ?? 0}
        </div>
      {/if}

      <button class="btn" disabled={loading} on:click={doImport}>
        {loading ? "导入中…" : "开始导入"}
      </button>
    </div>

    <div class="help">
      <strong>格式说明</strong>
      <ul>
        <li>根节点为数组，每项为一条笔记。</li>
        <li>每条笔记：<code>title</code>（字符串）、<code>content</code>（字符串）、<code>tags</code>（字符串数组，可选）。</li>
        <li>可从「导出」页面导出 JSON 后，取其中的 <code>notes</code> 数组粘贴到此处做迁移。</li>
      </ul>
    </div>
  </div>
</div>

<style>
  .wrap {
    padding: 24px 16px;
    max-width: 600px;
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
    margin-bottom: 16px;
  }
  .hint code {
    background: var(--panel-2);
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 12px;
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
  .textarea {
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 10px 12px;
    font-size: 13px;
    font-family: ui-monospace, monospace;
    resize: vertical;
    min-height: 120px;
  }
  .btn {
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 10px 16px;
    cursor: pointer;
    font-weight: 600;
  }
  .btn.ghost {
    border-color: var(--border);
    background: var(--panel);
    align-self: flex-start;
  }
  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
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
  .result {
    font-size: 13px;
    color: var(--muted);
  }
  .help {
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid var(--border);
    font-size: 13px;
    color: var(--muted);
  }
  .help ul {
    margin: 8px 0 0;
    padding-left: 20px;
  }
  .help code {
    background: var(--panel-2);
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 12px;
  }
</style>

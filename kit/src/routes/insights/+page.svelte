<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { goto } from '$app/navigation';

  let loading = false;
  let error = '';
  let notes = [];
  let insight = null;
  let timeRange = '30d';
  let selectedNote = '';
  let noteSummary = null;

  const timeRanges = [
    { value: '7d', label: '最近 7 天' },
    { value: '30d', label: '最近 30 天' },
    { value: '90d', label: '最近 3 个月' },
    { value: 'all', label: '全部' },
  ];

  onMount(async () => {
    await loadNotes();
  });

  async function loadNotes() {
    try {
      notes = await api.listNotes();
    } catch (e) {
      error = '加载笔记失败: ' + e.message;
    }
  }

  async function getInsight() {
    if (notes.length === 0) {
      error = '没有笔记可供分析';
      return;
    }

    loading = true;
    error = '';
    insight = null;

    try {
      const noteContents = notes.map(n => n.content);
      const res = await api.call('/api/insights', {
        method: 'POST',
        body: {
          notes: noteContents,
          time_range: timeRange,
        },
      });
      insight = res;
    } catch (e) {
      error = '获取洞察失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  async function summarizeNote() {
    if (!selectedNote) {
      error = '请先选择一条笔记';
      return;
    }

    loading = true;
    error = '';
    noteSummary = null;

    try {
      const res = await api.call('/api/summarize', {
        method: 'POST',
        body: {
          content: selectedNote,
        },
      });
      noteSummary = res;
    } catch (e) {
      error = '总结失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  async function batchSummarize() {
    if (notes.length === 0) {
      error = '没有笔记可供总结';
      return;
    }

    loading = true;
    error = '';

    try {
      const noteContents = notes.slice(0, 5).map(n => n.content);
      const res = await api.call('/api/summarize/batch', {
        method: 'POST',
        body: {
          notes: noteContents,
          limit: 5,
        },
      });

      // 显示批量总结结果
      insight = {
        summary: '批量总结',
        highlights: res.results.map(r => r.summary),
        tips: res.results.map((r, i) => `笔记 ${i + 1}: ${r.summary}`),
      };
    } catch (e) {
      error = '批量总结失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  function goBack() {
    goto('/');
  }

  // 扩展 api 对象以支持自定义端点
  async function customFetch(path, options = {}) {
    const token = localStorage.getItem('token') || '';
    const res = await fetch(`http://localhost:9000${path}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(options.headers || {}),
      },
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `请求失败 (${res.status})`);
    }

    return res.json();
  }

  // 覆盖 api.call
  api.call = customFetch;
</script>

<svelte:head>
  <title>洞察与总结 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">洞察 & 总结</div>
      <button class="link" on:click={goBack}>返回首页</button>
    </div>

    <div class="description">
      AI 驱动的笔记洞察分析，了解你的写作模式和内容主题。
      {#if !window.OPENAI_API_KEY}
        <span class="hint">（需配置 OPENAI_API_KEY 启用完整 AI 功能）</span>
      {/if}
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <div class="toolbar">
      <select bind:value={timeRange} class="select">
        {#each timeRanges as tr}
          <option value={tr.value}>{tr.label}</option>
        {/each}
      </select>

      <button class="btn" on:click={getInsight} disabled={loading || notes.length === 0}>
        {loading ? '分析中...' : '🔍 获取洞察'}
      </button>

      <button class="btn ghost" on:click={batchSummarize} disabled={loading || notes.length === 0}>
        📝 批量总结
      </button>
    </div>

    {#if insight}
      <div class="insight-result">
        <div class="section">
          <h3>📊 整体总结</h3>
          <div class="summary">{insight.summary || '暂无数据'}</div>
        </div>

        {#if insight.keywords && insight.keywords.length > 0}
          <div class="section">
            <h3>🏷️ 关键词</h3>
            <div class="tags">
              {#each insight.keywords as kw}
                <span class="tag">{kw}</span>
              {/each}
            </div>
          </div>
        {/if}

        {#if insight.categories && insight.categories.length > 0}
          <div class="section">
            <h3>📁 分类</h3>
            <div class="tags">
              {#each insight.categories as cat}
                <span class="tag category">{cat}</span>
              {/each}
            </div>
          </div>
        {/if}

        {#if insight.sentiment}
          <div class="section">
            <h3>😊 情感倾向</h3>
            <div class="sentiment">
              {#if insight.sentiment === 'positive'}
                <span class="positive">积极乐观</span>
              {:else if insight.sentiment === 'negative'}
                <span class="negative">有些消极</span>
              {:else}
                <span class="neutral">中性平和</span>
              {/if}
            </div>
          </div>
        {/if}

        {#if insight.trends && insight.trends.length > 0}
          <div class="section">
            <h3>📈 趋势</h3>
            <ul class="list">
              {#each insight.trends as trend}
                <li>{trend}</li>
              {/each}
            </ul>
          </div>
        {/if}

        {#if insight.tips && insight.tips.length > 0}
          <div class="section">
            <h3>💡 建议</h3>
            <ul class="list tips">
              {#each insight.tips as tip}
                <li>{tip}</li>
              {/each}
            </ul>
          </div>
        {/if}
      </div>
    {:else if !loading}
      <div class="empty">
        <p>点击「获取洞察」分析你的笔记</p>
        <p class="muted">共 {notes.length} 条笔记</p>
      </div>
    {/if}

    <hr class="divider" />

    <div class="section">
      <h3>📝 单条笔记总结</h3>

      <div class="toolbar">
        <select bind:value={selectedNote} class="select full">
          <option value="">选择一条笔记...</option>
          {#each notes as note}
            <option value={note.content}>
              {note.content?.substring(0, 50) || '(无内容)'}...
            </option>
          {/each}
        </select>

        <button class="btn" on:click={summarizeNote} disabled={loading || !selectedNote}>
          ✨ 总结
        </button>
      </div>

      {#if noteSummary}
        <div class="note-summary">
          <div class="summary">{noteSummary.summary}</div>

          {#if noteSummary.highlights && noteSummary.highlights.length > 0}
            <div class="sub-section">
              <h4>要点</h4>
              <ul class="list">
                {#each noteSummary.highlights as hl}
                  <li>{hl}</li>
                {/each}
              </ul>
            </div>
          {/if}

          {#if noteSummary.action_items && noteSummary.action_items.length > 0}
            <div class="sub-section">
              <h4>可执行任务</h4>
              <ul class="list action">
                {#each noteSummary.action_items as item}
                  <li>☑️ {item}</li>
                {/each}
              </ul>
            </div>
          {/if}
        </div>
      {/if}
    </div>

    <div class="stats-summary">
      <h3>📊 基础统计</h3>
      <div class="stats-grid">
        <div class="stat">
          <span class="num">{notes.length}</span>
          <span class="label">总笔记数</span>
        </div>
        <div class="stat">
          <span class="num">{notes.filter(n => n.tags?.length > 0).length}</span>
          <span class="label">有标签笔记</span>
        </div>
        <div class="stat">
          <span class="num">{notes.filter(n => n.pinned).length}</span>
          <span class="label">已固定</span>
        </div>
      </div>
    </div>
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
    max-width: 800px;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 24px;
  }
  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }
  .title {
    font-size: 24px;
    font-weight: 800;
  }
  .link {
    font-size: 14px;
    color: var(--accent);
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }
  .description {
    color: var(--muted);
    margin-bottom: 20px;
  }
  .hint {
    font-size: 12px;
    opacity: 0.7;
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 12px;
    margin-bottom: 16px;
    color: #f87171;
  }
  .toolbar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }
  .select {
    padding: 10px 14px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--panel-2);
    color: inherit;
    cursor: pointer;
    min-width: 140px;
  }
  .select.full {
    flex: 1;
    min-width: 200px;
  }
  .btn {
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 10px 18px;
    cursor: pointer;
    font-weight: 600;
  }
  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .btn.ghost {
    background: transparent;
  }
  .insight-result {
    background: var(--panel-2);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
  }
  .section {
    margin-bottom: 24px;
  }
  .section:last-child {
    margin-bottom: 0;
  }
  .section h3 {
    margin: 0 0 12px 0;
    font-size: 16px;
    font-weight: 600;
  }
  .summary {
    font-size: 16px;
    line-height: 1.6;
  }
  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .tag {
    border: 1px solid rgba(148, 163, 184, 0.24);
    border-radius: 999px;
    padding: 4px 12px;
    font-size: 13px;
    background: rgba(15, 23, 42, 0.08);
  }
  .tag.category {
    border-color: rgba(34, 197, 94, 0.4);
    background: rgba(34, 197, 94, 0.1);
  }
  .sentiment {
    font-size: 16px;
  }
  .positive { color: #22c55e; }
  .negative { color: #f87171; }
  .neutral { color: #94a3b8; }
  .list {
    margin: 0;
    padding-left: 20px;
  }
  .list li {
    margin: 6px 0;
    line-height: 1.5;
  }
  .tips li {
    color: #fbbf24;
  }
  .action li {
    color: #60a5fa;
  }
  .empty {
    text-align: center;
    padding: 40px;
    color: var(--muted);
  }
  .muted {
    font-size: 13px;
    opacity: 0.7;
  }
  .divider {
    border: none;
    border-top: 1px solid var(--border);
    margin: 24px 0;
  }
  .sub-section {
    margin-top: 16px;
    margin-bottom: 16px;
  }
  .sub-section h4 {
    margin: 0 0 8px 0;
    font-size: 14px;
    color: var(--muted);
  }
  .note-summary {
    background: var(--panel-2);
    border-radius: 12px;
    padding: 16px;
    margin-top: 16px;
  }
  .stats-summary {
    margin-top: 24px;
    padding-top: 24px;
    border-top: 1px solid var(--border);
  }
  .stats-summary h3 {
    margin: 0 0 16px 0;
    font-size: 16px;
  }
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
  }
  .stat {
    text-align: center;
    padding: 16px;
    background: var(--panel-2);
    border-radius: 12px;
  }
  .stat .num {
    display: block;
    font-size: 32px;
    font-weight: 800;
    color: var(--accent);
  }
  .stat .label {
    font-size: 13px;
    color: var(--muted);
  }
</style>

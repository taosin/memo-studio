<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { buildHeatmap, heatColor } from '$lib/heatmap.js';

  let input = '';
  let notes = [];
  let tags = [];

  let selectedTag = '';
  let searchQ = '';
  let loading = false;
  let error = '';

  let heat = { cells: [], max: 0 };

  function extractTags(text) {
    const matches = String(text || '').match(/#([\p{L}\p{N}_-]+)/gu) || [];
    return [...new Set(matches.map((m) => m.slice(1)))];
  }

  async function reload() {
    loading = true;
    error = '';
    try {
      const [ns, ts] = await Promise.all([api.listNotes(), api.listTags(true)]);
      notes = Array.isArray(ns) ? ns : [];
      tags = Array.isArray(ts) ? ts : [];
      heat = buildHeatmap(notes, 98);
    } catch (e) {
      error = e?.message || '加载失败';
    } finally {
      loading = false;
    }
  }

  async function submit() {
    const text = String(input || '').trim();
    if (!text) return;
    loading = true;
    error = '';
    try {
      const tgs = extractTags(text);
      await api.createNote({ content: text, tags: tgs });
      input = '';
      await reload();
    } catch (e) {
      error = e?.message || '保存失败';
    } finally {
      loading = false;
    }
  }

  async function doSearch() {
    const q = String(searchQ || '').trim();
    if (!q) {
      await reload();
      return;
    }
    loading = true;
    error = '';
    try {
      const ns = await api.search(q);
      notes = Array.isArray(ns) ? ns : [];
      heat = buildHeatmap(notes, 98);
    } catch (e) {
      error = e?.message || '搜索失败';
    } finally {
      loading = false;
    }
  }

  async function randomReview() {
    loading = true;
    error = '';
    try {
      const r = await api.randomReview({ limit: 1, tag: selectedTag || '' });
      if (Array.isArray(r) && r[0]) {
        alert((r[0].content || '').slice(0, 800));
      } else {
        alert('没有可回顾的笔记');
      }
    } catch (e) {
      error = e?.message || '回顾失败';
    } finally {
      loading = false;
    }
  }

  $: filtered = notes.filter((n) => {
    if (!selectedTag) return true;
    const ns = (n.tags || []).map((t) => t.name);
    return ns.includes(selectedTag);
  });

  onMount(reload);
</script>

<div class="grid">
  <aside class="sidebar">
    <div class="panel">
      <div class="panelTitle">标签</div>
      <button class:selected={!selectedTag} class="tag" on:click={() => (selectedTag = '')}>
        全部
      </button>
      {#each tags as t (t.id)}
        <button
          class:selected={selectedTag === t.name}
          class="tag"
          on:click={() => (selectedTag = t.name)}
          title={t.name}
        >
          <span class="dot" style="background:{t.color || 'rgba(34,197,94,0.9)'}" />
          <span class="name">{t.name}</span>
          {#if typeof t.note_count === 'number'}
            <span class="count">{t.note_count}</span>
          {/if}
        </button>
      {/each}
    </div>

    <div class="panel">
      <div class="panelTitle">热力图</div>
      <div class="heatmap" style="--cols: 14;">
        {#each heat.cells as c (c.date)}
          <div class="cell" title={`${c.date} · ${c.count}`} style={`background:${heatColor(c.count, heat.max)}`} />
        {/each}
      </div>
    </div>
  </aside>

  <section class="content">
    <div class="composer">
      <textarea
        class="input"
        bind:value={input}
        rows="3"
        placeholder="记录一条想法… 支持 #标签，例如：今天跑步了 #健康 #运动"
        on:keydown={(e) => {
          if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') submit();
        }}
      />
      <div class="actions">
        <div class="leftHint">Ctrl/Cmd + Enter 保存</div>
        <div class="btns">
          <button class="btn ghost" on:click={randomReview} disabled={loading}>随机回顾</button>
          <button class="btn" on:click={submit} disabled={loading}>保存</button>
        </div>
      </div>
    </div>

    <div class="toolbar">
      <input
        class="search"
        bind:value={searchQ}
        placeholder="全文搜索（FTS5）… 回车搜索"
        on:keydown={(e) => e.key === 'Enter' && doSearch()}
      />
      <button class="btn ghost" on:click={doSearch} disabled={loading}>搜索</button>
      <button class="btn ghost" on:click={reload} disabled={loading}>刷新</button>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if loading}
      <div class="muted">加载中…</div>
    {:else if filtered.length === 0}
      <div class="muted">暂无笔记</div>
    {:else}
      <div class="list">
        {#each filtered as n (n.id)}
          <article class="note">
            <div class="meta">
              <span class="date">{new Date(n.created_at).toLocaleString('zh-CN')}</span>
              <span class="tags">
                {#each n.tags || [] as tg (tg.id)}
                  <span class="pill" style="border-color:{tg.color || 'rgba(34,197,94,0.6)'}">{tg.name}</span>
                {/each}
              </span>
            </div>
            <div class="text">{n.content}</div>
          </article>
        {/each}
      </div>
    {/if}
  </section>
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: 260px 1fr;
    gap: 16px;
  }
  @media (max-width: 900px) {
    .grid {
      grid-template-columns: 1fr;
    }
  }

  .sidebar {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .panel {
    border: 1px solid rgba(148, 163, 184, 0.16);
    background: rgba(2, 6, 23, 0.35);
    border-radius: 12px;
    padding: 12px;
  }
  .panelTitle {
    font-size: 12px;
    color: rgba(148, 163, 184, 0.9);
    margin-bottom: 8px;
  }
  .tag {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 10px;
    border-radius: 10px;
    border: 1px solid transparent;
    background: transparent;
    color: inherit;
    cursor: pointer;
    text-align: left;
  }
  .tag:hover {
    background: rgba(148, 163, 184, 0.08);
  }
  .tag.selected {
    border-color: rgba(34, 197, 94, 0.45);
    background: rgba(34, 197, 94, 0.10);
  }
  .dot {
    width: 8px;
    height: 8px;
    border-radius: 999px;
    flex: 0 0 auto;
  }
  .name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .count {
    font-size: 12px;
    color: rgba(148, 163, 184, 0.95);
  }

  .heatmap {
    display: grid;
    grid-template-columns: repeat(var(--cols), 1fr);
    gap: 4px;
  }
  .cell {
    width: 100%;
    aspect-ratio: 1 / 1;
    border-radius: 4px;
    border: 1px solid rgba(148, 163, 184, 0.08);
  }

  .content {
    display: flex;
    flex-direction: column;
    gap: 12px;
    min-width: 0;
  }
  .composer {
    border: 1px solid rgba(148, 163, 184, 0.16);
    background: rgba(2, 6, 23, 0.35);
    border-radius: 12px;
    padding: 12px;
  }
  .input {
    width: 100%;
    resize: vertical;
    min-height: 90px;
    border-radius: 10px;
    border: 1px solid rgba(148, 163, 184, 0.18);
    background: rgba(15, 23, 42, 0.6);
    color: #e5e7eb;
    padding: 10px 12px;
    box-sizing: border-box;
    outline: none;
  }
  .input:focus {
    border-color: rgba(34, 197, 94, 0.55);
  }
  .actions {
    margin-top: 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
  }
  .leftHint {
    font-size: 12px;
    color: rgba(148, 163, 184, 0.9);
  }
  .btns {
    display: flex;
    gap: 8px;
  }
  .btn {
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: rgba(34, 197, 94, 0.16);
    color: #e5e7eb;
    padding: 8px 12px;
    cursor: pointer;
    font-weight: 600;
  }
  .btn.ghost {
    border-color: rgba(148, 163, 184, 0.22);
    background: rgba(148, 163, 184, 0.06);
    font-weight: 600;
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }

  .toolbar {
    display: flex;
    gap: 8px;
    align-items: center;
  }
  .search {
    flex: 1;
    min-width: 0;
    border-radius: 10px;
    border: 1px solid rgba(148, 163, 184, 0.18);
    background: rgba(15, 23, 42, 0.6);
    color: #e5e7eb;
    padding: 10px 12px;
    outline: none;
  }
  .search:focus {
    border-color: rgba(34, 197, 94, 0.55);
  }

  .list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .note {
    border: 1px solid rgba(148, 163, 184, 0.16);
    background: rgba(2, 6, 23, 0.30);
    border-radius: 12px;
    padding: 12px;
  }
  .meta {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    margin-bottom: 8px;
    color: rgba(148, 163, 184, 0.92);
    font-size: 12px;
  }
  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    justify-content: flex-end;
  }
  .pill {
    border: 1px solid rgba(148, 163, 184, 0.24);
    border-radius: 999px;
    padding: 2px 8px;
    color: rgba(229, 231, 235, 0.95);
    background: rgba(15, 23, 42, 0.35);
  }
  .text {
    white-space: pre-wrap;
    word-break: break-word;
    line-height: 1.6;
  }
  .muted {
    color: rgba(148, 163, 184, 0.9);
    padding: 10px 4px;
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.10);
    color: rgba(254, 202, 202, 0.95);
    border-radius: 12px;
    padding: 10px 12px;
  }
</style>


<script>
  import { onDestroy, onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { buildHeatmap, heatColor } from '$lib/heatmap.js';
  import { renderMiniMarkdown } from '$lib/miniMarkdown.js';
  import { notesStore, tagsStore } from '$lib/stores.js';
  import { goto } from '$app/navigation';

  let input = '';
  let baseNotes = [];
  let notes = []; // 当前展示数据（可能来自搜索）
  let tags = [];

  let selectedTag = '';
  let searchQ = '';
  let loading = false;
  let error = '';
  let toast = '';

  let heat = { cells: [], max: 0 };

  let inputEl;
  let showSidebar = true;
  let reviewOpen = false;
  let reviewText = '';
  let editOpen = false;
  let editId = null;
  let editText = '';
  let debounceTimer;
  let draftTimer;

  function extractTags(text) {
    const matches = String(text || '').match(/#([\p{L}\p{N}_-]+)/gu) || [];
    return [...new Set(matches.map((m) => m.slice(1)))];
  }

  function stripHtml(s) {
    return String(s || '').replace(/<[^>]*>/g, '').trim();
  }

  function setToast(msg) {
    toast = msg;
    if (!msg) return;
    setTimeout(() => {
      if (toast === msg) toast = '';
    }, 1800);
  }

  async function reload() {
    loading = true;
    error = '';
    try {
      const [ns, ts] = await Promise.all([api.listNotes(), api.listTags(true)]);
      baseNotes = Array.isArray(ns) ? ns : [];
      notes = baseNotes;
      tags = Array.isArray(ts) ? ts : [];
      notesStore.set(notes);
      tagsStore.set(tags);
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
      // 乐观更新：先塞一条到顶部，提升体感
      const optimistic = {
        id: `tmp_${Date.now()}`,
        content: text,
        created_at: new Date().toISOString(),
        tags: tgs.map((name) => ({ id: `tmp_${name}`, name, color: 'rgba(34,197,94,0.6)' }))
      };
      notes = [optimistic, ...notes];
      heat = buildHeatmap(notes, 98);

      await api.createNote({ content: text, tags: tgs });
      input = '';
      setToast('已保存');
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

  function scheduleSearch() {
    // 注意：不要用 `$:` 里读写 debounceTimer，会造成自触发循环导致闪烁/崩溃
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      doSearch();
    }, 300);
  }

  // ========== 图片上传 ==========
  let uploadLoading = false;

  async function uploadImage(file) {
    if (!file) return null;
    // 检查是否为图片
    if (!file.type.startsWith('image/')) {
      setToast('仅支持图片文件');
      return null;
    }
    // 检查大小（5MB）
    if (file.size > 5 * 1024 * 1024) {
      setToast('图片大小不能超过 5MB');
      return null;
    }
    uploadLoading = true;
    try {
      const res = await api.uploadResource(file);
      // 返回 Markdown 图片链接
      return `![${file.name}](/uploads/${res.path})`;
    } catch (e) {
      setToast(e?.message || '图片上传失败');
      return null;
    } finally {
      uploadLoading = false;
    }
  }

  async function handlePaste(ev) {
    const items = ev.clipboardData?.items;
    if (!items) return;
    for (const item of items) {
      if (item.type.startsWith('image/')) {
        ev.preventDefault();
        const file = item.getAsFile();
        const md = await uploadImage(file);
        if (md) {
          input = (input || '') + '\n' + md;
          setToast('图片已插入');
        }
        break;
      }
    }
  }

  async function handleDrop(ev) {
    ev.preventDefault();
    const files = ev.dataTransfer?.files;
    if (!files?.length) return;

    for (const file of files) {
      if (file.type.startsWith('image/')) {
        const md = await uploadImage(file);
        if (md) {
          input = (input || '') + '\n' + md;
        }
      }
    }
    if (input) setToast('图片已插入');
  }

  function handleDragOver(ev) {
    ev.preventDefault();
  }

  function clearSearch() {
    searchQ = '';
    scheduleSearch();
  }

  async function randomReview() {
    loading = true;
    error = '';
    try {
      const r = await api.randomReview({ limit: 1, tag: selectedTag || '' });
      if (Array.isArray(r) && r[0]) {
        reviewText = stripHtml(r[0].content || '').slice(0, 1200);
        reviewOpen = true;
      } else {
        setToast('没有可回顾的笔记');
      }
    } catch (e) {
      error = e?.message || '回顾失败';
    } finally {
      loading = false;
    }
  }

  function openEdit(note) {
    if (!note) return;
    const idStr = String(note.id ?? '');
    if (idStr.startsWith('tmp_')) return;
    editId = note.id;
    editText = stripHtml(note.content || '');
    editOpen = true;
  }

  async function saveEdit() {
    const text = String(editText || '').trim();
    if (!editId) return;
    if (!text) {
      setToast('内容不能为空');
      return;
    }
    loading = true;
    error = '';
    try {
      const tgs = extractTags(text);
      await api.updateNote(editId, { content: text, tags: tgs });
      setToast('已更新');
      editOpen = false;
      editId = null;
      await reload();
    } catch (e) {
      error = e?.message || '更新失败';
    } finally {
      loading = false;
    }
  }

  async function removeNote(noteId) {
    if (!confirm('确定删除这条笔记吗？')) return;
    loading = true;
    error = '';
    try {
      await api.deleteNote(noteId);
      setToast('已删除');
      await reload();
    } catch (e) {
      error = e?.message || '删除失败';
    } finally {
      loading = false;
    }
  }

  $: filtered = notes.filter((n) => {
    if (!selectedTag) return true;
    const ns = (n.tags || []).map((t) => t.name);
    return ns.includes(selectedTag);
  });

  onDestroy(() => {
    clearTimeout(debounceTimer);
    clearTimeout(draftTimer);
  });

  onMount(async () => {
    // 未登录：跳转登录页
    try {
      const t = localStorage.getItem('token') || '';
      if (!t) {
        await goto('/login');
        return;
      }
    } catch {
      await goto('/login');
      return;
    }

    // 恢复草稿
    try {
      const draft = localStorage.getItem('memo_draft_v1') || '';
      if (String(input || '').trim() === '' && String(draft).trim() !== '') {
        input = draft;
      }
    } catch {}

    await reload();
    inputEl?.focus();
    // 移动端默认收起侧边栏
    if (typeof window !== 'undefined' && window.innerWidth < 900) showSidebar = false;
  });
</script>

<div class="grid">
  <aside class="sidebar" class:hidden={!showSidebar}>
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
          <span class="dot" style="background:{t.color || 'rgba(34,197,94,0.9)'}"></span>
          <span class="name">{t.name}</span>
          {#if typeof t.note_count === 'number'}
            <span class="count">{t.note_count}</span>
          {/if}
        </button>
      {/each}
    </div>

    <div class="panel">
      <div class="panelTitle">热力图</div>
      <div class="heatmap">
        {#each heat.cells as c (c.date)}
          <div class="cell" title={`${c.date} · ${c.count}`} style={`background:${heatColor(c.count, heat.max)}`}></div>
        {/each}
      </div>
    </div>
  </aside>

  <section class="content">
    <div class="mobileBar">
      <button class="btn ghost" on:click={() => (showSidebar = !showSidebar)}>
        {showSidebar ? '收起侧栏' : '展开侧栏'}
      </button>
      <div class="chips">
        <button class:chipSelected={!selectedTag} class="chip" on:click={() => (selectedTag = '')}>全部</button>
        {#each tags.slice(0, 12) as t (t.id)}
          <button
            class:chipSelected={selectedTag === t.name}
            class="chip"
            on:click={() => (selectedTag = t.name)}
            title={t.name}
          >
            #{t.name}
          </button>
        {/each}
      </div>
    </div>

    <div class="composer">
      <textarea
        class="input"
        bind:value={input}
        rows="3"
        placeholder="写一条想法… #标签 · 粘贴/拖拽图片 · Ctrl+Enter 保存"
        bind:this={inputEl}
        on:input={() => {
          // 本地草稿保存（debounce）
          clearTimeout(draftTimer);
          draftTimer = setTimeout(() => {
            try {
              localStorage.setItem('memo_draft_v1', String(input || ''));
            } catch {}
          }, 250);
        }}
        on:keydown={(e) => {
          if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') submit();
        }}
        on:paste={handlePaste}
        on:drop={handleDrop}
        on:dragover={handleDragOver}
      ></textarea>
      <div class="preview">
        <div class="previewTitle">即时预览（极简 Markdown）</div>
        <div class="previewBody">{@html renderMiniMarkdown(input)}</div>
      </div>
      <div class="actions">
        <div class="leftHint">
          {#if uploadLoading}
            <span class="uploading">图片上传中…</span>
          {:else}
            Ctrl/Cmd + Enter 保存，可粘贴/拖拽图片
          {/if}
        </div>
        <div class="btns">
          <button class="btn ghost" on:click={randomReview} disabled={loading || uploadLoading}>随机回顾</button>
          <button class="btn" on:click={submit} disabled={loading || uploadLoading}>保存</button>
        </div>
      </div>
    </div>

    <div class="toolbar">
      <input
        class="search"
        bind:value={searchQ}
        placeholder="全文搜索（FTS5）… 输入即搜 / Enter 立即搜"
        on:input={scheduleSearch}
        on:keydown={(e) => e.key === 'Enter' && doSearch()}
      />
      <button class="btn ghost" on:click={clearSearch} disabled={loading}>清空</button>
      <button class="btn ghost" on:click={doSearch} disabled={loading}>搜索</button>
      <button class="btn ghost" on:click={reload} disabled={loading}>刷新</button>
    </div>

    {#if error}
      <div class="error" role="alert">{error}</div>
    {/if}
    {#if toast}
      <div class="toast" role="status">{toast}</div>
    {/if}

    {#if loading && baseNotes.length === 0}
      <div class="loadingState">
        <div class="loadingDots"><span></span><span></span><span></span></div>
        <p class="loadingText">加载中…</p>
      </div>
    {:else if loading && baseNotes.length > 0}
      <div class="list listLoading">
        {#each filtered as n (n.id)}
          <article class="note" on:dblclick={() => openEdit(n)} title="双击编辑">
            <div class="meta">
              <span class="date">{new Date(n.created_at).toLocaleString('zh-CN')}</span>
              <span class="tags">
                {#each n.tags || [] as tg (tg.id)}
                  <button
                    class="pill"
                    style="border-color:{tg.color || 'rgba(34,197,94,0.6)'}"
                    on:click={() => (selectedTag = tg.name)}
                    title="按该标签筛选"
                  >
                    {tg.name}
                  </button>
                {/each}
              </span>
            </div>
            <div class="text">{stripHtml(n.content)}</div>
            <div class="rowActions">
              <button class="miniBtn" on:click={() => navigator.clipboard?.writeText(stripHtml(n.content))}>复制</button>
              {#if String(n.id).startsWith('tmp_') === false}
                <button class="miniBtn" on:click={() => openEdit(n)}>编辑</button>
                <button class="miniBtn danger" on:click={() => removeNote(n.id)}>删除</button>
              {/if}
            </div>
          </article>
        {/each}
        <div class="loadingBar" aria-hidden="true"></div>
      </div>
    {:else if filtered.length === 0}
      <div class="emptyState">
        <div class="emptyIcon" aria-hidden="true">✍️</div>
        <h2 class="emptyTitle">还没有笔记</h2>
        <p class="emptyDesc">在上方输入框写一条想法，用 <kbd>Ctrl</kbd>+<kbd>Enter</kbd> 或点击「保存」即可</p>
        <p class="emptyTip">支持 #标签、粘贴/拖拽图片</p>
      </div>
    {:else}
      <div class="list">
        {#each filtered as n (n.id)}
          <article class="note" on:dblclick={() => openEdit(n)} title="双击编辑">
            <div class="meta">
              <span class="date">{new Date(n.created_at).toLocaleString('zh-CN')}</span>
              <span class="tags">
                {#each n.tags || [] as tg (tg.id)}
                  <button
                    class="pill"
                    style="border-color:{tg.color || 'rgba(34,197,94,0.6)'}"
                    on:click={() => (selectedTag = tg.name)}
                    title="按该标签筛选"
                  >
                    {tg.name}
                  </button>
                {/each}
              </span>
            </div>
            <div class="text">{stripHtml(n.content)}</div>
            <div class="rowActions">
              <button class="miniBtn" on:click={() => navigator.clipboard?.writeText(stripHtml(n.content))}>复制</button>
              {#if String(n.id).startsWith('tmp_') === false}
                <button class="miniBtn" on:click={() => openEdit(n)}>编辑</button>
                <button class="miniBtn danger" on:click={() => removeNote(n.id)}>删除</button>
              {/if}
            </div>
          </article>
        {/each}
      </div>
    {/if}
  </section>
</div>

{#if reviewOpen}
  <div
    class="overlay"
    role="button"
    tabindex="0"
    on:click={(e) => e.target === e.currentTarget && (reviewOpen = false)}
    on:keydown={(e) => e.key === 'Escape' && (reviewOpen = false)}
  >
    <div
      class="dialog"
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <div class="dialogTitle">随机回顾</div>
      <div class="dialogBody">{reviewText}</div>
      <div class="dialogActions">
        <button class="btn ghost" on:click={() => (reviewOpen = false)}>关闭</button>
        <button class="btn" on:click={randomReview} disabled={loading}>再来一条</button>
      </div>
    </div>
  </div>
{/if}

{#if editOpen}
  <div
    class="overlay"
    role="button"
    tabindex="0"
    on:click={(e) => e.target === e.currentTarget && (editOpen = false)}
    on:keydown={(e) => e.key === 'Escape' && (editOpen = false)}
  >
    <div
      class="dialog"
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <div class="dialogTitle">编辑笔记</div>
      <textarea
        class="dialogInput"
        bind:value={editText}
        rows="8"
        placeholder="修改内容… 支持 #标签"
        on:keydown={(e) => {
          if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') saveEdit();
        }}
      ></textarea>
      <div class="dialogActions">
        <button class="btn ghost" on:click={() => (editOpen = false)} disabled={loading}>取消</button>
        <button class="btn" on:click={saveEdit} disabled={loading}>保存修改</button>
      </div>
    </div>
  </div>
{/if}

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
  .sidebar.hidden {
    display: none;
  }
  .panel {
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 12px;
    padding: 12px;
  }
  .panelTitle {
    font-size: 12px;
    color: var(--muted);
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
    transition: background 0.15s ease, border-color 0.15s ease;
  }
  .tag:hover {
    background: rgba(148, 163, 184, 0.08);
  }
  .tag.selected {
    border-color: rgba(34, 197, 94, 0.45);
    background: var(--accent-soft);
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
    /* 固定 14 列，避免 iOS/Safari 对 repeat(var(--cols)) 兼容问题导致溢出 */
    grid-template-columns: repeat(14, minmax(0, 1fr));
    gap: 4px;
    max-width: 100%;
    /* 给子像素/边框留一点缓冲，避免最后一列被裁切 */
    padding: 1px;
    box-sizing: border-box;
    overflow: visible;
  }
  .cell {
    width: 100%;
    aspect-ratio: 1 / 1;
    border-radius: 4px;
    border: 1px solid rgba(148, 163, 184, 0.10);
    box-sizing: border-box;
  }

  .content {
    display: flex;
    flex-direction: column;
    gap: 12px;
    min-width: 0;
  }
  .mobileBar {
    display: none;
    gap: 10px;
    align-items: center;
  }
  .chips {
    display: flex;
    gap: 8px;
    overflow: auto;
    padding-bottom: 4px;
  }
  .chip {
    white-space: nowrap;
    border-radius: 999px;
    border: 1px solid var(--border);
    background: var(--panel);
    color: inherit;
    padding: 6px 10px;
    cursor: pointer;
    font-size: 12px;
  }
  .chipSelected {
    border-color: rgba(34, 197, 94, 0.45);
    background: var(--accent-soft);
  }
  @media (max-width: 900px) {
    .mobileBar {
      display: flex;
    }
  }

  .composer {
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 16px;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }
  .composer:focus-within {
    border-color: rgba(34, 197, 94, 0.35);
    box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.08);
  }
  .preview {
    margin-top: 12px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: rgba(15, 23, 42, 0.03);
    padding: 10px 12px;
  }
  :global(:root[data-theme="light"]) .preview {
    background: rgba(15, 23, 42, 0.04);
  }
  .previewTitle {
    font-size: 12px;
    color: var(--muted);
    margin-bottom: 6px;
  }
  .previewBody :global(p) {
    margin: 0 0 8px 0;
    line-height: 1.7;
  }
  .previewBody :global(ul) {
    margin: 0 0 8px 18px;
  }
  .previewBody :global(code) {
    padding: 1px 6px;
    border-radius: 8px;
    border: 1px solid var(--border);
    background: rgba(15, 23, 42, 0.06);
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
      "Courier New", monospace;
    font-size: 0.95em;
  }
  .previewBody :global(a) {
    text-decoration: underline;
  }
  .input {
    width: 100%;
    resize: vertical;
    min-height: 90px;
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 12px 14px;
    box-sizing: border-box;
    outline: none;
    font-size: 15px;
    line-height: 1.6;
    transition: border-color 0.2s ease;
  }
  :global(:root[data-theme="light"]) .input {
    background: rgba(15, 23, 42, 0.04);
  }
  .input:focus {
    border-color: rgba(34, 197, 94, 0.55);
  }
  .input::placeholder {
    color: var(--muted);
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
    color: var(--muted);
  }
  .uploading {
    color: var(--accent);
    font-weight: 600;
  }
  .btns {
    display: flex;
    gap: 8px;
  }
  .btn {
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 8px 14px;
    cursor: pointer;
    font-weight: 600;
    transition: background 0.15s ease, border-color 0.15s ease, opacity 0.15s ease;
  }
  .btn:hover:not(:disabled) {
    background: rgba(34, 197, 94, 0.22);
    border-color: rgba(34, 197, 94, 0.65);
  }
  .btn.ghost {
    border-color: var(--border);
    background: var(--panel);
    font-weight: 600;
  }
  .btn.ghost:hover:not(:disabled) {
    background: rgba(148, 163, 184, 0.1);
    border-color: rgba(148, 163, 184, 0.25);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }

  .toolbar {
    display: flex;
    gap: 8px;
    align-items: center;
    flex-wrap: wrap;
  }
  .search {
    flex: 1;
    min-width: 120px;
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 10px 12px;
    outline: none;
    font-size: 14px;
    transition: border-color 0.2s ease;
  }
  :global(:root[data-theme="light"]) .search {
    background: rgba(15, 23, 42, 0.04);
  }
  .search:focus {
    border-color: rgba(34, 197, 94, 0.55);
  }

  .list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .listLoading {
    position: relative;
  }
  .loadingBar {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(34, 197, 94, 0.4),
      transparent
    );
    background-size: 200% 100%;
    animation: loadingBar 1.2s ease-in-out infinite;
  }
  @keyframes loadingBar {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }
  .loadingState {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px;
    gap: 16px;
  }
  .loadingDots {
    display: flex;
    gap: 8px;
  }
  .loadingDots span {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--accent);
    opacity: 0.6;
    animation: loadingDot 1.4s ease-in-out infinite both;
  }
  .loadingDots span:nth-child(1) { animation-delay: 0s; }
  .loadingDots span:nth-child(2) { animation-delay: 0.2s; }
  .loadingDots span:nth-child(3) { animation-delay: 0.4s; }
  @keyframes loadingDot {
    0%, 80%, 100% { transform: scale(0.8); opacity: 0.5; }
    40% { transform: scale(1.2); opacity: 1; }
  }
  .loadingText {
    margin: 0;
    font-size: 14px;
    color: var(--muted);
  }
  .emptyState {
    text-align: center;
    padding: 48px 24px;
    max-width: 360px;
    margin: 0 auto;
  }
  .emptyIcon {
    font-size: 48px;
    margin-bottom: 12px;
    opacity: 0.9;
  }
  .emptyTitle {
    margin: 0 0 8px 0;
    font-size: 18px;
    font-weight: 600;
    color: inherit;
  }
  .emptyDesc {
    margin: 0 0 8px 0;
    font-size: 14px;
    color: var(--muted);
    line-height: 1.6;
  }
  .emptyTip {
    margin: 0;
    font-size: 13px;
    color: var(--muted);
    opacity: 0.9;
  }
  .emptyState kbd {
    display: inline-block;
    padding: 2px 6px;
    font-size: 12px;
    border-radius: 6px;
    border: 1px solid var(--border);
    background: var(--panel);
    font-family: inherit;
  }
  .note {
    border: 1px solid var(--border);
    background: var(--panel-2);
    border-radius: 12px;
    padding: 14px 16px;
    transition: border-color 0.2s ease, background 0.2s ease, transform 0.15s ease;
  }
  .note:hover {
    border-color: rgba(148, 163, 184, 0.24);
    background: var(--panel);
  }
  .note:active {
    transform: scale(0.998);
  }
  .meta {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    margin-bottom: 8px;
    color: var(--muted);
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
    color: inherit;
    background: rgba(15, 23, 42, 0.08);
    cursor: pointer;
  }
  .text {
    white-space: pre-wrap;
    word-break: break-word;
    line-height: 1.6;
  }
  .rowActions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
    justify-content: flex-end;
    opacity: 0.85;
    transition: opacity 0.2s ease;
  }
  .note:hover .rowActions {
    opacity: 1;
  }
  .miniBtn {
    border-radius: 8px;
    border: 1px solid var(--border);
    background: var(--panel);
    color: inherit;
    padding: 6px 10px;
    cursor: pointer;
    font-size: 12px;
    transition: background 0.15s ease, border-color 0.15s ease;
  }
  .miniBtn:hover {
    background: rgba(148, 163, 184, 0.12);
    border-color: rgba(148, 163, 184, 0.25);
  }
  .miniBtn.danger {
    border-color: rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.10);
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.10);
    color: inherit;
    border-radius: 12px;
    padding: 10px 12px;
  }
  .toast {
    position: sticky;
    top: 58px;
    z-index: 5;
    border: 1px solid rgba(34, 197, 94, 0.35);
    background: var(--accent-soft);
    border-radius: 12px;
    padding: 10px 14px;
    animation: toastIn 0.25s ease;
  }
  @keyframes toastIn {
    from {
      opacity: 0;
      transform: translateY(-8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
    z-index: 50;
  }
  .dialog {
    width: 100%;
    max-width: 640px;
    border-radius: 14px;
    border: 1px solid var(--border);
    background: var(--panel);
    padding: 14px;
  }
  .dialogTitle {
    font-weight: 700;
    margin-bottom: 10px;
  }
  .dialogBody {
    white-space: pre-wrap;
    line-height: 1.7;
    color: inherit;
    max-height: 60vh;
    overflow: auto;
    border: 1px solid var(--border);
    background: rgba(15, 23, 42, 0.06);
    border-radius: 12px;
    padding: 12px;
  }
  .dialogActions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 12px;
  }
  .dialogInput {
    width: 100%;
    resize: vertical;
    min-height: 160px;
    border-radius: 12px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 10px 12px;
    box-sizing: border-box;
    outline: none;
    line-height: 1.6;
    white-space: pre-wrap;
  }
  .dialogInput:focus {
    border-color: rgba(34, 197, 94, 0.55);
  }
</style>


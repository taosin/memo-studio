<script>
  import { onDestroy, onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { buildHeatmap, heatColor } from '$lib/heatmap.js';
  import { notesStore, tagsStore, showToast } from '$lib/stores.js';
  import { keyboardManager } from '$lib/keyboardManager.js';
  import { goto } from '$app/navigation';
  
  // 导入新组件
  import Toast from '$lib/components/Toast.svelte';
  import LoadingState from '$lib/components/LoadingState.svelte';
  import KeyboardHelp from '$lib/components/KeyboardHelp.svelte';
  import SearchBar from '$lib/components/SearchBar.svelte';

  let input = '';
  let baseNotes = [];
  let notes = []; // 当前展示数据（可能来自搜索）
  let tags = [];

  let selectedTag = '';
  let searchQ = '';
  let loading = false;
  let error = '';
  let showKeyboardHelp = false;
  let searchBarEl = null;

  let heat = { cells: [], max: 0 };

  let inputEl;
  let showSidebar = true;
  let reviewOpen = false;
  let reviewText = '';
  let editOpen = false;
  let editId = null;
  let editText = '';
  let editLoading = false;
  let debounceTimer;
  let draftTimer;

  function extractTags(text) {
    const matches = String(text || '').match(/#([\p{L}\p{N}_-]+)/gu) || [];
    return [...new Set(matches.map((m) => m.slice(1)))];
  }

  function stripHtml(s) {
    return String(s || '').replace(/<[^>]*>/g, '').trim();
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
    const tempId = `tmp_${Date.now()}_${Math.random().toString(36).slice(2, 7)}`;
    try {
      const tgs = extractTags(text);
      // 乐观更新：先塞一条到顶部，提升体感
      const optimistic = {
        id: tempId,
        content: text,
        created_at: new Date().toISOString(),
        tags: tgs.map((name) => ({ id: `tmp_${name}`, name, color: 'rgba(34,197,94,0.6)' }))
      };
      notes = [optimistic, ...notes];
      heat = buildHeatmap(notes, 98);

      await api.createNote({ content: text, tags: tgs });
      input = '';
      showToast('✅ 已保存', 'success');
      await reload();
    } catch (e) {
      // 回滚乐观更新
      notes = notes.filter(n => n.id !== tempId);
      heat = buildHeatmap(notes, 98);
      error = e?.message || '保存失败';
      showToast('❌ 保存失败: ' + error, 'error');
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
      showToast('❌ 仅支持图片文件', 'error');
      return null;
    }
    // 检查大小（5MB）
    if (file.size > 5 * 1024 * 1024) {
      showToast('❌ 图片大小不能超过 5MB', 'error');
      return null;
    }
    uploadLoading = true;
    try {
      const res = await api.uploadResource(file);
      // 返回 Markdown 图片链接 - 使用后端返回的 url 字段
      return `![${file.name}](${res.url || '/uploads/' + res.storage_path})`;
    } catch (e) {
      showToast(e?.message || '图片上传失败', 'error');
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
          showToast('📷 图片已插入', 'success');
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
    if (input) showToast('📷 图片已插入', 'success');
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
        showToast('ℹ️ 没有可回顾的笔记', 'info');
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
      showToast('⚠️ 内容不能为空', 'warning');
      return;
    }
    editLoading = true;
    error = '';
    try {
      const tgs = extractTags(text);
      await api.updateNote(editId, { content: text, tags: tgs });
      showToast('✅ 已更新', 'success');
      editOpen = false;
      editId = null;
      await reload();
    } catch (e) {
      error = e?.message || '更新失败';
      showToast('❌ 更新失败: ' + error, 'error');
    } finally {
      editLoading = false;
    }
  }

  async function removeNote(noteId) {
    if (!confirm('确定删除这条笔记吗？')) return;
    loading = true;
    error = '';
    try {
      await api.deleteNote(noteId);
      showToast('✅ 已删除', 'success');
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
    // 注意：键盘事件监听由 onMount 的返回函数清理
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

    // 注册键盘快捷键
    keyboardManager.register('ctrl+k', () => {
      searchBarEl?.focus();
    });

    keyboardManager.register('ctrl+enter', () => {
      submit();
    });

    keyboardManager.register('escape', () => {
      if (editOpen) {
        editOpen = false;
      } else if (reviewOpen) {
        reviewOpen = false;
      }
    });

    keyboardManager.register('?', () => {
      showKeyboardHelp = true;
    });

    keyboardManager.register('r', () => {
      reload();
    });

    keyboardManager.register('b', () => {
      showSidebar = !showSidebar;
    });

    // 监听键盘事件
    const handleKeydown = (e) => keyboardManager.handle(e);
    document.addEventListener('keydown', handleKeydown);

    return () => {
      document.removeEventListener('keydown', handleKeydown);
    };
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
      <div class="composerFooter">
        <div class="leftHint">
          {#if uploadLoading}
            <span class="uploading">图片上传中…</span>
          {:else}
            <span class="hintText">Ctrl + Enter 保存</span>
          {/if}
        </div>
        <div class="rightActions">
          <button class="iconBtn" on:click={randomReview} disabled={loading || uploadLoading} title="随机回顾">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
              <path d="M3 3v5h5"/>
              <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"/>
              <path d="M16 16h5v5"/>
            </svg>
          </button>
          <button class="submitBtn" on:click={submit} disabled={loading || uploadLoading}>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="22" y1="2" x2="11" y2="13"></line>
              <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <div class="toolbar">
      <SearchBar
        bind:value={searchQ}
        bind:inputElement={searchBarEl}
        on:input={scheduleSearch}
        on:search={doSearch}
        on:clear={clearSearch}
      />
      <button class="btn ghost" on:click={doSearch} disabled={loading}>搜索</button>
      <button class="btn ghost" on:click={reload} disabled={loading}>刷新</button>
    </div>

    {#if error}
      <div class="error" role="alert">{error}</div>
    {/if}

    {#if loading && baseNotes.length === 0}
      <LoadingState type="dots" text="加载中…" />
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
        <button class="btn ghost" on:click={() => (editOpen = false)} disabled={editLoading}>取消</button>
        <button class="btn" on:click={saveEdit} disabled={editLoading}>
          {editLoading ? '保存中...' : '保存修改'}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Toast 通知 -->
<Toast />

<!-- 快捷键帮助 -->
{#if showKeyboardHelp}
  <KeyboardHelp on:close={() => showKeyboardHelp = false} />
{/if}

<!-- 快捷键提示 -->
<div class="keyboardHint">
  <span class="hintText">按</span>
  <kbd class="hintKey">?</kbd>
  <span class="hintText">查看快捷键</span>
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
    border-radius: 16px;
    padding: 16px;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }
  .composer:focus-within {
    border-color: rgba(34, 197, 94, 0.35);
    box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.08);
  }
  .input {
    width: 100%;
    resize: none;
    min-height: 80px;
    border: none;
    background: transparent;
    color: var(--text);
    padding: 0;
    box-sizing: border-box;
    outline: none;
    font-size: 15px;
    line-height: 1.7;
    font-family: inherit;
  }
  .input::placeholder {
    color: var(--muted);
  }
  .composerFooter {
    margin-top: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding-top: 12px;
    border-top: 1px solid var(--border);
  }
  .leftHint {
    font-size: 12px;
    color: var(--muted);
  }
  .hintText {
    opacity: 0.7;
  }
  .uploading {
    color: var(--accent);
    font-weight: 500;
  }
  .rightActions {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .iconBtn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    color: var(--muted);
    cursor: pointer;
    border-radius: 8px;
    transition: background 0.15s ease, color 0.15s ease;
  }
  .iconBtn:hover:not(:disabled) {
    background: rgba(148, 163, 184, 0.12);
    color: var(--text);
  }
  .iconBtn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
  .submitBtn {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: var(--accent);
    color: white;
    cursor: pointer;
    border-radius: 10px;
    transition: transform 0.15s ease, opacity 0.15s ease;
  }
  .submitBtn:hover:not(:disabled) {
    transform: scale(1.05);
  }
  .submitBtn:active:not(:disabled) {
    transform: scale(0.95);
  }
  .submitBtn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
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

  .keyboardHint {
    position: fixed;
    bottom: 20px;
    right: 20px;
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    z-index: 100;
    opacity: 0.7;
    transition: opacity 0.2s ease;
  }

  .keyboardHint:hover {
    opacity: 1;
  }

  .hintText {
    font-size: 12px;
    color: var(--muted);
  }

  .hintKey {
    padding: 2px 6px;
    border: 1px solid var(--border);
    border-radius: 4px;
    background: rgba(15, 23, 42, 0.08);
    font-family: ui-monospace, monospace;
    font-size: 11px;
    font-weight: 600;
  }

  @media (max-width: 768px) {
    .keyboardHint {
      display: none;
    }
  }
</style>


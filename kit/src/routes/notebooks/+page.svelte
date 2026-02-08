<script>
  import { api } from "$lib/api.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let notebooks = [];
  let selectedNotebook = null;
  let notebookNotes = [];
  let loading = false;
  let notesLoading = false;
  let error = "";
  let toast = "";
  let newName = "";
  let newColor = "#22c55e";
  let createOpen = false;
  let editId = null;
  let editName = "";
  let editColor = "";

  const PRESET_COLORS = [
    "#22c55e",
    "#3b82f6",
    "#f59e0b",
    "#ef4444",
    "#8b5cf6",
    "#ec4899",
    "#06b6d4",
    "#84cc16",
  ];

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ""), 1500);
  }

  async function load() {
    loading = true;
    error = "";
    try {
      notebooks = await api.listNotebooks();
      if (!Array.isArray(notebooks)) notebooks = [];
    } catch (e) {
      error = e?.message || "加载失败";
      if (String(error).includes("401") || String(error).includes("未登录")) {
        await goto("/login");
      }
    } finally {
      loading = false;
    }
  }

  async function selectNotebook(nb) {
    selectedNotebook = nb;
    notebookNotes = [];
    if (!nb) return;
    notesLoading = true;
    try {
      notebookNotes = await api.listNotebookNotes(nb.id, 100, 0);
      if (!Array.isArray(notebookNotes)) notebookNotes = [];
    } catch (e) {
      setToast(e?.message || "加载笔记失败");
    } finally {
      notesLoading = false;
    }
  }

  async function createNotebook() {
    const name = String(newName || "").trim();
    if (!name) {
      setToast("请输入笔记本名称");
      return;
    }
    try {
      await api.createNotebook({
        name,
        color: newColor || PRESET_COLORS[0],
        sort_order: notebooks.length,
      });
      setToast("已创建");
      newName = "";
      newColor = PRESET_COLORS[0];
      createOpen = false;
      await load();
    } catch (e) {
      setToast(e?.message || "创建失败");
    }
  }

  function openEdit(nb) {
    editId = nb.id;
    editName = nb.name;
    editColor = nb.color || PRESET_COLORS[0];
  }

  function cancelEdit() {
    editId = null;
    editName = "";
    editColor = "";
  }

  async function saveEdit() {
    if (!editId) return;
    const name = String(editName || "").trim();
    if (!name) {
      setToast("笔记本名称不能为空");
      return;
    }
    try {
      await api.updateNotebook(editId, {
        name,
        color: editColor,
      });
      setToast("已保存");
      cancelEdit();
      await load();
      if (selectedNotebook?.id === editId) {
        selectedNotebook =
          notebooks.find((n) => n.id === editId) || selectedNotebook;
      }
    } catch (e) {
      setToast(e?.message || "保存失败");
    }
  }

  async function deleteNotebook(id) {
    if (!confirm("确定删除该笔记本吗？笔记不会被删除，仅解除与笔记本的关联。"))
      return;
    try {
      await api.deleteNotebook(id);
      setToast("已删除");
      if (selectedNotebook?.id === id) {
        selectedNotebook = null;
        notebookNotes = [];
      }
      cancelEdit();
      await load();
    } catch (e) {
      setToast(e?.message || "删除失败");
    }
  }

  function formatDate(d) {
    try {
      return new Date(d).toLocaleString("zh-CN");
    } catch {
      return d;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>笔记本 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">笔记本</div>
      <a class="link" href="/">返回首页</a>
    </div>
    <p class="hint">
      将笔记归类到不同笔记本中，便于按主题管理。同一篇笔记可以属于多个笔记本。
    </p>
    <div class="toolbar">
      <button class="btn" on:click={() => (createOpen = true)}
        >新建笔记本</button
      >
      <span class="muted">共 {notebooks.length} 个笔记本</span>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}
    {#if toast}
      <div class="toast">{toast}</div>
    {/if}

    {#if loading && notebooks.length === 0}
      <div class="muted">加载中…</div>
    {:else if notebooks.length === 0}
      <div class="empty">暂无笔记本，点击「新建笔记本」创建第一个。</div>
    {:else}
      <div class="notebookGrid">
        {#each notebooks as nb (nb.id)}
          <div
            class="notebookCard"
            class:selected={selectedNotebook?.id === nb.id}
            role="button"
            tabindex="0"
            on:click={() => selectNotebook(nb)}
            on:keydown={(e) => e.key === "Enter" && selectNotebook(nb)}
          >
            <span
              class="dot"
              style="background: {nb.color || 'rgba(34,197,94,0.6)'}"
              title={nb.color || ""}
            ></span>
            <span class="name">{nb.name}</span>
            {#if typeof nb.note_count === "number"}
              <span class="count">{nb.note_count} 条笔记</span>
            {/if}
            {#if editId === nb.id}
              <span
                class="editRow"
                role="group"
                on:click|stopPropagation
                on:keydown|stopPropagation
              >
                <input
                  class="input"
                  bind:value={editName}
                  placeholder="名称"
                  on:keydown={(e) => e.key === "Enter" && saveEdit()}
                />
                <div class="colorRow">
                  {#each PRESET_COLORS as c}
                    <button
                      type="button"
                      class="colorDot"
                      class:active={editColor === c}
                      style="background: {c}"
                      on:click={() => (editColor = c)}
                      title={c}
                    ></button>
                  {/each}
                </div>
                <div class="editActions">
                  <button class="mini" on:click={saveEdit}>保存</button>
                  <button class="mini" on:click={cancelEdit}>取消</button>
                </div>
              </span>
            {:else}
              <span
                class="actions"
                role="group"
                on:click|stopPropagation
                on:keydown|stopPropagation
              >
                <button class="mini" on:click={() => openEdit(nb)}>编辑</button>
                <button
                  class="mini danger"
                  on:click={() => deleteNotebook(nb.id)}>删除</button
                >
              </span>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>

  {#if selectedNotebook}
    <div class="card notesCard">
      <div class="row">
        <div class="title">
          <span
            class="dot small"
            style="background: {selectedNotebook.color ||
              'rgba(34,197,94,0.6)'}"
          ></span>
          {selectedNotebook.name} 中的笔记
        </div>
        <button class="btn ghost" on:click={() => selectNotebook(null)}
          >关闭</button
        >
      </div>
      {#if notesLoading}
        <div class="muted">加载中…</div>
      {:else if notebookNotes.length === 0}
        <div class="empty">该笔记本下暂无笔记。</div>
      {:else}
        <ul class="noteList">
          {#each notebookNotes as note (note.id)}
            <li class="noteItem">
              <a href="/?note={note.id}" class="noteLink">
                <span class="noteTitle">{note.title || "无标题"}</span>
                {#if note.tags?.length}
                  <span class="noteTags">
                    {#each note.tags as t}
                      <span
                        class="tag"
                        style="--tag-color: {t.color || '#22c55e'}"
                        >{t.name}</span
                      >
                    {/each}
                  </span>
                {/if}
                <span class="noteDate">{formatDate(note.updated_at)}</span>
              </a>
            </li>
          {/each}
        </ul>
      {/if}
    </div>
  {/if}
</div>

{#if createOpen}
  <div
    class="overlay"
    role="button"
    tabindex="0"
    on:click={(e) => e.target === e.currentTarget && (createOpen = false)}
    on:keydown={(e) => e.key === "Escape" && (createOpen = false)}
  >
    <div class="modal" role="dialog">
      <div class="modalTitle">新建笔记本</div>
      <label class="label" for="new-nb-name">名称</label>
      <input
        id="new-nb-name"
        class="input"
        aria-describedby="new-nb-name-hint"
        bind:value={newName}
        placeholder="例如：工作、学习、生活"
      />
      <label class="label" for="new-nb-color">颜色</label>
      <div class="colorRow" role="group" aria-label="选择颜色">
        {#each PRESET_COLORS as c}
          <button
            type="button"
            class="colorDot"
            class:active={newColor === c}
            style="background: {c}"
            on:click={() => (newColor = c)}
            title={c}
          ></button>
        {/each}
      </div>
      <input
        id="new-nb-color"
        type="text"
        class="input"
        bind:value={newColor}
        placeholder="#22c55e"
      />
      <div class="modalActions">
        <button class="btn ghost" on:click={() => (createOpen = false)}
          >取消</button
        >
        <button class="btn" on:click={createNotebook}>创建</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .wrap {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 24px 16px;
    max-width: 900px;
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
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .link {
    font-size: 12px;
    text-decoration: underline;
    color: inherit;
  }
  .hint {
    font-size: 13px;
    color: var(--muted);
    margin-bottom: 12px;
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
  .btn.ghost {
    border-color: var(--border);
    background: var(--panel);
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
  .notebookGrid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 12px;
  }
  .notebookCard {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 12px 14px;
    background: var(--panel-2);
    cursor: pointer;
  }
  .notebookCard:hover {
    filter: brightness(1.02);
  }
  .notebookCard.selected {
    border-color: var(--accent);
    background: var(--accent-soft);
  }
  .notebookCard .dot {
    width: 12px;
    height: 12px;
    border-radius: 999px;
    flex-shrink: 0;
  }
  .notebookCard .name {
    font-weight: 600;
    flex: 1;
    min-width: 0;
  }
  .notebookCard .count {
    font-size: 12px;
    color: var(--muted);
  }
  .notebookCard .editRow {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px solid var(--border);
  }
  .colorRow {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
  }
  .colorDot {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    border: 2px solid transparent;
    cursor: pointer;
    padding: 0;
  }
  .colorDot.active {
    border-color: var(--text);
    box-shadow: 0 0 0 1px var(--border);
  }
  .editActions {
    display: flex;
    gap: 8px;
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
  }
  .mini.danger {
    border-color: rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
  }
  .input {
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 8px 12px;
    font-size: 14px;
  }
  .notesCard .title .dot.small {
    width: 10px;
    height: 10px;
    border-radius: 50%;
  }
  .noteList {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .noteItem {
    border-bottom: 1px solid var(--border);
  }
  .noteItem:last-child {
    border-bottom: none;
  }
  .noteLink {
    display: block;
    padding: 12px 0;
    color: inherit;
    text-decoration: none;
  }
  .noteLink:hover {
    text-decoration: underline;
  }
  .noteTitle {
    font-weight: 600;
    display: block;
    margin-bottom: 4px;
  }
  .noteTags {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
    margin-bottom: 4px;
  }
  .tag {
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 6px;
    background: color-mix(in srgb, var(--tag-color) 25%, transparent);
    color: var(--tag-color);
  }
  .noteDate {
    font-size: 12px;
    color: var(--muted);
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
  .modal {
    width: 100%;
    max-width: 400px;
    border-radius: 14px;
    border: 1px solid var(--border);
    background: var(--panel);
    padding: 16px;
  }
  .modalTitle {
    font-weight: 700;
    margin-bottom: 12px;
  }
  .label {
    display: block;
    font-size: 12px;
    color: var(--muted);
    margin-top: 10px;
    margin-bottom: 6px;
  }
  .modalActions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 16px;
  }
</style>

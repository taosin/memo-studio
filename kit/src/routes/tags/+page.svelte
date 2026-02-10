<script>
  import { api } from "$lib/api.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let tags = [];
  let loading = false;
  let error = "";
  let toast = "";
  let newName = "";
  let newColor = "";
  let createOpen = false;
  let editId = null;
  let editName = "";
  let editColor = "";
  let mergeSourceId = null;
  let mergeTargetId = null;
  let mergeOpen = false;

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ""), 1500);
  }

  async function load() {
    loading = true;
    error = "";
    try {
      tags = await api.listTags(true);
      if (!Array.isArray(tags)) tags = [];
    } catch (e) {
      error = e?.message || "加载失败";
      if (String(error).includes("401") || String(error).includes("未登录")) {
        await goto("/login");
      }
    } finally {
      loading = false;
    }
  }

  async function createTag() {
    const name = String(newName || "").trim();
    if (!name) {
      setToast("请输入标签名称");
      return;
    }
    try {
      await api.createTag({ name, color: newColor || "" });
      setToast("已创建");
      newName = "";
      newColor = "";
      createOpen = false;
      await load();
    } catch (e) {
      setToast(e?.message || "创建失败");
    }
  }

  function openEdit(tag) {
    editId = tag.id;
    editName = tag.name;
    editColor = tag.color || "";
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
      setToast("标签名称不能为空");
      return;
    }
    try {
      await api.updateTag(editId, { name, color: editColor });
      setToast("已保存");
      cancelEdit();
      await load();
    } catch (e) {
      setToast(e?.message || "保存失败");
    }
  }

  async function deleteTag(id) {
    if (!confirm("确定删除该标签吗？关联的笔记将保留，仅移除标签。")) return;
    try {
      await api.deleteTag(id);
      setToast("已删除");
      if (editId === id) cancelEdit();
      await load();
    } catch (e) {
      setToast(e?.message || "删除失败");
    }
  }

  function openMerge() {
    mergeSourceId = null;
    mergeTargetId = null;
    mergeOpen = true;
  }

  function closeMerge() {
    mergeOpen = false;
    mergeSourceId = null;
    mergeTargetId = null;
  }

  async function doMerge() {
    const src = mergeSourceId != null ? Number(mergeSourceId) : null;
    const tgt = mergeTargetId != null ? Number(mergeTargetId) : null;
    if (src == null || tgt == null || src === tgt) {
      setToast("请选择两个不同的标签");
      return;
    }
    try {
      await api.mergeTags(src, tgt);
      setToast("合并成功");
      closeMerge();
      await load();
    } catch (e) {
      setToast(e?.message || "合并失败");
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>标签库 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">标签库</div>
      <a class="link" href="/">返回首页</a>
    </div>

    <div class="toolbar">
      <button class="btn" on:click={() => (createOpen = true)}>新建标签</button>
      <button class="btn ghost" on:click={openMerge}>合并标签</button>
      <span class="muted">共 {tags.length} 个标签</span>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}
    {#if toast}
      <div class="toast">{toast}</div>
    {/if}

    {#if loading && tags.length === 0}
      <div class="muted">加载中…</div>
    {:else if tags.length === 0}
      <div class="empty">
        暂无标签，在写笔记时使用 #标签名 即可自动创建，或点击「新建标签」。
      </div>
    {:else}
      <div class="list">
        {#each tags as tag (tag.id)}
          <div class="item">
            {#if editId === tag.id}
              <div class="editRow">
                <input
                  class="input"
                  bind:value={editName}
                  placeholder="标签名称"
                  on:keydown={(e) => e.key === "Enter" && saveEdit()}
                />
                <input
                  class="input color"
                  type="text"
                  bind:value={editColor}
                  placeholder="颜色 #hex"
                />
                <button class="mini" on:click={saveEdit}>保存</button>
                <button class="mini" on:click={cancelEdit}>取消</button>
              </div>
            {:else}
              <span
                class="dot"
                style="background: {tag.color || 'rgba(34,197,94,0.6)'}"
                title={tag.color || ""}
              ></span>
              <span class="name">{tag.name}</span>
              {#if typeof tag.note_count === "number"}
                <span class="count">{tag.note_count} 条笔记</span>
              {/if}
              <div class="actions">
                <button class="mini" on:click={() => openEdit(tag)}>编辑</button
                >
                <button class="mini danger" on:click={() => deleteTag(tag.id)}
                  >删除</button
                >
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
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
      <div class="modalTitle">新建标签</div>
      <label class="label" for="new-tag-name">名称</label>
      <input
        id="new-tag-name"
        class="input"
        bind:value={newName}
        placeholder="标签名称"
      />
      <label class="label" for="new-tag-color">颜色（可选，如 #22c55e）</label>
      <input
        id="new-tag-color"
        class="input"
        bind:value={newColor}
        placeholder="#22c55e"
      />
      <div class="modalActions">
        <button class="btn ghost" on:click={() => (createOpen = false)}
          >取消</button
        >
        <button class="btn" on:click={createTag}>创建</button>
      </div>
    </div>
  </div>
{/if}

{#if mergeOpen}
  <div
    class="overlay"
    role="button"
    tabindex="0"
    on:click={(e) => e.target === e.currentTarget && closeMerge()}
    on:keydown={(e) => e.key === "Escape" && closeMerge()}
  >
    <div class="modal" role="dialog">
      <div class="modalTitle">合并标签</div>
      <p class="hint">
        将「源标签」下的笔记全部归入「目标标签」，然后删除源标签。
      </p>
      <label class="label" for="merge-source">源标签（被合并）</label>
      <select id="merge-source" class="input" bind:value={mergeSourceId}>
        <option value={null}>请选择</option>
        {#each tags as t (t.id)}
          <option value={t.id}>{t.name}</option>
        {/each}
      </select>
      <label class="label" for="merge-target">目标标签（保留）</label>
      <select id="merge-target" class="input" bind:value={mergeTargetId}>
        <option value={null}>请选择</option>
        {#each tags as t (t.id)}
          <option value={t.id}>{t.name}</option>
        {/each}
      </select>
      <div class="modalActions">
        <button class="btn ghost" on:click={closeMerge}>取消</button>
        <button class="btn" on:click={doMerge}>合并</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .wrap {
    display: flex;
    justify-content: center;
    padding: 24px 16px;
  }
  .card {
    width: 100%;
    max-width: 640px;
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
  .list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .item {
    display: flex;
    align-items: center;
    gap: 10px;
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 10px 12px;
    background: var(--panel-2);
  }
  .dot {
    width: 12px;
    height: 12px;
    border-radius: 999px;
    flex-shrink: 0;
  }
  .name {
    font-weight: 600;
    flex: 1;
    min-width: 0;
  }
  .count {
    font-size: 12px;
    color: var(--muted);
  }
  .actions {
    display: flex;
    gap: 8px;
  }
  .editRow {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
    width: 100%;
  }
  .editRow .input {
    flex: 1;
    min-width: 100px;
  }
  .editRow .input.color {
    max-width: 120px;
  }
  .input {
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 8px 12px;
    font-size: 14px;
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
  .hint {
    font-size: 12px;
    color: var(--muted);
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

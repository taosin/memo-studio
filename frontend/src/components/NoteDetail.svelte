<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { createEventDispatcher } from 'svelte';

  export let noteId;
  const dispatch = createEventDispatcher();

  let note = null;
  let loading = true;
  let error = null;

  onMount(async () => {
    await loadNote();
  });

  async function loadNote() {
    try {
      loading = true;
      note = await api.getNote(noteId);
      error = null;
    } catch (err) {
      error = err.message;
      console.error('加载笔记失败:', err);
    } finally {
      loading = false;
    }
  }

  function handleEdit() {
    dispatch('edit', note);
  }

  function handleBack() {
    dispatch('back');
  }
</script>

<div class="note-detail-container">
  {#if loading}
    <div class="loading">加载中...</div>
  {:else if error}
    <div class="error">错误: {error}</div>
  {:else if note}
    <div class="note-detail">
      <div class="detail-header">
        <button class="btn btn-back" on:click={handleBack}>← 返回</button>
        <button class="btn btn-edit" on:click={handleEdit}>编辑</button>
      </div>

      <h1 class="detail-title">{note.title || '无标题'}</h1>
      
      <div class="detail-meta">
        <div class="tags">
          {#each note.tags || [] as tag}
            <span class="tag" style="background-color: {tag.color || '#4ECDC4'}20; color: {tag.color || '#4ECDC4'}">
              {tag.name}
            </span>
          {/each}
        </div>
        <span class="date">
          {new Date(note.created_at).toLocaleString('zh-CN')}
        </span>
      </div>

      <div class="detail-content">
        {note.content}
      </div>
    </div>
  {/if}
</div>

<style>
  .note-detail-container {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
  }

  .loading,
  .error {
    text-align: center;
    padding: 3rem;
    color: var(--text-secondary);
  }

  .note-detail {
    background-color: var(--card-bg);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: var(--shadow);
  }

  .detail-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2rem;
  }

  .btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 8px;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
    font-weight: 500;
  }

  .btn-back {
    background-color: var(--bg-secondary);
    color: var(--text-primary);
  }

  .btn-back:hover {
    background-color: var(--border-color);
  }

  .btn-edit {
    background-color: var(--accent-color);
    color: white;
  }

  .btn-edit:hover {
    background-color: var(--accent-hover);
  }

  .detail-title {
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
    line-height: 1.3;
  }

  .detail-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid var(--border-color);
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .tag {
    padding: 0.4rem 0.9rem;
    border-radius: 12px;
    font-size: 0.85rem;
    font-weight: 500;
  }

  .date {
    font-size: 0.9rem;
    color: var(--text-secondary);
  }

  .detail-content {
    color: var(--text-primary);
    line-height: 1.8;
    font-size: 1rem;
    white-space: pre-wrap;
    word-wrap: break-word;
  }

  @media (max-width: 768px) {
    .note-detail {
      padding: 1.5rem;
    }

    .detail-title {
      font-size: 1.5rem;
    }

    .detail-meta {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
    }
  }
</style>

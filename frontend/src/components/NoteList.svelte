<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import NoteCard from './NoteCard.svelte';
  import { api } from '../utils/api.js';

  const dispatch = createEventDispatcher();

  let notes = [];
  let loading = true;
  let error = null;

  onMount(async () => {
    await loadNotes();
  });

  async function loadNotes() {
    try {
      loading = true;
      notes = await api.getNotes();
      error = null;
    } catch (err) {
      error = err.message;
      console.error('加载笔记失败:', err);
    } finally {
      loading = false;
    }
  }

  function handleNoteClick(noteId) {
    dispatch('noteClick', noteId);
  }
</script>

<div class="note-list-container">
  {#if loading}
    <div class="loading">加载中...</div>
  {:else if error}
    <div class="error">错误: {error}</div>
  {:else if notes.length === 0}
    <div class="empty-state">
      <p>还没有笔记，创建第一个吧！</p>
    </div>
  {:else}
    <div class="notes-grid">
      {#each notes as note (note.id)}
        <NoteCard {note} on:click={() => handleNoteClick(note.id)} />
      {/each}
    </div>
  {/if}
</div>

<style>
  .note-list-container {
    width: 100%;
  }

  .loading,
  .error,
  .empty-state {
    text-align: center;
    padding: 3rem;
    color: var(--text-secondary);
  }

  .notes-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
  }

  @media (max-width: 768px) {
    .notes-grid {
      grid-template-columns: 1fr;
      gap: 1rem;
    }
  }
</style>

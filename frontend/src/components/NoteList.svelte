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

<div class="w-full">
  {#if loading}
    <div class="text-center py-12 text-muted-foreground">加载中...</div>
  {:else if error}
    <div class="text-center py-12 text-destructive">错误: {error}</div>
  {:else if notes.length === 0}
    <div class="text-center py-12 text-muted-foreground">
      <p>还没有笔记，创建第一个吧！</p>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each notes as note (note.id)}
        <NoteCard {note} on:click={() => handleNoteClick(note.id)} />
      {/each}
    </div>
  {/if}
</div>

<script>
  import { onMount } from 'svelte';
  import { themeStore } from './stores/theme.js';
  import NoteList from './components/NoteList.svelte';
  import NoteDetail from './components/NoteDetail.svelte';
  import NoteEditor from './components/NoteEditor.svelte';
  import ThemeToggle from './components/ThemeToggle.svelte';

  let currentView = 'list'; // 'list', 'detail', 'editor'
  let selectedNoteId = null;
  let editingNote = null;
  let listKey = 0; // 用于强制刷新列表

  $: theme = $themeStore;

  onMount(() => {
    // 应用主题
    document.documentElement.setAttribute('data-theme', theme);
  });

  $: {
    document.documentElement.setAttribute('data-theme', theme);
  }

  function handleNoteClick(noteId) {
    selectedNoteId = noteId;
    currentView = 'detail';
  }

  function handleNewNote() {
    editingNote = null;
    currentView = 'editor';
  }

  function handleEditNote(note) {
    editingNote = note;
    currentView = 'editor';
  }

  function handleBack() {
    currentView = 'list';
    selectedNoteId = null;
    editingNote = null;
  }

  function handleSave() {
    currentView = 'list';
    editingNote = null;
    listKey++; // 触发列表刷新
  }
</script>

<div class="app-container">
  <header class="header">
    <div class="header-content">
      <h1 class="logo" on:click={handleBack}>📝 Memo Studio</h1>
      <div class="header-actions">
        {#if currentView === 'list'}
          <button class="btn btn-primary" on:click={handleNewNote}>
            + 新建笔记
          </button>
        {/if}
        <ThemeToggle />
      </div>
    </div>
  </header>

  <main class="main-content">
    {#if currentView === 'list'}
      <NoteList key={listKey} on:noteClick={(e) => handleNoteClick(e.detail)} />
    {:else if currentView === 'detail'}
      <NoteDetail 
        noteId={selectedNoteId} 
        on:back={handleBack}
        on:edit={(e) => handleEditNote(e.detail)}
      />
    {:else if currentView === 'editor'}
      <NoteEditor 
        note={editingNote}
        on:save={handleSave}
        on:cancel={handleBack}
      />
    {/if}
  </main>
</div>

<style>
  .app-container {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: var(--bg-primary);
  }

  .header {
    background-color: var(--card-bg);
    border-bottom: 1px solid var(--border-color);
    padding: 1rem 0;
    position: sticky;
    top: 0;
    z-index: 100;
    box-shadow: var(--shadow);
  }

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .logo {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--text-primary);
    cursor: pointer;
    user-select: none;
  }

  .header-actions {
    display: flex;
    gap: 1rem;
    align-items: center;
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

  .btn-primary {
    background-color: var(--accent-color);
    color: white;
  }

  .btn-primary:hover {
    background-color: var(--accent-hover);
    transform: translateY(-1px);
  }

  .main-content {
    flex: 1;
    max-width: 1200px;
    width: 100%;
    margin: 0 auto;
    padding: 2rem 1.5rem;
  }

  @media (max-width: 768px) {
    .header-content {
      padding: 0 1rem;
    }

    .logo {
      font-size: 1.2rem;
    }

    .main-content {
      padding: 1rem;
    }

    .btn {
      padding: 0.4rem 0.8rem;
      font-size: 0.85rem;
    }
  }
</style>

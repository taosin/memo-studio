<script>
  import { onMount, onDestroy } from 'svelte';
  import { themeStore } from './stores/theme.js';
  import { authStore } from './stores/auth.js';
  import LoginPage from './components/LoginPage.svelte';
  import NoteList from './components/NoteList.svelte';
  import NoteDetail from './components/NoteDetail.svelte';
  import NoteEditor from './components/NoteEditor.svelte';
  import ProfilePage from './components/ProfilePage.svelte';
  import ThemeToggle from './components/ThemeToggle.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import { api } from './utils/api.js';

  let currentView = 'list';
  let selectedNoteId = null;
  let editingNote = null;
  let listKey = 0;
  let notes = [];
  let previousView = '';

  $: isAuthenticated = $authStore.isAuthenticated;

  function handleAuthSuccess() {
    currentView = 'list';
  }

  onMount(() => {
    window.addEventListener('auth-success', handleAuthSuccess);
    
    if ($authStore.isAuthenticated) {
      verifyToken();
    }
  });

  onDestroy(() => {
    window.removeEventListener('auth-success', handleAuthSuccess);
  });

  async function verifyToken() {
    try {
      const user = await api.getCurrentUser();
      authStore.setUser(user);
    } catch (err) {
      authStore.logout();
    }
  }

  function handleNoteClick(noteId) {
    previousView = currentView;
    selectedNoteId = noteId;
    currentView = 'detail';
  }

  function handleNewNote() {
    previousView = currentView;
    editingNote = null;
    currentView = 'editor';
  }

  function handleEditNote(note) {
    previousView = currentView;
    editingNote = note;
    currentView = 'editor';
  }

  function handleBack() {
    currentView = 'list';
    selectedNoteId = null;
    editingNote = null;
  }

  function handleProfile() {
    previousView = currentView;
    currentView = 'profile';
  }

  function handleLogout() {
    authStore.logout();
    currentView = 'login';
  }

  function handleSave() {
    currentView = 'list';
    editingNote = null;
    listKey++;
  }

  async function handleQuickEdit(noteId) {
    try {
      const note = await api.getNote(noteId);
      handleEditNote(note);
    } catch (err) {
      console.error('获取笔记失败:', err);
    }
  }
</script>

{#if !isAuthenticated}
  <LoginPage />
{:else}
  <div class="min-h-screen flex flex-col bg-background">
    <header class="sticky top-0 z-50 w-full border-b bg-card/80 backdrop-blur-md">
      <div class="container mx-auto px-4">
        <div class="flex h-14 sm:h-16 items-center justify-between">
          <button
            class="text-xl sm:text-2xl font-bold cursor-pointer select-none bg-transparent border-none p-0 text-left flex items-center gap-2 hover:opacity-80 transition-opacity"
            on:click={handleBack}
          >
            <span class="text-2xl">📝</span>
            <span class="hidden sm:inline bg-gradient-to-r from-primary to-primary-light bg-clip-text text-transparent">
              Memo Studio
            </span>
          </button>
          <div class="flex items-center gap-2 sm:gap-4">
            {#if currentView === 'list'}
              <Button on:click={handleNewNote} size="sm" class="text-xs sm:text-sm shadow-sm hover:shadow-md transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-1">
                  <line x1="12" y1="5" x2="12" y2="19"></line>
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
                新建
              </Button>
              <Button variant="ghost" size="sm" on:click={handleProfile} class="hover:bg-accent transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
              </Button>
            {/if}
            <ThemeToggle />
          </div>
        </div>
      </div>
    </header>

    <main class="flex-1 container mx-auto px-4 py-4 max-w-[1400px]">
      {#if currentView === 'list'}
        <NoteList 
          key={listKey} 
          on:noteClick={(e) => handleNoteClick(e.detail)}
          onQuickEdit={handleQuickEdit}
        />
      {:else if currentView === 'detail'}
        <div class="animate-fade-in">
          <NoteDetail 
            noteId={selectedNoteId} 
            on:back={handleBack}
            on:edit={(e) => handleEditNote(e.detail)}
            on:deleted={() => {
              listKey++;
              handleBack();
            }}
          />
        </div>
      {:else if currentView === 'editor'}
        <div class="animate-scale-in">
          <NoteEditor 
            note={editingNote}
            on:save={handleSave}
            on:cancel={handleBack}
          />
        </div>
      {:else if currentView === 'profile'}
        <div class="animate-fade-in">
          <ProfilePage on:logout={handleLogout} />
        </div>
      {/if}
    </main>
  </div>
{/if}

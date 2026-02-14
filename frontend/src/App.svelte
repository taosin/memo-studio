<script>
  import { onMount, onDestroy } from 'svelte';
  import { themeStore } from './stores/theme.js';
  import { authStore } from './stores/auth.js';
  import LoginPage from './components/LoginPage.svelte';
  import NoteList from './components/NoteList.svelte';
  import NoteDetail from './components/NoteDetail.svelte';
  import FlomoEditor from './components/FlomoEditor.svelte';
  import ProfilePage from './components/ProfilePage.svelte';
  import ThemeToggle from './components/ThemeToggle.svelte';
  import KeyboardManager from './components/KeyboardManager.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import { api } from './utils/api.js';
  import { exportData } from './utils/exportImport.js';
  import { saveDraft } from './utils/backup.js';

  let currentView = 'list';
  let selectedNoteId = null;
  let editingNote = null;
  let listKey = 0;
  let showEditor = false;
  let editorMode = 'create';
  let viewMode = 'timeline';
  let sidebarCollapsed = false;
  let keyboardManager;

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
    selectedNoteId = noteId;
    currentView = 'detail';
  }

  function handleNewNote() {
    editingNote = null;
    editorMode = 'create';
    showEditor = true;
  }

  function handleEditNote(note) {
    editingNote = note;
    editorMode = 'edit';
    showEditor = true;
  }

  function handleBack() {
    currentView = 'list';
    selectedNoteId = null;
    editingNote = null;
  }

  function handleProfile() {
    currentView = 'profile';
  }

  function handleLogout() {
    authStore.logout();
    currentView = 'login';
  }

  function handleSave() {
    showEditor = false;
    editingNote = null;
    listKey++;
    if (currentView === 'detail') {
      currentView = 'list';
      selectedNoteId = null;
    }
  }

  function handleEditorCancel() {
    showEditor = false;
    editingNote = null;
  }

  async function handleQuickEdit(noteId) {
    try {
      const note = await api.getNote(noteId);
      handleEditNote(note);
    } catch (err) {
      console.error('获取笔记失败:', err);
    }
  }

  // 键盘事件处理
  function handleFocusSearch() {
    // 触发搜索栏聚焦
    const searchInput = document.querySelector('.search-bar input');
    if (searchInput) searchInput.focus();
  }

  function handleSaveEditor() {
    // 触发编辑器保存
    const saveBtn = document.querySelector('.flomo-editor .save-btn');
    if (saveBtn) saveBtn.click();
  }

  function handleCloseAll() {
    if (showEditor) {
      handleEditorCancel();
    } else if (currentView !== 'list') {
      handleBack();
    }
  }

  function handleNewNoteKey() {
    handleNewNote();
  }

  function handleEdit() {
    // 编辑当前选中的笔记
    if (selectedNoteId) {
      handleQuickEdit(selectedNoteId);
    }
  }

  function handleDelete() {
    // 删除当前笔记
    if (selectedNoteId && confirm('确定删除此笔记吗？')) {
      api.deleteNote(selectedNoteId).then(() => {
        listKey++;
        handleBack();
      });
    }
  }

  function handleNavigate(e) {
    // 笔记列表导航
    const direction = e.detail.direction;
    const notes = document.querySelectorAll('.note-card, .note-item');
    if (notes.length === 0) return;
    
    const active = document.activeElement;
    let index = -1;
    notes.forEach((note, i) => {
      if (note === active || note.contains(active)) index = i;
    });
    
    if (direction === 'down') {
      index = Math.min(index + 1, notes.length - 1);
    } else {
      index = Math.max(index - 1, 0);
    }
    
    if (notes[index]) {
      notes[index].focus();
      if (notes[index].scrollIntoView) {
        notes[index].scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
    }
  }

  function handleChangeView(e) {
    viewMode = e.detail.mode;
  }

  function handleToggleSidebar() {
    sidebarCollapsed = !sidebarCollapsed;
  }

  function handleShowTags() {
    // 显示标签面板
    const tagButton = document.querySelector('.tag-tree-toggle');
    if (tagButton) tagButton.click();
  }

  function handleFocusTagSearch() {
    const tagSearch = document.querySelector('.tag-search input');
    if (tagSearch) tagSearch.focus();
  }

  function handleImport() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.json,.md,.txt';
    fileInput.onchange = async (e) => {
      const file = e.target.files[0];
      if (file) {
        try {
          const text = await file.text();
          const data = JSON.parse(text);
          if (data.notes && Array.isArray(data.notes)) {
            for (const note of data.notes) {
              await api.createNote(note);
            }
            listKey++;
            alert(`成功导入 ${data.notes.length} 条笔记`);
          }
        } catch (err) {
          alert('导入失败: ' + err.message);
        }
      }
    };
    fileInput.click();
  }
</script>

{#if !isAuthenticated}
  <LoginPage />
{:else}
  <div class="min-h-screen flex flex-col bg-background">
    <!-- 键盘管理器 -->
    <KeyboardManager
      bind:this={keyboardManager}
      on:focusSearch={handleFocusSearch}
      on:saveEditor={handleSaveEditor}
      on:closeAll={handleCloseAll}
      on:newNote={handleNewNoteKey}
      on:edit={handleEdit}
      on:delete={handleDelete}
      on:navigate={handleNavigate}
      on:changeView={handleChangeView}
      on:toggleSidebar={handleToggleSidebar}
      on:showTags={handleShowTags}
      on:focusTagSearch={handleFocusTagSearch}
      on:import={handleImport}
    />

    <header class="sticky top-0 z-40 w-full border-b bg-card/80 backdrop-blur-md transition-all duration-300">
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
            <Button variant="ghost" size="sm" on:click={handleProfile} class="hover:bg-accent transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </Button>
            <ThemeToggle />
          </div>
        </div>
      </div>
    </header>

    <main class="flex-1 container mx-auto px-4 py-4 max-w-[1400px] pb-24">
      {#if currentView === 'list'}
        <NoteList 
          key={listKey} 
          {viewMode}
          {sidebarCollapsed}
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
      {:else if currentView === 'profile'}
        <div class="animate-fade-in">
          <ProfilePage on:logout={handleLogout} />
        </div>
      {/if}
    </main>

    <!-- Flomo 风格底部编辑器 -->
    {#if showEditor}
      <FlomoEditor 
        note={editingNote}
        mode={editorMode}
        on:save={handleSave}
        on:cancel={handleEditorCancel}
      />
    {:else}
      <!-- 底部浮动按钮 -->
      <div class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 animate-fade-in">
        <button
          on:click={handleNewNote}
          class="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-primary to-primary-light text-primary-foreground rounded-full shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 hover:scale-105 active:scale-95 transition-all duration-300 group"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:rotate-90 transition-transform duration-300">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          <span class="font-medium">记录灵感</span>
        </button>
      </div>
    {/if}
  </div>
{/if}

<style>
  :global(.note-card:focus, .note-item:focus) {
    outline: 2px solid hsl(var(--primary));
    outline-offset: 2px;
  }
</style>

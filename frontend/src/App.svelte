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

  let currentView = 'list'; // 'list', 'detail', 'editor', 'profile'
  let selectedNoteId = null;
  let editingNote = null;
  let listKey = 0; // 用于强制刷新列表
  let notes = []; // 用于快速编辑

  $: isAuthenticated = $authStore.isAuthenticated;

  function handleAuthSuccess() {
    // 登录成功后刷新页面或跳转到列表页
    currentView = 'list';
  }

  onMount(() => {
    // 监听登录成功事件
    window.addEventListener('auth-success', handleAuthSuccess);
    
    // 检查是否已登录
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
      // Token 无效，清除认证信息
      authStore.logout();
    }
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

  function handleProfile() {
    currentView = 'profile';
  }

  function handleLogout() {
    authStore.logout();
    currentView = 'login';
  }

  function handleSave() {
    currentView = 'list';
    editingNote = null;
    listKey++; // 触发列表刷新
  }

  async function handleQuickEdit(noteId) {
    // 快速编辑：直接进入编辑模式
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
    <header class="sticky top-0 z-50 w-full border-b bg-card">
      <div class="container mx-auto px-4">
        <div class="flex h-14 sm:h-16 items-center justify-between">
          <button
            class="text-xl sm:text-2xl font-semibold cursor-pointer select-none bg-transparent border-none p-0 text-left"
            on:click={handleBack}
          >
            📝 Memo Studio
          </button>
          <div class="flex items-center gap-2 sm:gap-4">
            {#if currentView === 'list'}
              <Button on:click={handleNewNote} size="sm" class="text-xs sm:text-sm">+ 新建</Button>
              <Button variant="ghost" size="sm" on:click={handleProfile}>
                👤
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
        <NoteDetail 
          noteId={selectedNoteId} 
          on:back={handleBack}
          on:edit={(e) => handleEditNote(e.detail)}
          on:deleted={() => {
            listKey++;
            handleBack();
          }}
        />
      {:else if currentView === 'editor'}
        <NoteEditor 
          note={editingNote}
          on:save={handleSave}
          on:cancel={handleBack}
        />
      {:else if currentView === 'profile'}
        <ProfilePage on:logout={handleLogout} />
      {/if}
    </main>
  </div>
{/if}

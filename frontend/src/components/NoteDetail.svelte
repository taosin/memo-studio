<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { createEventDispatcher } from 'svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  export let noteId;
  const dispatch = createEventDispatcher();

  let note = null;
  let loading = true;
  let error = null;
  let isHovered = false;

  onMount(async () => {
    await loadNote();
  });

  async function loadNote() {
    try {
      loading = true;
      const loadedNote = await api.getNote(noteId);
      if (typeof loadedNote.content !== 'string') {
        console.warn('NoteDetail - content 不是字符串:', typeof loadedNote.content, loadedNote.content);
        loadedNote.content = typeof loadedNote.content === 'object' 
          ? JSON.stringify(loadedNote.content) 
          : String(loadedNote.content || '');
      }
      note = loadedNote;
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

  async function handleDelete() {
    if (!confirm('确定要删除这条笔记吗？此操作不可恢复。')) {
      return;
    }

    try {
      await api.deleteNote(noteId);
      dispatch('deleted');
      dispatch('back');
    } catch (err) {
      alert('删除失败: ' + err.message);
    }
  }

  function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric',
      weekday: 'long',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function getRelativeTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return '今天';
    if (diffDays === 1) return '昨天';
    if (diffDays < 7) return `${diffDays}天前`;
    return formatDate(dateString);
  }
</script>

<div class="w-full max-w-4xl mx-auto animate-fade-in">
  <!-- 加载状态 -->
  {#if loading}
    <div class="flex flex-col items-center justify-center py-20">
      <div class="w-12 h-12 border-4 border-primary/20 border-t-primary rounded-full animate-spin mb-4"></div>
      <p class="text-muted-foreground">加载中...</p>
    </div>
  {:else if error}
    <!-- 错误状态 -->
    <div class="flex flex-col items-center justify-center py-20 text-center animate-fade-in">
      <div class="w-16 h-16 bg-destructive/10 rounded-full flex items-center justify-center mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-destructive">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
      </div>
      <h3 class="text-lg font-semibold mb-2">加载失败</h3>
      <p class="text-muted-foreground mb-4">{error}</p>
      <Button on:click={loadNote}>重新加载</Button>
    </div>
  {:else if note}
    <!-- 笔记详情 -->
    <div class="space-y-6">
      <!-- 顶部导航 -->
      <div class="flex items-center justify-between">
        <button 
          on:click={handleBack}
          class="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors group"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform">
            <path d="M19 12H5M12 19l-7-7 7-7"/>
          </svg>
          <span>返回列表</span>
        </button>
        
        <div class="flex items-center gap-3">
          <Button 
            variant="outline" 
            on:click={handleEdit}
            class="shadow-sm"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            编辑
          </Button>
          <Button 
            variant="destructive" 
            on:click={handleDelete}
            class="shadow-sm"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              <line x1="10" y1="11" x2="10" y2="17"/>
              <line x1="14" y1="11" x2="14" y2="17"/>
            </svg>
            删除
          </Button>
        </div>
      </div>

      <!-- 笔记卡片 -->
      <Card 
        class="border-2 shadow-lg overflow-hidden {isHovered ? 'shadow-xl' : ''}"
        on:mouseenter={() => isHovered = true}
        on:mouseleave={() => isHovered = false}
      >
        <!-- 头部信息 -->
        <div class="px-6 py-4 border-b bg-gradient-to-r from-primary/5 to-transparent">
          <!-- 标题 -->
          <h1 class="text-3xl font-bold mb-4 text-foreground">
            {note.title || '无标题'}
          </h1>
          
          <!-- 元信息 -->
          <div class="flex flex-wrap items-center gap-4">
            <!-- 时间信息 -->
            <div class="flex items-center gap-2 text-sm text-muted-foreground">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              <span>{getRelativeTime(note.created_at)}</span>
              <span class="text-xs opacity-50">创建</span>
            </div>
            
            {#if note.updated_at && note.updated_at !== note.created_at}
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
                <span>{getRelativeTime(note.updated_at)}</span>
                <span class="text-xs opacity-50">更新</span>
              </div>
            {/if}
          </div>

          <!-- 标签 -->
          {#if note.tags && note.tags.length > 0}
            <div class="flex flex-wrap gap-2 mt-4">
              {#each note.tags as tag}
                <Badge 
                  class="px-3 py-1 text-sm font-medium"
                  style="background-color: {tag.color || '#4ECDC4'}20; border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="currentColor" class="mr-1">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                  {tag.name}
                </Badge>
              {/each}
            </div>
          {/if}
        </div>

        <!-- 内容区域 -->
        <CardContent class="p-6">
          <div class="prose prose-lg dark:prose-invert max-w-none break-words">
            {@html note.content}
          </div>
        </CardContent>
      </Card>

      <!-- 底部操作提示 -->
      <div class="text-center text-sm text-muted-foreground">
        <p>💡 双击卡片可快速编辑</p>
      </div>
    </div>
  {/if}
</div>

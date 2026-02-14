<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import NoteCard from './NoteCard.svelte';
  import TagTree from './TagTree.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import { api } from '../utils/api.js';

  const dispatch = createEventDispatcher();
  export let onQuickEdit = null;

  let notes = [];
  let filteredNotes = [];
  let loading = true;
  let error = null;
  let searchQuery = '';
  let selectedTags = [];
  let viewMode = 'timeline';
  let selectedNoteIds = new Set();
  let randomNote = null;
  let showRandomNote = false;
  let sidebarCollapsed = false;
  let showFilters = false;

  onMount(async () => {
    await Promise.all([loadNotes(), loadRandomNote()]);
  });

  async function loadNotes() {
    try {
      loading = true;
      const data = await api.getNotes();
      notes = Array.isArray(data) ? data : [];
      filterNotes();
      error = null;
    } catch (err) {
      error = err.message || '加载笔记失败';
      notes = [];
      filteredNotes = [];
    } finally {
      loading = false;
    }
  }

  async function loadRandomNote() {
    try {
      const response = await fetch('/api/review/random', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      if (response.ok) {
        const data = await response.json();
        if (Array.isArray(data) && data.length > 0) {
          randomNote = data[0];
        }
      }
    } catch (err) {
      // 静默处理
    }
  }

  function filterNotes() {
    if (!Array.isArray(notes)) {
      filteredNotes = [];
      return;
    }
    
    let filtered = [...notes];

    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(note => 
        (note.title || '').toLowerCase().includes(query) ||
        (note.content || '').replace(/<[^>]*>/g, '').toLowerCase().includes(query)
      );
    }

    if (selectedTags.length > 0) {
      filtered = filtered.filter(note => {
        const noteTagIds = (note.tags || []).map(t => t.id);
        return selectedTags.some(tagId => noteTagIds.includes(tagId));
      });
    }

    filteredNotes = filtered;
  }

  function toggleNoteSelection(noteId) {
    if (selectedNoteIds.has(noteId)) {
      selectedNoteIds.delete(noteId);
    } else {
      selectedNoteIds.add(noteId);
    }
    selectedNoteIds = selectedNoteIds;
  }

  async function handleBatchDelete() {
    if (selectedNoteIds.size === 0) return;
    if (!confirm(`确定删除这 ${selectedNoteIds.size} 条笔记吗？`)) return;
    
    try {
      await api.deleteNotes(Array.from(selectedNoteIds));
      selectedNoteIds.clear();
      await loadNotes();
    } catch (err) {
      alert('删除失败');
    }
  }

  function handleSearch(e) {
    searchQuery = e.target.value;
    filterNotes();
  }

  function handleTagSelect(tag) {
    if (selectedTags.includes(tag.id)) {
      selectedTags = selectedTags.filter(id => id !== tag.id);
    } else {
      selectedTags = [...selectedTags, tag.id];
    }
    filterNotes();
  }

  function handleViewModeChange(mode) {
    viewMode = mode;
  }

  function handleNoteClick(noteId) {
    dispatch('noteClick', noteId);
  }

  function handleNoteDoubleClick(noteId) {
    if (onQuickEdit) onQuickEdit(noteId);
  }

  function clearFilters() {
    searchQuery = '';
    selectedTags = [];
    filterNotes();
  }

  function formatGroupDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = now - date;
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return '今天';
    if (diffDays === 1) return '昨天';
    if (diffDays < 7) return `${diffDays}天前`;
    return date.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' });
  }

  $: groupedNotes = (() => {
    const groups = {};
    filteredNotes.forEach(note => {
      const date = new Date(note.created_at).toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
      if (!groups[date]) groups[date] = [];
      groups[date].push(note);
    });
    
    const sortedDates = Object.keys(groups).sort((a, b) => {
      return new Date(b) - new Date(a);
    });
    
    return sortedDates.reduce((acc, key) => {
      acc[key] = groups[key];
      return acc;
    }, {});
  })();

  $: activeFilters = (searchQuery.trim() ? 1 : 0) + selectedTags.length;
</script>

<div class="flex min-h-screen">
  <!-- 侧边栏 -->
  <aside class="w-64 flex-shrink-0 border-r bg-card/50 transition-all duration-300 {sidebarCollapsed ? 'w-16' : ''}">
    <div class="sticky top-0 h-screen flex flex-col p-4">
      <!-- Logo / 品牌 -->
      <div class="flex items-center gap-3 mb-6">
        <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary to-primary-light flex items-center justify-center shadow-lg shadow-primary/20">
          <span class="text-xl">📝</span>
        </div>
        {#if !sidebarCollapsed}
          <span class="font-bold bg-gradient-to-r from-primary to-primary-light bg-clip-text text-transparent">Memo</span>
        {/if}
      </div>

      <!-- 搜索栏 -->
      {#if !sidebarCollapsed}
        <div class="relative mb-4">
          <input
            type="text"
            placeholder="搜索..."
            value={searchQuery}
            on:input={handleSearch}
            class="w-full h-10 pl-10 pr-4 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all"
          />
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </div>
      {/if}

      <!-- 标签管理 -->
      {#if !sidebarCollapsed}
        <div class="flex-1 overflow-y-auto">
          <TagTree 
            {selectedTags} 
            on:tagSelect={handleTagSelect}
          />
        </div>
      {/if}

      <!-- 底部操作 -->
      <div class="pt-4 border-t">
        {#if !sidebarCollapsed}
          <button 
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-muted-foreground hover:bg-accent hover:text-foreground transition-colors"
            on:click={() => sidebarCollapsed = !sidebarCollapsed}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="15 18 9 12 15 6"></polyline>
            </svg>
            <span>收起侧边栏</span>
          </button>
        {:else}
          <button 
            class="w-full flex justify-center py-2 rounded-lg text-muted-foreground hover:bg-accent hover:text-foreground transition-colors"
            on:click={() => sidebarCollapsed = !sidebarCollapsed}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
          </button>
        {/if}
      </div>
    </div>
  </aside>

  <!-- 主内容区 -->
  <main class="flex-1 min-w-0 pb-32">
    <!-- 顶部栏 -->
    <div class="sticky top-0 z-30 bg-background/80 backdrop-blur-md border-b">
      <div class="flex items-center justify-between px-6 py-4">
        <!-- 视图切换 -->
        <div class="flex items-center gap-1 bg-secondary/50 rounded-full p-1">
          <button
            class="px-4 py-1.5 rounded-full text-sm font-medium transition-all"
            class:bg-primary={viewMode === 'timeline'}
            class:text-primary-foreground={viewMode === 'timeline'}
            class:text-muted-foreground={viewMode !== 'timeline'}
            on:click={() => handleViewModeChange('timeline')}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="inline mr-1">
              <line x1="8" y1="6" x2="21" y2="6"></line>
              <line x1="8" y1="12" x2="21" y2="12"></line>
              <line x1="8" y1="18" x2="21" y2="18"></line>
              <line x1="3" y1="6" x2="3.01" y2="6"></line>
              <line x1="3" y1="12" x2="3.01" y2="12"></line>
              <line x1="3" y1="18" x2="3.01" y2="18"></line>
            </svg>
            信息流
          </button>
          <button
            class="px-4 py-1.5 rounded-full text-sm font-medium transition-all"
            class:bg-primary={viewMode === 'waterfall'}
            class:text-primary-foreground={viewMode === 'waterfall'}
            class:text-muted-foreground={viewMode !== 'waterfall'}
            on:click={() => handleViewModeChange('waterfall')}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="inline mr-1">
              <rect x="3" y="3" width="7" height="7"></rect>
              <rect x="14" y="3" width="7" height="7"></rect>
              <rect x="14" y="14" width="7" height="7"></rect>
              <rect x="3" y="14" width="7" height="7"></rect>
            </svg>
            卡片
          </button>
        </div>

        <!-- 右侧操作 -->
        <div class="flex items-center gap-3">
          {#if activeFilters > 0}
            <span class="text-sm text-muted-foreground">
              {filteredNotes.length} 条结果
            </span>
            <Button variant="ghost" size="sm" on:click={clearFilters}>
              清除筛选
            </Button>
          {:else}
            <span class="text-sm text-muted-foreground">
              {filteredNotes.length} 条笔记
            </span>
          {/if}

          {#if selectedNoteIds.size > 0}
            <span class="text-sm text-primary font-medium">
              已选 {selectedNoteIds.size}
            </span>
            <Button variant="destructive" size="sm" on:click={handleBatchDelete}>
              删除
            </Button>
          {/if}
        </div>
      </div>
    </div>

    <!-- 内容区 -->
    <div class="px-6 py-6">
      <!-- 加载状态 -->
      {#if loading}
        <div class="space-y-8">
          {#each Array(3) as _, i}
            <div class="animate-pulse">
              <div class="h-6 w-24 bg-muted rounded-full mb-4"></div>
              <div class="space-y-3">
                <div class="h-40 bg-muted rounded-2xl"></div>
                <div class="h-40 bg-muted rounded-2xl"></div>
              </div>
            </div>
          {/each}
        </div>

      <!-- 错误状态 -->
      {:else if error}
        <div class="flex flex-col items-center justify-center py-20 text-center">
          <div class="w-16 h-16 bg-destructive/10 rounded-full flex items-center justify-center mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-destructive">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
          </div>
          <h3 class="text-lg font-semibold mb-2">加载失败</h3>
          <p class="text-muted-foreground mb-4">{error}</p>
          <Button on:click={loadNotes}>重试</Button>
        </div>

      <!-- 空状态 -->
      {:else if filteredNotes.length === 0}
        <div class="flex flex-col items-center justify-center py-20 text-center animate-fade-in">
          <div class="w-24 h-24 bg-primary/5 rounded-full flex items-center justify-center mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="text-primary">
              <path d="M12 19l7-7 3 3-7 7-3-3z"></path>
              <path d="M18 13l-1.5-7.5L2 2l3.5 14.5L13 18l5-5z"></path>
              <path d="M2 2l7.586 7.586"></path>
              <circle cx="11" cy="11" r="2"></circle>
            </svg>
          </div>
          <h3 class="text-xl font-semibold mb-2">开始记录灵感</h3>
          <p class="text-muted-foreground mb-6 max-w-sm">
            点击底部的 + 按钮，记录你的第一条笔记
          </p>
          <div class="flex items-center gap-6 text-sm text-muted-foreground">
            <span class="flex items-center gap-1">
              <kbd class="px-2 py-0.5 bg-muted rounded text-xs">#</kbd>
              添加标签
            </span>
            <span class="flex items-center gap-1">
              <kbd class="px-2 py-0.5 bg-muted rounded text-xs">Ctrl</kbd>
              <kbd class="px-2 py-0.5 bg-muted rounded text-xs">Enter</kbd>
              快速保存
            </span>
          </div>
        </div>

      <!-- 笔记列表 -->
      {:else}
        <!-- 时间线视图 -->
        {#if viewMode === 'timeline'}
          <div class="space-y-1">
            {#each Object.entries(groupedNotes) as [date, dateNotes], index}
              <div class="animate-fade-in" style="animation-delay: {index * 80}ms">
                <!-- 日期标题 -->
                <div class="flex items-center gap-4 mb-6 mt-8 first:mt-0">
                  <div class="flex-1 h-px bg-gradient-to-r from-transparent via-border to-transparent"></div>
                  <div class="flex items-center gap-2 px-4 py-1.5 bg-primary/5 rounded-full">
                    <span class="text-sm font-medium text-primary">{formatGroupDate(date)}</span>
                    <span class="text-xs text-muted-foreground">({dateNotes.length})</span>
                  </div>
                  <div class="flex-1 h-px bg-gradient-to-r from-transparent via-border to-transparent"></div>
                </div>

                <!-- 笔记卡片流 -->
                <div class="space-y-3">
                  {#each dateNotes as note, noteIndex (note.id)}
                    <div 
                      class="group relative animate-slide-up"
                      style="animation-delay: {noteIndex * 40}ms"
                    >
                      <!-- 选择指示器 -->
                      <div 
                        class="absolute -left-3 top-4 w-8 h-8 rounded-full border-2 border-border bg-card flex items-center justify-center cursor-pointer opacity-0 group-hover:opacity-100 transition-all duration-200 shadow-sm"
                        class:opacity-100={selectedNoteIds.has(note.id)}
                        class:bg-primary={selectedNoteIds.has(note.id)}
                        class:border-primary={selectedNoteIds.has(note.id)}
                        on:click|stopPropagation={() => toggleNoteSelection(note.id)}
                      >
                        {#if selectedNoteIds.has(note.id)}
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" class="text-primary-foreground">
                            <polyline points="20 6 9 17 4 12"></polyline>
                          </svg>
                        {/if}
                      </div>

                      <div class="pl-8">
                        <NoteCard 
                          {note} 
                          on:click={() => handleNoteClick(note.id)}
                          on:doubleClick={() => handleNoteDoubleClick(note.id)}
                        />
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            {/each}
          </div>

        <!-- 瀑布流视图 -->
        {:else}
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {#each filteredNotes as note, index (note.id)}
              <div 
                class="relative group animate-fade-in"
                style="animation-delay: {index * 30}ms"
              >
                <!-- 选择指示器 -->
                <div 
                  class="absolute -left-2 top-2 w-6 h-6 rounded-full border-2 border-border bg-card flex items-center justify-center cursor-pointer opacity-0 group-hover:opacity-100 transition-all duration-200 z-10"
                  class:opacity-100={selectedNoteIds.has(note.id)}
                  class:bg-primary={selectedNoteIds.has(note.id)}
                  class:border-primary={selectedNoteIds.has(note.id)}
                  on:click|stopPropagation={() => toggleNoteSelection(note.id)}
                >
                  {#if selectedNoteIds.has(note.id)}
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" class="text-primary-foreground">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                  {/if}
                </div>

                <NoteCard 
                  {note} 
                  on:click={() => handleNoteClick(note.id)}
                  on:doubleClick={() => handleNoteDoubleClick(note.id)}
                />
              </div>
            {/each}
          </div>
        {/if}
      {/if}
    </div>
  </main>
</div>

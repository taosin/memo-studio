<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import NoteCard from './NoteCard.svelte';
  import TagTree from './TagTree.svelte';
  import SearchBar from './SearchBar.svelte';
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
  let collapsedGroups = new Set();
  let selectedNoteIds = new Set();
  let randomNote = null;
  let showRandomNote = false;

  onMount(async () => {
    await loadNotes();
    await loadRandomNote();
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
      console.error('加载笔记失败:', err);
    } finally {
      loading = false;
    }
  }

  async function loadRandomNote() {
    try {
      const response = await fetch('/api/review/random', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      if (response.ok) {
        const data = await response.json();
        if (Array.isArray(data) && data.length > 0) {
          randomNote = data[0];
        }
      }
    } catch (err) {
      console.error('加载随机笔记失败:', err);
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
    
    if (!confirm(`确定要删除选中的 ${selectedNoteIds.size} 条笔记吗？`)) {
      return;
    }

    try {
      await api.deleteNotes(Array.from(selectedNoteIds));
      selectedNoteIds.clear();
      await loadNotes();
    } catch (err) {
      alert('删除失败: ' + err.message);
    }
  }

  function handleSearch(query) {
    searchQuery = query;
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
    collapsedGroups.clear();
  }

  function handleNoteClick(noteId) {
    dispatch('noteClick', noteId);
  }

  function handleNoteDoubleClick(noteId) {
    if (onQuickEdit) {
      onQuickEdit(noteId);
    } else {
      handleNoteClick(noteId);
    }
  }

  function handleTagClick(tag, event) {
    event.stopPropagation();
    handleTagSelect(tag);
  }

  function toggleGroup(date) {
    if (collapsedGroups.has(date)) {
      collapsedGroups.delete(date);
    } else {
      collapsedGroups.add(date);
    }
    collapsedGroups = collapsedGroups;
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
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(note);
    });
    
    const sortedDates = Object.keys(groups).sort((a, b) => {
      const dateA = groups[a][0]?.created_at || '';
      const dateB = groups[b][0]?.created_at || '';
      return new Date(dateB) - new Date(dateA);
    });
    
    return sortedDates.reduce((acc, key) => {
      acc[key] = groups[key];
      return acc;
    }, {});
  })();

  $: {
    filterNotes();
  }

  function clearFilters() {
    searchQuery = '';
    selectedTags = [];
    filterNotes();
  }
</script>

<div class="flex flex-col lg:flex-row gap-6">
  <!-- 左侧边栏 -->
  <aside class="w-full lg:w-64 flex-shrink-0">
    <div class="sticky top-24 space-y-4">
      <!-- 搜索 -->
      <SearchBar value={searchQuery} on:search={(e) => handleSearch(e.detail)} />

      <!-- 快捷筛选 -->
      {#if selectedTags.length > 0 || searchQuery}
        <div class="flex items-center justify-between p-3 bg-primary/5 rounded-lg">
          <span class="text-sm text-muted-foreground">
            {filteredNotes.length} 条结果
          </span>
          <button 
            class="text-sm text-primary hover:underline"
            on:click={clearFilters}
          >
            清除筛选
          </button>
        </div>
      {/if}

      <!-- 标签树 -->
      <TagTree 
        {selectedTags} 
        on:tagSelect={(e) => handleTagSelect(e.detail)} 
      />

      <!-- 随机回顾 -->
      {#if randomNote && !showRandomNote}
        <div class="p-4 bg-gradient-to-br from-primary/10 to-primary/5 rounded-xl border border-primary/20">
          <div class="flex items-center gap-2 mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-primary">
              <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
              <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
            </svg>
            <span class="text-sm font-medium text-primary">今日回顾</span>
          </div>
          <p class="text-sm text-muted-foreground line-clamp-2 mb-3">
            {randomNote.title || randomNote.content?.substring(0, 50)}...
          </p>
          <Button variant="outline" size="sm" class="w-full" on:click={() => showRandomNote = true}>
            查看详情
          </Button>
        </div>
      {/if}
    </div>
  </aside>

  <!-- 主内容区 -->
  <div class="flex-1 min-w-0">
    <!-- 顶部操作栏 -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <!-- 视图切换 -->
        <div class="flex bg-secondary/50 rounded-lg p-1">
          <button
            class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
            class:bg-primary={viewMode === 'timeline'}
            class:text-primary-foreground={viewMode === 'timeline'}
            class:text-muted-foreground={viewMode !== 'timeline'}
            class:hover:bg-accent={viewMode !== 'timeline'}
            on:click={() => handleViewModeChange('timeline')}
          >
            时间线
          </button>
          <button
            class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
            class:bg-primary={viewMode === 'waterfall'}
            class:text-primary-foreground={viewMode === 'waterfall'}
            class:text-muted-foreground={viewMode !== 'waterfall'}
            class:hover:bg-accent={viewMode !== 'waterfall'}
            on:click={() => handleViewModeChange('waterfall')}
          >
            卡片
          </button>
        </div>
        
        {#if selectedNoteIds.size > 0}
          <span class="text-sm text-muted-foreground">
            已选 {selectedNoteIds.size} 条
          </span>
        {/if}
      </div>

      <div class="flex items-center gap-2">
        {#if selectedNoteIds.size > 0}
          <Button variant="destructive" size="sm" on:click={handleBatchDelete}>
            删除
          </Button>
        {/if}
      </div>
    </div>

    <!-- 加载状态 -->
    {#if loading}
      <div class="space-y-6">
        {#each Array(3) as _, i}
          <div class="animate-pulse">
            <div class="h-6 w-32 bg-muted rounded mb-4"></div>
            <div class="space-y-3">
              <div class="h-32 bg-muted rounded-xl"></div>
              <div class="h-32 bg-muted rounded-xl"></div>
            </div>
          </div>
        {/each}
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
        <Button on:click={loadNotes}>重试</Button>
      </div>
    {:else if filteredNotes.length === 0}
      <!-- 空状态 -->
      <div class="flex flex-col items-center justify-center py-20 text-center animate-fade-in">
        <div class="w-24 h-24 bg-primary/10 rounded-full flex items-center justify-center mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-primary">
            <path d="M12 19l7-7 3 3-7 7-3-3z"></path>
            <path d="M18 13l-1.5-7.5L2 2l3.5 14.5L13 18l5-5z"></path>
            <path d="M2 2l7.586 7.586"></path>
            <circle cx="11" cy="11" r="2"></circle>
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-2">开始记录你的灵感</h3>
        <p class="text-muted-foreground mb-6 max-w-sm">
          点击底部的"记录灵感"按钮，开始你的第一篇笔记
        </p>
        <div class="flex items-center gap-4 text-sm text-muted-foreground">
          <span class="flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            支持 #标签
          </span>
          <span class="flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            全文搜索
          </span>
          <span class="flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            每日回顾
          </span>
        </div>
      </div>
    {:else}
      <!-- 时间线视图 -->
      {#if viewMode === 'timeline'}
        <div class="space-y-8">
          {#each Object.entries(groupedNotes) as [date, dateNotes], index}
            <div class="animate-fade-in" style="animation-delay: {index * 50}ms">
              <!-- 日期分组标题 -->
              <div class="flex items-center gap-3 mb-4">
                <div class="flex-1 h-px bg-border"></div>
                <span class="text-sm font-medium text-muted-foreground px-3">
                  {formatGroupDate(date)}
                </span>
                <div class="flex-1 h-px bg-border"></div>
              </div>

              <!-- 该日期的笔记列表 -->
              <div class="space-y-3">
                {#each dateNotes as note, noteIndex (note.id)}
                  <div 
                    class="relative animate-slide-up"
                    style="animation-delay: {noteIndex * 30}ms"
                  >
                    <!-- 选择checkbox -->
                    <div 
                      class="absolute -left-3 top-4 w-6 h-6 rounded-full border-2 border-border bg-card flex items-center justify-center cursor-pointer hover:border-primary transition-colors z-10"
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

                    <div class="ml-6">
                      <NoteCard 
                        {note} 
                        on:click={() => handleNoteClick(note.id)}
                        on:doubleClick={() => handleNoteDoubleClick(note.id)}
                        on:tagClick={(e) => handleTagClick(e.detail.tag, e.detail.event)}
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
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each filteredNotes as note, index (note.id)}
            <div 
              class="relative animate-fade-in"
              style="animation-delay: {index * 30}ms"
            >
              <!-- 选择checkbox -->
              <div 
                class="absolute -left-2 top-2 w-6 h-6 rounded-full border-2 border-border bg-card flex items-center justify-center cursor-pointer hover:border-primary transition-colors z-10"
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
                on:tagClick={(e) => handleTagClick(e.detail.tag, e.detail.event)}
              />
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</div>

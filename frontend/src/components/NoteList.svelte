<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import NoteCard from './NoteCard.svelte';
  import TagTree from './TagTree.svelte';
  import SearchBar from './SearchBar.svelte';
  import StatsOverview from './StatsOverview.svelte';
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
  let sidebarCollapsed = false;
  let mobileMenuOpen = false;

  // 响应式侧边栏控制
  function checkMobile() {
    if (typeof window !== 'undefined') {
      mobileMenuOpen = window.innerWidth < 768;
      sidebarCollapsed = window.innerWidth < 768;
    }
  }

  onMount(() => {
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  });

  onMount(async () => {
    await loadNotes();
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

  function filterNotes() {
    if (!Array.isArray(notes)) {
      filteredNotes = [];
      return;
    }
    
    let filtered = notes;

    // 如果有搜索条件才过滤
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(note => 
        (note.title || '').toLowerCase().includes(query) ||
        (note.content || '').replace(/<[^>]*>/g, '').toLowerCase().includes(query)
      );
    }

    // 如果有标签筛选才过滤
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
    searchQuery = e.detail;
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
    const diffTime = Math.abs(now - date);
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
    
    const sortedDates = Object.keys(groups).sort((a, b) => new Date(b) - new Date(a));
    
    return sortedDates.reduce((acc, key) => {
      acc[key] = groups[key];
      return acc;
    }, {});
  })();

  $: activeFilters = (searchQuery.trim() ? 1 : 0) + selectedTags.length;
</script>

<div class="flex min-h-screen">
  <!-- 移动端侧边栏遮罩 -->
  {#if mobileMenuOpen && !sidebarCollapsed}
    <div 
      class="fixed inset-0 bg-black/50 z-30 md:hidden"
      on:click={() => sidebarCollapsed = true}
      on:keydown={(e) => e.key === 'Escape' && (sidebarCollapsed = true)}
      role="button"
      tabindex="0"
    ></div>
  {/if}

  <!-- 侧边栏 -->
  <aside 
    class="w-64 flex-shrink-0 border-r bg-card/50 transition-all duration-300 fixed md:relative z-40 h-full {sidebarCollapsed ? 'w-0 md:w-16 -translate-x-full md:translate-x-0' : ''} {mobileMenuOpen ? 'translate-x-0' : ''}"
  >
    <div class="sticky top-0 h-screen flex flex-col p-4 overflow-hidden">
      <!-- Logo / 品牌 -->
      <div class="flex items-center gap-3 mb-6">
        <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary to-primary-light flex items-center justify-center shadow-lg shadow-primary/20 flex-shrink-0">
          <span class="text-xl">📝</span>
        </div>
        {#if !sidebarCollapsed}
          <span class="font-bold bg-gradient-to-r from-primary to-primary-light bg-clip-text text-transparent whitespace-nowrap">Memo</span>
        {/if}
      </div>

      <!-- 统计概览 -->
      {#if !sidebarCollapsed}
        <div class="mb-4">
          <StatsOverview />
        </div>
      {/if}

      <!-- 搜索栏 -->
      {#if !sidebarCollapsed}
        <div class="mb-4">
          <SearchBar on:search={handleSearch} />
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
            aria-label="展开侧边栏"
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
      <div class="flex items-center justify-between px-4 sm:px-6 py-3 sm:py-4">
        <!-- 左侧：移动端菜单按钮 + 视图切换 -->
        <div class="flex items-center gap-2">
          <!-- 移动端菜单按钮 -->
          <button
            class="md:hidden p-2 rounded-lg hover:bg-accent transition-colors"
            on:click={() => sidebarCollapsed = !sidebarCollapsed}
            aria-label="打开菜单"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="3" y1="12" x2="21" y2="12"></line>
              <line x1="3" y1="6" x2="21" y2="6"></line>
              <line x1="3" y1="18" x2="21" y2="18"></line>
            </svg>
          </button>
          
          <!-- 视图切换 -->
          <div class="flex items-center gap-1 bg-secondary/50 rounded-full p-0.5 sm:p-1">
            <button
              class="px-2 sm:px-4 py-1.5 rounded-full text-xs sm:text-sm font-medium transition-all"
              class:bg-primary={viewMode === 'timeline'}
              class:text-primary-foreground={viewMode === 'timeline'}
              class:text-muted-foreground={viewMode !== 'timeline'}
              on:click={() => handleViewModeChange('timeline')}
            >
              <span class="hidden sm:inline">信息流</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="sm:hidden inline">
                <line x1="8" y1="6" x2="21" y2="6"></line>
                <line x1="8" y1="12" x2="21" y2="12"></line>
                <line x1="8" y1="18" x2="21" y2="18"></line>
                <line x1="3" y1="6" x2="3.01" y2="6"></line>
                <line x1="3" y1="12" x2="3.01" y2="12"></line>
                <line x1="3" y1="18" x2="3.01" y2="18"></line>
              </svg>
            </button>
            <button
              class="px-2 sm:px-4 py-1.5 rounded-full text-xs sm:text-sm font-medium transition-all"
              class:bg-primary={viewMode === 'waterfall'}
              class:text-primary-foreground={viewMode === 'waterfall'}
              class:text-muted-foreground={viewMode !== 'waterfall'}
              on:click={() => handleViewModeChange('waterfall')}
            >
              <span class="hidden sm:inline">卡片</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="sm:hidden inline">
                <rect x="3" y="3" width="7" height="7"></rect>
                <rect x="14" y="3" width="7" height="7"></rect>
                <rect x="14" y="14" width="7" height="7"></rect>
                <rect x="3" y="14" width="7" height="7"></rect>
              </svg>
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
              <div class="h-6 w-20 sm:w-24 bg-muted rounded-full mb-4"></div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="h-36 sm:h-40 bg-muted/50 rounded-2xl"></div>
                <div class="h-36 sm:h-40 bg-muted/50 rounded-2xl hidden sm:block"></div>
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
        <div class="flex flex-col items-center justify-center py-16 sm:py-20 text-center animate-fade-in">
          <!-- 动态插图 - 使用emoji动画 -->
          <div class="relative mb-8">
            <div class="w-28 h-28 sm:w-32 sm:h-32 bg-gradient-to-br from-primary/10 to-primary/5 rounded-full flex items-center justify-center animate-pulse">
              <span class="text-5xl sm:text-6xl">✨</span>
            </div>
            <!-- 浮动装饰 -->
            <div class="absolute -top-2 -right-2 w-8 h-8 bg-amber-100 rounded-full flex items-center justify-center animate-bounce" style="animation-delay: 0.2s">
              <span class="text-lg">💡</span>
            </div>
            <div class="absolute -bottom-1 -left-1 w-6 h-6 bg-purple-100 rounded-full flex items-center justify-center animate-bounce" style="animation-delay: 0.4s">
              <span class="text-sm">📝</span>
            </div>
          </div>
          
          <h3 class="text-xl sm:text-2xl font-bold mb-3 bg-gradient-to-r from-primary to-primary-light bg-clip-text text-transparent">
            {activeFilters > 0 ? '没有找到匹配的笔记' : '开始记录灵感'}
          </h3>
          <p class="text-muted-foreground mb-8 max-w-sm mx-auto">
            {activeFilters > 0 ? '试试调整筛选条件，或者记录一条新笔记' : '点击底部的 ✏️ 按钮，记录你的第一条灵感'}
          </p>
          
          {#if activeFilters > 0}
            <Button on:click={clearFilters} class="mb-4">
              清除筛选条件
            </Button>
          {:else}
            <div class="flex flex-wrap justify-center gap-4 sm:gap-6 text-sm text-muted-foreground">
              <span class="flex items-center gap-2 bg-secondary/50 px-3 py-2 rounded-lg">
                <kbd class="px-2 py-0.5 bg-background rounded text-xs shadow-sm">#</kbd>
                添加标签
              </span>
              <span class="flex items-center gap-2 bg-secondary/50 px-3 py-2 rounded-lg">
                <kbd class="px-2 py-0.5 bg-background rounded text-xs shadow-sm">/</kbd>
                快速搜索
              </span>
              <span class="flex items-center gap-2 bg-secondary/50 px-3 py-2 rounded-lg">
                <kbd class="px-2 py-0.5 bg-background rounded text-xs shadow-sm">⌘N</kbd>
                新建笔记
              </span>
            </div>
          {/if}
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
                      <button 
                        type="button"
                        class="absolute -left-3 top-4 w-8 h-8 rounded-full border-2 border-border bg-card flex items-center justify-center cursor-pointer opacity-0 group-hover:opacity-100 transition-all duration-200 shadow-sm"
                        class:opacity-100={selectedNoteIds.has(note.id)}
                        class:bg-primary={selectedNoteIds.has(note.id)}
                        class:border-primary={selectedNoteIds.has(note.id)}
                        on:click|stopPropagation={() => toggleNoteSelection(note.id)}
                        aria-label={selectedNoteIds.has(note.id) ? '取消选择' : '选择笔记'}
                      >
                        {#if selectedNoteIds.has(note.id)}
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" class="text-primary-foreground">
                            <polyline points="20 6 9 17 4 12"></polyline>
                          </svg>
                        {/if}
                      </button>

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
                <button 
                  type="button"
                  class="absolute -left-2 top-2 w-6 h-6 rounded-full border-2 border-border bg-card flex items-center justify-center cursor-pointer opacity-0 group-hover:opacity-100 transition-all duration-200 z-10"
                  class:opacity-100={selectedNoteIds.has(note.id)}
                  class:bg-primary={selectedNoteIds.has(note.id)}
                  class:border-primary={selectedNoteIds.has(note.id)}
                  on:click|stopPropagation={() => toggleNoteSelection(note.id)}
                  aria-label={selectedNoteIds.has(note.id) ? '取消选择' : '选择笔记'}
                >
                  {#if selectedNoteIds.has(note.id)}
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" class="text-primary-foreground">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                  {/if}
                </button>

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

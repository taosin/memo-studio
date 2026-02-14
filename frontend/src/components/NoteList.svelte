<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import NoteCard from './NoteCard.svelte';
  import TagTree from './TagTree.svelte';
  import SearchBar from './SearchBar.svelte';
  import AdvancedSearch from './AdvancedSearch.svelte';
  import ViewModeToggle from './ViewModeToggle.svelte';
  import ExportDialog from './ExportDialog.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import { api } from '../utils/api.js';

  const dispatch = createEventDispatcher();

  export let onQuickEdit = null;
  export let onDelete = null;

  let notes = [];
  let filteredNotes = [];
  let loading = true;
  let error = null;
  let searchQuery = '';
  let selectedTags = [];
  let viewMode = 'waterfall';
  let collapsedGroups = new Set();
  let showAdvancedSearch = false;
  let showExportDialog = false;
  let searchHistory = [];
  let selectedNoteIds = new Set();
  let searchFilters = {
    keyword: '',
    tags: [],
    dateFrom: '',
    dateTo: ''
  };

  onMount(async () => {
    await loadNotes();
    const saved = localStorage.getItem('searchHistory');
    if (saved) {
      try {
        searchHistory = JSON.parse(saved);
      } catch (e) {
        console.error('加载搜索历史失败:', e);
      }
    }
  });

  async function loadNotes() {
    try {
      loading = true;
      const data = await api.getNotes();
      notes = Array.isArray(data) ? data : [];
      if (!Array.isArray(data)) {
        console.warn('API 返回的数据不是数组:', data);
      }
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

  function filterNotes() {
    if (!Array.isArray(notes)) {
      console.warn('notes 不是数组:', notes);
      filteredNotes = [];
      return;
    }
    let filtered = [...notes];

    if (searchFilters.keyword) {
      const query = searchFilters.keyword.toLowerCase();
      filtered = filtered.filter(note => 
        (note.title || '').toLowerCase().includes(query) ||
        (note.content || '').replace(/<[^>]*>/g, '').toLowerCase().includes(query)
      );
    }

    const tagsToFilter = searchFilters.tags.length > 0 ? searchFilters.tags : selectedTags;
    if (tagsToFilter.length > 0) {
      filtered = filtered.filter(note => {
        const noteTagIds = (note.tags || []).map(t => t.id);
        return tagsToFilter.some(tagId => noteTagIds.includes(tagId));
      });
    }

    if (searchFilters.dateFrom) {
      const fromDate = new Date(searchFilters.dateFrom);
      filtered = filtered.filter(note => new Date(note.created_at) >= fromDate);
    }
    if (searchFilters.dateTo) {
      const toDate = new Date(searchFilters.dateTo);
      toDate.setHours(23, 59, 59, 999);
      filtered = filtered.filter(note => new Date(note.created_at) <= toDate);
    }

    if (!searchFilters.keyword && searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(note => 
        (note.title || '').toLowerCase().includes(query) ||
        (note.content || '').replace(/<[^>]*>/g, '').toLowerCase().includes(query)
      );
    }

    filteredNotes = filtered;
  }

  function handleAdvancedSearch(e) {
    const filters = e.detail;
    searchFilters = filters;
    
    if (filters.keyword || filters.tags.length > 0) {
      const historyItem = {
        keyword: filters.keyword,
        tags: filters.tags,
        dateFrom: filters.dateFrom,
        dateTo: filters.dateTo,
        timestamp: new Date().toISOString()
      };
      searchHistory = [historyItem, ...searchHistory.filter(h => 
        h.keyword !== historyItem.keyword || 
        JSON.stringify(h.tags) !== JSON.stringify(historyItem.tags)
      )].slice(0, 10);
      localStorage.setItem('searchHistory', JSON.stringify(searchHistory));
    }
    
    filterNotes();
  }

  function handleClearSearch() {
    searchFilters = { keyword: '', tags: [], dateFrom: '', dateTo: '' };
    filterNotes();
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
      if (onDelete) onDelete();
    } catch (err) {
      alert('删除失败: ' + err.message);
    }
  }

  function handleSearch(query) {
    searchQuery = query;
    filterNotes();
  }

  function handleTagSelect(tag) {
    const tagId = tag.id;
    if (selectedTags.includes(tagId)) {
      selectedTags = selectedTags.filter(id => id !== tagId);
    } else {
      selectedTags = [...selectedTags, tagId];
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

  $: groupedNotes = (() => {
    if (viewMode !== 'timeline') return {};
    
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

  // 骨架屏组件
  function SkeletonCard() {
    return `
      <div class="bg-card border border-border/60 rounded-lg p-4 animate-pulse">
        <div class="h-5 bg-muted rounded w-3/4 mb-3"></div>
        <div class="h-4 bg-muted rounded w-full mb-2"></div>
        <div class="h-4 bg-muted rounded w-5/6 mb-2"></div>
        <div class="h-4 bg-muted rounded w-4/6 mb-4"></div>
        <div class="flex gap-2">
          <div class="h-5 bg-muted rounded w-16"></div>
          <div class="h-5 bg-muted rounded w-16"></div>
        </div>
      </div>
    `;
  }
</script>

<div class="flex flex-col md:flex-row gap-4">
  <!-- 左侧边栏 -->
  <aside class="w-full md:w-64 flex-shrink-0 md:block">
    <div class="sticky top-24 space-y-4">
      <div class="hidden md:block">
        <SearchBar value={searchQuery} on:search={(e) => handleSearch(e.detail)} />
      </div>
      <div class="hidden lg:block">
        <TagTree {selectedTags} on:tagSelect={(e) => handleTagSelect(e.detail)} />
      </div>
    </div>
  </aside>

  <!-- 主内容区 -->
  <div class="flex-1 min-w-0">
    <!-- 移动端搜索栏和标签树 -->
    <div class="md:hidden mb-4 space-y-3">
      <SearchBar value={searchQuery} on:search={(e) => handleSearch(e.detail)} />
      <div class="lg:hidden">
        <TagTree {selectedTags} on:tagSelect={(e) => handleTagSelect(e.detail)} />
      </div>
    </div>

    <!-- 工具栏 -->
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-4">
      <div class="flex items-center gap-3">
        <ViewModeToggle mode={viewMode} on:change={(e) => handleViewModeChange(e.detail)} />
        <Button
          variant="outline"
          size="sm"
          on:click={() => showAdvancedSearch = !showAdvancedSearch}
        >
          {showAdvancedSearch ? '收起' : '高级搜索'}
        </Button>
        <Button
          variant="outline"
          size="sm"
          on:click={() => showExportDialog = true}
        >
          导出
        </Button>
        {#if selectedNoteIds.size > 0}
          <Button
            variant="destructive"
            size="sm"
            on:click={handleBatchDelete}
          >
            删除选中 ({selectedNoteIds.size})
          </Button>
        {/if}
      </div>
      <div class="text-sm text-muted-foreground">
        共 {filteredNotes.length} 条笔记
      </div>
    </div>

    <!-- 高级搜索 -->
    <AdvancedSearch
      visible={showAdvancedSearch}
      {searchHistory}
      on:search={handleAdvancedSearch}
      on:clear={handleClearSearch}
      on:close={() => showAdvancedSearch = false}
    />

    <!-- 导出对话框 -->
    <ExportDialog
      visible={showExportDialog}
      selectedNotes={Array.from(selectedNoteIds).map(id => notes.find(n => n.id === id)).filter(Boolean)}
      on:close={() => showExportDialog = false}
      on:exported={() => {
        selectedNoteIds.clear();
      }}
    />

    {#if loading}
      <!-- 加载骨架屏 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {#each Array(8) as _, i}
          <div class="bg-card border border-border/60 rounded-lg p-4 animate-pulse" key={i}>
            <div class="h-5 bg-muted rounded w-3/4 mb-3"></div>
            <div class="h-4 bg-muted rounded w-full mb-2"></div>
            <div class="h-4 bg-muted rounded w-5/6 mb-2"></div>
            <div class="h-4 bg-muted rounded w-4/6 mb-4"></div>
            <div class="flex gap-2">
              <div class="h-5 bg-muted rounded w-16"></div>
              <div class="h-5 bg-muted rounded w-16"></div>
            </div>
          </div>
        {/each}
      </div>
    {:else if error}
      <!-- 错误状态 -->
      <div class="flex flex-col items-center justify-center py-16 text-center animate-fade-in">
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
      <div class="flex flex-col items-center justify-center py-16 text-center animate-fade-in">
        <div class="w-20 h-20 bg-primary/10 rounded-full flex items-center justify-center mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-primary">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="12" y1="18" x2="12" y2="12"></line>
            <line x1="9" y1="15" x2="15" y2="15"></line>
          </svg>
        </div>
        <h3 class="text-lg font-semibold mb-2">还没有笔记</h3>
        <p class="text-muted-foreground mb-4">点击右上角按钮创建你的第一篇笔记</p>
      </div>
    {:else if viewMode === 'waterfall'}
      <!-- 瀑布流模式 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {#each filteredNotes as note, index (note.id)}
          <div 
            class="relative animate-fade-in" 
            style="animation-delay: {index * 50}ms"
          >
            {#if selectedNoteIds.has(note.id)}
              <div class="absolute top-2 right-2 z-10 bg-primary text-primary-foreground rounded-full w-6 h-6 flex items-center justify-center text-xs shadow-md">
                ✓
              </div>
            {/if}
            <NoteCard 
              {note} 
              on:click={() => handleNoteClick(note.id)}
              on:doubleClick={() => handleNoteDoubleClick(note.id)}
              on:tagClick={(e) => handleTagClick(e.detail.tag, e.detail.event)}
            />
            <input
              type="checkbox"
              class="absolute top-2 left-2 z-10 w-5 h-5 cursor-pointer"
              checked={selectedNoteIds.has(note.id)}
              on:change={() => toggleNoteSelection(note.id)}
              on:click|stopPropagation
            />
          </div>
        {/each}
      </div>
    {:else}
      <!-- Timeline 模式 -->
      <div class="space-y-6">
        {#each Object.entries(groupedNotes) as [date, dateNotes], index}
          <div 
            class="border-l-2 border-primary/20 pl-4 animate-slide-down"
            style="animation-delay: {index * 100}ms"
          >
            <div 
              class="flex items-center gap-2 mb-4 cursor-pointer hover:text-primary transition-colors"
              role="button"
              tabindex="0"
              on:click={() => toggleGroup(date)}
              on:keydown={(e) => e.key === 'Enter' && toggleGroup(date)}
            >
              <span class="text-xs transition-transform duration-200" class:rotate-90={!collapsedGroups.has(date)}>
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="9 18 15 12 9 6"></polyline>
                </svg>
              </span>
              <h3 class="text-lg font-semibold">{date}</h3>
              <span class="text-xs text-muted-foreground">({dateNotes.length})</span>
            </div>
            {#if !collapsedGroups.has(date)}
              <div class="space-y-3 ml-4">
                {#each dateNotes as note, noteIndex (note.id)}
                  <div 
                    class="relative animate-fade-in"
                    style="animation-delay: {noteIndex * 50}ms"
                  >
                    {#if selectedNoteIds.has(note.id)}
                      <div class="absolute top-2 right-2 z-10 bg-primary text-primary-foreground rounded-full w-6 h-6 flex items-center justify-center text-xs shadow-md">
                        ✓
                      </div>
                    {/if}
                    <NoteCard 
                      {note} 
                      on:click={() => handleNoteClick(note.id)}
                      on:doubleClick={() => handleNoteDoubleClick(note.id)}
                      on:tagClick={(e) => handleTagClick(e.detail.tag, e.detail.event)}
                    />
                    <input
                      type="checkbox"
                      class="absolute top-2 left-2 z-10 w-5 h-5 cursor-pointer"
                      checked={selectedNoteIds.has(note.id)}
                      on:change={() => toggleNoteSelection(note.id)}
                      on:click|stopPropagation
                    />
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

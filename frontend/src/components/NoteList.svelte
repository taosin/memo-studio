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
  let viewMode = 'waterfall'; // 'waterfall' | 'timeline'
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
    // 加载搜索历史
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
      notes = await api.getNotes();
      filterNotes();
      error = null;
    } catch (err) {
      error = err.message;
      console.error('加载笔记失败:', err);
    } finally {
      loading = false;
    }
  }

  function filterNotes() {
    let filtered = [...notes];

    // 高级搜索过滤
    if (searchFilters.keyword) {
      const query = searchFilters.keyword.toLowerCase();
      filtered = filtered.filter(note => 
        (note.title || '').toLowerCase().includes(query) ||
        (note.content || '').replace(/<[^>]*>/g, '').toLowerCase().includes(query)
      );
    }

    // 按标签过滤
    const tagsToFilter = searchFilters.tags.length > 0 ? searchFilters.tags : selectedTags;
    if (tagsToFilter.length > 0) {
      filtered = filtered.filter(note => {
        const noteTagIds = (note.tags || []).map(t => t.id);
        return tagsToFilter.some(tagId => noteTagIds.includes(tagId));
      });
    }

    // 按日期范围过滤
    if (searchFilters.dateFrom) {
      const fromDate = new Date(searchFilters.dateFrom);
      filtered = filtered.filter(note => new Date(note.created_at) >= fromDate);
    }
    if (searchFilters.dateTo) {
      const toDate = new Date(searchFilters.dateTo);
      toDate.setHours(23, 59, 59, 999);
      filtered = filtered.filter(note => new Date(note.created_at) <= toDate);
    }

    // 简单搜索（如果高级搜索未使用）
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
    
    // 保存到搜索历史
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
      )].slice(0, 10); // 保留最近10条
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
    selectedNoteIds = selectedNoteIds; // 触发更新
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
      // 如果没有提供快速编辑函数，则触发普通点击
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
    collapsedGroups = collapsedGroups; // 触发更新
  }

  // 按日期分组（用于timeline模式）
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
    
    // 按日期排序（需要按实际日期对象排序）
    const sortedDates = Object.keys(groups).sort((a, b) => {
      // 从日期字符串中提取第一个笔记的创建时间
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
      <div class="text-center py-12 text-muted-foreground">加载中...</div>
    {:else if error}
      <div class="text-center py-12 text-destructive">错误: {error}</div>
    {:else if filteredNotes.length === 0}
      <div class="text-center py-12 text-muted-foreground">
        <p>没有找到笔记</p>
      </div>
    {:else if viewMode === 'waterfall'}
      <!-- 瀑布流模式 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {#each filteredNotes as note (note.id)}
          <div class="relative">
            {#if selectedNoteIds.has(note.id)}
              <div class="absolute top-2 right-2 z-10 bg-primary text-primary-foreground rounded-full w-6 h-6 flex items-center justify-center text-xs">
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
              class="absolute top-2 left-2 z-10"
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
        {#each Object.entries(groupedNotes) as [date, dateNotes]}
          <div class="border-l-2 border-primary/20 pl-4">
            <div 
              class="flex items-center gap-2 mb-4 cursor-pointer hover:text-primary transition-colors"
              role="button"
              tabindex="0"
              on:click={() => toggleGroup(date)}
              on:keydown={(e) => e.key === 'Enter' && toggleGroup(date)}
            >
              <span class="text-xs">{collapsedGroups.has(date) ? '▶' : '▼'}</span>
              <h3 class="text-lg font-semibold">{date}</h3>
              <span class="text-xs text-muted-foreground">({dateNotes.length})</span>
            </div>
            {#if !collapsedGroups.has(date)}
              <div class="space-y-3 ml-6">
                {#each dateNotes as note (note.id)}
                  <div class="relative">
                    {#if selectedNoteIds.has(note.id)}
                      <div class="absolute top-2 right-2 z-10 bg-primary text-primary-foreground rounded-full w-6 h-6 flex items-center justify-center text-xs">
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
                      class="absolute top-2 left-2 z-10"
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

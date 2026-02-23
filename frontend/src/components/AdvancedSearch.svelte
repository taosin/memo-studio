<script>
  import { onMount, onDestroy } from 'svelte';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let showHelp = false;
  let searchQuery = '';
  let dateFrom = '';
  let dateTo = '';
  let selectedTags = [];
  let sortBy = 'created_at';
  let sortOrder = 'desc';

  const shortcuts = [
    { key: 'n', desc: '新建笔记' },
    { key: '/', desc: '聚焦搜索' },
    { key: 'Esc', desc: '关闭弹窗' },
    { key: '⌘K', desc: '打开搜索' },
    { key: '?', desc: '显示快捷键' },
  ];

  function handleKeydown(e) {
    // ? 显示快捷键帮助
    if (e.key === '?' && !e.target.matches('input, textarea')) {
      e.preventDefault();
      showHelp = !showHelp;
    }
    // Escape 关闭
    if (e.key === 'Escape') {
      showHelp = false;
    }
  }

  onMount(() => {
    document.addEventListener('keydown', handleKeydown);
  });

  onDestroy(() => {
    document.removeEventListener('keydown', handleKeydown);
  });

  function handleSearch() {
    dispatch('search', { keyword: searchQuery, tags: selectedTags, dateFrom, dateTo, sortBy, sortOrder });
  }

  function handleClear() {
    searchQuery = '';
    dateFrom = '';
    dateTo = '';
    selectedTags = [];
    dispatch('clear');
  }
</script>

{#if showHelp}
  <!-- 遮罩层 -->
  <div 
    class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center animate-fade-in"
    on:click={() => showHelp = false}
    on:keydown={(e) => e.key === 'Escape' && (showHelp = false)}
    role="button"
    tabindex="-1"
  >
    <!-- 弹窗 -->
    <div 
      class="bg-popover border border-border rounded-2xl shadow-2xl p-6 max-w-md w-full mx-4 animate-scale-in"
      on:click|stopPropagation
      role="dialog"
    >
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-semibold">键盘快捷键</h3>
        <button 
          class="p-1 rounded-lg hover:bg-accent transition-colors"
          on:click={() => showHelp = false}
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <div class="space-y-3">
        {#each shortcuts as shortcut}
          <div class="flex items-center justify-between py-2 border-b border-border/50 last:border-0">
            <span class="text-sm text-muted-foreground">{shortcut.desc}</span>
            <kbd class="px-3 py-1 bg-muted rounded-lg text-sm font-mono">{shortcut.key}</kbd>
          </div>
        {/each}
      </div>

      <p class="text-xs text-muted-foreground mt-6 text-center">
        按 <kbd class="px-1.5 py-0.5 bg-muted rounded">Esc</kbd> 关闭此窗口
      </p>
    </div>
  </div>
{/if}

<!-- 高级搜索表单 -->
<div class="bg-card border border-border rounded-xl p-4 space-y-4">
  <div class="flex items-center gap-2 text-sm font-medium text-muted-foreground">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
    </svg>
    高级搜索
  </div>

  <!-- 关键词 -->
  <div>
    <label class="text-sm text-muted-foreground mb-2 block">关键词</label>
    <input
      type="text"
      bind:value={searchQuery}
      placeholder="搜索标题和内容..."
      class="w-full h-10 px-3 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
    />
  </div>

  <!-- 日期范围 -->
  <div class="grid grid-cols-2 gap-4">
    <div>
      <label class="text-sm text-muted-foreground mb-2 block">开始日期</label>
      <input
        type="date"
        bind:value={dateFrom}
        class="w-full h-10 px-3 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
      />
    </div>
    <div>
      <label class="text-sm text-muted-foreground mb-2 block">结束日期</label>
      <input
        type="date"
        bind:value={dateTo}
        class="w-full h-10 px-3 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
      />
    </div>
  </div>

  <!-- 排序 -->
  <div class="grid grid-cols-2 gap-4">
    <div>
      <label class="text-sm text-muted-foreground mb-2 block">排序方式</label>
      <select
        bind:value={sortBy}
        class="w-full h-10 px-3 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
      >
        <option value="created_at">创建时间</option>
        <option value="updated_at">更新时间</option>
        <option value="title">标题</option>
      </select>
    </div>
    <div>
      <label class="text-sm text-muted-foreground mb-2 block">排序顺序</label>
      <select
        bind:value={sortOrder}
        class="w-full h-10 px-3 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
      >
        <option value="desc">降序（最新优先）</option>
        <option value="asc">升序（最早优先）</option>
      </select>
    </div>
  </div>

  <!-- 操作按钮 -->
  <div class="flex items-center justify-end gap-2 pt-2">
    <button 
      class="px-4 py-2 rounded-lg text-sm text-muted-foreground hover:bg-accent transition-colors"
      on:click={handleClear}
    >
      清除
    </button>
    <button 
      class="px-4 py-2 rounded-lg text-sm bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
      on:click={handleSearch}
    >
      应用筛选
    </button>
  </div>
</div>

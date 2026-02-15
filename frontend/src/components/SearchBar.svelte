<script>
  import { onMount, onDestroy } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import AdvancedSearch from './AdvancedSearch.svelte';
  import Button from '$lib/components/ui/button/button.svelte';

  export let value = '';
  const dispatch = createEventDispatcher();

  let inputValue = value;
  let isFocused = false;
  let showAdvanced = false;
  let showSuggestions = false;
  let suggestions = [];
  let recentSearches = [];

  $: {
    if (inputValue !== value) {
      inputValue = value;
    }
  }

  onMount(() => {
    // 加载最近搜索
    const saved = localStorage.getItem('recentSearches');
    if (saved) {
      try {
        recentSearches = JSON.parse(saved).slice(0, 5);
      } catch (e) {}
    }

    // 全局快捷键
    document.addEventListener('keydown', handleGlobalKeydown);
  });

  onDestroy(() => {
    document.removeEventListener('keydown', handleGlobalKeydown);
  });

  function handleGlobalKeydown(e) {
    // Cmd/Ctrl + K 打开搜索
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault();
      focusSearch();
    }
    // Escape 关闭高级搜索
    if (e.key === 'Escape' && showAdvanced) {
      showAdvanced = false;
    }
  }

  function focusSearch() {
    const input = document.querySelector('.search-input');
    if (input) input.focus();
  }

  function handleInput() {
    dispatch('search', inputValue);
    
    if (inputValue.trim()) {
      showSuggestions = true;
      suggestions = generateSuggestions(inputValue);
    } else {
      showSuggestions = false;
    }
  }

  function generateSuggestions(query) {
    // 基于查询生成建议（实际应该从API获取）
    return [
      { type: 'tag', label: `#${query}` },
      { type: 'text', label: `搜索"${query}"` },
    ];
  }

  function handleSelectSuggestion(suggestion) {
    if (suggestion.type === 'tag') {
      inputValue = suggestion.label;
    } else {
      inputValue = suggestion.label.replace(/"/g, '');
    }
    handleSubmit();
    showSuggestions = false;
  }

  function handleSubmit() {
    dispatch('search', inputValue);
    
    // 保存到最近搜索
    if (inputValue.trim() && !recentSearches.includes(inputValue.trim())) {
      recentSearches = [inputValue.trim(), ...recentSearches.slice(0, 4)];
      localStorage.setItem('recentSearches', JSON.stringify(recentSearches));
    }
    
    showSuggestions = false;
  }

  function handleClear() {
    inputValue = '';
    dispatch('search', '');
    showSuggestions = false;
    focusSearch();
  }

  function handleRecentSearch(term) {
    inputValue = term;
    handleSubmit();
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') {
      handleSubmit();
    }
    if (e.key === 'Escape') {
      showSuggestions = false;
      showAdvanced = false;
    }
  }
</script>

<div class="relative w-full search-container">
  <!-- 搜索框 -->
  <div 
    class="relative flex items-center gap-2 {isFocused ? 'animate-scale-in' : ''}"
  >
    <!-- 搜索图标 -->
    <div class="absolute left-4 text-muted-foreground">
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"></circle>
        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
      </svg>
    </div>

    <!-- 输入框 -->
    <input
      type="text"
      placeholder="搜索笔记、标签..."
      bind:value={inputValue}
      on:input={handleInput}
      on:focus={() => { isFocused = true; if (inputValue) showSuggestions = true; }}
      on:blur={() => setTimeout(() => { isFocused = false; showSuggestions = false; }, 200)}
      on:keydown={handleKeydown}
      class="search-input w-full h-12 pl-12 pr-24 rounded-xl border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all placeholder:text-muted-foreground/50"
    />

    <!-- 快捷键提示 -->
    {#if !isFocused && !inputValue}
      <div class="absolute right-4 flex items-center gap-1 text-xs text-muted-foreground/50 pointer-events-none">
        <kbd class="px-1.5 py-0.5 bg-muted rounded">⌘</kbd>
        <kbd class="px-1.5 py-0.5 bg-muted rounded">K</kbd>
      </div>
    {/if}

    <!-- 清除按钮 -->
    {#if inputValue}
      <button
        class="absolute right-16 p-1 text-muted-foreground hover:text-foreground transition-colors"
        on:click={handleClear}
        aria-label="清除搜索"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    {/if}

    <!-- 高级搜索按钮 -->
    <Button 
      variant="ghost" 
      size="sm" 
      class="absolute right-2"
      on:click={() => showAdvanced = !showAdvanced}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
      </svg>
    </Button>
  </div>

  <!-- 建议下拉框 -->
  {#if showSuggestions && (suggestions.length > 0 || recentSearches.length > 0)}
    <div class="absolute top-full left-0 right-0 mt-2 p-2 bg-popover border border-border rounded-xl shadow-xl z-50 animate-fade-in">
      <!-- 最近搜索 -->
      {#if recentSearches.length > 0 && !inputValue}
        <div class="mb-2">
          <p class="px-3 py-1.5 text-xs text-muted-foreground">最近搜索</p>
          {#each recentSearches as term}
            <button
              class="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm hover:bg-accent transition-colors text-left"
              on:click={() => handleRecentSearch(term)}
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
              {term}
            </button>
          {/each}
        </div>
      {/if}

      <!-- 搜索建议 -->
      {#if suggestions.length > 0}
        <div>
          <p class="px-3 py-1.5 text-xs text-muted-foreground">建议</p>
          {#each suggestions as suggestion}
            <button
              class="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm hover:bg-accent transition-colors text-left"
              on:click={() => handleSelectSuggestion(suggestion)}
            >
              {#if suggestion.type === 'tag'}
                <span class="text-primary">#</span>
              {:else}
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                </svg>
              {/if}
              {suggestion.label}
            </button>
          {/each}
        </div>
      {/if}
    </div>
  {/if}

  <!-- 高级搜索弹窗 -->
  {#if showAdvanced}
    <AdvancedSearch 
      on:close={() => showAdvanced = false}
    />
  {/if}
</div>

<style>
  .search-input {
    box-shadow: 0 0 0 0 transparent;
    transition: all 0.2s ease;
  }
  
  .search-input:focus {
    box-shadow: 0 0 0 2px hsl(var(--primary) / 0.1);
  }
  
  kbd {
    font-family: inherit;
    font-size: 11px;
  }
</style>

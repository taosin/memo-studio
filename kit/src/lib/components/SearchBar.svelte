<script>
  import { createEventDispatcher, onMount } from 'svelte';
  
  const dispatch = createEventDispatcher();
  
  export let value = '';
  export let placeholder = '全文搜索（FTS5）...';
  export let inputElement = null;
  
  let isFocused = false;
  let searchHistory = [];
  let showHistory = false;

  onMount(() => {
    // 加载搜索历史
    try {
      const history = localStorage.getItem('search_history');
      if (history) {
        searchHistory = JSON.parse(history).slice(0, 5);
      }
    } catch {}
  });

  function handleInput(e) {
    value = e.target.value;
    dispatch('input', value);
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      saveToHistory(value);
      dispatch('search', value);
      showHistory = false;
    } else if (e.key === 'Escape') {
      e.target.blur();
      showHistory = false;
    }
  }

  function handleFocus() {
    isFocused = true;
    if (searchHistory.length > 0 && !value) {
      showHistory = true;
    }
  }

  function handleBlur() {
    isFocused = false;
    setTimeout(() => {
      showHistory = false;
    }, 200);
  }

  function handleClear() {
    value = '';
    dispatch('clear');
    inputElement?.focus();
  }

  function selectHistory(item) {
    value = item;
    dispatch('search', value);
    showHistory = false;
  }

  function saveToHistory(query) {
    if (!query || !query.trim()) return;
    
    try {
      const newHistory = [query, ...searchHistory.filter(h => h !== query)].slice(0, 5);
      searchHistory = newHistory;
      localStorage.setItem('search_history', JSON.stringify(newHistory));
    } catch {}
  }

  function clearHistory() {
    searchHistory = [];
    try {
      localStorage.removeItem('search_history');
    } catch {}
    showHistory = false;
  }
</script>

<div class="searchBarWrap">
  <div class="searchBar" class:focused={isFocused}>
    <div class="searchIcon">🔍</div>
    <input
      bind:this={inputElement}
      type="text"
      class="searchInput"
      {placeholder}
      value={value}
      on:input={handleInput}
      on:keydown={handleKeydown}
      on:focus={handleFocus}
      on:blur={handleBlur}
    />
    {#if value}
      <button class="clearBtn" on:click={handleClear} aria-label="清空">
        ✕
      </button>
    {:else}
      <kbd class="shortcutHint">⌘K</kbd>
    {/if}
  </div>

  {#if showHistory && searchHistory.length > 0}
    <div class="historyDropdown">
      <div class="historyHeader">
        <span>搜索历史</span>
        <button class="clearHistoryBtn" on:click={clearHistory}>清空</button>
      </div>
      {#each searchHistory as item}
        <button class="historyItem" on:click={() => selectHistory(item)}>
          <span class="historyIcon">🕐</span>
          <span class="historyText">{item}</span>
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .searchBarWrap {
    position: relative;
    width: 100%;
  }

  .searchBar {
    position: relative;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 12px;
    border: 1px solid var(--border-2);
    border-radius: 10px;
    background: rgba(15, 23, 42, 0.06);
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }

  :global(:root[data-theme="light"]) .searchBar {
    background: rgba(15, 23, 42, 0.04);
  }

  .searchBar.focused {
    border-color: rgba(34, 197, 94, 0.55);
    box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.08);
  }

  .searchIcon {
    font-size: 16px;
    flex-shrink: 0;
    opacity: 0.7;
  }

  .searchInput {
    flex: 1;
    border: none;
    background: transparent;
    color: var(--text);
    padding: 10px 0;
    font-size: 14px;
    outline: none;
  }

  .searchInput::placeholder {
    color: var(--muted);
  }

  .clearBtn {
    flex-shrink: 0;
    border: none;
    background: transparent;
    color: var(--muted);
    cursor: pointer;
    padding: 4px;
    font-size: 14px;
    border-radius: 4px;
    transition: background 0.15s ease, color 0.15s ease;
  }

  .clearBtn:hover {
    background: rgba(148, 163, 184, 0.12);
    color: var(--text);
  }

  .shortcutHint {
    flex-shrink: 0;
    padding: 3px 8px;
    border: 1px solid var(--border);
    border-radius: 4px;
    background: rgba(15, 23, 42, 0.06);
    font-family: ui-monospace, monospace;
    font-size: 11px;
    color: var(--muted);
    pointer-events: none;
  }

  .historyDropdown {
    position: absolute;
    top: calc(100% + 8px);
    left: 0;
    right: 0;
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 10px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    z-index: 100;
    overflow: hidden;
    animation: slideDown 0.2s ease;
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-4px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .historyHeader {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    border-bottom: 1px solid var(--border);
    font-size: 12px;
    color: var(--muted);
  }

  .clearHistoryBtn {
    border: none;
    background: transparent;
    color: var(--muted);
    cursor: pointer;
    padding: 2px 6px;
    font-size: 11px;
    border-radius: 4px;
    transition: background 0.15s ease, color 0.15s ease;
  }

  .clearHistoryBtn:hover {
    background: rgba(148, 163, 184, 0.12);
    color: var(--text);
  }

  .historyItem {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    border: none;
    background: transparent;
    color: var(--text);
    text-align: left;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  .historyItem:hover {
    background: rgba(148, 163, 184, 0.08);
  }

  .historyIcon {
    font-size: 14px;
    opacity: 0.6;
  }

  .historyText {
    flex: 1;
    font-size: 14px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  @media (max-width: 600px) {
    .shortcutHint {
      display: none;
    }
  }
</style>

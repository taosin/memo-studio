<script>
  import { createEventDispatcher } from 'svelte';
  import { shortcuts, formatShortcut } from '../keyboardManager.js';

  const dispatch = createEventDispatcher();
  
  let searchQuery = '';

  const categories = {
    navigation: '🧭 导航',
    editor: '✏️ 编辑',
    note: '📝 笔记',
    view: '👁️ 视图',
    help: '❓ 帮助'
  };

  $: groupedShortcuts = Object.entries(shortcuts).reduce((acc, [key, val]) => {
    const cat = val.category || 'other';
    if (!acc[cat]) acc[cat] = [];
    acc[cat].push({ key, ...val });
    return acc;
  }, {});

  $: filteredGroups = Object.entries(groupedShortcuts)
    .map(([cat, items]) => ({
      category: cat,
      name: categories[cat] || cat,
      shortcuts: items.filter(s => 
        !searchQuery || 
        s.desc.toLowerCase().includes(searchQuery.toLowerCase()) ||
        s.key.toLowerCase().includes(searchQuery.toLowerCase())
      )
    }))
    .filter(g => g.shortcuts.length > 0);

  function handleClose() {
    dispatch('close');
  }

  function handleOverlayClick(e) {
    if (e.target === e.currentTarget) {
      handleClose();
    }
  }
</script>

<div 
  class="overlay" 
  on:click={handleOverlayClick}
  on:keydown={(e) => e.key === 'Escape' && handleClose()}
  role="button"
  tabindex="0"
>
  <div class="dialog" role="dialog" aria-modal="true" on:click|stopPropagation>
    <div class="header">
      <h2 class="title">⌨️ 键盘快捷键</h2>
      <button class="closeBtn" on:click={handleClose} aria-label="关闭">
        ✕
      </button>
    </div>

    <div class="searchBox">
      <input 
        type="text" 
        placeholder="搜索快捷键..." 
        bind:value={searchQuery}
        class="searchInput"
      />
    </div>

    <div class="content">
      {#if filteredGroups.length > 0}
        {#each filteredGroups as group}
          <div class="group">
            <h3 class="groupTitle">{group.name}</h3>
            <div class="shortcutList">
              {#each group.shortcuts as shortcut}
                <div class="shortcutItem">
                  <span class="desc">{shortcut.desc}</span>
                  <kbd class="key">{formatShortcut(shortcut.key)}</kbd>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      {:else}
        <div class="empty">
          <p>未找到匹配的快捷键</p>
        </div>
      {/if}
    </div>

    <div class="footer">
      <p class="hint">按 <kbd>Esc</kbd> 关闭此面板</p>
    </div>
  </div>
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
    z-index: 9998;
    animation: fadeIn 0.2s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .dialog {
    width: 100%;
    max-width: 640px;
    max-height: 80vh;
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 16px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s ease;
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px;
    border-bottom: 1px solid var(--border);
  }

  .title {
    margin: 0;
    font-size: 20px;
    font-weight: 700;
    color: var(--text);
  }

  .closeBtn {
    border: none;
    background: transparent;
    color: var(--muted);
    cursor: pointer;
    padding: 8px;
    font-size: 20px;
    border-radius: 8px;
    transition: background 0.15s ease, color 0.15s ease;
  }

  .closeBtn:hover {
    background: rgba(148, 163, 184, 0.12);
    color: var(--text);
  }

  .searchBox {
    padding: 16px 24px;
    border-bottom: 1px solid var(--border);
  }

  .searchInput {
    width: 100%;
    padding: 10px 14px;
    border: 1px solid var(--border-2);
    border-radius: 10px;
    background: rgba(15, 23, 42, 0.06);
    color: var(--text);
    font-size: 14px;
    outline: none;
    transition: border-color 0.2s ease;
  }

  :global(:root[data-theme="light"]) .searchInput {
    background: rgba(15, 23, 42, 0.04);
  }

  .searchInput:focus {
    border-color: rgba(34, 197, 94, 0.55);
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 16px 24px;
  }

  .group {
    margin-bottom: 24px;
  }

  .group:last-child {
    margin-bottom: 0;
  }

  .groupTitle {
    margin: 0 0 12px 0;
    font-size: 14px;
    font-weight: 600;
    color: var(--muted);
  }

  .shortcutList {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .shortcutItem {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 12px;
    border-radius: 10px;
    background: rgba(15, 23, 42, 0.04);
    border: 1px solid transparent;
    transition: border-color 0.15s ease, background 0.15s ease;
  }

  .shortcutItem:hover {
    border-color: var(--border);
    background: rgba(15, 23, 42, 0.06);
  }

  .desc {
    font-size: 14px;
    color: var(--text);
  }

  .key, kbd {
    display: inline-flex;
    align-items: center;
    padding: 4px 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: rgba(15, 23, 42, 0.08);
    font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
    font-size: 12px;
    font-weight: 600;
    color: var(--text);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .footer {
    padding: 16px 24px;
    border-top: 1px solid var(--border);
    background: rgba(15, 23, 42, 0.02);
  }

  .hint {
    margin: 0;
    text-align: center;
    font-size: 12px;
    color: var(--muted);
  }

  .empty {
    text-align: center;
    padding: 48px 24px;
    color: var(--muted);
  }

  @media (max-width: 600px) {
    .dialog {
      max-height: 90vh;
      border-radius: 12px;
    }

    .header, .searchBox, .content, .footer {
      padding-left: 16px;
      padding-right: 16px;
    }

    .title {
      font-size: 18px;
    }
  }
</style>

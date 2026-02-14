<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import KeyboardShortcuts from './KeyboardShortcuts.svelte';
  import { exportData } from '../utils/exportImport.js';
  import { saveDraft, getDrafts } from '../utils/backup.js';

  const dispatch = createEventDispatcher();

  let showShortcuts = false;
  let searchFocused = false;
  let editorFocused = false;
  let inputFocused = false;

  const handlers = new Map();

  function handleKeydown(e) {
    // 忽略输入中的快捷键（除了 Ctrl+Enter, Ctrl+S）
    const isInput = e.target.matches('input, textarea, [contenteditable]');
    
    // Ctrl/Cmd + K - 搜索
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault();
      dispatch('focusSearch');
      return;
    }

    // ? - 显示帮助
    if (e.key === '?' && !isInput) {
      e.preventDefault();
      showShortcuts = true;
      return;
    }

    // Escape - 关闭弹窗
    if (e.key === 'Escape') {
      if (showShortcuts) {
        showShortcuts = false;
      } else {
        dispatch('closeAll');
      }
      return;
    }

    // Ctrl/Cmd + S - 保存
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
      e.preventDefault();
      dispatch('save');
      return;
    }

    // Ctrl/Cmd + Enter - 保存编辑器
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      if (isInput || editorFocused) {
        e.preventDefault();
        dispatch('saveEditor');
      }
      return;
    }

    // N - 新建笔记
    if (e.key === 'n' && !isInput && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      dispatch('newNote');
      return;
    }

    // E - 编辑
    if (e.key === 'e' && !isInput && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      dispatch('edit');
      return;
    }

    // D - 删除
    if (e.key === 'd' && !isInput && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      dispatch('delete');
      return;
    }

    // J/Down - 下移
    if ((e.key === 'j' || e.key === 'ArrowDown') && !isInput) {
      e.preventDefault();
      dispatch('navigate', { direction: 'down' });
      return;
    }

    // K/Up - 上移
    if ((e.key === 'k' || e.key === 'ArrowUp') && !isInput) {
      e.preventDefault();
      dispatch('navigate', { direction: 'up' });
      return;
    }

    // 1 - 信息流视图
    if (e.key === '1' && !isInput && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      dispatch('changeView', { mode: 'timeline' });
      return;
    }

    // 2 - 卡片视图
    if (e.key === '2' && !isInput && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      dispatch('changeView', { mode: 'waterfall' });
      return;
    }

    // B - 切换侧边栏
    if (e.key === 'b' && !isInput && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      dispatch('toggleSidebar');
      return;
    }

    // T - 显示标签
    if (e.key === 't' && !isInput && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      dispatch('showTags');
      return;
    }

    // / - 搜索标签
    if (e.key === '/' && !isInput) {
      e.preventDefault();
      dispatch('focusTagSearch');
      return;
    }

    // Ctrl/Cmd + E - 导出
    if ((e.ctrlKey || e.metaKey) && e.key === 'e') {
      e.preventDefault();
      exportData('markdown');
      return;
    }

    // Ctrl/Cmd + I - 导入
    if ((e.ctrlKey || e.metaKey) && e.key === 'i') {
      e.preventDefault();
      dispatch('import');
      return;
    }
  }

  // 注册编辑器内的快捷键处理
  export function registerEditorHandler(callback) {
    handlers.set('editor', callback);
  }

  export function unregisterEditorHandler() {
    handlers.delete('editor');
  }

  // 检测焦点状态
  function handleFocusIn(e) {
    const tag = e.target.tagName.toLowerCase();
    if (tag === 'input' || tag === 'textarea' || e.target.isContentEditable) {
      if (e.target.closest('.note-editor') || e.target.closest('[data-editor]')) {
        editorFocused = true;
      } else {
        inputFocused = true;
      }
    }
    searchFocused = e.target.closest('.search-bar') !== null;
  }

  function handleFocusOut(e) {
    const related = e.relatedTarget;
    if (!related || !e.target.contains(related)) {
      setTimeout(() => {
        editorFocused = false;
        inputFocused = false;
        searchFocused = false;
      }, 100);
    }
  }

  onMount(() => {
    document.addEventListener('keydown', handleKeydown);
    document.addEventListener('focusin', handleFocusIn);
    document.addEventListener('focusout', handleFocusOut);
    
    return () => {
      document.removeEventListener('keydown', handleKeydown);
      document.removeEventListener('focusin', handleFocusIn);
      document.removeEventListener('focusout', handleFocusOut);
    };
  });
</script>

<svelte:window />

<!-- 快捷键帮助弹窗 -->
{#if showShortcuts}
  <KeyboardShortcuts on:close={() => showShortcuts = false} />
{/if}

<!-- 快捷键提示浮层 -->
<div class="fixed bottom-4 right-4 z-40">
  <div class="flex items-center gap-2 px-3 py-2 rounded-full bg-card/80 backdrop-blur-sm border border-border/50 shadow-lg">
    <span class="text-xs text-muted-foreground">按</span>
    <kbd class="px-2 py-0.5 text-xs bg-muted rounded">?</kbd>
    <span class="text-xs text-muted-foreground">查看快捷键</span>
  </div>
</div>

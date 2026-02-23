import { writable } from 'svelte/store';

export const keyboardState = writable({
  showHelp: false,
  activeContext: 'global' // global, editor, search
});

class KeyboardManager {
  constructor() {
    this.handlers = new Map();
    this.context = 'global';
    this.enabled = true;
  }

  setContext(context) {
    this.context = context;
  }

  register(key, handler, context = 'global') {
    const keyString = `${context}:${key}`;
    if (!this.handlers.has(keyString)) {
      this.handlers.set(keyString, []);
    }
    this.handlers.get(keyString).push(handler);
  }

  unregister(key, context = 'global') {
    const keyString = `${context}:${key}`;
    this.handlers.delete(keyString);
  }

  handle(event) {
    if (!this.enabled) return false;

    const isInput = event.target.matches('input, textarea, [contenteditable]');
    let key = event.key.toLowerCase();
    
    // Build modifier prefix
    const modifiers = [];
    if (event.ctrlKey || event.metaKey) modifiers.push('ctrl');
    if (event.shiftKey) modifiers.push('shift');
    if (event.altKey) modifiers.push('alt');
    
    if (modifiers.length > 0) {
      key = `${modifiers.join('+')}+${key}`;
    }

    // Try context-specific handlers first
    const contextKey = `${this.context}:${key}`;
    if (this.handlers.has(contextKey)) {
      const handlers = this.handlers.get(contextKey);
      for (const handler of handlers) {
        if (handler(event) !== false) {
          event.preventDefault();
          return true;
        }
      }
    }

    // Try global handlers (except when in input for single character keys)
    // Allow Ctrl/Cmd/Alt combinations even in input fields
    const isModifierKey = key.includes('+');
    if (!isInput || isModifierKey) {
      const globalKey = `global:${key}`;
      if (this.handlers.has(globalKey)) {
        const handlers = this.handlers.get(globalKey);
        for (const handler of handlers) {
          if (handler(event) !== false) {
            event.preventDefault();
            return true;
          }
        }
      }
    }

    return false;
  }

  enable() {
    this.enabled = true;
  }

  disable() {
    this.enabled = false;
  }
}

export const keyboardManager = new KeyboardManager();

// Default shortcuts
export const shortcuts = {
  'ctrl+k': { desc: '聚焦搜索', category: 'navigation' },
  'ctrl+enter': { desc: '保存笔记', category: 'editor' },
  'escape': { desc: '关闭对话框', category: 'navigation' },
  'n': { desc: '新建笔记', category: 'note' },
  'e': { desc: '编辑笔记', category: 'note' },
  'd': { desc: '删除笔记', category: 'note' },
  'j': { desc: '下一条', category: 'navigation' },
  'k': { desc: '上一条', category: 'navigation' },
  '?': { desc: '显示快捷键', category: 'help' },
  '/': { desc: '快速搜索', category: 'navigation' },
  'b': { desc: '切换侧边栏', category: 'view' },
  'r': { desc: '刷新列表', category: 'navigation' },
};

export function formatShortcut(key) {
  return key
    .split('+')
    .map(k => {
      if (k === 'ctrl') return '⌘';
      if (k === 'shift') return '⇧';
      if (k === 'alt') return '⌥';
      return k.toUpperCase();
    })
    .join(' ');
}

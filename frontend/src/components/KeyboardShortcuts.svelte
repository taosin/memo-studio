<script>
  import { onMount, onDestroy } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/ui/card/card-title.svelte';

  export let onClose;

  let searchQuery = '';
  let filteredShortcuts = [];

  const shortcutGroups = [
    {
      name: '📝 笔记操作',
      shortcuts: [
        { keys: ['n', 'c'], description: '新建笔记', category: 'global' },
        { keys: ['e'], description: '编辑当前笔记', category: 'note' },
        { keys: ['d'], description: '删除当前笔记', category: 'note' },
        { keys: ['Enter'], description: '展开笔记详情', category: 'note' },
        { keys: ['Escape'], description: '关闭弹窗/编辑器', category: 'global' },
      ]
    },
    {
      name: '🔍 搜索与导航',
      shortcuts: [
        { keys: ['Ctrl', 'k'], description: '聚焦搜索栏', category: 'global' },
        { keys: ['/'], description: '搜索标签', category: 'search' },
        { keys: ['?'], description: '显示快捷键帮助', category: 'global' },
        { keys: ['j', '↓'], description: '向下选择', category: 'navigation' },
        { keys: ['k', '↑'], description: '向上选择', category: 'navigation' },
      ]
    },
    {
      name: '🏷️ 标签操作',
      shortcuts: [
        { keys: ['#'], description: '在编辑器中添加标签', category: 'editor' },
        { keys: ['t'], description: '显示标签列表', category: 'global' },
      ]
    },
    {
      name: '💾 编辑器',
      shortcuts: [
        { keys: ['Ctrl', 'Enter'], description: '保存笔记', category: 'editor' },
        { keys: ['Ctrl', 's'], description: '保存笔记', category: 'editor' },
        { keys: ['Tab'], description: '插入标签补全', category: 'editor' },
      ]
    },
    {
      name: '📋 列表操作',
      shortcuts: [
        { keys: ['a'], description: '全选笔记', category: 'list' },
        { keys: ['x'], description: '多选笔记', category: 'list' },
        { keys: ['m'], description: '移动笔记', category: 'list' },
      ]
    },
    {
      name: '🎨 视图切换',
      shortcuts: [
        { keys: ['1'], description: '信息流视图', category: 'view' },
        { keys: ['2'], description: '卡片视图', category: 'view' },
        { keys: ['b'], description: '收起/展开侧边栏', category: 'view' },
        { keys: ['f'], description: '全屏阅读', category: 'view' },
      ]
    },
    {
      name: '🔐 隐私与数据',
      shortcuts: [
        { keys: ['Ctrl', '\\'], description: '锁定应用', category: 'privacy' },
        { keys: ['Ctrl', 'e'], description: '导出数据', category: 'data' },
        { keys: ['Ctrl', 'i'], description: '导入数据', category: 'data' },
      ]
    }
  ];

  $: {
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filteredShortcuts = shortcutGroups.map(group => ({
        ...group,
        shortcuts: group.shortcuts.filter(s => 
          s.description.toLowerCase().includes(query) ||
          s.keys.some(k => k.toLowerCase().includes(query))
        )
      })).filter(group => group.shortcuts.length > 0);
    } else {
      filteredShortcuts = shortcutGroups;
    }
  }

  function formatKeys(keys) {
    return keys.map(k => {
      if (k === 'Ctrl') return '⌘';
      if (k === 'Meta') return '⌘';
      if (k === 'Shift') return '⇧';
      if (k === 'Alt') return '⌥';
      return k.toUpperCase();
    }).join(' ');
  }

  onMount(() => {
    // ESC 关闭
    const handleKeydown = (e) => {
      if (e.key === 'Escape' && !e.target.closest('input, textarea')) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });
</script>

<div class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4" on:click={onClose}>
  <Card class="w-full max-w-2xl max-h-[80vh] overflow-hidden" on:click|stopPropagation>
    <CardHeader class="pb-3 border-b">
      <div class="flex items-center justify-between">
        <CardTitle class="flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
            <line x1="8" y1="21" x2="16" y2="21"/>
            <line x1="12" y1="17" x2="12" y2="21"/>
          </svg>
          键盘快捷键
        </CardTitle>
        <button 
          class="text-muted-foreground hover:text-foreground transition-colors"
          on:click={onClose}
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      
      <!-- 搜索框 -->
      <div class="relative mt-3">
        <input
          type="text"
          placeholder="搜索快捷键..."
          bind:value={searchQuery}
          class="w-full px-4 py-2 pl-10 rounded-lg border border-border bg-background focus:outline-none focus:ring-2 focus:ring-primary/50"
        />
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
      </div>
    </CardHeader>

    <CardContent class="overflow-y-auto max-h-[60vh] p-4">
      <div class="space-y-6">
        {#each filteredShortcuts as group}
          <div>
            <h4 class="text-sm font-medium text-muted-foreground mb-3">{group.name}</h4>
            <div class="grid grid-cols-1 gap-2">
              {#each group.shortcuts as shortcut}
                <div class="flex items-center justify-between py-2 border-b border-border/50 last:border-0">
                  <span class="text-sm">{shortcut.description}</span>
                  <div class="flex gap-1">
                    {#each shortcut.keys as key}
                      <kbd class="px-2 py-1 text-xs font-medium bg-muted rounded-md border border-border">
                        {key === ' ' ? 'Space' : key}
                      </kbd>
                    {/each}
                  </div>
                </div>
              {/each}
            </div>
          </div>
        {/each}
        
        {#if filteredShortcuts.length === 0}
          <div class="text-center py-8 text-muted-foreground">
            <p>未找到匹配的快捷键</p>
          </div>
        {/if}
      </div>
    </CardContent>

    <!-- 底部提示 -->
    <div class="p-3 border-t bg-muted/30 text-center text-xs text-muted-foreground">
      按 <kbd class="px-1.5 py-0.5 bg-muted rounded">Esc</kbd> 关闭此面板
    </div>
  </Card>
</div>

<style>
  kbd {
    font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
  }
</style>

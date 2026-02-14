<script>
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { api } from '../utils/api.js';
  import Button from '$lib/components/ui/button/button.svelte';

  export let note = null;
  export let mode = 'create'; // 'create' | 'edit'
  
  const dispatch = createEventDispatcher();

  let content = '';
  let title = '';
  let tags = '';
  let allTags = [];
  let loading = false;
  let textareaRef;
  let showTags = false;
  let tagInputFocused = false;

  // 自动调整 textarea 高度
  function autoResize() {
    if (textareaRef) {
      textareaRef.style.height = 'auto';
      textareaRef.style.height = Math.min(textareaRef.scrollHeight, 200) + 'px';
    }
  }

  onMount(async () => {
    if (note) {
      title = note.title || '';
      content = typeof note.content === 'string' ? note.content : '';
      tags = (note.tags || []).map(t => t.name).join(', ');
    }
    await loadTags();
    
    // 自动聚焦
    setTimeout(() => {
      if (textareaRef) {
        textareaRef.focus();
        autoResize();
      }
    }, 100);
  });

  async function loadTags() {
    try {
      allTags = await api.getTags();
    } catch (err) {
      console.error('加载标签失败:', err);
    }
  }

  function handleInput(e) {
    content = e.target.value;
    autoResize();
    
    // 检测 # 标签触发
    const lastWord = content.split(/[\s#]/).pop();
    if (lastWord && !tagInputFocused) {
      showTags = lastWord.length > 0;
    }
  }

  function handleKeydown(e) {
    // Ctrl + Enter 保存
    if (e.ctrlKey && e.key === 'Enter') {
      handleSave();
    }
    // Escape 取消
    if (e.key === 'Escape') {
      dispatch('cancel');
    }
  }

  function handleTagClick(tagName) {
    const lastHashIndex = content.lastIndexOf('#');
    if (lastHashIndex >= 0) {
      // 替换当前输入的标签
      const before = content.substring(0, lastHashIndex);
      content = before + '#' + tagName + ' ';
    } else {
      content += ' #' + tagName + ' ';
    }
    showTags = false;
    if (textareaRef) textareaRef.focus();
    autoResize();
  }

  async function handleSave() {
    const safeContent = content.trim();
    const safeTitle = title.trim();
    
    // 提取内容中的标签
    const tagMatches = safeContent.match(/#([\w\u4e00-\u9fa5]+)/g) || [];
    const contentTags = tagMatches.map(t => t.substring(1));
    const manualTags = (tags || '').split(',').map(t => t.trim()).filter(t => t);
    const tagList = [...new Set([...manualTags, ...contentTags])];

    if (!safeContent && !safeTitle) {
      return; // 至少要有内容
    }

    loading = true;
    try {
      if (mode === 'edit' && note?.id) {
        await api.updateNote(note.id, safeTitle, safeContent, tagList);
      } else {
        await api.createNote(safeTitle, safeContent, tagList);
      }
      dispatch('save');
    } catch (err) {
      console.error('保存失败:', err);
      alert('保存失败: ' + (err.message || '未知错误'));
    } finally {
      loading = false;
    }
  }

  function handleCancel() {
    if (content.trim() || title.trim()) {
      if (!confirm('确定放弃此次记录吗？')) {
        return;
      }
    }
    dispatch('cancel');
  }

  // 格式化相对时间
  function formatTime() {
    const now = new Date();
    const hour = now.getHours();
    if (hour >= 5 && hour < 12) return '早上好';
    if (hour >= 12 && hour < 14) return '中午好';
    if (hour >= 14 && hour < 18) return '下午好';
    if (hour >= 18 && hour < 22) return '晚上好';
    return '夜深了';
  }
</script>

<!-- 遮罩层 -->
<div 
  class="fixed inset-0 bg-black/20 backdrop-blur-sm z-50 animate-fade-in"
  on:click={handleCancel}
  on:keydown={(e) => e.key === 'Escape' && handleCancel()}
  role="button"
  tabindex="-1"
></div>

<!-- 编辑器面板 -->
<div class="fixed bottom-0 left-0 right-0 z-50 animate-slide-up">
  <div class="max-w-3xl mx-auto bg-card border-t border-border rounded-t-2xl shadow-2xl overflow-hidden">
    <!-- 头部 -->
    <div class="flex items-center justify-between px-4 py-3 border-b border-border/50">
      <div class="flex items-center gap-2 text-sm text-muted-foreground">
        <span>{formatTime()}</span>
        <span class="w-1 h-1 rounded-full bg-muted-foreground/50"></span>
        <span class="text-xs">{new Date().toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', weekday: 'short' })}</span>
      </div>
      <div class="flex items-center gap-2">
        {#if loading}
          <span class="text-sm text-muted-foreground">保存中...</span>
        {/if}
        <Button variant="ghost" size="sm" on:click={handleCancel} disabled={loading}>
          取消
        </Button>
        <Button 
          variant="gradient" 
          size="sm" 
          on:click={handleSave} 
          disabled={loading || (!content.trim() && !title.trim())}
          loading={loading}
        >
          {mode === 'edit' ? '更新' : '保存'}
        </Button>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="p-4 space-y-3">
      <!-- 标题输入（可选） -->
      <input
        type="text"
        bind:value={title}
        placeholder="标题（可选）"
        class="w-full bg-transparent border-none outline-none text-lg font-medium placeholder:text-muted-foreground/50"
        disabled={loading}
      />

      <!-- 主要输入框 -->
      <div class="relative">
        <textarea
          bind:this={textareaRef}
          value={content}
          on:input={handleInput}
          on:keydown={handleKeydown}
          on:focus={() => tagInputFocused = true}
          on:blur={() => setTimeout(() => tagInputFocused = false, 200)}
          placeholder="记录你的想法... (支持 #标签)"
          class="w-full min-h-[80px] max-h-[200px] bg-transparent border-none outline-none resize-none text-base leading-relaxed placeholder:text-muted-foreground/50"
          disabled={loading}
        ></textarea>

        <!-- 标签建议下拉框 -->
        {#if showTags && allTags.length > 0}
          <div class="absolute bottom-full left-0 mb-2 p-2 bg-popover border border-border rounded-lg shadow-xl animate-fade-in max-h-48 overflow-y-auto">
            <p class="text-xs text-muted-foreground mb-2 px-2">选择标签</p>
            <div class="flex flex-wrap gap-1">
              {#each allTags as tag}
                <button
                  class="px-3 py-1 rounded-full text-sm transition-all hover:scale-105"
                  style="background-color: {tag.color || '#4ECDC4'}20; border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
                  on:click={() => handleTagClick(tag.name)}
                >
                  #{tag.name}
                </button>
              {/each}
            </div>
          </div>
        {/if}
      </div>

      <!-- 底部工具栏 -->
      <div class="flex items-center justify-between pt-2 border-t border-border/30">
        <div class="flex items-center gap-4 text-muted-foreground">
          <!-- 快捷提示 -->
          <span class="text-xs flex items-center gap-1">
            <kbd class="px-1.5 py-0.5 bg-muted rounded text-xs">#</kbd>
            添加标签
          </span>
          <span class="text-xs flex items-center gap-1 hidden sm:inline">
            <kbd class="px-1.5 py-0.5 bg-muted rounded text-xs">Ctrl+Enter</kbd>
            保存
          </span>
        </div>
        
        <!-- 字数统计 -->
        <span class="text-xs text-muted-foreground">
          {content.length} 字
        </span>
      </div>
    </div>
  </div>
</div>

<style>
  textarea::-webkit-scrollbar {
    width: 4px;
  }
  
  textarea::-webkit-scrollbar-thumb {
    background-color: hsl(var(--muted-foreground) / 0.3);
    border-radius: 4px;
  }
  
  textarea::-webkit-scrollbar-track {
    background: transparent;
  }
</style>

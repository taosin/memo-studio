<script>
  import { createEventDispatcher } from 'svelte';

  export let note;
  const dispatch = createEventDispatcher();

  let isHovered = false;

  function getContentPreview(content) {
    if (!content) return '';
    
    // 统一转换为字符串并去除 HTML 标签
    const textContent = String(content)
      .replace(/<[^>]*>/g, '')
      .trim();
    
    return textContent.length > 150 
      ? textContent.substring(0, 150) + '...' 
      : textContent;
  }

  function handleClick() {
    dispatch('click');
  }

  function handleDoubleClick() {
    dispatch('doubleClick');
  }

  function handleTagClick(tag, event) {
    event.stopPropagation();
    dispatch('tagClick', { tag, event });
  }

  function formatTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffMinutes = Math.floor(diffTime / (1000 * 60));
    const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

    if (diffMinutes < 1) return '刚刚';
    if (diffMinutes < 60) return `${diffMinutes}分钟前`;
    if (diffHours < 24) return `${diffHours}小时前`;
    if (diffDays === 1) return '昨天';
    if (diffDays < 7) return `${diffDays}天前`;
    return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
  }
</script>

<div
  class="cursor-pointer transition-all duration-300 ease-out group"
  on:mouseenter={() => isHovered = true}
  on:mouseleave={() => isHovered = false}
  on:click={handleClick}
  on:dblclick={handleDoubleClick}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
>
  <div 
    class="relative p-5 rounded-2xl transition-all duration-300"
    class:bg-card={!isHovered}
    class:bg-card/80={isHovered}
    class:shadow-sm={!isHovered}
    class:shadow-lg={isHovered}
    class:shadow-primary/5={isHovered}
    class:border={!isHovered}
    class:border-border/40={!isHovered}
    class:border-primary/20={isHovered}
  >
    <!-- 标题 -->
    {#if note.title}
      <h3 class="font-semibold text-foreground mb-3 leading-snug transition-colors group-hover:text-primary">
        {note.title}
      </h3>
    {/if}

    <!-- 内容预览 -->
    {#if getContentPreview(note.content)}
      <p class="text-muted-foreground text-sm leading-relaxed line-clamp-3 mb-4">
        {getContentPreview(note.content)}
      </p>
    {/if}

    <!-- 标签 -->
    {#if note.tags && note.tags.length > 0}
      <div class="flex flex-wrap gap-1.5 mb-4">
        {#each note.tags.slice(0, 5) as tag}
          <button
            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium cursor-pointer transition-all duration-200 hover:scale-105"
            style="background-color: {tag.color || '#4ECDC4'}15; color: {tag.color || '#4ECDC4'}"
            on:click={(e) => handleTagClick(tag, e)}
          >
            #{tag.name}
          </button>
        {/each}
        {#if note.tags.length > 5}
          <span class="text-xs text-muted-foreground self-center">+{note.tags.length - 5}</span>
        {/if}
      </div>
    {/if}

    <!-- 底部信息 -->
    <div class="flex items-center justify-between pt-3 border-t border-border/30">
      <div class="flex items-center gap-2">
        <span class="text-xs text-muted-foreground/60">
          {formatTime(note.created_at)}
        </span>
      </div>
      
      <div class="flex items-center gap-2 text-muted-foreground/40">
        {#if note.tags && note.tags.length > 0}
          <div class="flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
            <span class="text-xs">{note.tags.length}</span>
          </div>
        {/if}
        
        <!-- 悬停时显示的操作提示 -->
        <span class="text-xs opacity-0 group-hover:opacity-100 transition-opacity">
          点击查看
        </span>
      </div>
    </div>

    <!-- 悬停时的装饰线条 -->
    <div 
      class="absolute left-0 top-0 bottom-0 w-1 rounded-l-2xl bg-gradient-to-b from-primary to-primary/50 opacity-0 group-hover:opacity-100 transition-all duration-300"
    ></div>
    
    <!-- 双击编辑提示 -->
    <div class="absolute right-3 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity text-xs text-muted-foreground bg-background/80 px-2 py-1 rounded">
      双击编辑
    </div>
  </div>
</div>

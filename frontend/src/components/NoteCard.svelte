<script>
  import { createEventDispatcher } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  export let note;
  const dispatch = createEventDispatcher();

  let isHovered = false;

  function getContentPreview(content) {
    if (typeof content !== 'string' && content !== null && content !== undefined) {
      console.warn('NoteCard - content 不是字符串:', typeof content, content);
    }

    let safeContent = '';
    if (typeof content === 'string') {
      safeContent = content;
    } else if (content === null || content === undefined) {
      safeContent = '';
    } else {
      try {
        if (typeof content === 'object') {
          safeContent = JSON.stringify(content);
        } else {
          safeContent = String(content);
        }
      } catch (e) {
        safeContent = '';
      }
    }

    const textContent = safeContent.replace(/<[^>]*>/g, '').trim();
    if (textContent.length > 120) {
      return textContent.substring(0, 120) + '...';
    }
    return textContent || '无内容';
  }

  function handleClick() {
    dispatch('click');
  }

  function handleDoubleClick() {
    dispatch('doubleClick');
  }

  function handleTagClick(tag, event) {
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
  class:shadow-lg={isHovered}
  class:shadow-primary/10={isHovered}
  class:-translate-y-0.5={isHovered}
  class:scale-[1.01]={isHovered}
  on:mouseenter={() => isHovered = true}
  on:mouseleave={() => isHovered = false}
  on:click={handleClick}
  on:dblclick={handleDoubleClick}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
>
  <Card class="h-full border-border/40 hover:border-primary/30 hover:shadow-md transition-all duration-300 bg-card/50 backdrop-blur-sm">
    <CardContent class="p-4 space-y-3">
      <!-- 标题 -->
      {#if note.title}
        <h3 class="font-semibold text-foreground line-clamp-1 group-hover:text-primary transition-colors">
          {note.title}
        </h3>
      {/if}

      <!-- 内容预览 -->
      <p class="text-muted-foreground text-sm leading-relaxed line-clamp-3">
        {getContentPreview(note.content)}
      </p>

      <!-- 标签 -->
      {#if note.tags && note.tags.length > 0}
        <div class="flex flex-wrap gap-1.5 pt-1">
          {#each note.tags.slice(0, 4) as tag}
            <span
              class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium cursor-pointer hover:scale-105 transition-transform"
              style="background-color: {tag.color || '#4ECDC4'}15; color: {tag.color || '#4ECDC4'}"
              role="button"
              tabindex="0"
              on:click|stopPropagation={(e) => handleTagClick(tag, e)}
              on:keydown={(e) => e.key === 'Enter' && handleTagClick(tag, e)}
            >
              #{tag.name}
            </span>
          {/each}
          {#if note.tags.length > 4}
            <span class="text-xs text-muted-foreground self-center">+{note.tags.length - 4}</span>
          {/if}
        </div>
      {/if}

      <!-- 底部信息 -->
      <div class="flex items-center justify-between pt-2 border-t border-border/30">
        <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground/70">
            {formatTime(note.created_at)}
          </span>
        </div>
        
        {#if note.tags}
          <div class="flex items-center gap-1 text-muted-foreground/50">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
            <span class="text-xs">{note.tags.length}</span>
          </div>
        {/if}
      </div>
    </CardContent>
  </Card>
</div>

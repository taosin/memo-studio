<script>
  import { createEventDispatcher } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import CardFooter from '$lib/components/ui/card/card-footer.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  export let note;
  const dispatch = createEventDispatcher();

  let isHovered = false;

  // 安全地提取笔记内容的纯文本预览
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
        console.error('转换 content 失败:', e);
        safeContent = '';
      }
    }

    const textContent = safeContent.replace(/<[^>]*>/g, '').trim();
    if (textContent.length > 150) {
      return textContent.substring(0, 150) + '...';
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

  function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 0) {
      return '今天';
    } else if (diffDays === 1) {
      return '昨天';
    } else if (diffDays < 7) {
      return `${diffDays}天前`;
    } else {
      return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
    }
  }
</script>

<div
  class="cursor-pointer transition-all duration-300 ease-out group"
  class:shadow-lg={isHovered}
  class:-translate-y-1={isHovered}
  class:scale-[1.02]={isHovered}
  on:mouseenter={() => isHovered = true}
  on:mouseleave={() => isHovered = false}
  on:click={handleClick}
  on:dblclick={handleDoubleClick}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
>
  <Card class="h-full border-border/60 hover:border-primary hover:shadow-md transition-all duration-300">
    <CardContent class="p-4">
      <h3 class="text-lg font-semibold mb-2 text-card-foreground line-clamp-1 group-hover:text-primary transition-colors">
        {note.title || '无标题'}
      </h3>
      <p class="text-muted-foreground text-sm line-clamp-3 mb-3 leading-relaxed">
        {getContentPreview(note.content)}
      </p>

      {#if note.tags && note.tags.length > 0}
        <div class="flex flex-wrap gap-1.5 mb-3">
          {#each note.tags.slice(0, 3) as tag}
            <span
              class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium cursor-pointer hover:scale-105 transition-transform"
              style="background-color: {tag.color || '#4ECDC4'}20; border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
              role="button"
              tabindex="0"
              on:click|stopPropagation={(e) => handleTagClick(tag, e)}
              on:keydown={(e) => e.key === 'Enter' && handleTagClick(tag, e)}
            >
              #{tag.name}
            </span>
          {/each}
          {#if note.tags.length > 3}
            <span class="text-xs text-muted-foreground self-center">+{note.tags.length - 3}</span>
          {/if}
        </div>
      {/if}
    </CardContent>
    <CardFooter class="flex justify-between items-center pt-0 pb-3 px-4 border-t border-border/50">
      <div class="flex items-center gap-2">
        <span class="text-xs text-muted-foreground">
          {formatDate(note.created_at)}
        </span>
      </div>
      <div class="flex items-center gap-1 text-muted-foreground">
        {#if note.tags}
          <span class="text-xs flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
            {note.tags.length}
          </span>
        {/if}
      </div>
    </CardFooter>
  </Card>
</div>

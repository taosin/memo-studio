<script>
  import { createEventDispatcher } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import CardFooter from '$lib/components/ui/card/card-footer.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  
  export let note;
  const dispatch = createEventDispatcher();
  
  // 安全地提取笔记内容的纯文本预览
  function getContentPreview(content) {
    // 调试：检查 content 的类型
    if (typeof content !== 'string' && content !== null && content !== undefined) {
      console.warn('NoteCard - content 不是字符串:', typeof content, content);
    }
    
    // 处理各种可能的类型
    let safeContent = '';
    if (typeof content === 'string') {
      safeContent = content;
    } else if (content === null || content === undefined) {
      safeContent = '';
    } else {
      // 如果是对象或其他类型，尝试转换
      try {
        if (typeof content === 'object') {
          // 如果是对象，可能是错误存储的数据，尝试提取文本
          safeContent = JSON.stringify(content);
        } else {
          safeContent = String(content);
        }
      } catch (e) {
        console.error('转换 content 失败:', e);
        safeContent = '';
      }
    }
    
    // 如果是 HTML，提取纯文本；否则直接使用
    const textContent = safeContent.replace(/<[^>]*>/g, '').trim();
    // 截取前 150 个字符
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
</script>

<div
  class="cursor-pointer transition-all hover:shadow-lg hover:-translate-y-0.5"
  on:click={handleClick}
  on:dblclick={handleDoubleClick}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
>
  <Card class="hover:border-primary">
    <CardContent class="p-3">
      <h3 class="text-lg font-semibold mb-2 text-card-foreground">
        {note.title || '无标题'}
      </h3>
      <p class="text-muted-foreground text-sm line-clamp-3 mb-3">
        {getContentPreview(note.content)}
      </p>
    </CardContent>
    <CardFooter class="flex justify-between items-center pt-0 pb-3 px-3 border-t">
      <div class="flex flex-wrap gap-1.5">
        {#each note.tags || [] as tag}
          <span
            class="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold cursor-pointer hover:bg-accent transition-colors"
            style="border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
            role="button"
            tabindex="0"
            on:click|stopPropagation={(e) => handleTagClick(tag, e)}
            on:keydown={(e) => e.key === 'Enter' && handleTagClick(tag, e)}
          >
            {tag.name}
          </span>
        {/each}
      </div>
      <span class="text-xs text-muted-foreground">
        {new Date(note.created_at).toLocaleDateString('zh-CN')}
      </span>
    </CardFooter>
  </Card>
</div>

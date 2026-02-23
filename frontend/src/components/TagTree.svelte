<script>
  import { createEventDispatcher } from 'svelte';
  import { api } from '../utils/api.js';
  import { onMount } from 'svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  export let selectedTags = [];
  const dispatch = createEventDispatcher();

  let tags = [];
  let collapsed = false;

  onMount(async () => {
    await loadTags();
  });

  async function loadTags() {
    try {
      const allTagsData = await api.getTags();
      // 计算每个标签的使用次数
      const notes = await api.getNotes();
      tags = allTagsData.map(tag => {
        const count = notes.filter(note => 
          (note.tags || []).some(t => t.id === tag.id)
        ).length;
        return { ...tag, count };
      });
    } catch (err) {
      console.error('加载标签失败:', err);
    }
  }

  function toggleCollapse() {
    collapsed = !collapsed;
  }

  function handleTagClick(tag) {
    dispatch('tagSelect', tag);
  }

  function isTagSelected(tagId) {
    return selectedTags.includes(tagId);
  }
</script>

<div class="w-full">
  <button 
    type="button"
    class="flex items-center justify-between p-3 cursor-pointer hover:bg-accent rounded-md transition-colors w-full text-left"
    on:click={toggleCollapse}
  >
    <h3 class="font-semibold text-sm">标签</h3>
    <span class="text-xs text-muted-foreground">
      {collapsed ? '▶' : '▼'}
    </span>
  </button>
  
  {#if !collapsed}
    <div class="space-y-1 mt-2">
      {#each tags as tag}
        <button
          type="button"
          class="flex items-center gap-2 p-2 rounded-md cursor-pointer transition-colors w-full text-left {isTagSelected(tag.id) ? 'bg-primary/10 border border-primary' : 'hover:bg-accent'}"
          on:click={() => handleTagClick(tag)}
        >
          <Badge
            variant="outline"
            class="text-xs"
            style="border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
          >
            {tag.name}
          </Badge>
          <span class="text-xs text-muted-foreground ml-auto">
            {tag.count || 0}
          </span>
        </button>
      {/each}
    </div>
  {/if}
</div>

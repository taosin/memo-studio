<script>
  import { createEventDispatcher } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import CardFooter from '$lib/components/ui/card/card-footer.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  
  export let note;
  const dispatch = createEventDispatcher();
  
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

<Card 
  class="cursor-pointer transition-all hover:shadow-lg hover:-translate-y-0.5 hover:border-primary"
  on:click={handleClick}
  on:dblclick={handleDoubleClick}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
>
  <CardContent class="p-3">
    <h3 class="text-lg font-semibold mb-2 text-card-foreground">
      {note.title || '无标题'}
    </h3>
    <p class="text-muted-foreground text-sm line-clamp-3 mb-3">
      {(note.content || '').substring(0, 150)}{(note.content || '').length > 150 ? '...' : ''}
    </p>
  </CardContent>
  <CardFooter class="flex justify-between items-center pt-0 pb-3 px-3 border-t">
    <div class="flex flex-wrap gap-1.5">
      {#each note.tags || [] as tag}
        <Badge 
          variant="outline" 
          class="cursor-pointer hover:bg-accent"
          style="border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
          on:click={(e) => handleTagClick(tag, e)}
        >
          {tag.name}
        </Badge>
      {/each}
    </div>
    <span class="text-xs text-muted-foreground">
      {new Date(note.created_at).toLocaleDateString('zh-CN')}
    </span>
  </CardFooter>
</Card>

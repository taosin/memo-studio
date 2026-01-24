<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { createEventDispatcher } from 'svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/ui/card/card-title.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import CardFooter from '$lib/components/ui/card/card-footer.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  export let noteId;
  const dispatch = createEventDispatcher();

  let note = null;
  let loading = true;
  let error = null;

  onMount(async () => {
    await loadNote();
  });

  async function loadNote() {
    try {
      loading = true;
      note = await api.getNote(noteId);
      error = null;
    } catch (err) {
      error = err.message;
      console.error('加载笔记失败:', err);
    } finally {
      loading = false;
    }
  }

  function handleEdit() {
    dispatch('edit', note);
  }

  function handleBack() {
    dispatch('back');
  }

  async function handleDelete() {
    if (!confirm('确定要删除这条笔记吗？')) {
      return;
    }

    try {
      await api.deleteNote(noteId);
      dispatch('deleted');
      dispatch('back');
    } catch (err) {
      alert('删除失败: ' + err.message);
    }
  }
</script>

<div class="w-full max-w-3xl mx-auto">
  {#if loading}
    <div class="text-center py-12 text-muted-foreground">加载中...</div>
  {:else if error}
    <div class="text-center py-12 text-destructive">错误: {error}</div>
  {:else if note}
    <div class="flex justify-between mb-6">
      <Button variant="outline" on:click={handleBack}>← 返回</Button>
      <div class="flex gap-2">
        <Button variant="outline" on:click={handleEdit}>编辑</Button>
        <Button variant="destructive" on:click={handleDelete}>删除</Button>
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle class="text-3xl">{note.title || '无标题'}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-3 p-3">
        <div class="flex flex-wrap justify-between items-center gap-4 pb-3 border-b">
          <div class="flex flex-wrap gap-2">
            {#each note.tags || [] as tag}
              <Badge 
                variant="outline"
                style="border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
              >
                {tag.name}
              </Badge>
            {/each}
          </div>
          <span class="text-sm text-muted-foreground">
            {new Date(note.created_at).toLocaleString('zh-CN')}
          </span>
        </div>
        <div class="prose prose-sm dark:prose-invert max-w-none break-words" innerHTML={note.content}></div>
      </CardContent>
    </Card>
  {/if}
</div>

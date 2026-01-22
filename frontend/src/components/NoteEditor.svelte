<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { createEventDispatcher } from 'svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Input from '$lib/components/ui/input/input.svelte';
  import Textarea from '$lib/components/ui/textarea/textarea.svelte';
  import Label from '$lib/components/ui/label/label.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  export let note = null; // 如果存在则是编辑模式，否则是新建模式
  const dispatch = createEventDispatcher();

  let title = '';
  let content = '';
  let tags = '';
  let allTags = [];
  let loading = false;

  onMount(async () => {
    if (note) {
      // 编辑模式
      title = note.title || '';
      content = note.content || '';
      tags = (note.tags || []).map(t => t.name).join(',');
    }
    await loadTags();
  });

  async function loadTags() {
    try {
      allTags = await api.getTags();
    } catch (err) {
      console.error('加载标签失败:', err);
    }
  }

  async function handleSave() {
    if (!content.trim()) {
      alert('内容不能为空');
      return;
    }

    loading = true;
    try {
      const tagList = tags.split(',').map(t => t.trim()).filter(t => t);
      
      if (note) {
        // 编辑模式 - 这里可以后续实现更新接口
        await api.createNote(title, content, tagList);
      } else {
        // 新建模式
        await api.createNote(title, content, tagList);
      }
      
      dispatch('save');
    } catch (err) {
      alert('保存失败: ' + err.message);
      console.error('保存笔记失败:', err);
    } finally {
      loading = false;
    }
  }

  function handleCancel() {
    dispatch('cancel');
  }
</script>

<div class="w-full max-w-3xl mx-auto">
  <div class="flex justify-between mb-6">
    <Button variant="outline" on:click={handleCancel}>← 取消</Button>
    <Button on:click={handleSave} disabled={loading}>
      {loading ? '保存中...' : '保存'}
    </Button>
  </div>

  <Card>
    <CardContent class="p-6 space-y-6">
      <div>
        <Input
          className="text-2xl font-semibold"
          placeholder="标题（可选）"
          bind:value={title}
        />
      </div>

      <div class="space-y-2">
        <Label>标签（用逗号分隔）</Label>
        <Input
          placeholder="例如: 工作, 学习, 生活"
          bind:value={tags}
        />
        {#if allTags.length > 0}
          <div class="flex flex-wrap gap-2 items-center mt-2">
            <span class="text-sm text-muted-foreground">常用标签:</span>
            {#each allTags as tag}
              <Badge
                variant="outline"
                class="cursor-pointer hover:bg-accent"
                style="border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
                on:click={() => {
                  const tagList = tags.split(',').map(t => t.trim()).filter(t => t);
                  if (!tagList.includes(tag.name)) {
                    tags = tagList.length > 0 ? tags + ', ' + tag.name : tag.name;
                  }
                }}
              >
                {tag.name}
              </Badge>
            {/each}
          </div>
        {/if}
      </div>

      <div>
        <Textarea
          className="min-h-[400px] text-base"
          placeholder="开始记录你的想法..."
          bind:value={content}
        />
      </div>
    </CardContent>
  </Card>
</div>

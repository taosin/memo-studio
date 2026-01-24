<script>
  import { createEventDispatcher } from 'svelte';
  import { onMount } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Input from '$lib/components/ui/input/input.svelte';
  import Label from '$lib/components/ui/label/label.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import { api } from '../utils/api.js';

  export let visible = false;
  export let searchHistory = [];
  const dispatch = createEventDispatcher();

  let keyword = '';
  let selectedTags = [];
  let dateFrom = '';
  let dateTo = '';
  let allTags = [];

  onMount(async () => {
    await loadTags();
  });

  async function loadTags() {
    try {
      allTags = await api.getTags();
    } catch (err) {
      console.error('加载标签失败:', err);
    }
  }

  function handleSearch() {
    dispatch('search', {
      keyword,
      tags: selectedTags,
      dateFrom,
      dateTo
    });
  }

  function handleClear() {
    keyword = '';
    selectedTags = [];
    dateFrom = '';
    dateTo = '';
    dispatch('clear');
  }

  function toggleTag(tagId) {
    if (selectedTags.includes(tagId)) {
      selectedTags = selectedTags.filter(id => id !== tagId);
    } else {
      selectedTags = [...selectedTags, tagId];
    }
  }

  function selectHistoryItem(item) {
    keyword = item.keyword || '';
    selectedTags = item.tags || [];
    dateFrom = item.dateFrom || '';
    dateTo = item.dateTo || '';
    handleSearch();
  }
</script>

{#if visible}
  <Card class="mb-4">
    <CardContent class="p-3 space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="font-semibold">高级搜索</h3>
        <Button variant="ghost" size="sm" on:click={() => dispatch('close')}>✕</Button>
      </div>

      <div>
        <Label>关键词</Label>
        <Input
          placeholder="搜索标题或内容"
          bind:value={keyword}
        />
      </div>

      <div>
        <Label>标签筛选</Label>
        <div class="flex flex-wrap gap-2 mt-2">
          {#each allTags as tag}
            <Badge
              variant={selectedTags.includes(tag.id) ? 'default' : 'outline'}
              class="cursor-pointer"
              style={selectedTags.includes(tag.id) ? '' : `border-color: ${tag.color}; color: ${tag.color}`}
              on:click={() => toggleTag(tag.id)}
            >
              {tag.name}
            </Badge>
          {/each}
        </div>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <Label>开始日期</Label>
          <Input type="date" bind:value={dateFrom} />
        </div>
        <div>
          <Label>结束日期</Label>
          <Input type="date" bind:value={dateTo} />
        </div>
      </div>

      {#if searchHistory.length > 0}
        <div>
          <Label>搜索历史</Label>
          <div class="flex flex-wrap gap-2 mt-2">
            {#each searchHistory.slice(0, 5) as item, index}
              <Button
                variant="outline"
                size="sm"
                on:click={() => selectHistoryItem(item)}
              >
                {item.keyword || '历史搜索 ' + (index + 1)}
              </Button>
            {/each}
          </div>
        </div>
      {/if}

      <div class="flex gap-2">
        <Button on:click={handleSearch} class="flex-1">搜索</Button>
        <Button variant="outline" on:click={handleClear}>清空</Button>
      </div>
    </CardContent>
  </Card>
{/if}

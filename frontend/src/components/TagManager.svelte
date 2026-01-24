<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import { api } from '../utils/api.js';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/ui/card/card-title.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Input from '$lib/components/ui/input/input.svelte';
  import Label from '$lib/components/ui/label/label.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  const dispatch = createEventDispatcher();

  let tags = [];
  let editingTag = null;
  let showEditDialog = false;
  let editName = '';
  let editColor = '';
  let mergeSource = null;
  let showMergeDialog = false;

  onMount(async () => {
    await loadTags();
  });

  async function loadTags() {
    try {
      tags = await api.getTags();
    } catch (err) {
      console.error('加载标签失败:', err);
    }
  }

  function startEdit(tag) {
    editingTag = tag;
    editName = tag.name;
    editColor = tag.color;
    showEditDialog = true;
  }

  function startMerge(tag) {
    mergeSource = tag;
    showMergeDialog = true;
  }

  async function saveEdit() {
    if (!editName.trim()) {
      alert('标签名称不能为空');
      return;
    }

    try {
      await api.updateTag(editingTag.id, editName, editColor);
      showEditDialog = false;
      editingTag = null;
      await loadTags();
      dispatch('updated');
    } catch (err) {
      alert('更新失败: ' + err.message);
    }
  }

  async function handleDelete(tag) {
    if (!confirm(`确定要删除标签 "${tag.name}" 吗？这将从所有笔记中移除该标签。`)) {
      return;
    }

    try {
      await api.deleteTag(tag.id);
      await loadTags();
      dispatch('updated');
    } catch (err) {
      alert('删除失败: ' + err.message);
    }
  }

  async function handleMerge(targetTag) {
    if (!confirm(`确定要将标签 "${mergeSource.name}" 合并到 "${targetTag.name}" 吗？`)) {
      return;
    }

    try {
      await api.mergeTags(mergeSource.id, targetTag.id);
      showMergeDialog = false;
      mergeSource = null;
      await loadTags();
      dispatch('updated');
    } catch (err) {
      alert('合并失败: ' + err.message);
    }
  }

  function cancelEdit() {
    showEditDialog = false;
    editingTag = null;
    editName = '';
    editColor = '';
  }

  function cancelMerge() {
    showMergeDialog = false;
    mergeSource = null;
  }
</script>

<div class="w-full">
  <div class="mb-4">
    <h3 class="text-lg font-semibold">标签管理</h3>
    <p class="text-sm text-muted-foreground">编辑、删除或合并标签</p>
  </div>

  <div class="space-y-3">
    {#each tags as tag}
      <Card>
        <CardContent class="p-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <Badge
                variant="outline"
                style="border-color: {tag.color}; color: {tag.color}"
              >
                {tag.name}
              </Badge>
              <span class="text-sm text-muted-foreground">
                使用 {tag.count || 0} 次
              </span>
            </div>
            <div class="flex gap-2">
              <Button variant="outline" size="sm" on:click={() => startEdit(tag)}>
                编辑
              </Button>
              <Button variant="outline" size="sm" on:click={() => startMerge(tag)}>
                合并
              </Button>
              <Button variant="destructive" size="sm" on:click={() => handleDelete(tag)}>
                删除
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    {/each}
  </div>

  <!-- 编辑对话框 -->
  {#if showEditDialog && editingTag}
    <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" on:click={cancelEdit}>
      <Card class="w-full max-w-md" on:click|stopPropagation>
        <CardHeader>
          <CardTitle>编辑标签</CardTitle>
        </CardHeader>
        <CardContent class="p-3 space-y-4">
          <div>
            <Label>标签名称</Label>
            <Input bind:value={editName} />
          </div>
          <div>
            <Label>标签颜色</Label>
            <div class="flex gap-2 mt-2">
              <Input type="color" bind:value={editColor} class="w-20 h-10" />
              <Input bind:value={editColor} placeholder="#4ECDC4" />
            </div>
          </div>
          <div class="flex gap-2">
            <Button on:click={saveEdit} class="flex-1">保存</Button>
            <Button variant="outline" on:click={cancelEdit}>取消</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  {/if}

  <!-- 合并对话框 -->
  {#if showMergeDialog && mergeSource}
    <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" on:click={cancelMerge}>
      <Card class="w-full max-w-md" on:click|stopPropagation>
        <CardHeader>
          <CardTitle>合并标签</CardTitle>
        </CardHeader>
        <CardContent class="p-3 space-y-4">
          <p class="text-sm text-muted-foreground">
            选择要将 "<strong>{mergeSource.name}</strong>" 合并到的目标标签：
          </p>
          <div class="space-y-2 max-h-60 overflow-y-auto">
            {#each tags.filter(t => t.id !== mergeSource.id) as tag}
              <Button
                variant="outline"
                class="w-full justify-start"
                on:click={() => handleMerge(tag)}
              >
                <Badge
                  variant="outline"
                  style="border-color: {tag.color}; color: {tag.color}"
                >
                  {tag.name}
                </Badge>
              </Button>
            {/each}
          </div>
          <Button variant="outline" on:click={cancelMerge} class="w-full">取消</Button>
        </CardContent>
      </Card>
    </div>
  {/if}
</div>

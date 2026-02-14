<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import { getTemplates, useTemplate, saveTemplate, deleteTemplate } from '../utils/smartTags.js';

  const dispatch = createEventDispatcher();

  let templates = [];
  let selectedTemplate = null;
  let showCreateForm = false;
  let newTemplate = {
    name: '',
    icon: '📝',
    content: '',
    tags: []
  };

  onMount(() => {
    templates = getTemplates();
  });

  function selectTemplate(template) {
    selectedTemplate = template;
    const noteData = useTemplate(template.id);
    dispatch('select', noteData);
  }

  function handleCreateTemplate() {
    if (newTemplate.name.trim() && newTemplate.content.trim()) {
      const saved = saveTemplate(newTemplate);
      templates = getTemplates();
      newTemplate = { name: '', icon: '📝', content: '', tags: [] };
      showCreateForm = false;
      selectTemplate(saved);
    }
  }

  function handleDeleteTemplate(e, templateId) {
    e.stopPropagation();
    if (confirm('确定删除此模板吗？')) {
      deleteTemplate(templateId);
      templates = getTemplates();
      if (selectedTemplate?.id === templateId) {
        selectedTemplate = null;
      }
    }
  }
</script>

<div class="space-y-4">
  <!-- 标题 -->
  <div class="flex items-center justify-between">
    <h3 class="text-lg font-semibold">选择模板</h3>
    <Button variant="ghost" size="sm" on:click={() => showCreateForm = !showCreateForm}>
      {showCreateForm ? '取消' : '+ 自定义'}
    </Button>
  </div>

  <!-- 创建模板表单 -->
  {#if showCreateForm}
    <Card class="border-primary/30">
      <CardContent class="p-4 space-y-3">
        <div class="flex gap-2">
          <input
            type="text"
            placeholder="模板名称"
            bind:value={newTemplate.name}
            class="flex-1 px-3 py-2 rounded-lg border border-border bg-background"
          />
          <select 
            bind:value={newTemplate.icon}
            class="px-3 py-2 rounded-lg border border-border bg-background"
          >
            {#each ['📝', '💡', '📅', '📚', '📦', '🔄', '🎯', '💰', '🏃', '🍎'] as emoji}
              <option value={emoji}>{emoji}</option>
            {/each}
          </select>
        </div>
        <textarea
          placeholder="模板内容（可以使用 {{date}} 等变量）"
          bind:value={newTemplate.content}
          rows="4"
          class="w-full px-3 py-2 rounded-lg border border-border bg-background resize-none"
        ></textarea>
        <div class="flex justify-end">
          <Button size="sm" on:click={handleCreateTemplate}>保存模板</Button>
        </div>
      </CardContent>
    </Card>
  {/if}

  <!-- 模板列表 -->
  <div class="grid grid-cols-2 gap-3">
    {#each templates as template}
      <button
        class="group relative p-4 rounded-xl border border-border bg-card hover:border-primary/50 transition-all text-left {selectedTemplate?.id === template.id ? 'ring-2 ring-primary border-primary' : ''}"
        on:click={() => selectTemplate(template)}
      >
        <!-- 删除按钮 -->
        {#if template.id.startsWith('custom-')}
          <button
            class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity text-destructive hover:bg-destructive/10 rounded-full p-1"
            on:click={(e) => handleDeleteTemplate(e, template.id)}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        {/if}

        <!-- 图标 -->
        <div class="text-2xl mb-2">{template.icon}</div>
        
        <!-- 名称 -->
        <div class="font-medium text-sm">{template.name}</div>
        
        <!-- 标签预览 -->
        <div class="flex gap-1 mt-2 flex-wrap">
          {#each (template.tags || []).slice(0, 3) as tag}
            <span class="text-xs px-1.5 py-0.5 rounded bg-primary/10 text-primary">
              #{tag}
            </span>
          {/each}
        </div>
      </button>
    {/each}
  </div>

  <!-- 使用提示 -->
  {#if selectedTemplate}
    <div class="text-sm text-muted-foreground bg-accent/50 rounded-lg p-3">
      <p>💡 提示：内容中的 <code class="text-primary">{`{{date}}`}</code> 会自动替换为当前日期</p>
    </div>
  {/if}
</div>

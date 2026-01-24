<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { createEventDispatcher } from 'svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Input from '$lib/components/ui/input/input.svelte';
  import Label from '$lib/components/ui/label/label.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import RichTextEditor from './RichTextEditor.svelte';

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
      // 调试：检查 note.content 的类型
      console.log('NoteEditor onMount - note.content type:', typeof note.content, 'value:', note.content);
      
      title = String(note.title || '');
      // 如果内容是HTML，直接使用；否则转换为HTML
      // 处理可能的对象类型
      let noteContent = '';
      if (typeof note.content === 'string') {
        noteContent = note.content;
      } else if (note.content !== null && note.content !== undefined) {
        // 如果是对象，尝试转换为字符串
        try {
          noteContent = typeof note.content === 'object' ? JSON.stringify(note.content) : String(note.content);
        } catch (e) {
          console.error('转换 content 失败:', e);
          noteContent = '';
        }
      }
      
      if (noteContent && !noteContent.includes('<')) {
        // 纯文本，转换为HTML（保留换行）
        content = noteContent.replace(/\n/g, '<br>');
      } else {
        content = noteContent;
      }
      tags = (note.tags || []).map(t => t.name).join(',');
    } else {
      // 新建模式，确保初始值为空字符串
      title = '';
      content = '';
      tags = '';
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

  function handleContentChange(e) {
    // 确保 content 始终是字符串类型
    content = String(e.detail || '');
  }

  async function handleSave() {
    // 确保 content 和 title 是字符串（使用 String() 强制转换，处理 null/undefined）
    const safeContent = String(content || '');
    const safeTitle = String(title || '');
    
    // 从富文本内容中提取纯文本用于验证
    const textContent = safeContent.replace(/<[^>]*>/g, '').trim();
    const titleText = safeTitle.trim();
    
    if (!textContent && !titleText) {
      alert('标题和内容不能同时为空');
      return;
    }

    loading = true;
    try {
      // 从内容中提取标签（#标签格式，支持中文）
      const tagMatches = safeContent.match(/#([\w\u4e00-\u9fa5]+)/g) || [];
      const contentTags = tagMatches.map(match => match.substring(1));
      
      // 合并手动输入的标签和内容中的标签
      const manualTags = (tags || '').split(',').map(t => t.trim()).filter(t => t);
      const tagList = [...new Set([...manualTags, ...contentTags])];
      
      // 确保至少有一个字段不为空（后端会验证）
      const finalTitle = titleText || '';
      const finalContent = safeContent.trim() || '';
      
      console.log('保存笔记:', { 
        mode: note && note.id ? 'edit' : 'create',
        title: finalTitle,
        titleLength: finalTitle.length,
        content: finalContent.substring(0, 100), // 只显示前100个字符
        contentLength: finalContent.length,
        contentType: typeof finalContent,
        tags: tagList 
      });
      
      if (note && note.id) {
        // 编辑模式 - 使用更新接口
        const result = await api.updateNote(note.id, finalTitle, finalContent, tagList);
        console.log('更新成功:', result);
      } else {
        // 新建模式
        const result = await api.createNote(finalTitle, finalContent, tagList);
        console.log('创建成功:', result);
      }
      
      dispatch('save');
    } catch (err) {
      console.error('保存笔记失败 - 详细错误:', err);
      console.error('请求数据:', { 
        title: safeTitle.trim(), 
        contentLength: safeContent.trim().length,
        hasNote: !!note,
        noteId: note?.id 
      });
      alert('保存失败: ' + (err.message || '未知错误'));
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
    <CardContent class="p-3 space-y-4">
      <div>
        <Input
          className="text-2xl font-semibold"
          placeholder="标题（可选）"
          bind:value={title}
        />
      </div>

      <div class="space-y-2">
        <Label>标签（用逗号分隔，或在内容中使用 #标签）</Label>
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
        <Label class="mb-2 block">内容（支持富文本，输入 # 选择标签，输入 @ 引用笔记）</Label>
        <RichTextEditor
          value={content}
          placeholder="开始记录你的想法... 输入 # 选择标签，输入 @ 引用笔记"
          on:input={handleContentChange}
        />
      </div>
    </CardContent>
  </Card>
</div>

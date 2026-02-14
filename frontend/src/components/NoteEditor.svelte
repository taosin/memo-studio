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

  export let note = null;
  const dispatch = createEventDispatcher();

  let title = '';
  let content = '';
  let tags = '';
  let allTags = [];
  let loading = false;
  let showTagSuggestions = false;

  onMount(async () => {
    if (note) {
      console.log('NoteEditor onMount - note.content type:', typeof note.content, 'value:', note.content);
      
      title = String(note.title || '');
      let noteContent = '';
      if (typeof note.content === 'string') {
        noteContent = note.content;
      } else if (note.content !== null && note.content !== undefined) {
        try {
          noteContent = typeof note.content === 'object' ? JSON.stringify(note.content) : String(note.content);
        } catch (e) {
          console.error('转换 content 失败:', e);
          noteContent = '';
        }
      }
      
      if (noteContent && !noteContent.includes('<')) {
        content = noteContent.replace(/\n/g, '<br>');
      } else {
        content = noteContent;
      }
      tags = (note.tags || []).map(t => t.name).join(',');
    } else {
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
    content = String(e.detail || '');
  }

  async function handleSave() {
    const safeContent = String(content || '');
    const safeTitle = String(title || '');
    const textContent = safeContent.replace(/<[^>]*>/g, '').trim();
    const titleText = safeTitle.trim();
    
    if (!textContent && !titleText) {
      alert('标题和内容不能同时为空');
      return;
    }

    loading = true;
    try {
      const tagMatches = safeContent.match(/#([\w\u4e00-\u9fa5]+)/g) || [];
      const contentTags = tagMatches.map(match => match.substring(1));
      const manualTags = (tags || '').split(',').map(t => t.trim()).filter(t => t);
      const tagList = [...new Set([...manualTags, ...contentTags])];
      
      const finalTitle = titleText || '';
      const finalContent = safeContent.trim() || '';
      
      console.log('保存笔记:', { 
        mode: note && note.id ? 'edit' : 'create',
        title: finalTitle,
        contentLength: finalContent.length,
        tags: tagList 
      });
      
      if (note && note.id) {
        const result = await api.updateNote(note.id, finalTitle, finalContent, tagList);
        console.log('更新成功:', result);
      } else {
        const result = await api.createNote(finalTitle, finalContent, tagList);
        console.log('创建成功:', result);
      }
      
      dispatch('save');
    } catch (err) {
      console.error('保存笔记失败:', err);
      alert('保存失败: ' + (err.message || '未知错误'));
    } finally {
      loading = false;
    }
  }

  function handleCancel() {
    dispatch('cancel');
  }

  function addTag(tagName) {
    const tagList = tags.split(',').map(t => t.trim()).filter(t => t);
    if (!tagList.includes(tagName)) {
      tags = tagList.length > 0 ? tags + ', ' + tagName : tagName;
    }
    showTagSuggestions = false;
  }

  function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<div class="w-full max-w-4xl mx-auto">
  <!-- 顶部导航 -->
  <div class="flex items-center justify-between mb-6">
    <button 
      on:click={handleCancel}
      class="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors group"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform">
        <path d="M19 12H5M12 19l-7-7 7-7"/>
      </svg>
      <span>返回</span>
    </button>
    <div class="flex items-center gap-3">
      {#if note}
        <span class="text-sm text-muted-foreground hidden sm:inline">
          最后修改: {formatDate(note.updated_at || note.created_at)}
        </span>
      {/if}
      <Button 
        variant="gradient" 
        on:click={handleSave} 
        loading={loading}
        class="shadow-md"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2">
          <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
          <polyline points="17 21 17 13 7 13 7 21"/>
          <polyline points="7 3 7 8 15 8"/>
        </svg>
        保存
      </Button>
    </div>
  </div>

  <!-- 编辑器卡片 -->
  <Card class="border-2 shadow-lg animate-scale-in">
    <CardContent class="p-6 space-y-6">
      <!-- 标题输入 -->
      <div class="relative">
        <input
          type="text"
          bind:value={title}
          placeholder="无标题笔记"
          class="w-full text-3xl font-bold bg-transparent border-none outline-none placeholder:text-muted-foreground/50 focus:ring-0"
        />
      </div>

      <!-- 标签管理 -->
      <div class="space-y-3">
        <div class="flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-muted-foreground">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
          <Label class="text-sm font-medium">标签</Label>
        </div>
        
        <div class="relative">
          <input
            type="text"
            bind:value={tags}
            placeholder="添加标签，用逗号分隔..."
            class="w-full h-10 px-4 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all"
            on:focus={() => showTagSuggestions = true}
            on:blur={() => setTimeout(() => showTagSuggestions = false, 200)}
          />
          
          <!-- 标签建议下拉框 -->
          {#if showTagSuggestions && allTags.length > 0}
            <div class="absolute top-full left-0 right-0 mt-2 p-3 bg-card border border-border rounded-lg shadow-xl z-50 animate-fade-in">
              <p class="text-xs text-muted-foreground mb-2">点击添加常用标签：</p>
              <div class="flex flex-wrap gap-2">
                {#each allTags as tag}
                  <button
                    class="px-3 py-1 rounded-full text-sm font-medium transition-all hover:scale-105"
                    style="background-color: {tag.color || '#4ECDC4'}20; border-color: {tag.color || '#4ECDC4'}; color: {tag.color || '#4ECDC4'}"
                    on:click={() => addTag(tag.name)}
                  >
                    + {tag.name}
                  </button>
                {/each}
              </div>
            </div>
          {/if}
        </div>

        <!-- 已选标签预览 -->
        {#if tags}
          <div class="flex flex-wrap gap-2 mt-3">
            {#each tags.split(',').map(t => t.trim()).filter(t => t) as tag}
              <span 
                class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium animate-fade-in"
                style="background-color: #4ECDC420; border-color: #4ECDC4; color: #4ECDC4"
              >
                #{tag}
                <button 
                  class="ml-1 hover:opacity-70 transition-opacity"
                  on:click={() => {
                    const tagList = tags.split(',').map(t => t.trim()).filter(t => t && t !== tag);
                    tags = tagList.join(', ');
                  }}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </span>
            {/each}
          </div>
        {/if}
      </div>

      <!-- 分隔线 -->
      <hr class="border-border" />

      <!-- 富文本编辑器 -->
      <div class="space-y-3">
        <div class="flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-muted-foreground">
            <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
            <line x1="10" y1="9" x2="8" y2="9"/>
          </svg>
          <Label class="text-sm font-medium">内容</Label>
          <span class="text-xs text-muted-foreground ml-auto">支持 Markdown 语法</span>
        </div>
        
        <div class="border-2 border-dashed border-border rounded-lg overflow-hidden hover:border-primary/50 transition-colors">
          <RichTextEditor
            value={content}
            placeholder="开始记录你的想法..."
            on:input={handleContentChange}
          />
        </div>
      </div>
    </CardContent>
  </Card>

  <!-- 底部提示 -->
  <div class="mt-6 text-center text-sm text-muted-foreground">
    <p>💡 提示：输入 # 即可快速添加标签</p>
  </div>
</div>

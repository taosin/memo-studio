<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { createEventDispatcher } from 'svelte';

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

<div class="note-editor-container">
  <div class="editor-header">
    <button class="btn btn-back" on:click={handleCancel}>← 取消</button>
    <button class="btn btn-save" on:click={handleSave} disabled={loading}>
      {loading ? '保存中...' : '保存'}
    </button>
  </div>

  <div class="editor-content">
    <input
      type="text"
      class="title-input"
      placeholder="标题（可选）"
      bind:value={title}
    />

    <div class="tags-section">
      <label class="label">标签（用逗号分隔）</label>
      <input
        type="text"
        class="tags-input"
        placeholder="例如: 工作, 学习, 生活"
        bind:value={tags}
      />
      {#if allTags.length > 0}
        <div class="suggested-tags">
          <span class="label">常用标签:</span>
          {#each allTags as tag}
            <button
              class="tag-suggestion"
              style="background-color: {tag.color || '#4ECDC4'}20; color: {tag.color || '#4ECDC4'}"
              on:click={() => {
                const tagList = tags.split(',').map(t => t.trim()).filter(t => t);
                if (!tagList.includes(tag.name)) {
                  tags = tagList.length > 0 ? tags + ', ' + tag.name : tag.name;
                }
              }}
            >
              {tag.name}
            </button>
          {/each}
        </div>
      {/if}
    </div>

    <textarea
      class="content-textarea"
      placeholder="开始记录你的想法..."
      bind:value={content}
    ></textarea>
  </div>
</div>

<style>
  .note-editor-container {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
  }

  .editor-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2rem;
  }

  .btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 8px;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
    font-weight: 500;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-back {
    background-color: var(--bg-secondary);
    color: var(--text-primary);
  }

  .btn-back:hover {
    background-color: var(--border-color);
  }

  .btn-save {
    background-color: var(--accent-color);
    color: white;
  }

  .btn-save:hover:not(:disabled) {
    background-color: var(--accent-hover);
  }

  .editor-content {
    background-color: var(--card-bg);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: var(--shadow);
  }

  .title-input {
    width: 100%;
    padding: 1rem;
    font-size: 1.5rem;
    font-weight: 600;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    background-color: var(--bg-primary);
    color: var(--text-primary);
    margin-bottom: 1.5rem;
    font-family: inherit;
  }

  .title-input:focus {
    outline: none;
    border-color: var(--accent-color);
  }

  .tags-section {
    margin-bottom: 1.5rem;
  }

  .label {
    display: block;
    font-size: 0.9rem;
    color: var(--text-secondary);
    margin-bottom: 0.5rem;
    font-weight: 500;
  }

  .tags-input {
    width: 100%;
    padding: 0.75rem;
    font-size: 0.9rem;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    background-color: var(--bg-primary);
    color: var(--text-primary);
    font-family: inherit;
    margin-bottom: 0.75rem;
  }

  .tags-input:focus {
    outline: none;
    border-color: var(--accent-color);
  }

  .suggested-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    align-items: center;
  }

  .tag-suggestion {
    padding: 0.3rem 0.7rem;
    border: none;
    border-radius: 12px;
    font-size: 0.8rem;
    cursor: pointer;
    transition: transform 0.2s ease;
  }

  .tag-suggestion:hover {
    transform: scale(1.05);
  }

  .content-textarea {
    width: 100%;
    min-height: 400px;
    padding: 1rem;
    font-size: 1rem;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    background-color: var(--bg-primary);
    color: var(--text-primary);
    font-family: inherit;
    line-height: 1.6;
    resize: vertical;
  }

  .content-textarea:focus {
    outline: none;
    border-color: var(--accent-color);
  }

  @media (max-width: 768px) {
    .editor-content {
      padding: 1.5rem;
    }

    .title-input {
      font-size: 1.2rem;
      padding: 0.8rem;
    }

    .content-textarea {
      min-height: 300px;
    }
  }
</style>

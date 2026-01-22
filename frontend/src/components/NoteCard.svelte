<script>
  import { createEventDispatcher } from 'svelte';
  
  export let note;
  const dispatch = createEventDispatcher();
  
  function handleClick() {
    dispatch('click');
  }
</script>

<div class="note-card" on:click={handleClick} role="button" tabindex="0" on:keydown={(e) => e.key === 'Enter' && handleClick()}>
  <h3 class="note-title">{note.title || '无标题'}</h3>
  <p class="note-preview">{(note.content || '').substring(0, 150)}{(note.content || '').length > 150 ? '...' : ''}</p>
  <div class="note-footer">
    <div class="tags">
      {#each note.tags || [] as tag}
        <span class="tag" style="background-color: {tag.color || '#4ECDC4'}20; color: {tag.color || '#4ECDC4'}">
          {tag.name}
        </span>
      {/each}
    </div>
    <span class="note-date">
      {new Date(note.created_at).toLocaleDateString('zh-CN')}
    </span>
  </div>
</div>

<style>
  .note-card {
    background-color: var(--card-bg);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    padding: 1.5rem;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: var(--shadow);
  }

  .note-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    border-color: var(--accent-color);
  }

  .note-title {
    font-size: 1.2rem;
    font-weight: 600;
    margin-bottom: 0.75rem;
    color: var(--text-primary);
    line-height: 1.4;
  }

  .note-preview {
    color: var(--text-secondary);
    line-height: 1.6;
    margin-bottom: 1rem;
    font-size: 0.9rem;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .note-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border-color);
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .tag {
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .note-date {
    font-size: 0.8rem;
    color: var(--text-secondary);
  }

  @media (max-width: 768px) {
    .note-card {
      padding: 1rem;
    }

    .note-title {
      font-size: 1.1rem;
    }
  }
</style>

<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';

  let stats = {
    totalNotes: 0,
    totalTags: 0,
    todayNotes: 0,
    weekNotes: 0
  };
  let loading = true;

  onMount(async () => {
    await loadStats();
  });

  async function loadStats() {
    try {
      // 获取所有笔记
      const notes = await api.getNotes().catch(() => []);
      const tags = await api.getTags().catch(() => []);
      
      const now = new Date();
      const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
      const weekAgo = new Date(today);
      weekAgo.setDate(weekAgo.getDate() - 7);

      const todayNotesList = notes.filter(n => new Date(n.created_at) >= today);
      const weekNotesList = notes.filter(n => new Date(n.created_at) >= weekAgo);

      stats = {
        totalNotes: notes.length,
        totalTags: tags.length,
        todayNotes: todayNotesList.length,
        weekNotes: weekNotesList.length
      };
    } catch (err) {
      console.error('加载统计失败:', err);
    } finally {
      loading = false;
    }
  }
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
  <!-- 总笔记数 -->
  <div class="bg-gradient-to-br from-primary/5 to-primary/10 rounded-2xl p-4 border border-primary/10">
    <div class="flex items-center justify-between mb-2">
      <span class="text-sm text-muted-foreground">总笔记</span>
      <div class="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-primary">
          <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/>
          <polyline points="14 2 14 8 20 8"/>
        </svg>
      </div>
    </div>
    {#if loading}
      <div class="h-8 w-16 bg-primary/20 rounded animate-pulse"></div>
    {:else}
      <div class="text-2xl font-bold text-primary">{stats.totalNotes}</div>
    {/if}
    <div class="text-xs text-muted-foreground mt-1">篇</div>
  </div>

  <!-- 今日笔记 -->
  <div class="bg-gradient-to-br from-green-500/5 to-green-500/10 rounded-2xl p-4 border border-green-500/10">
    <div class="flex items-center justify-between mb-2">
      <span class="text-sm text-muted-foreground">今日</span>
      <div class="w-8 h-8 rounded-full bg-green-500/20 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-green-500">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
      </div>
    </div>
    {#if loading}
      <div class="h-8 w-12 bg-green-500/20 rounded animate-pulse"></div>
    {:else}
      <div class="text-2xl font-bold text-green-500">{stats.todayNotes}</div>
    {/if}
    <div class="text-xs text-muted-foreground mt-1">篇新增</div>
  </div>

  <!-- 本周笔记 -->
  <div class="bg-gradient-to-br from-orange-500/5 to-orange-500/10 rounded-2xl p-4 border border-orange-500/10">
    <div class="flex items-center justify-between mb-2">
      <span class="text-sm text-muted-foreground">本周</span>
      <div class="w-8 h-8 rounded-full bg-orange-500/20 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-orange-500">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
          <line x1="16" y1="2" x2="16" y2="6"/>
          <line x1="8" y1="2" x2="8" y2="6"/>
          <line x1="3" y1="10" x2="21" y2="10"/>
        </svg>
      </div>
    </div>
    {#if loading}
      <div class="h-8 w-12 bg-orange-500/20 rounded animate-pulse"></div>
    {:else}
      <div class="text-2xl font-bold text-orange-500">{stats.weekNotes}</div>
    {/if}
    <div class="text-xs text-muted-foreground mt-1">篇新增</div>
  </div>

  <!-- 标签数 -->
  <div class="bg-gradient-to-br from-purple-500/5 to-purple-500/10 rounded-2xl p-4 border border-purple-500/10">
    <div class="flex items-center justify-between mb-2">
      <span class="text-sm text-muted-foreground">标签</span>
      <div class="w-8 h-8 rounded-full bg-purple-500/20 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-purple-500">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
        </svg>
      </div>
    </div>
    {#if loading}
      <div class="h-8 w-12 bg-purple-500/20 rounded animate-pulse"></div>
    {:else}
      <div class="text-2xl font-bold text-purple-500">{stats.totalTags}</div>
    {/if}
    <div class="text-xs text-muted-foreground mt-1">个标签</div>
  </div>
</div>

<style>
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
  }
  
  .animate-pulse {
    animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }
</style>

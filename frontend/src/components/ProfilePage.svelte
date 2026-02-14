<script>
  import { createEventDispatcher } from 'svelte';
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/ui/card/card-title.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Heatmap from './Heatmap.svelte';
  import TagManager from './TagManager.svelte';
  import PrivacySettings from './PrivacySettings.svelte';
  import { themeStore } from '../stores/theme.js';
  import { authStore } from '../stores/auth.js';

  const dispatch = createEventDispatcher();

  let activeTab = 'detail'; // 'detail', 'settings', 'stats', 'tags', 'privacy'
  let notes = [];
  let tags = [];
  let stats = {
    totalNotes: 0,
    totalTags: 0,
    totalWords: 0,
    avgNotesPerDay: 0
  };

  onMount(async () => {
    await loadData();
    calculateStats();
  });

  async function loadData() {
    try {
      [notes, tags] = await Promise.all([
        api.getNotes(),
        api.getTags()
      ]);
      calculateStats();
    } catch (err) {
      console.error('加载数据失败:', err);
    }
  }

  function calculateStats() {
    stats.totalNotes = notes.length;
    stats.totalTags = tags.length;
    
    stats.totalWords = notes.reduce((sum, note) => {
      const text = (note.content || '').replace(/<[^>]*>/g, '');
      return sum + text.length;
    }, 0);

    if (notes.length > 0) {
      const firstNote = notes[notes.length - 1];
      const daysDiff = Math.ceil(
        (new Date() - new Date(firstNote.created_at)) / (1000 * 60 * 60 * 24)
      );
      stats.avgNotesPerDay = daysDiff > 0 ? (notes.length / daysDiff).toFixed(2) : 0;
    }
  }

  function handleLogout() {
    if (confirm('确定要退出吗？')) {
      dispatch('logout');
    }
  }

  function handleThemeChange() {
    $themeStore = $themeStore === 'light' ? 'dark' : 'light';
  }
</script>

<div class="w-full max-w-4xl mx-auto">
  <div class="mb-4">
    <h2 class="text-2xl font-bold">个人信息</h2>
  </div>

  <!-- 标签页 -->
  <div class="flex gap-1 mb-4 p-1 bg-card/50 rounded-lg">
    {#each [
      { id: 'detail', label: '👤 详情', icon: '' },
      { id: 'settings', label: '⚙️ 设置', icon: '' },
      { id: 'stats', label: '📊 统计', icon: '' },
      { id: 'tags', label: '🏷️ 标签', icon: '' },
      { id: 'privacy', label: '🔒 隐私', icon: '' }
    ] as tab}
      <button
        class="flex-1 px-4 py-2 rounded-md text-sm font-medium transition-all {activeTab === tab.id ? 'bg-primary text-primary-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground hover:bg-accent'}"
        on:click={() => activeTab = tab.id}
      >
        {tab.label}
      </button>
    {/each}
  </div>

  <!-- 个人详情 -->
  {#if activeTab === 'detail'}
    <Card>
      <CardHeader>
        <CardTitle>个人详情</CardTitle>
      </CardHeader>
      <CardContent class="p-3 space-y-4">
        <div class="flex items-center gap-4">
          <div class="w-16 h-16 rounded-full bg-gradient-to-br from-primary to-primary-light flex items-center justify-center text-2xl shadow-lg">
            👤
          </div>
          <div>
            <h3 class="text-lg font-semibold">{$authStore.user?.username || '用户'}</h3>
            <p class="text-sm text-muted-foreground">{$authStore.user?.email || '未设置邮箱'}</p>
          </div>
        </div>
        <div class="space-y-2 pt-4 border-t">
          <div class="flex justify-between">
            <span class="text-muted-foreground">注册时间</span>
            <span>{$authStore.user?.created_at ? new Date($authStore.user.created_at).toLocaleDateString('zh-CN') : '-'}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">笔记总数</span>
            <span>{stats.totalNotes}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">标签总数</span>
            <span>{stats.totalTags}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  {/if}

  <!-- 偏好设置 -->
  {#if activeTab === 'settings'}
    <Card>
      <CardHeader>
        <CardTitle>偏好设置</CardTitle>
      </CardHeader>
      <CardContent class="p-3 space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h4 class="font-medium">主题模式</h4>
            <p class="text-sm text-muted-foreground">切换明暗主题</p>
          </div>
          <Button variant="outline" on:click={handleThemeChange}>
            {$themeStore === 'light' ? '🌙 暗色' : '☀️ 亮色'}
          </Button>
        </div>
      </CardContent>
    </Card>
  {/if}

  <!-- 记录统计 -->
  {#if activeTab === 'stats'}
    <div class="space-y-4">
      <!-- 统计卡片 -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <Card class="bg-gradient-to-br from-primary/5 to-primary/10 border-primary/20">
          <CardContent class="p-3 text-center">
            <div class="text-2xl font-bold text-primary">{stats.totalNotes}</div>
            <div class="text-sm text-muted-foreground mt-1">笔记总数</div>
          </CardContent>
        </Card>
        <Card class="bg-gradient-to-br from-green-500/5 to-green-500/10 border-green-500/20">
          <CardContent class="p-3 text-center">
            <div class="text-2xl font-bold text-green-500">{stats.totalTags}</div>
            <div class="text-sm text-muted-foreground mt-1">标签总数</div>
          </CardContent>
        </Card>
        <Card class="bg-gradient-to-br from-orange-500/5 to-orange-500/10 border-orange-500/20">
          <CardContent class="p-3 text-center">
            <div class="text-2xl font-bold text-orange-500">{stats.totalWords.toLocaleString()}</div>
            <div class="text-sm text-muted-foreground mt-1">总字数</div>
          </CardContent>
        </Card>
        <Card class="bg-gradient-to-br from-purple-500/5 to-purple-500/10 border-purple-500/20">
          <CardContent class="p-3 text-center">
            <div class="text-2xl font-bold text-purple-500">{stats.avgNotesPerDay}</div>
            <div class="text-sm text-muted-foreground mt-1">日均笔记</div>
          </CardContent>
        </Card>
      </div>

      <!-- 热力图 -->
      <Card>
        <CardContent class="p-3">
          <Heatmap />
        </CardContent>
      </Card>
    </div>
  {/if}

  <!-- 标签管理 -->
  {#if activeTab === 'tags'}
    <TagManager on:updated={loadData} />
  {/if}

  <!-- 隐私设置 -->
  {#if activeTab === 'privacy'}
    <PrivacySettings onLogout={handleLogout} />
  {/if}

  <!-- 退出按钮 -->
  <div class="mt-6 flex justify-end">
    <Button variant="destructive" on:click={handleLogout}>
      退出登录
    </Button>
  </div>
</div>

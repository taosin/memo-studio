<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { goto } from '$app/navigation';

  let loading = false;
  let error = '';
  let locationStats = [];
  let notesByLocation = {};
  let selectedLocation = '';
  let notes = [];

  onMount(async () => {
    await loadStats();
  });

  async function loadStats() {
    loading = true;
    error = '';
    try {
      const res = await api.call('/api/locations/stats', { method: 'GET' });
      locationStats = res.locations || [];
    } catch (e) {
      error = '加载统计失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  async function selectLocation(loc) {
    selectedLocation = loc;
    loading = true;
    notes = [];
    try {
      const res = await api.call('/api/notes?location=' + encodeURIComponent(loc), { method: 'GET' });
      notes = res.notes || [];
    } catch (e) {
      error = '加载笔记失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  async function detectAndSave(noteId) {
    loading = true;
    try {
      const res = await api.call('/api/memos/' + noteId + '/detect-and-save', { method: 'POST' });
      if (res.success) {
        await loadStats();
        if (selectedLocation) {
          await selectLocation(selectedLocation);
        }
      }
    } catch (e) {
      error = '检测失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  async function updateLocation(noteId, loc, lat, lng) {
    loading = true;
    try {
      await api.call('/api/memos/' + noteId + '/location', {
        method: 'PUT',
        body: { location: loc, latitude: lat, longitude: lng }
      });
      await loadStats();
      if (selectedLocation) {
        await selectLocation(selectedLocation);
      }
    } catch (e) {
      error = '更新失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  function goBack() {
    goto('/');
  }

  // 扩展 api 对象以支持自定义端点
  async function customFetch(path, options = {}) {
    const token = localStorage.getItem('token') || '';
    const res = await fetch(`http://localhost:9000${path}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(options.headers || {}),
      },
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `请求失败 (${res.status})`);
    }

    return res.json();
  }

  api.call = customFetch;
</script>

<svelte:head>
  <title>位置 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">📍 位置洞察</div>
      <button class="link" on:click={goBack}>返回首页</button>
    </div>

    <div class="description">
      查看按地点分类的笔记，了解你的足迹分布。
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <div class="grid">
      <!-- 位置统计 -->
      <div class="section">
        <h3>🏙️ 地点分布</h3>
        {#if loading && locationStats.length === 0}
          <div class="loading">加载中...</div>
        {:else if locationStats.length === 0}
          <div class="empty">
            <p>暂无位置数据</p>
            <p class="muted">创建笔记时提及地点（如"在北京"），系统会自动识别</p>
          </div>
        {:else}
          <div class="location-list">
            {#each locationStats as stat}
              <button
                class="location-item"
                class:selected={selectedLocation === stat.location}
                on:click={() => selectLocation(stat.location)}
              >
                <span class="icon">📍</span>
                <span class="name">{stat.location}</span>
                <span class="count">{stat.count} 篇</span>
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <!-- 笔记列表 -->
      <div class="section">
        <h3>
          {#if selectedLocation}
            📝 {selectedLocation} 的笔记
          {:else}
            📝 选择一个地点查看
          {/if}
        </h3>

        {#if loading && notes.length === 0}
          <div class="loading">加载中...</div>
        {:else if !selectedLocation}
          <div class="empty">
            <p>点击左侧地点查看相关笔记</p>
          </div>
        {:else if notes.length === 0}
          <div class="empty">
            <p>该地点暂无笔记</p>
          </div>
        {:else}
          <div class="notes-list">
            {#each notes as note}
              <div class="note-card">
                <div class="note-content">
                  {note.content?.substring(0, 150) || '(无内容)'}
                  {#if note.content?.length > 150}
                    ...
                  {/if}
                </div>
                <div class="note-meta">
                  <span class="date">{new Date(note.created_at).toLocaleDateString('zh-CN')}</span>
                  <span class="location-tag">📍 {note.location || '未知'}</span>
                </div>
                <div class="note-actions">
                  {#if !note.location}
                    <button class="mini-btn" on:click={() => detectAndSave(note.id)}>
                      🔍 识别位置
                    </button>
                  {/if}
                  <button class="mini-btn" on:click={() => goto('/note/' + note.id)}>
                    查看详情
                  </button>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <div class="tips">
      <h4>💡 使用提示</h4>
      <ul>
        <li>📝 在笔记中提及地点（如"今天在北京"）</li>
        <li>🔍 系统会自动识别并标记位置</li>
        <li>📊 在此页面查看足迹分布</li>
      </ul>
    </div>
  </div>
</div>

<style>
  .wrap {
    display: flex;
    justify-content: center;
    padding: 24px 16px;
  }
  .card {
    width: 100%;
    max-width: 1000px;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 24px;
  }
  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }
  .title {
    font-size: 24px;
    font-weight: 800;
  }
  .link {
    font-size: 14px;
    color: var(--accent);
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }
  .description {
    color: var(--muted);
    margin-bottom: 20px;
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 12px;
    margin-bottom: 16px;
    color: #f87171;
  }
  .grid {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 24px;
  }
  @media (max-width: 700px) {
    .grid {
      grid-template-columns: 1fr;
    }
  }
  .section {
    background: var(--panel-2);
    border-radius: 12px;
    padding: 16px;
  }
  .section h3 {
    margin: 0 0 16px 0;
    font-size: 16px;
    font-weight: 600;
  }
  .loading {
    text-align: center;
    padding: 20px;
    color: var(--muted);
  }
  .empty {
    text-align: center;
    padding: 24px;
    color: var(--muted);
  }
  .muted {
    font-size: 13px;
    opacity: 0.7;
  }
  .location-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .location-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--panel);
    cursor: pointer;
    text-align: left;
    transition: all 0.2s;
  }
  .location-item:hover {
    border-color: var(--accent);
  }
  .location-item.selected {
    border-color: var(--accent);
    background: var(--accent-soft);
  }
  .location-item .icon {
    font-size: 18px;
  }
  .location-item .name {
    flex: 1;
    font-weight: 500;
  }
  .location-item .count {
    font-size: 13px;
    color: var(--muted);
  }
  .notes-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .note-card {
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 14px;
  }
  .note-content {
    font-size: 14px;
    line-height: 1.5;
    margin-bottom: 10px;
    white-space: pre-wrap;
    word-break: break-word;
  }
  .note-meta {
    display: flex;
    gap: 12px;
    font-size: 12px;
    color: var(--muted);
    margin-bottom: 10px;
  }
  .location-tag {
    background: rgba(34, 197, 94, 0.1);
    padding: 2px 8px;
    border-radius: 999px;
    color: #22c55e;
  }
  .note-actions {
    display: flex;
    gap: 8px;
  }
  .mini-btn {
    padding: 6px 12px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--panel-2);
    color: inherit;
    font-size: 12px;
    cursor: pointer;
  }
  .mini-btn:hover {
    border-color: var(--accent);
  }
  .tips {
    margin-top: 24px;
    padding-top: 20px;
    border-top: 1px solid var(--border);
    color: var(--muted);
  }
  .tips h4 {
    margin: 0 0 12px 0;
    font-size: 14px;
  }
  .tips ul {
    margin: 0;
    padding-left: 20px;
  }
  .tips li {
    margin: 8px 0;
    font-size: 13px;
  }
</style>

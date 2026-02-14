<script>
  import { onMount } from 'svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/card-title.svelte';

  export let onLogout;

  let privacySettings = {
    enableEncryption: true,
    autoLock: false,
    autoLockTimeout: 5,
    clearClipboard: true,
    clearClipboardTimeout: 30,
    privateBrowsing: false
  };

  let securityLogs = [];
  let loading = false;

  onMount(() => {
    loadSettings();
    loadSecurityLogs();
  });

  function loadSettings() {
    try {
      const saved = localStorage.getItem('privacy_settings');
      if (saved) {
        privacySettings = { ...privacySettings, ...JSON.parse(saved) };
      }
    } catch (e) {
      console.error('加载设置失败:', e);
    }
  }

  function saveSettings() {
    localStorage.setItem('privacy_settings', JSON.stringify(privacySettings));
    addLog('设置已更新');
  }

  function addLog(message) {
    const log = {
      id: Date.now(),
      time: new Date().toISOString(),
      message
    };
    securityLogs = [log, ...securityLogs.slice(0, 49)];
    localStorage.setItem('security_logs', JSON.stringify(securityLogs));
  }

  function loadSecurityLogs() {
    try {
      const saved = localStorage.getItem('security_logs');
      if (saved) {
        securityLogs = JSON.parse(saved);
      }
    } catch {
      securityLogs = [];
    }
  }

  function formatTime(isoString) {
    const date = new Date(isoString);
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function exportData() {
    // 触发数据导出
    window.dispatchEvent(new CustomEvent('export-data'));
    addLog('触发数据导出');
  }

  function clearAllData() {
    if (confirm('确定要清除所有本地数据吗？此操作不可恢复！')) {
      if (confirm('再次确认：真的要清除所有数据吗？')) {
        localStorage.clear();
        addLog('清除所有本地数据');
        alert('已清除所有数据，请重新登录');
        if (onLogout) onLogout();
      }
    }
  }

  function regenerateEncryptionKey() {
    if (confirm('重新生成密钥将导致已加密数据无法解密。确定继续吗？')) {
      localStorage.removeItem('memo_encryption_key');
      addLog('重新生成加密密钥');
      alert('加密密钥已重新生成。请注意：这可能导致旧数据无法解密。');
    }
  }
</script>

<div class="max-w-2xl mx-auto space-y-6">
  <!-- 标题 -->
  <div class="text-center mb-8">
    <h2 class="text-2xl font-bold bg-gradient-to-r from-primary to-primary-light bg-clip-text text-transparent">
      🔒 隐私与安全
    </h2>
    <p class="text-muted-foreground mt-2">
      管理您的数据安全和隐私设置
    </p>
  </div>

  <!-- 加密设置 -->
  <Card>
    <CardHeader>
      <CardTitle class="flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-primary">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
        <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
      </svg>
        加密设置
      </CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <p class="font-medium">本地数据加密</p>
          <p class="text-sm text-muted-foreground">使用 AES-256 加密本地存储的笔记</p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input 
            type="checkbox" 
            bind:checked={privacySettings.enableEncryption}
            on:change={saveSettings}
            class="sr-only peer"
          />
          <div class="w-11 h-6 bg-muted-foreground/20 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
        </label>
      </div>

      {#if privacySettings.enableEncryption}
        <div class="flex items-center justify-between pt-4 border-t">
          <div>
            <p class="font-medium">重新生成加密密钥</p>
            <p class="text-sm text-muted-foreground">可能导致旧数据无法解密</p>
          </div>
          <Button variant="outline" size="sm" on:click={regenerateEncryptionKey}>
            重新生成
          </Button>
        </div>
      {/if}
    </CardContent>
  </Card>

  <!-- 自动锁定 -->
  <Card>
    <CardHeader>
      <CardTitle class="flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-orange-500">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
        自动锁定
      </CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <p class="font-medium">离开时自动锁定</p>
          <p class="text-sm text-muted-foreground">页面闲置一段时间后清除敏感数据</p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input 
            type="checkbox" 
            bind:checked={privacySettings.autoLock}
            on:change={saveSettings}
            class="sr-only peer"
          />
          <div class="w-11 h-6 bg-muted-foreground/20 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-orange-500"></div>
        </label>
      </div>

      {#if privacySettings.autoLock}
        <div class="flex items-center justify-between">
          <p class="text-sm">锁定超时</p>
          <select 
            bind:value={privacySettings.autoLockTimeout}
            on:change={saveSettings}
            class="px-3 py-1.5 rounded-lg border border-border bg-background"
          >
            <option value={1}>1 分钟</option>
            <option value={5}>5 分钟</option>
            <option value={15}>15 分钟</option>
            <option value={30}>30 分钟</option>
          </select>
        </div>
      {/if}
    </CardContent>
  </Card>

  <!-- 剪贴板清理 -->
  <Card>
    <CardHeader>
      <CardTitle class="flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-purple-500">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
        </svg>
        剪贴板清理
      </CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <p class="font-medium">自动清除剪贴板</p>
          <p class="text-sm text-muted-foreground">复制敏感内容后自动清除</p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input 
            type="checkbox" 
            bind:checked={privacySettings.clearClipboard}
            on:change={saveSettings}
            class="sr-only peer"
          />
          <div class="w-11 h-6 bg-muted-foreground/20 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-500"></div>
        </label>
      </div>

      {#if privacySettings.clearClipboard}
        <div class="flex items-center justify-between">
          <p class="text-sm">清除前等待</p>
          <select 
            bind:value={privacySettings.clearClipboardTimeout}
            on:change={saveSettings}
            class="px-3 py-1.5 rounded-lg border border-border bg-background"
          >
            <option value={10}>10 秒</option>
            <option value={30}>30 秒</option>
            <option value={60}>1 分钟</option>
          </select>
        </div>
      {/if}
    </CardContent>
  </Card>

  <!-- 数据管理 -->
  <Card class="border-destructive/20">
    <CardHeader>
      <CardTitle class="flex items-center gap-2 text-destructive">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 6h18"/>
          <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
          <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
        </svg>
        危险区域
      </CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <p class="font-medium">导出所有数据</p>
          <p class="text-sm text-muted-foreground">下载 JSON 格式的备份</p>
        </div>
        <Button variant="outline" on:click={exportData}>
          导出
        </Button>
      </div>

      <div class="flex items-center justify-between pt-4 border-t">
        <div>
          <p class="font-medium text-destructive">清除所有本地数据</p>
          <p class="text-sm text-muted-foreground">此操作不可恢复</p>
        </div>
        <Button variant="destructive" on:click={clearAllData}>
          清除
        </Button>
      </div>
    </CardContent>
  </Card>

  <!-- 安全日志 -->
  <Card>
    <CardHeader>
      <CardTitle class="flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-blue-500">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
          <line x1="16" y1="13" x2="8" y2="13"/>
          <line x1="16" y1="17" x2="8" y2="17"/>
          <polyline points="10 9 9 9 8 9"/>
        </svg>
        操作日志
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="space-y-2 max-h-60 overflow-y-auto">
        {#each securityLogs.slice(0, 20) as log}
          <div class="flex items-center justify-between text-sm py-2 border-b border-border/50 last:border-0">
            <span>{log.message}</span>
            <span class="text-muted-foreground">{formatTime(log.time)}</span>
          </div>
        {:else}
          <p class="text-center text-muted-foreground py-4">暂无日志</p>
        {/each}
      </div>
    </CardContent>
  </Card>
</div>

<style>
  select {
    cursor: pointer;
  }
</style>

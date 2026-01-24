<script>
  import { onMount } from 'svelte';
  import { authStore } from '../stores/auth.js';
  import { api } from '../utils/api.js';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/ui/card/card-title.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Input from '$lib/components/ui/input/input.svelte';
  import Label from '$lib/components/ui/label/label.svelte';
  import Button from '$lib/components/ui/button/button.svelte';

  let username = '';
  let password = '';
  let email = '';
  let loading = false;
  let error = '';
  let isRegister = false;
  let showPassword = false;

  async function handleLogin() {
    if (!username.trim() || !password.trim()) {
      error = '请输入用户名和密码';
      return;
    }

    loading = true;
    error = '';

    try {
      const response = await api.login(username, password);
      authStore.login(response.token, response.user);
      window.dispatchEvent(new CustomEvent('auth-success'));
    } catch (err) {
      error = err.message || '登录失败';
    } finally {
      loading = false;
    }
  }

  async function handleRegister() {
    if (!username.trim() || !password.trim()) {
      error = '请输入用户名和密码';
      return;
    }

    if (username.length < 3) {
      error = '用户名长度至少为3位';
      return;
    }

    if (password.length < 6) {
      error = '密码长度至少为6位';
      return;
    }

    loading = true;
    error = '';

    try {
      const response = await api.register(username, password, email);
      authStore.login(response.token, response.user);
      window.dispatchEvent(new CustomEvent('auth-success'));
    } catch (err) {
      error = err.message || '注册失败';
    } finally {
      loading = false;
    }
  }

  function toggleMode() {
    isRegister = !isRegister;
    error = '';
    password = '';
    email = '';
    showPassword = false;
  }

  function handleSubmit() {
    if (isRegister) {
      handleRegister();
    } else {
      handleLogin();
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-background via-background to-muted/20 px-4 py-12 relative overflow-hidden">
  <!-- 装饰性背景元素 -->
  <div class="absolute inset-0 overflow-hidden pointer-events-none">
    <div class="absolute -top-40 -right-40 w-80 h-80 bg-primary/5 rounded-full blur-3xl"></div>
    <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-primary/5 rounded-full blur-3xl"></div>
  </div>

  <div class="w-full max-w-md relative z-10">
    <Card class="border-2 shadow-xl backdrop-blur-sm bg-card/95">
      <CardHeader class="space-y-6 pb-6">
        <!-- Logo 和标题 -->
        <div class="flex flex-col items-center space-y-3">
          <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-primary/60 flex items-center justify-center shadow-lg transform transition-transform hover:scale-105">
            <span class="text-3xl">📝</span>
          </div>
          <div class="text-center space-y-1">
            <CardTitle class="text-3xl font-bold tracking-tight">
              {isRegister ? '创建账号' : '欢迎回来'}
            </CardTitle>
            <p class="text-sm text-muted-foreground">
              {isRegister ? '开始使用 Memo Studio' : '登录以继续使用 Memo Studio'}
            </p>
          </div>
        </div>
      </CardHeader>

      <CardContent class="space-y-5 px-6 pb-6">
        <!-- 错误提示 -->
        {#if error}
          <div class="p-4 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm flex items-start gap-2 animate-in fade-in slide-in-from-top-2">
            <svg class="w-5 h-5 mt-0.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{error}</span>
          </div>
        {/if}

        <form on:submit|preventDefault={handleSubmit} class="space-y-5">
          <!-- 用户名输入 -->
          <div class="space-y-2">
            <Label for="username" class="text-sm font-medium">
              用户名
            </Label>
            <div class="relative">
              <div class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <Input
                id="username"
                placeholder="请输入用户名"
                bind:value={username}
                disabled={loading}
                class="pl-10 h-11"
                autocomplete="username"
              />
            </div>
          </div>

          <!-- 密码输入 -->
          <div class="space-y-2">
            <Label for="password" class="text-sm font-medium">
              密码
            </Label>
            <div class="relative">
              <div class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              {#if showPassword}
                <input
                  id="password"
                  type="text"
                  placeholder="请输入密码"
                  bind:value={password}
                  disabled={loading}
                  autocomplete={isRegister ? 'new-password' : 'current-password'}
                  class="flex h-11 w-full rounded-md border border-input bg-background pl-10 pr-10 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  on:keydown={(e) => e.key === 'Enter' && handleSubmit()}
                />
              {:else}
                <input
                  id="password"
                  type="password"
                  placeholder="请输入密码"
                  bind:value={password}
                  disabled={loading}
                  autocomplete={isRegister ? 'new-password' : 'current-password'}
                  class="flex h-11 w-full rounded-md border border-input bg-background pl-10 pr-10 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  on:keydown={(e) => e.key === 'Enter' && handleSubmit()}
                />
              {/if}
              <button
                type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                on:click={() => showPassword = !showPassword}
                disabled={loading}
              >
                {#if showPassword}
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                {:else}
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                {/if}
              </button>
            </div>
          </div>

          <!-- 邮箱输入（注册时显示） -->
          {#if isRegister}
            <div class="space-y-2 animate-in fade-in slide-in-from-top-2">
              <Label for="email" class="text-sm font-medium">
                邮箱 <span class="text-muted-foreground font-normal">(可选)</span>
              </Label>
              <div class="relative">
                <div class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
                <input
                  id="email"
                  type="email"
                  placeholder="请输入邮箱"
                  bind:value={email}
                  disabled={loading}
                  autocomplete="email"
                  class="flex h-11 w-full rounded-md border border-input bg-background pl-10 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
              </div>
            </div>
          {/if}

          <!-- 提交按钮 -->
          <Button
            type="submit"
            class="w-full h-11 text-base font-semibold shadow-md hover:shadow-lg transition-all duration-200"
            disabled={loading}
          >
            {#if loading}
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {isRegister ? '注册中...' : '登录中...'}
            {:else}
              {isRegister ? '创建账号' : '登录'}
            {/if}
          </Button>
        </form>

        <!-- 切换登录/注册 -->
        <div class="pt-4 border-t">
          <div class="text-center text-sm">
            {#if isRegister}
              <span class="text-muted-foreground">已有账号？</span>
              <button
                type="button"
                class="ml-1.5 font-medium text-primary hover:text-primary/80 transition-colors underline-offset-4 hover:underline"
                on:click={toggleMode}
                disabled={loading}
              >
                立即登录
              </button>
            {:else}
              <span class="text-muted-foreground">还没有账号？</span>
              <button
                type="button"
                class="ml-1.5 font-medium text-primary hover:text-primary/80 transition-colors underline-offset-4 hover:underline"
                on:click={toggleMode}
                disabled={loading}
              >
                立即注册
              </button>
            {/if}
          </div>
        </div>

        <!-- 提示信息 -->
        {#if !isRegister}
          <div class="pt-2">
            <div class="flex items-center justify-center gap-2 text-xs text-muted-foreground">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>首次使用请先注册账号</span>
            </div>
          </div>
        {/if}
      </CardContent>
    </Card>
  </div>
</div>

<style>
  @keyframes fade-in {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes slide-in-from-top-2 {
    from {
      transform: translateY(-0.5rem);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .animate-in {
    animation: fade-in 0.2s ease-out;
  }

  .fade-in {
    animation: fade-in 0.2s ease-out;
  }

  .slide-in-from-top-2 {
    animation: slide-in-from-top-2 0.2s ease-out;
  }
</style>

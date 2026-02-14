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
  let usernameFocused = false;
  let passwordFocused = false;

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
      error = err.message || '登录失败，请检查用户名和密码';
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
      error = err.message || '注册失败，请稍后重试';
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

  // 回车键快速提交
  function handleKeydown(e) {
    if (e.key === 'Enter' && !loading) {
      handleSubmit();
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-background via-background to-primary/5 px-4 py-12 relative overflow-hidden">
  <!-- 动态背景装饰 -->
  <div class="absolute inset-0 overflow-hidden pointer-events-none">
    <div class="absolute -top-40 -right-40 w-96 h-96 bg-primary/5 rounded-full blur-3xl animate-pulse-soft"></div>
    <div class="absolute -bottom-40 -left-40 w-96 h-96 bg-primary/5 rounded-full blur-3xl animate-pulse-soft" style="animation-delay: 1s;"></div>
    <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-primary/3 rounded-full blur-3xl"></div>
  </div>

  <div class="w-full max-w-md relative z-10">
    <!-- Logo 和标题 -->
    <div class="text-center mb-8 animate-fade-in">
      <div class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-primary to-primary-light shadow-xl shadow-primary/25 mb-4 transform hover:scale-105 transition-transform duration-300">
        <span class="text-4xl">📝</span>
      </div>
      <h1 class="text-3xl font-bold bg-gradient-to-r from-primary to-primary-light bg-clip-text text-transparent">
        Memo Studio
      </h1>
      <p class="text-muted-foreground mt-2">记录美好生活</p>
    </div>

    <Card class="border-2 shadow-2xl backdrop-blur-sm bg-card/95 animate-slide-up">
      <CardHeader class="space-y-6 pb-6">
        <div class="text-center space-y-2">
          <CardTitle class="text-2xl font-bold">
            {isRegister ? '创建账号' : '欢迎回来'}
          </CardTitle>
          <p class="text-sm text-muted-foreground">
            {isRegister ? '立即开始记录你的灵感' : '登录以继续使用 Memo Studio'}
          </p>
        </div>
      </CardHeader>

      <CardContent class="px-6 pb-6">
        <!-- 错误提示 -->
        {#if error}
          <div class="mb-6 p-4 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm flex items-start gap-3 animate-fade-in">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mt-0.5 flex-shrink-0">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <span>{error}</span>
          </div>
        {/if}

        <!-- 表单 -->
        <form on:submit|preventDefault={handleSubmit} class="space-y-5">
          <!-- 用户名输入 -->
          <div class="space-y-2">
            <Label for="username" class="text-sm font-medium flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              用户名
            </Label>
            <div class="relative group">
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none z-10 transition-colors group-focus-within:text-primary">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <input
                id="username"
                type="text"
                placeholder="请输入用户名"
                bind:value={username}
                disabled={loading}
                autocomplete="username"
                on:keydown={handleKeydown}
                class="flex h-12 w-full rounded-lg border border-border bg-background pl-12 pr-4 text-sm placeholder:text-muted-foreground/50 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>
          </div>

          <!-- 密码输入 -->
          <div class="space-y-2">
            <Label for="password" class="text-sm font-medium flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              密码
            </Label>
            <div class="relative group">
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none z-10 transition-colors group-focus-within:text-primary">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input
                id="password"
                type={showPassword ? 'text' : 'password'}
                placeholder="请输入密码"
                bind:value={password}
                disabled={loading}
                autocomplete={isRegister ? 'new-password' : 'current-password'}
                on:keydown={handleKeydown}
                class="flex h-12 w-full rounded-lg border border-border bg-background pl-12 pr-12 text-sm placeholder:text-muted-foreground/50 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-50"
              />
              <button
                type="button"
                class="absolute right-4 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors p-1"
                on:click={() => showPassword = !showPassword}
                disabled={loading}
              >
                {#if showPassword}
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                {:else}
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                {/if}
              </button>
            </div>
          </div>

          <!-- 邮箱输入（注册时显示） -->
          {#if isRegister}
            <div class="space-y-2 animate-fade-in">
              <Label for="email" class="text-sm font-medium flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                  <polyline points="22,6 12,13 2,6"/>
                </svg>
                邮箱 <span class="text-muted-foreground font-normal">(可选)</span>
              </Label>
              <div class="relative group">
                <div class="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none z-10 transition-colors group-focus-within:text-primary">
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
                <input
                  id="email"
                  type="email"
                  placeholder="用于找回密码"
                  bind:value={email}
                  disabled={loading}
                  autocomplete="email"
                  class="flex h-12 w-full rounded-lg border border-border bg-background pl-12 pr-4 text-sm placeholder:text-muted-foreground/50 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-50"
                />
              </div>
            </div>
          {/if}

          <!-- 提交按钮 -->
          <Button
            type="submit"
            variant="gradient"
            class="w-full h-12 text-base font-semibold shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-all duration-300 mt-6"
            loading={loading}
          >
            {#if loading}
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              {isRegister ? '注册中...' : '登录中...'}
            {:else}
              {isRegister ? '创建账号' : '登录'}
            {/if}
          </Button>
        </form>

        <!-- 切换登录/注册 -->
        <div class="mt-6 pt-6 border-t border-border">
          <div class="text-center text-sm">
            <span class="text-muted-foreground">
              {isRegister ? '已有账号？' : '还没有账号？'}
            </span>
            <button
              type="button"
              class="ml-2 font-semibold text-primary hover:text-primary/80 transition-colors"
              on:click={toggleMode}
              disabled={loading}
            >
              {isRegister ? '立即登录' : '立即注册'}
            </button>
          </div>
        </div>

        <!-- 特性提示 -->
        <div class="mt-6 grid grid-cols-3 gap-4 text-center">
          <div class="flex flex-col items-center gap-1">
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-primary">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
            </div>
            <span class="text-xs text-muted-foreground">安全</span>
          </div>
          <div class="flex flex-col items-center gap-1">
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-primary">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
            </div>
            <span class="text-xs text-muted-foreground">简单</span>
          </div>
          <div class="flex flex-col items-center gap-1">
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-primary">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </div>
            <span class="text-xs text-muted-foreground">快速</span>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</div>

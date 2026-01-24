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
      // 触发登录成功事件，由父组件处理
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
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-background px-4 py-12">
  <Card class="w-full max-w-md">
    <CardHeader>
      <CardTitle class="text-2xl text-center">
        {isRegister ? '注册' : '登录'}
      </CardTitle>
    </CardHeader>
    <CardContent class="p-3 space-y-4">
      {#if error}
        <div class="p-3 rounded-md bg-destructive/10 text-destructive text-sm">
          {error}
        </div>
      {/if}

      <div>
        <Label>用户名</Label>
        <Input
          placeholder="请输入用户名"
          bind:value={username}
          disabled={loading}
        />
      </div>

      <div>
        <Label>密码</Label>
        <Input
          type="password"
          placeholder="请输入密码"
          bind:value={password}
          disabled={loading}
          on:keydown={(e) => e.key === 'Enter' && (isRegister ? handleRegister() : handleLogin())}
        />
      </div>

      {#if isRegister}
        <div>
          <Label>邮箱（可选）</Label>
          <Input
            type="email"
            placeholder="请输入邮箱"
            bind:value={email}
            disabled={loading}
          />
        </div>
      {/if}

      <Button
        class="w-full"
        on:click={isRegister ? handleRegister : handleLogin}
        disabled={loading}
      >
        {loading ? '处理中...' : (isRegister ? '注册' : '登录')}
      </Button>

      <div class="text-center text-sm text-muted-foreground">
        {#if isRegister}
          <span>已有账号？</span>
          <button
            class="text-primary hover:underline"
            on:click={toggleMode}
            disabled={loading}
          >
            立即登录
          </button>
        {:else}
          <span>还没有账号？</span>
          <button
            class="text-primary hover:underline"
            on:click={toggleMode}
            disabled={loading}
          >
            立即注册
          </button>
        {/if}
      </div>

      {#if !isRegister}
        <div class="text-xs text-muted-foreground text-center pt-4 border-t">
          <p>首次使用请先注册账号</p>
        </div>
      {/if}
    </CardContent>
  </Card>
</div>

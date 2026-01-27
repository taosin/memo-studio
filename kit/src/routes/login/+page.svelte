<script>
  import { api } from '$lib/api.js';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let username = 'admin';
  let password = '';
  let loading = false;
  let error = '';
  const uid = `u_${Math.random().toString(16).slice(2)}`;
  const pid = `p_${Math.random().toString(16).slice(2)}`;

  onMount(() => {
    // 已登录则跳回首页
    try {
      const t = localStorage.getItem('token');
      if (t) goto('/');
    } catch {}
  });

  async function submit() {
    loading = true;
    error = '';
    try {
      const res = await api.login(username.trim(), password);
      try {
        localStorage.setItem('token', res.token);
        localStorage.setItem('user', JSON.stringify(res.user || {}));
      } catch {}
      await goto('/');
    } catch (e) {
      error = e?.message || '登录失败';
    } finally {
      loading = false;
    }
  }
</script>

<div class="wrap">
  <div class="card">
    <div class="title">登录</div>
    <div class="hint">默认管理员用户名：admin（密码请查看后端启动日志，或设置 MEMO_ADMIN_PASSWORD）</div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <label class="label" for={uid}>用户名</label>
    <input id={uid} class="input" bind:value={username} autocomplete="username" />

    <label class="label" for={pid}>密码</label>
    <input id={pid} class="input" type="password" bind:value={password} autocomplete="current-password" />

    <button class="btn" on:click={submit} disabled={loading}>
      {loading ? '登录中…' : '登录'}
    </button>
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
    max-width: 420px;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 16px;
  }
  .title {
    font-weight: 800;
    margin-bottom: 6px;
  }
  .hint {
    color: var(--muted);
    font-size: 12px;
    margin-bottom: 12px;
  }
  .label {
    display: block;
    margin-top: 10px;
    margin-bottom: 6px;
    color: var(--muted);
    font-size: 12px;
  }
  .input {
    width: 100%;
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 10px 12px;
    box-sizing: border-box;
    outline: none;
  }
  .btn {
    width: 100%;
    margin-top: 14px;
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 10px 12px;
    cursor: pointer;
    font-weight: 700;
  }
  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 10px 12px;
    margin: 8px 0 10px;
  }
</style>


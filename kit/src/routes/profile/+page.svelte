<script>
  import { api } from '$lib/api.js';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let me = null;
  let username = '';
  let email = '';
  let old_password = '';
  let new_password = '';
  let loading = false;
  let error = '';
  let toast = '';
  const idUser = `id_${Math.random().toString(16).slice(2)}`;
  const idEmail = `em_${Math.random().toString(16).slice(2)}`;
  const idOld = `op_${Math.random().toString(16).slice(2)}`;
  const idNew = `np_${Math.random().toString(16).slice(2)}`;

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ''), 1500);
  }

  async function load() {
    loading = true;
    error = '';
    try {
      me = await api.me();
      username = me?.username || '';
      email = me?.email || '';
    } catch (e) {
      error = e?.message || '加载失败';
      if (String(error).includes('401') || String(error).includes('未认证')) {
        await goto('/login');
      }
    } finally {
      loading = false;
    }
  }

  async function saveProfile() {
    loading = true;
    error = '';
    try {
      me = await api.updateMe({ username, email });
      try {
        localStorage.setItem('user', JSON.stringify(me || {}));
      } catch {}
      setToast('已保存');
    } catch (e) {
      error = e?.message || '保存失败';
    } finally {
      loading = false;
    }
  }

  async function changePassword() {
    loading = true;
    error = '';
    try {
      await api.changePassword({ old_password, new_password });
      old_password = '';
      new_password = '';
      setToast('密码已修改');
    } catch (e) {
      error = e?.message || '修改失败';
    } finally {
      loading = false;
    }
  }

  async function logout() {
    await api.logout();
    await goto('/login');
  }

  onMount(load);
</script>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">个人信息</div>
      <button class="mini" on:click={logout}>登出</button>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}
    {#if toast}
      <div class="toast">{toast}</div>
    {/if}

    {#if loading && !me}
      <div class="muted">加载中…</div>
    {:else}
      <div class="section">
        <div class="sectionTitle">基础信息</div>
        <label class="label" for={idUser}>用户名</label>
        <input id={idUser} class="input" bind:value={username} />
        <label class="label" for={idEmail}>邮箱（可选）</label>
        <input id={idEmail} class="input" bind:value={email} />
        <button class="btn" on:click={saveProfile} disabled={loading}>保存</button>
      </div>

      <div class="section">
        <div class="sectionTitle">修改密码</div>
        <label class="label" for={idOld}>旧密码</label>
        <input id={idOld} class="input" type="password" bind:value={old_password} />
        <label class="label" for={idNew}>新密码</label>
        <input id={idNew} class="input" type="password" bind:value={new_password} />
        <button class="btn" on:click={changePassword} disabled={loading}>修改密码</button>
      </div>

      {#if me?.is_admin}
        <div class="section">
          <div class="sectionTitle">管理员</div>
          <a class="link" href="/admin/users">用户管理</a>
        </div>
      {/if}
    {/if}
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
    max-width: 700px;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 16px;
  }
  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    margin-bottom: 10px;
  }
  .title {
    font-weight: 800;
  }
  .mini {
    border-radius: 10px;
    border: 1px solid var(--border);
    background: var(--panel);
    color: inherit;
    padding: 6px 10px;
    cursor: pointer;
    font-size: 12px;
  }
  .section {
    border-top: 1px solid var(--border);
    padding-top: 12px;
    margin-top: 12px;
  }
  .sectionTitle {
    font-weight: 700;
    margin-bottom: 10px;
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
    margin-top: 12px;
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
  .toast {
    border: 1px solid rgba(34, 197, 94, 0.35);
    background: var(--accent-soft);
    border-radius: 12px;
    padding: 10px 12px;
    margin: 8px 0 10px;
  }
  .muted {
    color: var(--muted);
    padding: 10px 0;
  }
  .link {
    display: inline-block;
    text-decoration: underline;
    color: inherit;
  }
</style>


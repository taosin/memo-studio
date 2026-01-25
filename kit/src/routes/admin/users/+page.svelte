<script>
  import { api } from '$lib/api.js';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let users = [];
  let loading = false;
  let error = '';
  let toast = '';

  let newUsername = '';
  let newPassword = '';
  let newEmail = '';
  let newIsAdmin = false;

  function setToast(s) {
    toast = s;
    setTimeout(() => (toast = ''), 1500);
  }

  async function load() {
    loading = true;
    error = '';
    try {
      users = await api.adminListUsers();
    } catch (e) {
      error = e?.message || '加载失败';
      if (String(error).includes('未认证')) {
        await goto('/login');
      }
    } finally {
      loading = false;
    }
  }

  async function createUser() {
    loading = true;
    error = '';
    try {
      await api.adminCreateUser({
        username: newUsername,
        password: newPassword,
        email: newEmail,
        is_admin: newIsAdmin
      });
      newUsername = '';
      newPassword = '';
      newEmail = '';
      newIsAdmin = false;
      setToast('已创建');
      await load();
    } catch (e) {
      error = e?.message || '创建失败';
    } finally {
      loading = false;
    }
  }

  async function saveUser(u) {
    loading = true;
    error = '';
    try {
      await api.adminUpdateUser(u.id, {
        username: u.username,
        email: u.email || '',
        is_admin: !!u.is_admin
      });
      setToast('已保存');
      await load();
    } catch (e) {
      error = e?.message || '保存失败';
    } finally {
      loading = false;
    }
  }

  async function delUser(id) {
    if (!confirm('确定删除该用户吗？')) return;
    loading = true;
    error = '';
    try {
      await api.adminDeleteUser(id);
      setToast('已删除');
      await load();
    } catch (e) {
      error = e?.message || '删除失败';
    } finally {
      loading = false;
    }
  }

  onMount(load);
</script>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">用户管理</div>
      <a class="link" href="/profile">返回个人信息</a>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}
    {#if toast}
      <div class="toast">{toast}</div>
    {/if}

    <div class="section">
      <div class="sectionTitle">新增用户</div>
      <div class="grid">
        <input class="input" placeholder="用户名" bind:value={newUsername} />
        <input class="input" type="password" placeholder="密码" bind:value={newPassword} />
        <input class="input" placeholder="邮箱(可选)" bind:value={newEmail} />
        <label class="chk">
          <input type="checkbox" bind:checked={newIsAdmin} />
          管理员
        </label>
        <button class="btn" on:click={createUser} disabled={loading}>创建</button>
      </div>
    </div>

    <div class="section">
      <div class="sectionTitle">用户列表</div>
      {#if loading && users.length === 0}
        <div class="muted">加载中…</div>
      {:else}
        <div class="table">
          <div class="tr head">
            <div>ID</div>
            <div>用户名</div>
            <div>邮箱</div>
            <div>管理员</div>
            <div>操作</div>
          </div>
          {#each users as u (u.id)}
            <div class="tr">
              <div class="id">{u.id}</div>
              <div><input class="cell" bind:value={u.username} /></div>
              <div><input class="cell" bind:value={u.email} /></div>
              <div>
                <input type="checkbox" bind:checked={u.is_admin} />
              </div>
              <div class="ops">
                <button class="mini" on:click={() => saveUser(u)} disabled={loading}>保存</button>
                <button class="mini danger" on:click={() => delUser(u.id)} disabled={loading}>删除</button>
              </div>
            </div>
          {/each}
        </div>
      {/if}
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
    max-width: 980px;
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
  .link {
    text-decoration: underline;
    color: inherit;
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
  .grid {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr auto auto;
    gap: 10px;
    align-items: center;
  }
  @media (max-width: 900px) {
    .grid {
      grid-template-columns: 1fr;
    }
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
  .chk {
    display: flex;
    gap: 8px;
    align-items: center;
    color: var(--muted);
    font-size: 12px;
  }
  .btn {
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
  .table {
    display: grid;
    gap: 8px;
  }
  .tr {
    display: grid;
    grid-template-columns: 70px 1fr 1fr 90px 160px;
    gap: 10px;
    align-items: center;
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 10px 12px;
    background: var(--panel-2);
  }
  .tr.head {
    background: transparent;
    border: none;
    padding: 0 2px;
    color: var(--muted);
    font-size: 12px;
  }
  .id {
    color: var(--muted);
    font-size: 12px;
  }
  .cell {
    width: 100%;
    border-radius: 10px;
    border: 1px solid var(--border-2);
    background: rgba(15, 23, 42, 0.06);
    color: inherit;
    padding: 8px 10px;
    box-sizing: border-box;
    outline: none;
  }
  .ops {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
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
  .mini.danger {
    border-color: rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
  }
  .muted {
    color: var(--muted);
    padding: 10px 0;
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
</style>


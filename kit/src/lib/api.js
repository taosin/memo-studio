const API_BASE = '/api';

function getToken() {
  try {
    return localStorage.getItem('token') || '';
  } catch {
    return '';
  }
}

function requireToken() {
  const t = getToken();
  if (!t) throw new Error('未登录');
  return t;
}

async function jsonFetch(path, options = {}) {
  const token = getToken();
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options.headers || {})
    },
    ...options
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `请求失败(${res.status})`);
  }
  return res.json();
}

export const api = {
  async login(username, password) {
    return jsonFetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    });
  },
  async logout() {
    // 无需请求后端：清 token 即可
    try {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    } catch {}
    return true;
  },
  async me() {
    return jsonFetch('/users/me', { method: 'GET' });
  },
  async updateMe({ username, email }) {
    return jsonFetch('/users/me', {
      method: 'PUT',
      body: JSON.stringify({ username, email })
    });
  },
  async changePassword({ old_password, new_password }) {
    return jsonFetch('/users/me/password', {
      method: 'PUT',
      body: JSON.stringify({ old_password, new_password })
    });
  },
  async adminListUsers() {
    return jsonFetch('/users', { method: 'GET' });
  },
  async adminCreateUser(payload) {
    return jsonFetch('/users', { method: 'POST', body: JSON.stringify(payload) });
  },
  async adminUpdateUser(id, payload) {
    return jsonFetch(`/users/${id}`, { method: 'PUT', body: JSON.stringify(payload) });
  },
  async adminDeleteUser(id) {
    return jsonFetch(`/users/${id}`, { method: 'DELETE' });
  },
  async listNotes() {
    // 迁移到 memos（需要登录）
    requireToken();
    return jsonFetch('/memos');
  },
  async createNote({ content, tags = [] }) {
    requireToken();
    return jsonFetch('/memos', {
      method: 'POST',
      body: JSON.stringify({ title: '', content, tags, pinned: false, content_type: 'markdown', resource_ids: [] })
    });
  },
  async updateNote(id, { content, tags = [] }) {
    requireToken();
    return jsonFetch(`/memos/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ title: '', content, tags, pinned: false, content_type: 'markdown', resource_ids: [] })
    });
  },
  async deleteNote(id) {
    requireToken();
    return jsonFetch(`/memos/${id}`, { method: 'DELETE' });
  },
  async search(q) {
    requireToken();
    const qp = new URLSearchParams({ q: q || '', limit: '50', offset: '0' });
    return jsonFetch(`/memos?${qp.toString()}`, { method: 'GET' });
  },
  async listTags(withCount = true) {
    // 标签仍走公共接口（如需按用户隔离，可再加 tag 统计接口）
    const qp = new URLSearchParams(withCount ? { withCount: '1' } : {});
    return jsonFetch(`/tags?${qp.toString()}`, { method: 'GET' });
  },
  async randomReview({ limit = 1, tag = '', days = 0 } = {}) {
    requireToken();
    const qp = new URLSearchParams({
      limit: String(limit),
      tag: tag || '',
      days: String(days || 0)
    });
    return jsonFetch(`/review/random?${qp.toString()}`, { method: 'GET' });
  }
};

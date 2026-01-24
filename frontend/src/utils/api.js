// 使用真实 API，连接后端 SQLite 数据库
const USE_MOCK = false;

const API_BASE = '/api';

// 获取认证 token
function getAuthToken() {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token');
  }
  return null;
}

// 获取带认证头的 fetch
function fetchWithAuth(url, options = {}) {
  const token = getAuthToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  return fetch(url, {
    ...options,
    headers,
  });
}

const realApi = {
  // 认证相关
  async login(username, password) {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || '登录失败');
    }
    return await response.json();
  },

  async register(username, password, email = '') {
    const response = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password, email }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || '注册失败');
    }
    return await response.json();
  },

  async getCurrentUser() {
    const response = await fetchWithAuth(`${API_BASE}/auth/me`);
    if (!response.ok) {
      throw new Error('获取用户信息失败');
    }
    return await response.json();
  },

  // 笔记相关
  async getNotes() {
    const response = await fetchWithAuth(`${API_BASE}/notes`);
    if (!response.ok) {
      throw new Error('获取笔记列表失败');
    }
    return await response.json();
  },

  async getNote(id) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`);
    if (!response.ok) {
      throw new Error('获取笔记失败');
    }
    return await response.json();
  },

  async createNote(title, content, tags) {
    const response = await fetchWithAuth(`${API_BASE}/notes`, {
      method: 'POST',
      body: JSON.stringify({
        title,
        content,
        tags,
      }),
    });
    if (!response.ok) {
      throw new Error('创建笔记失败');
    }
    return await response.json();
  },

  async getTags() {
    const response = await fetchWithAuth(`${API_BASE}/tags`);
    if (!response.ok) {
      throw new Error('获取标签列表失败');
    }
    return await response.json();
  },

  async updateNote(id, title, content, tags) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`, {
      method: 'PUT',
      body: JSON.stringify({
        title,
        content,
        tags,
      }),
    });
    if (!response.ok) {
      throw new Error('更新笔记失败');
    }
    return await response.json();
  },

  async deleteNote(id) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error('删除笔记失败');
    }
    return await response.json();
  },

  async deleteNotes(ids) {
    const response = await fetchWithAuth(`${API_BASE}/notes/batch`, {
      method: 'DELETE',
      body: JSON.stringify({ ids }),
    });
    if (!response.ok) {
      throw new Error('批量删除笔记失败');
    }
    return await response.json();
  },

  async updateTag(id, name, color) {
    const response = await fetchWithAuth(`${API_BASE}/tags/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name, color }),
    });
    if (!response.ok) {
      throw new Error('更新标签失败');
    }
    return await response.json();
  },

  async deleteTag(id) {
    const response = await fetchWithAuth(`${API_BASE}/tags/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error('删除标签失败');
    }
    return await response.json();
  },

  async mergeTags(sourceId, targetId) {
    const response = await fetchWithAuth(`${API_BASE}/tags/merge`, {
      method: 'POST',
      body: JSON.stringify({ sourceId, targetId }),
    });
    if (!response.ok) {
      throw new Error('合并标签失败');
    }
    return await response.json();
  },
};

// 根据 USE_MOCK 选择使用 mock 数据还是真实 API
export const api = USE_MOCK ? (() => {
  // 如果需要使用 mock，取消下面的注释
  // return mockApi;
  throw new Error('Mock API 未启用，请设置 USE_MOCK = true 并导入 mockApi');
})() : realApi;

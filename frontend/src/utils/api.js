// API 配置
const USE_MOCK = false;
const API_BASE = '/api/v1'; // 使用新的 API 版本

// ===== 认证拦截器 =====
let authInterceptors = [];

export function addAuthInterceptor(fn) {
  authInterceptors.push(fn);
}

export function removeAuthInterceptor(fn) {
  authInterceptors = authInterceptors.filter(f => f !== fn);
}

function handleAuthError() {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    // 触发认证错误事件，让应用知道需要重新登录
    window.dispatchEvent(new CustomEvent('auth-expired'));
  }
}

// ===== 获取 token =====
function getAuthToken() {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token');
  }
  return null;
}

// ===== 统一错误处理 =====
function handleApiError(response, customMessage) {
  if (response.status === 401) {
    handleAuthError();
    throw new Error('登录已过期，请重新登录');
  }
  if (response.status === 404) {
    throw new Error(customMessage || '资源不存在');
  }
  if (response.status === 429) {
    throw new Error('请求过于频繁，请稍后再试');
  }
  if (response.status >= 400) {
    const error = response.json ? response.json().catch(() => ({})) : Promise.resolve({});
    throw new Error(customMessage || error.error || `请求失败 (${response.status})`);
  }
  return null;
}

// ===== 带认证的 fetch =====
async function fetchWithAuth(url, options = {}) {
  const token = getAuthToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  // 调用拦截器
  for (const interceptor of authInterceptors) {
    const result = interceptor({ url, options, headers });
    if (result === false) {
      throw new Error('请求被拦截器取消');
    }
  }

  return fetch(url, {
    ...options,
    headers,
  });
}

// ===== Content 清理工具 =====
function cleanContent(content) {
  if (typeof content === 'string') {
    if (content === '[object Object]' || content === '[object object]') {
      return '';
    }
    // 检查是否是 JSON 字符串化的对象
    if (content.startsWith('{') || content.startsWith('[')) {
      try {
        const parsed = JSON.parse(content);
        if (typeof parsed === 'object' && parsed !== null) {
          return parsed.content || parsed.text || parsed.value || JSON.stringify(parsed);
        }
      } catch (e) {
        // 不是有效的 JSON，保持原样
      }
    }
    return content;
  }
  if (content !== null && content !== undefined) {
    if (typeof content === 'object') {
      return content.content || content.text || content.value || JSON.stringify(content);
    }
    return String(content);
  }
  return '';
}

function cleanNote(note) {
  return {
    ...note,
    content: cleanContent(note.content),
    title: typeof note.title === 'string' ? note.title : (note.title ? String(note.title) : ''),
  };
}

// ===== API 实现 =====
const realApi = {
  // 认证相关
  async login(username, password) {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: '登录失败' }));
      throw new Error(error.error || '登录失败');
    }
    return await response.json();
  },

  async register(username, password, email = '') {
    const response = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password, email }),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: '注册失败' }));
      throw new Error(error.error || '注册失败');
    }
    return await response.json();
  },

  async getCurrentUser() {
    const response = await fetchWithAuth(`${API_BASE}/auth/me`);
    if (!response.ok) {
      handleAuthError();
      throw new Error('登录已过期，请重新登录');
    }
    return await response.json();
  },

  // 笔记相关
  async getNotes() {
    const response = await fetchWithAuth(`${API_BASE}/notes`);
    if (!response.ok) {
      handleAuthError();
      throw new Error('获取笔记列表失败');
    }
    const data = await response.json();
    return Array.isArray(data) ? data.map(cleanNote) : [];
  },

  async getNote(id) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`);
    if (!response.ok) {
      handleAuthError();
      if (response.status === 404) throw new Error('笔记不存在');
      throw new Error('获取笔记失败');
    }
    const note = await response.json();
    return cleanNote(note);
  },

  async createNote(title, content, tags) {
    const response = await fetchWithAuth(`${API_BASE}/notes`, {
      method: 'POST',
      body: JSON.stringify({ title, content, tags }),
    });
    if (!response.ok) {
      handleAuthError();
      if (response.status === 400) {
        const error = await response.json().catch(() => ({ error: '' }));
        throw new Error(error.error || '标题和内容不能同时为空');
      }
      throw new Error('创建笔记失败');
    }
    return await response.json();
  },

  async updateNote(id, title, content, tags) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ title, content, tags }),
    });
    if (!response.ok) {
      handleAuthError();
      if (response.status === 404) throw new Error('笔记不存在');
      throw new Error('更新笔记失败');
    }
    return await response.json();
  },

  async deleteNote(id) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      handleAuthError();
      throw new Error('删除笔记失败');
    }
    return await response.json();
  },

  async deleteNotes(ids) {
    if (!Array.isArray(ids) || ids.length === 0) {
      throw new Error('请选择要删除的笔记');
    }
    const response = await fetchWithAuth(`${API_BASE}/notes/batch`, {
      method: 'DELETE',
      body: JSON.stringify({ ids }),
    });
    if (!response.ok) {
      handleAuthError();
      throw new Error('批量删除笔记失败');
    }
    return await response.json();
  },

  // 标签相关
  async getTags() {
    const response = await fetchWithAuth(`${API_BASE}/tags`);
    if (!response.ok) {
      handleAuthError();
      throw new Error('获取标签列表失败');
    }
    const data = await response.json();
    return Array.isArray(data) ? data : [];
  },

  async createTag(name, color) {
    const response = await fetchWithAuth(`${API_BASE}/tags`, {
      method: 'POST',
      body: JSON.stringify({ name, color }),
    });
    if (!response.ok) {
      handleAuthError();
      const error = await response.json().catch(() => ({ error: '创建标签失败' }));
      throw new Error(error.error || '创建标签失败');
    }
    return await response.json();
  },

  async updateTag(id, name, color) {
    if (!name || name.trim() === '') {
      throw new Error('标签名称不能为空');
    }
    const response = await fetchWithAuth(`${API_BASE}/tags/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name, color }),
    });
    if (!response.ok) {
      handleAuthError();
      if (response.status === 404) throw new Error('标签不存在');
      throw new Error('更新标签失败');
    }
    return await response.json();
  },

  async deleteTag(id) {
    const response = await fetchWithAuth(`${API_BASE}/tags/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      handleAuthError();
      throw new Error('删除标签失败');
    }
    return await response.json();
  },

  async mergeTags(sourceId, targetId) {
    if (!sourceId || !targetId) {
      throw new Error('请选择要合并的标签');
    }
    if (sourceId === targetId) {
      throw new Error('不能合并相同的标签');
    }
    const response = await fetchWithAuth(`${API_BASE}/tags/merge`, {
      method: 'POST',
      body: JSON.stringify({ sourceId, targetId }),
    });
    if (!response.ok) {
      handleAuthError();
      throw new Error('合并标签失败');
    }
    return await response.json();
  },

  // 搜索
  async searchNotes(query) {
    const response = await fetchWithAuth(`${API_BASE}/search?q=${encodeURIComponent(query)}`);
    if (!response.ok) {
      handleAuthError();
      throw new Error('搜索笔记失败');
    }
    const data = await response.json();
    return Array.isArray(data) ? data.map(cleanNote) : [];
  },
};

// ===== 导出 =====
export const api = USE_MOCK ? (() => {
  throw new Error('Mock API 未启用');
})() : realApi;

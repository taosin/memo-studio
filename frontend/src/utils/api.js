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
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '获取用户信息失败');
    }
    return await response.json();
  },

  // 笔记相关
  async getNotes() {
    const response = await fetchWithAuth(`${API_BASE}/notes`);
    if (!response.ok) {
      // 如果是 401 未授权，可能是 token 过期
      if (response.status === 401) {
        // 清除 token，让用户重新登录
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      throw new Error('获取笔记列表失败');
    }
    const data = await response.json();
    // 确保返回的是数组，并清理每个笔记的 content 字段
    if (!Array.isArray(data)) {
      return [];
    }
    // 清理每个笔记的数据
    return data.map(note => {
      let cleanContent = '';
      if (typeof note.content === 'string') {
        cleanContent = note.content;
        // 检查是否是 JSON 字符串化的对象
        if (cleanContent.startsWith('{') || cleanContent.startsWith('[')) {
          try {
            const parsed = JSON.parse(cleanContent);
            // 如果是对象，尝试提取文本内容
            if (typeof parsed === 'object') {
              cleanContent = parsed.content || parsed.text || parsed.value || JSON.stringify(parsed);
            }
          } catch (e) {
            // 不是有效的 JSON，保持原样
          }
        }
      } else if (note.content !== null && note.content !== undefined) {
        // 如果是对象，尝试提取或转换
        if (typeof note.content === 'object') {
          cleanContent = note.content.content || note.content.text || note.content.value || JSON.stringify(note.content);
        } else {
          cleanContent = String(note.content);
        }
      }
      
      return {
        ...note,
        content: cleanContent,
        title: typeof note.title === 'string' ? note.title : (note.title ? String(note.title) : '')
      };
    });
  },

  async getNote(id) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`);
    if (!response.ok) {
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      if (response.status === 404) {
        throw new Error('笔记不存在');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '获取笔记失败');
    }
    const note = await response.json();
    // 清理 content 字段
    let cleanContent = '';
    if (typeof note.content === 'string') {
      cleanContent = note.content;
      // 检查是否是 JSON 字符串化的对象
      if (cleanContent.startsWith('{') || cleanContent.startsWith('[')) {
        try {
          const parsed = JSON.parse(cleanContent);
          // 如果是对象，尝试提取文本内容
          if (typeof parsed === 'object') {
            cleanContent = parsed.content || parsed.text || parsed.value || JSON.stringify(parsed);
          }
        } catch (e) {
          // 不是有效的 JSON，保持原样
        }
      }
    } else if (note.content !== null && note.content !== undefined) {
      // 如果是对象，尝试提取或转换
      if (typeof note.content === 'object') {
        cleanContent = note.content.content || note.content.text || note.content.value || JSON.stringify(note.content);
      } else {
        cleanContent = String(note.content);
      }
    }
    
    // 确保 content 和 title 是字符串
    return {
      ...note,
      content: cleanContent,
      title: typeof note.title === 'string' ? note.title : (note.title ? String(note.title) : '')
    };
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
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '创建笔记失败');
    }
    return await response.json();
  },

  async getTags() {
    const response = await fetchWithAuth(`${API_BASE}/tags`);
    if (!response.ok) {
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '获取标签列表失败');
    }
    const data = await response.json();
    // 确保返回的是数组
    return Array.isArray(data) ? data : [];
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
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      if (response.status === 404) {
        throw new Error('笔记不存在');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '更新笔记失败');
    }
    return await response.json();
  },

  async deleteNote(id) {
    const response = await fetchWithAuth(`${API_BASE}/notes/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      if (response.status === 404) {
        throw new Error('笔记不存在');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '删除笔记失败');
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
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '批量删除笔记失败');
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
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      if (response.status === 404) {
        throw new Error('标签不存在');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '更新标签失败');
    }
    return await response.json();
  },

  async deleteTag(id) {
    const response = await fetchWithAuth(`${API_BASE}/tags/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      if (response.status === 404) {
        throw new Error('标签不存在');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '删除标签失败');
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
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        throw new Error('登录已过期，请重新登录');
      }
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || '合并标签失败');
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

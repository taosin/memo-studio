// 使用真实 API，连接后端 SQLite 数据库
const USE_MOCK = false;

const API_BASE = '/api';

const realApi = {
  async getNotes() {
    const response = await fetch(`${API_BASE}/notes`);
    if (!response.ok) {
      throw new Error('获取笔记列表失败');
    }
    return await response.json();
  },

  async getNote(id) {
    const response = await fetch(`${API_BASE}/notes/${id}`);
    if (!response.ok) {
      throw new Error('获取笔记失败');
    }
    return await response.json();
  },

  async createNote(title, content, tags) {
    const response = await fetch(`${API_BASE}/notes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
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
    const response = await fetch(`${API_BASE}/tags`);
    if (!response.ok) {
      throw new Error('获取标签列表失败');
    }
    return await response.json();
  },

  async updateNote(id, title, content, tags) {
    const response = await fetch(`${API_BASE}/notes/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
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
    const response = await fetch(`${API_BASE}/notes/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error('删除笔记失败');
    }
    return await response.json();
  },

  async deleteNotes(ids) {
    const response = await fetch(`${API_BASE}/notes/batch`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ ids }),
    });
    if (!response.ok) {
      throw new Error('批量删除笔记失败');
    }
    return await response.json();
  },

  async updateTag(id, name, color) {
    const response = await fetch(`${API_BASE}/tags/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, color }),
    });
    if (!response.ok) {
      throw new Error('更新标签失败');
    }
    return await response.json();
  },

  async deleteTag(id) {
    const response = await fetch(`${API_BASE}/tags/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error('删除标签失败');
    }
    return await response.json();
  },

  async mergeTags(sourceId, targetId) {
    const response = await fetch(`${API_BASE}/tags/merge`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
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

// 使用 mock 数据，不依赖后端
import { mockApi } from './mockData.js';

// 如果需要切换到真实 API，将 USE_MOCK 设置为 false
const USE_MOCK = true;

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
};

// 根据 USE_MOCK 选择使用 mock 数据还是真实 API
export const api = USE_MOCK ? mockApi : realApi;

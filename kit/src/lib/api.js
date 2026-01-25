const API_BASE = '/api';

async function jsonFetch(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
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
  async listNotes() {
    return jsonFetch('/notes');
  },
  async createNote({ content, tags = [] }) {
    return jsonFetch('/notes', {
      method: 'POST',
      body: JSON.stringify({ title: '', content, tags })
    });
  },
  async updateNote(id, { content, tags = [] }) {
    return jsonFetch(`/notes/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ title: '', content, tags })
    });
  },
  async deleteNote(id) {
    return jsonFetch(`/notes/${id}`, { method: 'DELETE' });
  },
  async search(q) {
    const qp = new URLSearchParams({ q: q || '', limit: '50', offset: '0' });
    return jsonFetch(`/search?${qp.toString()}`, { method: 'GET' });
  },
  async listTags(withCount = true) {
    const qp = new URLSearchParams(withCount ? { withCount: '1' } : {});
    return jsonFetch(`/tags?${qp.toString()}`, { method: 'GET' });
  },
  async randomReview({ limit = 1, tag = '', days = 0 } = {}) {
    const qp = new URLSearchParams({
      limit: String(limit),
      tag: tag || '',
      days: String(days || 0)
    });
    return jsonFetch(`/review/random?${qp.toString()}`, { method: 'GET' });
  }
};


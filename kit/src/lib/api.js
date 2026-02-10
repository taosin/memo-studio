const API_BASE = "/api";

function getToken() {
  try {
    return localStorage.getItem("token") || "";
  } catch {
    return "";
  }
}

function requireToken() {
  const t = getToken();
  if (!t) throw new Error("未登录");
  return t;
}

async function jsonFetch(path, options = {}) {
  // always attempt to read token via getToken (works in Node tests where window is undefined but global.localStorage is mocked)
  const token = getToken();
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options.headers || {}),
    },
    ...options,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `请求失败(${res.status})`);
  }
  return res.json();
}

export const api = {
  async login(username, password) {
    return jsonFetch("/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    });
  },
  async logout() {
    // 无需请求后端：清 token 即可
    try {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    } catch {}
    return true;
  },
  async me() {
    return jsonFetch("/users/me", { method: "GET" });
  },
  async updateMe({ username, email }) {
    return jsonFetch("/users/me", {
      method: "PUT",
      body: JSON.stringify({ username, email }),
    });
  },
  async changePassword({ old_password, new_password }) {
    return jsonFetch("/users/me/password", {
      method: "PUT",
      body: JSON.stringify({ old_password, new_password }),
    });
  },
  async adminListUsers() {
    return jsonFetch("/users", { method: "GET" });
  },
  async adminCreateUser(payload) {
    return jsonFetch("/users", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
  async adminUpdateUser(id, payload) {
    return jsonFetch(`/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    });
  },
  async adminDeleteUser(id) {
    return jsonFetch(`/users/${id}`, { method: "DELETE" });
  },
  async listNotes() {
    // 迁移到 memos（需要登录）
    requireToken();
    return jsonFetch("/memos");
  },
  async createNote({ content, tags = [] }) {
    requireToken();
    return jsonFetch("/memos", {
      method: "POST",
      body: JSON.stringify({
        title: "",
        content,
        tags,
        pinned: false,
        content_type: "markdown",
        resource_ids: [],
      }),
    });
  },
  async updateNote(id, { content, tags = [] }) {
    requireToken();
    return jsonFetch(`/memos/${id}`, {
      method: "PUT",
      body: JSON.stringify({
        title: "",
        content,
        tags,
        pinned: false,
        content_type: "markdown",
        resource_ids: [],
      }),
    });
  },
  async deleteNote(id) {
    requireToken();
    return jsonFetch(`/memos/${id}`, { method: "DELETE" });
  },
  async search(q) {
    requireToken();
    const qp = new URLSearchParams({ q: q || "", limit: "50", offset: "0" });
    return jsonFetch(`/memos?${qp.toString()}`, { method: "GET" });
  },
  async listTags(withCount = true) {
    const qp = new URLSearchParams(withCount ? { withCount: "1" } : {});
    return jsonFetch(`/tags?${qp.toString()}`, { method: "GET" });
  },
  async createTag({ name, color }) {
    requireToken();
    return jsonFetch("/tags", {
      method: "POST",
      body: JSON.stringify({ name: name || "", color: color || "" }),
    });
  },
  async updateTag(id, { name, color }) {
    requireToken();
    return jsonFetch(`/tags/${id}`, {
      method: "PUT",
      body: JSON.stringify({ name: name || "", color: color || "" }),
    });
  },
  async deleteTag(id) {
    requireToken();
    return jsonFetch(`/tags/${id}`, { method: "DELETE" });
  },
  async mergeTags(sourceId, targetId) {
    requireToken();
    return jsonFetch("/tags/merge", {
      method: "POST",
      body: JSON.stringify({ sourceId, targetId }),
    });
  },
  async listResources(limit = 20, offset = 0) {
    requireToken();
    const qp = new URLSearchParams({
      limit: String(limit),
      offset: String(offset),
    });
    return jsonFetch(`/resources?${qp.toString()}`, { method: "GET" });
  },
  async deleteResource(id) {
    requireToken();
    return jsonFetch(`/resources/${id}`, { method: "DELETE" });
  },
  async uploadResource(file) {
    requireToken();
    const form = new FormData();
    form.append("file", file);
    const token = getToken();
    const res = await fetch(`${API_BASE}/resources`, {
      method: "POST",
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: form,
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `上传失败(${res.status})`);
    }
    return res.json();
  },
  async randomReview({ limit = 1, tag = "", days = 0 } = {}) {
    requireToken();
    const qp = new URLSearchParams({
      limit: String(limit),
      tag: tag || "",
      days: String(days || 0),
    });
    return jsonFetch(`/review/random?${qp.toString()}`, { method: "GET" });
  },
  // 笔记本
  async listNotebooks() {
    requireToken();
    return jsonFetch("/notebooks", { method: "GET" });
  },
  async getNotebook(id) {
    requireToken();
    return jsonFetch(`/notebooks/${id}`, { method: "GET" });
  },
  async createNotebook({ name, color, sort_order = 0 }) {
    requireToken();
    return jsonFetch("/notebooks", {
      method: "POST",
      body: JSON.stringify({
        name: name || "",
        color: color || "",
        sort_order,
      }),
    });
  },
  async updateNotebook(id, { name, color, sort_order }) {
    requireToken();
    return jsonFetch(`/notebooks/${id}`, {
      method: "PUT",
      body: JSON.stringify({ name, color, sort_order }),
    });
  },
  async deleteNotebook(id) {
    requireToken();
    return jsonFetch(`/notebooks/${id}`, { method: "DELETE" });
  },
  async listNotebookNotes(id, limit = 50, offset = 0) {
    requireToken();
    const qp = new URLSearchParams({
      limit: String(limit),
      offset: String(offset),
    });
    return jsonFetch(`/notebooks/${id}/notes?${qp.toString()}`, {
      method: "GET",
    });
  },
  // 统计
  async getStats() {
    requireToken();
    return jsonFetch("/stats", { method: "GET" });
  },
  // 导出：format=json 返回 JSON 对象；format=markdown 返回 { blob, filename }
  async exportNotes(format = "json", limit = 500) {
    requireToken();
    const qp = new URLSearchParams({
      format: format === "markdown" ? "markdown" : "json",
      limit: String(limit),
    });
    const token = getToken();
    const res = await fetch(`${API_BASE}/export?${qp.toString()}`, {
      method: "GET",
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `导出失败(${res.status})`);
    }
    if (format === "markdown") {
      const blob = await res.blob();
      const disposition = res.headers.get("Content-Disposition") || "";
      const m = disposition.match(/filename=(.+)/);
      const filename = m ? m[1].replace(/^["']|["']$/g, "") : "memo-export.md";
      return { blob, filename };
    }
    return res.json();
  },
  // 导入
  async importNotes(notes) {
    requireToken();
    return jsonFetch("/import", {
      method: "POST",
      body: JSON.stringify({ notes: notes || [] }),
    });
  },
};

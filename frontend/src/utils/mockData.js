// Mock 数据
let mockNotes = [
  {
    id: 1,
    title: '欢迎使用 Memo Studio',
    content: '这是一个简洁美观的笔记应用，支持标签管理和明暗主题切换。你可以在这里记录你的想法、学习笔记、工作计划等。\n\n开始创建你的第一个笔记吧！',
    tags: [
      { id: 1, name: '欢迎', color: '#4ECDC4' },
      { id: 2, name: '指南', color: '#FF6B6B' }
    ],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  },
  {
    id: 2,
    title: 'Svelte 学习笔记',
    content: 'Svelte 是一个构建用户界面的框架。与 React 和 Vue 不同，Svelte 在构建时将应用编译为高度优化的 JavaScript 代码。\n\n主要特点：\n- 无虚拟 DOM\n- 更小的打包体积\n- 更好的性能\n- 简洁的语法',
    tags: [
      { id: 3, name: '学习', color: '#45B7D1' },
      { id: 4, name: '前端', color: '#98D8C8' }
    ],
    created_at: new Date(Date.now() - 86400000).toISOString(),
    updated_at: new Date(Date.now() - 86400000).toISOString()
  },
  {
    id: 3,
    title: '今日待办',
    content: '1. 完成项目重构\n2. 测试新功能\n3. 更新文档\n4. 代码审查',
    tags: [
      { id: 5, name: '工作', color: '#F7DC6F' },
      { id: 6, name: '待办', color: '#BB8FCE' }
    ],
    created_at: new Date(Date.now() - 172800000).toISOString(),
    updated_at: new Date(Date.now() - 172800000).toISOString()
  }
];

let mockTags = [
  { id: 1, name: '欢迎', color: '#4ECDC4', created_at: new Date().toISOString(), count: 1 },
  { id: 2, name: '指南', color: '#FF6B6B', created_at: new Date().toISOString(), count: 1 },
  { id: 3, name: '学习', color: '#45B7D1', created_at: new Date().toISOString(), count: 1 },
  { id: 4, name: '前端', color: '#98D8C8', created_at: new Date().toISOString(), count: 1 },
  { id: 5, name: '工作', color: '#F7DC6F', created_at: new Date().toISOString(), count: 1 },
  { id: 6, name: '待办', color: '#BB8FCE', created_at: new Date().toISOString(), count: 1 },
  { id: 7, name: '生活', color: '#85C1E2', created_at: new Date().toISOString(), count: 0 }
];

let nextNoteId = 4;
let nextTagId = 8;

// 模拟网络延迟
const delay = (ms = 300) => new Promise(resolve => setTimeout(resolve, ms));

export const mockApi = {
  async getNotes() {
    await delay();
    return [...mockNotes].sort((a, b) => 
      new Date(b.created_at) - new Date(a.created_at)
    );
  },

  async getNote(id) {
    await delay();
    const note = mockNotes.find(n => n.id === parseInt(id));
    if (!note) {
      throw new Error('笔记不存在');
    }
    return { ...note };
  },

  async createNote(title, content, tags = []) {
    await delay();
    
    // 创建或获取标签
    const noteTags = tags.map(tagName => {
      let tag = mockTags.find(t => t.name === tagName);
      if (!tag) {
        // 创建新标签
        const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A', '#98D8C8', '#F7DC6F', '#BB8FCE', '#85C1E2'];
        const hash = tagName.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
        tag = {
          id: nextTagId++,
          name: tagName,
          color: colors[hash % colors.length],
          created_at: new Date().toISOString()
        };
        mockTags.push(tag);
      }
      return tag;
    });

    const newNote = {
      id: nextNoteId++,
      title: title || '无标题',
      content,
      tags: noteTags,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };

    mockNotes.unshift(newNote); // 添加到开头
    updateTagCounts();
    return { ...newNote };
  },

  async updateNote(id, title, content, tags = []) {
    await delay();
    
    const noteIndex = mockNotes.findIndex(n => n.id === parseInt(id));
    if (noteIndex === -1) {
      throw new Error('笔记不存在');
    }

    // 创建或获取标签
    const noteTags = tags.map(tagName => {
      let tag = mockTags.find(t => t.name === tagName);
      if (!tag) {
        const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A', '#98D8C8', '#F7DC6F', '#BB8FCE', '#85C1E2'];
        const hash = tagName.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
        tag = {
          id: nextTagId++,
          name: tagName,
          color: colors[hash % colors.length],
          created_at: new Date().toISOString()
        };
        mockTags.push(tag);
      }
      return tag;
    });

    const note = mockNotes[noteIndex];
    note.title = title || '无标题';
    note.content = content;
    note.tags = noteTags;
    note.updated_at = new Date().toISOString();

    updateTagCounts();
    return { ...note };
  },

  async deleteNote(id) {
    await delay();
    
    const noteIndex = mockNotes.findIndex(n => n.id === parseInt(id));
    if (noteIndex === -1) {
      throw new Error('笔记不存在');
    }

    mockNotes.splice(noteIndex, 1);
    updateTagCounts();
    return { success: true };
  },

  async deleteNotes(ids) {
    await delay();
    
    ids.forEach(id => {
      const noteIndex = mockNotes.findIndex(n => n.id === parseInt(id));
      if (noteIndex !== -1) {
        mockNotes.splice(noteIndex, 1);
      }
    });

    updateTagCounts();
    return { success: true, deleted: ids.length };
  },

  async getTags() {
    await delay();
    updateTagCounts();
    return [...mockTags];
  },

  async updateTag(id, name, color) {
    await delay();
    
    const tag = mockTags.find(t => t.id === parseInt(id));
    if (!tag) {
      throw new Error('标签不存在');
    }

    tag.name = name;
    if (color) tag.color = color;
    tag.updated_at = new Date().toISOString();

    // 更新所有使用该标签的笔记
    mockNotes.forEach(note => {
      const noteTag = note.tags.find(t => t.id === tag.id);
      if (noteTag) {
        noteTag.name = name;
        if (color) noteTag.color = color;
      }
    });

    return { ...tag };
  },

  async deleteTag(id) {
    await delay();
    
    const tagIndex = mockTags.findIndex(t => t.id === parseInt(id));
    if (tagIndex === -1) {
      throw new Error('标签不存在');
    }

    // 从所有笔记中移除该标签
    mockNotes.forEach(note => {
      note.tags = note.tags.filter(t => t.id !== parseInt(id));
    });

    mockTags.splice(tagIndex, 1);
    return { success: true };
  },

  async mergeTags(sourceId, targetId) {
    await delay();
    
    const sourceTag = mockTags.find(t => t.id === parseInt(sourceId));
    const targetTag = mockTags.find(t => t.id === parseInt(targetId));
    
    if (!sourceTag || !targetTag) {
      throw new Error('标签不存在');
    }

    // 将所有使用源标签的笔记改为使用目标标签
    mockNotes.forEach(note => {
      const hasSource = note.tags.some(t => t.id === parseInt(sourceId));
      const hasTarget = note.tags.some(t => t.id === parseInt(targetId));
      
      if (hasSource && !hasTarget) {
        note.tags = note.tags.map(t => 
          t.id === parseInt(sourceId) ? targetTag : t
        );
      } else if (hasSource && hasTarget) {
        note.tags = note.tags.filter(t => t.id !== parseInt(sourceId));
      }
    });

    // 删除源标签
    await this.deleteTag(sourceId);
    updateTagCounts();
    
    return { success: true };
  }
};

// 更新标签计数
function updateTagCounts() {
  mockTags.forEach(tag => {
    tag.count = mockNotes.filter(note => 
      (note.tags || []).some(t => t.id === tag.id)
    ).length;
  });
}

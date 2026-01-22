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
  { id: 1, name: '欢迎', color: '#4ECDC4', created_at: new Date().toISOString() },
  { id: 2, name: '指南', color: '#FF6B6B', created_at: new Date().toISOString() },
  { id: 3, name: '学习', color: '#45B7D1', created_at: new Date().toISOString() },
  { id: 4, name: '前端', color: '#98D8C8', created_at: new Date().toISOString() },
  { id: 5, name: '工作', color: '#F7DC6F', created_at: new Date().toISOString() },
  { id: 6, name: '待办', color: '#BB8FCE', created_at: new Date().toISOString() },
  { id: 7, name: '生活', color: '#85C1E2', created_at: new Date().toISOString() }
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
    return { ...newNote };
  },

  async getTags() {
    await delay();
    return [...mockTags];
  }
};

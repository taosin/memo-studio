// ===== 标签智能推荐 =====

import { api } from './api.js';

// 标签使用统计
let tagUsageStats = {};

// 加载标签使用统计
export async function loadTagStats() {
  try {
    const notes = await api.getNotes();
    const stats = {};
    
    notes.forEach(note => {
      (note.tags || []).forEach(tag => {
        const tagId = tag.id || tag;
        if (!stats[tagId]) {
          stats[tagId] = { count: 0, recentNotes: [] };
        }
        stats[tagId].count++;
        stats[tagId].recentNotes.push(note.id);
      });
    });
    
    tagUsageStats = stats;
    return stats;
  } catch (e) {
    console.error('加载标签统计失败:', e);
    return {};
  }
}

// 获取最常使用的标签
export function getFrequentlyUsedTags(tags, limit = 5) {
  const sorted = Object.entries(tagUsageStats)
    .map(([id, stats]) => ({ id, ...stats }))
    .sort((a, b) => b.count - a.count)
    .slice(0, limit);
  
  return sorted.filter(s => tags.find(t => t.id === s.id));
}

// 基于内容推荐标签
export function recommendTags(content, availableTags) {
  const recommendations = [];
  const contentLower = content.toLowerCase();
  
  // 关键词到标签的映射
  const keywordMap = {
    '工作': ['工作', '任务', '项目'],
    '学习': ['学习', '笔记', '知识'],
    '灵感': ['灵感', '想法', '创意'],
    '代码': ['代码', '开发', '编程'],
    '读书': ['阅读', '书籍', '读书'],
    '旅行': ['旅行', '出行', '旅游'],
    '健康': ['健康', '运动', '健身'],
    '财务': ['财务', '金钱', '消费'],
    '会议': ['会议', '讨论', '沟通'],
    '计划': ['计划', '目标', '待办'],
    '回顾': ['回顾', '总结', '复盘'],
    '重要': ['重要', '紧急', '优先级'],
    '待办': ['待办', 'todo', '任务']
  };
  
  // 分析内容并推荐标签
  for (const [tagName, keywords] of Object.entries(keywordMap)) {
    const matchingTag = availableTags.find(t => 
      t.name.toLowerCase().includes(tagName.toLowerCase())
    );
    
    if (matchingTag) {
      const matchCount = keywords.filter(kw => 
        contentLower.includes(kw.toLowerCase())
      ).length;
      
      if (matchCount > 0) {
        recommendations.push({
          tag: matchingTag,
          score: matchCount * 10
        });
      }
    }
  }
  
  // 基于已有标签推荐相关标签
  const existingTagIds = [];
  const contentTags = content.match(/#[\u4e00-\u9fa5a-zA-Z0-9_]+/g) || [];
  contentTags.forEach(tag => {
    const tagName = tag.slice(1);
    const matched = availableTags.find(t => 
      t.name.toLowerCase() === tagName.toLowerCase()
    );
    if (matched) existingTagIds.push(matched.id);
  });
  
  // 推荐与已有标签相关的标签
  existingTagIds.forEach(tagId => {
    Object.entries(tagUsageStats)
      .filter(([id]) => id !== tagId)
      .forEach(([id, stats]) => {
        if (stats.recentNotes.some(noteId => 
          stats.recentNotes.includes(noteId) && 
          existingTagIds.some(eid => 
            tagUsageStats[eid]?.recentNotes.includes(noteId)
          )
        )) {
          const tag = availableTags.find(t => t.id === id);
          if (tag) {
            recommendations.push({ tag, score: 5 });
          }
        }
      });
  });
  
  // 去重并排序
  const seen = new Set();
  return recommendations
    .filter(r => {
      if (seen.has(r.tag.id)) return false;
      seen.add(r.tag.id);
      return true;
    })
    .sort((a, b) => b.score - a.score)
    .slice(0, 5);
}

// ===== 智能模板 =====

const defaultTemplates = [
  {
    id: 'daily-note',
    name: '📝 每日笔记',
    icon: '📝',
    content: `# 今日笔记 - {{date}}

## 今日完成
- [ ] 

## 待办事项
- [ ] 

## 想法与灵感


## 明日计划

`,
    tags: ['日记', '每日']
  },
  {
    id: 'meeting-note',
    name: '📅 会议记录',
    icon: '📅',
    content: `# {{title}}

## 基本信息
- **时间**: {{date}}
- **参与人员**: 
- **地点**: 

## 会议议程


## 讨论要点


## 决议事项


## 待办事项


## 下次会议

`,
    tags: ['会议', '工作']
  },
  {
    id: 'idea-note',
    name: '💡 灵感记录',
    icon: '💡',
    content: `# 灵感 - {{date}}

## 核心想法


## 应用场景


## 优缺点分析
**优点**:
- 

**缺点**:
- 

## 下一步行动


## 相关链接

`,
    tags: ['灵感', '创意']
  },
  {
    id: 'book-note',
    name: '📚 读书笔记',
    icon: '📚',
    content: `# {{bookName}}

## 书籍信息
- **作者**: 
- **阅读日期**: {{date}}
- **评分**: ⭐⭐⭐⭐⭐

## 核心观点


## 精彩摘录


## 个人感悟


## 推荐章节


`,
    tags: ['阅读', '学习']
  },
  {
    id: 'project-note',
    name: '📦 项目记录',
    icon: '📦',
    content: `# {{projectName}}

## 项目概述


## 目标


## 当前进度


## 问题与挑战


## 下一步计划


## 资源链接

`,
    tags: ['项目', '工作']
  },
  {
    id: 'review-note',
    name: '🔄 回顾总结',
    icon: '🔄',
    content: `# {{period}} 回顾

## 完成的事项


## 未完成的事项


## 收获与成长


## 待改进的地方


## 下阶段目标

`,
    tags: ['回顾', '总结']
  }
];

// 获取模板列表
export function getTemplates() {
  try {
    const saved = localStorage.getItem('note_templates');
    if (saved) {
      const custom = JSON.parse(saved);
      return [...defaultTemplates, ...custom];
    }
  } catch (e) {
    console.error('加载模板失败:', e);
  }
  return defaultTemplates;
}

// 保存自定义模板
export function saveTemplate(template) {
  const templates = getTemplates().filter(t => t.id.startsWith('custom-'));
  templates.push({ ...template, id: `custom-${Date.now()}` });
  localStorage.setItem('note_templates', JSON.stringify(templates));
  return templates[templates.length - 1];
}

// 删除自定义模板
export function deleteTemplate(templateId) {
  if (templateId.startsWith('custom-')) {
    const templates = getTemplates().filter(t => t.id !== templateId);
    localStorage.setItem('note_templates', JSON.stringify(
      templates.filter(t => t.id.startsWith('custom-'))
    ));
  }
}

// 使用模板创建笔记
export function useTemplate(templateId, extraData = {}) {
  const template = getTemplates().find(t => t.id === templateId);
  if (!template) return null;
  
  const now = new Date();
  const dateStr = now.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
  
  const replacements = {
    '{{date}}': dateStr,
    '{{title}}': extraData.title || '',
    '{{period}}': extraData.period || `${now.getFullYear()}年${now.getMonth() + 1}月`,
    '{{bookName}}': extraData.bookName || '',
    '{{projectName}}': extraData.projectName || ''
  };
  
  let content = template.content;
  Object.entries(replacements).forEach(([key, value]) => {
    content = content.replace(new RegExp(key, 'g'), value);
  });
  
  return {
    title: extraData.title || template.name,
    content,
    tags: template.tags
  };
}

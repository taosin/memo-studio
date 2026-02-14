// ===== 多格式数据导出 =====

import { api } from './api.js';

// 导出为 Markdown
export function exportAsMarkdown(notes, tags) {
  let content = '# Memo Studio 导出\n\n';
  content += `导出时间: ${new Date().toLocaleString('zh-CN')}\n`;
  content += `笔记总数: ${notes.length}\n\n`;
  content += '---\n\n';
  
  // 按日期分组
  const grouped = {};
  notes.forEach(note => {
    const date = new Date(note.created_at).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
    if (!grouped[date]) grouped[date] = [];
    grouped[date].push(note);
  });
  
  // 生成内容
  Object.entries(grouped).forEach(([date, dateNotes]) => {
    content += `## ${date}\n\n`;
    
    dateNotes.forEach(note => {
      content += `### ${note.title || '无标题'}\n\n`;
      
      // 标签
      if (note.tags && note.tags.length > 0) {
        content += note.tags.map(t => `#${t.name}`).join(' ') + '\n\n';
      }
      
      // 内容
      const plainContent = (note.content || '').replace(/<[^>]*>/g, '');
      content += plainContent + '\n\n';
      
      // 元信息
      content += `> 创建于: ${new Date(note.created_at).toLocaleString('zh-CN')}\n`;
      if (note.updated_at !== note.created_at) {
        content += `> 更新于: ${new Date(note.updated_at).toLocaleString('zh-CN')}\n`;
      }
      content += '\n---\n\n';
    });
  });
  
  return content;
}

// 导出为 HTML
export function exportAsHTML(notes, tags) {
  const styles = `
    <style>
      * { box-sizing: border-box; margin: 0; padding: 0; }
      body { 
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        max-width: 800px; margin: 0 auto; padding: 40px 20px;
        line-height: 1.6; color: #333;
      }
      h1 { font-size: 2em; margin-bottom: 10px; }
      .meta { color: #666; font-size: 0.9em; margin-bottom: 30px; }
      h2 { font-size: 1.5em; margin: 30px 0 15px; border-bottom: 2px solid #eee; padding-bottom: 10px; }
      h3 { font-size: 1.2em; margin: 20px 0 10px; }
      .note { background: #f9f9f9; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
      .tags { margin-bottom: 10px; }
      .tag { 
        display: inline-block; background: #e8f4fd; color: #1a73e8;
        padding: 2px 8px; border-radius: 12px; font-size: 0.85em; margin-right: 5px;
      }
      .timestamp { color: #999; font-size: 0.85em; margin-top: 10px; }
      hr { border: none; border-top: 1px solid #eee; margin: 30px 0; }
    </style>
  `;
  
  let html = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Memo Studio 导出</title>
  ${styles}
</head>
<body>
  <h1>📝 Memo Studio 导出</h1>
  <div class="meta">
    <p>导出时间: ${new Date().toLocaleString('zh-CN')}</p>
    <p>笔记总数: ${notes.length}</p>
  </div>
`;
  
  const grouped = {};
  notes.forEach(note => {
    const date = new Date(note.created_at).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
    if (!grouped[date]) grouped[date] = [];
    grouped[date].push(note);
  });
  
  Object.entries(grouped).forEach(([date, dateNotes]) => {
    html += `<h2>📅 ${date}</h2>`;
    
    dateNotes.forEach(note => {
      const plainContent = (note.content || '').replace(/<[^>]*>/g, '');
      
      html += `<div class="note">
        <h3>${note.title || '无标题'}</h3>
        ${note.tags && note.tags.length > 0 ? 
          `<div class="tags">${note.tags.map(t => `<span class="tag">#${t.name}</span>`).join('')}</div>` : ''}
        <div class="content">${plainContent.replace(/\n/g, '<br>')}</div>
        <div class="timestamp">
          创建: ${new Date(note.created_at).toLocaleString('zh-CN')}
        </div>
      </div>`;
    });
  });
  
  html += '</body></html>';
  return html;
}

// 导出为纯文本
export function exportAsPlainText(notes, tags) {
  let content = `Memo Studio 导出
${'='.repeat(50)}
导出时间: ${new Date().toLocaleString('zh-CN')}
笔记总数: ${notes.length}

`;
  
  notes.forEach((note, index) => {
    content += `${'─'.repeat(50)}\n`;
    content += `[${index + 1}] ${note.title || '无标题'}\n`;
    
    if (note.tags && note.tags.length > 0) {
      content += `标签: ${note.tags.map(t => '#' + t.name).join(' ')}\n`;
    }
    
    content += `\n${(note.content || '').replace(/<[^>]*>/g, '')}\n`;
    content += `\n创建: ${new Date(note.created_at).toLocaleString('zh-CN')}\n`;
  });
  
  return content;
}

// 导出为 CSV（适合数据分析）
export function exportAsCSV(notes, tags) {
  const headers = ['ID', '标题', '内容', '标签', '创建时间', '更新时间'];
  const rows = notes.map(note => [
    note.id,
    `"${(note.title || '').replace(/"/g, '""')}"`,
    `"${(note.content || '').replace(/<[^>]*>/g, '').replace(/"/g, '""')}"`,
    `"${(note.tags || []).map(t => t.name).join('; ')}"`,
    note.created_at,
    note.updated_at
  ]);
  
  return [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
}

// ===== 下载工具 =====

function downloadFile(content, filename, mimeType) {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

// 主导出函数
export async function exportData(format = 'markdown') {
  try {
    const [notes, tags] = await Promise.all([
      api.getNotes(),
      api.getTags()
    ]);
    
    const dateStr = new Date().toISOString().split('T')[0];
    const filename = `memo-studio-${dateStr}`;
    
    switch (format) {
      case 'markdown':
        downloadFile(
          exportAsMarkdown(notes, tags),
          `${filename}.md`,
          'text/markdown'
        );
        break;
        
      case 'html':
        downloadFile(
          exportAsHTML(notes, tags),
          `${filename}.html`,
          'text/html'
        );
        break;
        
      case 'plain':
        downloadFile(
          exportAsPlainText(notes, tags),
          `${filename}.txt`,
          'text/plain'
        );
        break;
        
      case 'csv':
        downloadFile(
          exportAsCSV(notes, tags),
          `${filename}.csv`,
          'text/csv'
        );
        break;
        
      case 'json':
        const jsonData = {
          exportDate: new Date().toISOString(),
          version: '1.0',
          notes,
          tags
        };
        downloadFile(
          JSON.stringify(jsonData, null, 2),
          `${filename}.json`,
          'application/json'
        );
        break;
        
      default:
        throw new Error('不支持的格式');
    }
    
    return true;
  } catch (e) {
    console.error('导出失败:', e);
    throw e;
  }
}

// ===== 导入工具 =====

// 解析导入的文件
export async function parseImportFile(file) {
  const content = await file.text();
  const ext = file.name.split('.').pop().toLowerCase();
  
  switch (ext) {
    case 'json':
      return JSON.parse(content);
    case 'md':
    case 'txt':
      return parseMarkdown(content);
    default:
      throw new Error('不支持的文件格式');
  }
}

// 解析 Markdown 导入
function parseMarkdown(content) {
  const notes = [];
  const lines = content.split('\n');
  let currentNote = null;
  
  lines.forEach(line => {
    // 检测笔记开始
    if (line.startsWith('### ')) {
      if (currentNote) notes.push(currentNote);
      currentNote = {
        title: line.replace('### ', '').trim(),
        content: '',
        tags: []
      };
    } else if (currentNote) {
      // 解析标签
      if (line.match(/^#[\u4e00-\u9fa5a-zA-Z0-9_]+/)) {
        const tags = line.match(/#[\u4e00-\u9fa5a-zA-Z0-9_]+/g);
        currentNote.tags = tags.map(t => ({
          name: t.replace('#', ''),
          color: '#3b82f6'
        }));
      } else {
        currentNote.content += line + '\n';
      }
    }
  });
  
  if (currentNote) notes.push(currentNote);
  return { notes };
}

// 从解析结果创建笔记
export async function createNotesFromImport(data, apiClient) {
  if (!data.notes || !Array.isArray(data.notes)) {
    throw new Error('无效的导入数据');
  }
  
  const created = [];
  for (const note of data.notes) {
    try {
      const createdNote = await apiClient.createNote({
        title: note.title || '导入笔记',
        content: note.content || '',
        tags: note.tags || []
      });
      created.push(createdNote);
    } catch (e) {
      console.warn('创建笔记失败:', note.title, e);
    }
  }
  
  return created;
}

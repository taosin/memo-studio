// 数据导出工具

export function exportToMarkdown(notes) {
  let markdown = '# 笔记导出\n\n';
  markdown += `导出时间: ${new Date().toLocaleString('zh-CN')}\n`;
  markdown += `共 ${notes.length} 条笔记\n\n`;
  markdown += '---\n\n';

  notes.forEach((note, index) => {
    markdown += `## ${note.title || '无标题'}\n\n`;
    
    // 标签
    if (note.tags && note.tags.length > 0) {
      const tags = note.tags.map(t => `\`${t.name}\``).join(' ');
      markdown += `标签: ${tags}\n\n`;
    }

    // 日期
    markdown += `创建时间: ${new Date(note.created_at).toLocaleString('zh-CN')}\n\n`;

    // 内容
    const content = note.content.replace(/<[^>]*>/g, '').trim();
    markdown += `${content}\n\n`;
    markdown += '---\n\n';
  });

  return markdown;
}

export function exportToJSON(notes) {
  const data = {
    exportTime: new Date().toISOString(),
    total: notes.length,
    notes: notes.map(note => ({
      id: note.id,
      title: note.title,
      content: note.content.replace(/<[^>]*>/g, '').trim(),
      tags: (note.tags || []).map(t => ({ name: t.name, color: t.color })),
      created_at: note.created_at,
      updated_at: note.updated_at
    }))
  };

  return JSON.stringify(data, null, 2);
}

export function exportToCSV(notes) {
  const headers = ['标题', '内容', '标签', '创建时间', '更新时间'];
  const rows = notes.map(note => {
    const content = note.content.replace(/<[^>]*>/g, '').replace(/"/g, '""').trim();
    const tags = (note.tags || []).map(t => t.name).join(';');
    return [
      `"${note.title || '无标题'}"`,
      `"${content}"`,
      `"${tags}"`,
      `"${new Date(note.created_at).toLocaleString('zh-CN')}"`,
      `"${new Date(note.updated_at).toLocaleString('zh-CN')}"`
    ].join(',');
  });

  return [headers.join(','), ...rows].join('\n');
}

export function downloadFile(content, filename, mimeType) {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

export async function exportNotes(notes, format = 'markdown') {
  let content = '';
  let filename = '';
  let mimeType = '';

  switch (format) {
    case 'markdown':
      content = exportToMarkdown(notes);
      filename = `notes-${new Date().toISOString().split('T')[0]}.md`;
      mimeType = 'text/markdown';
      break;
    case 'json':
      content = exportToJSON(notes);
      filename = `notes-${new Date().toISOString().split('T')[0]}.json`;
      mimeType = 'application/json';
      break;
    case 'csv':
      content = exportToCSV(notes);
      filename = `notes-${new Date().toISOString().split('T')[0]}.csv`;
      mimeType = 'text/csv';
      break;
    default:
      throw new Error('不支持的导出格式');
  }

  downloadFile(content, filename, mimeType);
}

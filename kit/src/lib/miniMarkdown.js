function escapeHtml(s) {
  return String(s ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

function safeLink(url) {
  url = String(url ?? '').trim();
  if (!url) return '';
  // 仅允许 http/https/mailto，避免 javascript: 注入
  if (/^(https?:\/\/|mailto:)/i.test(url)) return url;
  return '';
}

function inline(mdLine) {
  // 在 escape 之后做简单替换，避免 XSS
  let s = escapeHtml(mdLine);

  // 链接 [text](url)
  s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_, text, url) => {
    const href = safeLink(url);
    if (!href) return escapeHtml(text);
    return `<a href="${escapeHtml(href)}" target="_blank" rel="noreferrer noopener">${escapeHtml(text)}</a>`;
  });

  // 加粗 **text**
  s = s.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
  // 斜体 *text*（简单版）
  s = s.replace(/\*([^*]+)\*/g, '<em>$1</em>');

  // 行内代码 `code`
  s = s.replace(/`([^`]+)`/g, '<code>$1</code>');

  return s;
}

export function renderMiniMarkdown(md) {
  const lines = String(md ?? '').split(/\r?\n/);
  let html = '';
  let inList = false;

  function closeList() {
    if (inList) {
      html += '</ul>';
      inList = false;
    }
  }

  for (const raw of lines) {
    const line = raw ?? '';
    const m = line.match(/^\s*-\s+(.*)$/);
    if (m) {
      if (!inList) {
        html += '<ul>';
        inList = true;
      }
      html += `<li>${inline(m[1])}</li>`;
      continue;
    }

    closeList();

    if (line.trim() === '') {
      // 空行：分段
      html += '<div style="height:8px"></div>';
      continue;
    }

    html += `<p>${inline(line)}</p>`;
  }

  closeList();
  return html;
}


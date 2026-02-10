#!/usr/bin/env node
/**
 * 生成软件著作权用：源程序前30页、后30页（每页50行）
 * 用法: node scripts/gen-source-pages.js
 * 输出: docs/源程序前30页.txt, docs/源程序后30页.txt
 */

const fs = require('fs');
const path = require('path');

const ROOT = path.join(__dirname, '..');

const BACKEND_FILES = [
  'backend/main.go',
  'backend/database/database.go',
  'backend/utils/jwt.go',
  'backend/middleware/auth.go',
  'backend/handlers/auth.go',
  'backend/handlers/notes.go',
  'backend/handlers/memos.go',
  'backend/models/user.go',
  'backend/models/note.go',
  'backend/models/memo_query.go',
  'backend/models/notebook.go',
  'backend/models/resource.go',
  'backend/models/stats.go',
  'backend/handlers/users.go',
  'backend/handlers/notebooks.go',
  'backend/handlers/resources.go',
  'backend/handlers/export.go',
  'backend/handlers/import_handler.go',
  'backend/handlers/review.go',
  'backend/handlers/search.go',
  'backend/handlers/stats.go',
];

const KIT_FILES = [
  'kit/src/lib/api.js',
  'kit/src/lib/stores.js',
  'kit/src/lib/heatmap.js',
  'kit/src/lib/miniMarkdown.js',
  'kit/src/lib/theme.js',
  'kit/src/routes/+layout.js',
  'kit/src/routes/+layout.svelte',
  'kit/src/routes/+page.svelte',
  'kit/src/routes/login/+page.svelte',
  'kit/src/routes/profile/+page.svelte',
  'kit/src/routes/admin/users/+page.svelte',
  'kit/src/routes/export/+page.svelte',
  'kit/src/routes/help/+page.svelte',
  'kit/src/routes/import/+page.svelte',
  'kit/src/routes/notebooks/+page.svelte',
  'kit/src/routes/resources/+page.svelte',
  'kit/src/routes/settings/+page.svelte',
  'kit/src/routes/stats/+page.svelte',
  'kit/src/routes/tags/+page.svelte',
];

function readAllLines() {
  const lines = [];
  for (const f of [...BACKEND_FILES, ...KIT_FILES]) {
    const filePath = path.join(ROOT, f);
    if (!fs.existsSync(filePath)) continue;
    const content = fs.readFileSync(filePath, 'utf8');
    const fileLines = content.split(/\r?\n/);
    for (const line of fileLines) {
      lines.push(line);
    }
  }
  return lines;
}

function writePages(filename, lineArray, takeFirst) {
  const total = lineArray.length;
  const LINES_PER_PAGE = 50;
  const PAGES = 30;
  const needLines = LINES_PER_PAGE * PAGES; // 1500

  let slice;
  if (takeFirst) {
    slice = lineArray.slice(0, needLines);
  } else {
    slice = lineArray.slice(-needLines);
  }

  const out = [];
  out.push('Memo Studio 笔记应用软件 V1.0 源程序');
  out.push(takeFirst ? '（前 30 页，每页 50 行）' : '（后 30 页，每页 50 行）');
  out.push('');

  for (let p = 0; p < PAGES; p++) {
    const pageNum = takeFirst ? p + 1 : 31 + p;
    out.push('----- 第 ' + pageNum + ' 页 / 共 60 页 -----');
    const start = p * LINES_PER_PAGE;
    const chunk = slice.slice(start, start + LINES_PER_PAGE);
    for (const line of chunk) {
      out.push(line);
    }
    out.push('');
  }

  const outPath = path.join(ROOT, 'docs', filename);
  fs.mkdirSync(path.dirname(outPath), { recursive: true });
  fs.writeFileSync(outPath, out.join('\n'), 'utf8');
  console.log('Written:', outPath, '(' + slice.length + ' lines, 30 pages)');
}

const allLines = readAllLines();
console.log('Total source lines:', allLines.length);

writePages('源程序前30页.txt', allLines, true);
writePages('源程序后30页.txt', allLines, false);

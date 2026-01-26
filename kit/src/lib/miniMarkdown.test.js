import test from 'node:test';
import assert from 'node:assert/strict';

import { renderMiniMarkdown } from './miniMarkdown.js';

test('renderMiniMarkdown escapes HTML to prevent XSS', () => {
  const html = renderMiniMarkdown(`<img src=x onerror=alert(1)>`);
  assert.ok(html.includes('&lt;img'));
  assert.ok(!html.includes('<img'));
});

test('renderMiniMarkdown blocks javascript: links', () => {
  const html = renderMiniMarkdown(`[x](javascript:alert(1))`);
  assert.ok(!html.includes('href="javascript:'));
});

test('renderMiniMarkdown renders bold/list/link/code', () => {
  const md = [
    '**b**',
    '- item',
    '[t](https://example.com)',
    '`code`'
  ].join('\n');
  const html = renderMiniMarkdown(md);
  assert.ok(html.includes('<strong>b</strong>'));
  assert.ok(html.includes('<ul>') && html.includes('<li>'));
  assert.ok(html.includes('href="https://example.com"'));
  assert.ok(html.includes('<code>code</code>'));
});


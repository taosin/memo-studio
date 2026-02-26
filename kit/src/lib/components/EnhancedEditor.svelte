<script>
  import { createEventDispatcher, onMount, tick } from 'svelte';
  import VoiceInput from './VoiceInput.svelte';
  import EmojiPicker from './EmojiPicker.svelte';
  import AIAssist from './AIAssist.svelte';

  export let value = '';
  export let placeholder = '现在的想法是...';
  export let disabled = false;
  export let uploadLoading = false;

  const dispatch = createEventDispatcher();

  let editorEl;
  let renderTimer;
  let isComposing = false;
  let suppressRender = false;

  // 弹窗状态
  let showVoice = false;
  let showEmoji = false;
  let showAI = false;

  // ========== 光标位置管理 ==========
  function getCaretOffset() {
    const sel = window.getSelection();
    if (!sel.rangeCount || !editorEl) return 0;
    const range = sel.getRangeAt(0).cloneRange();
    range.selectNodeContents(editorEl);
    range.setEnd(sel.getRangeAt(0).startContainer, sel.getRangeAt(0).startOffset);
    return range.toString().length;
  }

  function setCaretOffset(offset) {
    if (!editorEl) return;
    const sel = window.getSelection();
    const range = document.createRange();

    let currentOffset = 0;
    let found = false;

    function walk(node) {
      if (found) return;
      if (node.nodeType === Node.TEXT_NODE) {
        const len = node.textContent.length;
        if (currentOffset + len >= offset) {
          range.setStart(node, offset - currentOffset);
          range.collapse(true);
          found = true;
          return;
        }
        currentOffset += len;
      } else {
        // 处理 <br> 和块级元素的换行
        if (node.nodeName === 'BR') {
          if (currentOffset + 1 >= offset) {
            range.setStartAfter(node);
            range.collapse(true);
            found = true;
            return;
          }
          currentOffset += 1;
        }
        for (const child of node.childNodes) {
          walk(child);
          if (found) return;
        }
        // 块级元素后有隐含换行
        if (node !== editorEl && node.nodeType === Node.ELEMENT_NODE) {
          const display = window.getComputedStyle(node).display;
          if (display === 'block' || display === 'list-item') {
            // 不额外加 - innerText 已经包含了
          }
        }
      }
    }

    walk(editorEl);

    if (!found) {
      // 光标放到末尾
      range.selectNodeContents(editorEl);
      range.collapse(false);
    }

    sel.removeAllRanges();
    sel.addRange(range);
  }

  // ========== 实时行内格式化渲染 ==========
  function escapeHtml(s) {
    return String(s ?? '')
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#39;');
  }

  function formatInline(text) {
    let s = escapeHtml(text);

    // 图片 ![alt](url)
    s = s.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, (_, alt, url) => {
      if (/^(https?:\/\/)/i.test(url.trim())) {
        return `<img src="${escapeHtml(url.trim())}" alt="${escapeHtml(alt)}" style="max-width:100%;border-radius:8px;margin:4px 0;" loading="lazy" />`;
      }
      return `<img src="${escapeHtml(url.trim())}" alt="${escapeHtml(alt)}" style="max-width:100%;border-radius:8px;margin:4px 0;" loading="lazy" />`;
    });

    // 链接 [text](url)
    s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_, text, url) => {
      const href = /^(https?:\/\/|mailto:)/i.test(url.trim()) ? url.trim() : '';
      if (!href) return escapeHtml(text);
      return `<a href="${escapeHtml(href)}" target="_blank" rel="noreferrer noopener" style="color:var(--accent);text-decoration:underline;">${text}</a>`;
    });

    // 粗体 **text**
    s = s.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
    // 斜体 *text*
    s = s.replace(/(?<!\*)\*([^*]+)\*(?!\*)/g, '<em>$1</em>');
    // 行内代码 `code`
    s = s.replace(/`([^`]+)`/g, '<code style="background:rgba(148,163,184,0.15);padding:2px 6px;border-radius:4px;font-family:monospace;font-size:0.9em;">$1</code>');
    // 标签 #tag
    s = s.replace(/(^|\s)#([\p{L}\p{N}_-]+)/gu, '$1<span style="color:var(--accent);font-weight:500;">#$2</span>');

    return s;
  }

  function renderContent(raw) {
    const lines = String(raw ?? '').split('\n');
    let html = '';
    let inList = false;

    for (const line of lines) {
      const trimmed = line.trim();

      // 无序列表 - item
      const ulMatch = trimmed.match(/^-\s+(.*)$/);
      if (ulMatch) {
        if (!inList) { html += '<ul style="margin:4px 0;padding-left:20px;">'; inList = true; }
        html += `<li>${formatInline(ulMatch[1])}</li>`;
        continue;
      }

      // 有序列表 1. item
      const olMatch = trimmed.match(/^\d+\.\s+(.*)$/);
      if (olMatch) {
        if (!inList) { html += '<ul style="margin:4px 0;padding-left:20px;">'; inList = true; }
        html += `<li>${formatInline(olMatch[1])}</li>`;
        continue;
      }

      if (inList) { html += '</ul>'; inList = false; }

      // 任务列表 - [ ] task / - [x] task
      const taskMatch = trimmed.match(/^-\s+\[([ x])\]\s+(.*)$/i);
      if (taskMatch) {
        const checked = taskMatch[1].toLowerCase() === 'x';
        html += `<div style="display:flex;align-items:center;gap:6px;margin:2px 0;"><span style="opacity:${checked ? '0.5' : '1'};${checked ? 'text-decoration:line-through;' : ''}">${checked ? '&#9745;' : '&#9744;'} ${formatInline(taskMatch[2])}</span></div>`;
        continue;
      }

      // 引用 > quote
      if (trimmed.startsWith('> ')) {
        html += `<div style="border-left:3px solid var(--accent);padding-left:12px;color:var(--muted);margin:4px 0;">${formatInline(trimmed.slice(2))}</div>`;
        continue;
      }

      // 分隔线 ---
      if (/^-{3,}$/.test(trimmed)) {
        html += '<hr style="border:none;border-top:1px solid var(--border);margin:8px 0;" />';
        continue;
      }

      // 标题 # heading
      const headingMatch = trimmed.match(/^(#{1,3})\s+(.*)$/);
      if (headingMatch) {
        const level = headingMatch[1].length;
        const sizes = { 1: '1.4em', 2: '1.2em', 3: '1.05em' };
        html += `<div style="font-size:${sizes[level]};font-weight:700;margin:6px 0;">${formatInline(headingMatch[2])}</div>`;
        continue;
      }

      // 空行
      if (trimmed === '') {
        html += '<br>';
        continue;
      }

      // 普通段落
      html += `<div>${formatInline(line)}</div>`;
    }

    if (inList) html += '</ul>';
    return html;
  }

  // ========== 内容同步 ==========
  function getPlainText() {
    if (!editorEl) return '';
    // 用 innerText 获取换行正确的纯文本
    return editorEl.innerText || '';
  }

  function syncValueFromEditor() {
    const text = getPlainText();
    // 去掉末尾多余换行
    value = text.replace(/\n$/, '');
    dispatch('input');
  }

  function renderToEditor() {
    if (!editorEl || isComposing || suppressRender) return;
    const caretPos = getCaretOffset();
    const html = renderContent(value);
    // 仅在内容实际变化时更新 DOM
    if (editorEl.innerHTML !== html) {
      editorEl.innerHTML = html || '';
      setCaretOffset(caretPos);
    }
  }

  function scheduleRender() {
    clearTimeout(renderTimer);
    renderTimer = setTimeout(renderToEditor, 120);
  }

  // ========== 事件处理 ==========
  function handleInput() {
    syncValueFromEditor();
    scheduleRender();
  }

  function handleCompositionStart() {
    isComposing = true;
  }

  function handleCompositionEnd() {
    isComposing = false;
    syncValueFromEditor();
    scheduleRender();
  }

  function handleKeydown(e) {
    // Ctrl/Cmd + Enter 提交
    if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
      e.preventDefault();
      dispatch('submit');
      return;
    }

    // Enter 键：插入纯换行而不是 <div>
    if (e.key === 'Enter' && !e.shiftKey && !e.metaKey && !e.ctrlKey) {
      e.preventDefault();
      document.execCommand('insertLineBreak');
      syncValueFromEditor();
      return;
    }

    // 格式快捷键
    if (e.metaKey || e.ctrlKey) {
      switch (e.key.toLowerCase()) {
        case 'b':
          e.preventDefault();
          wrapSelection('**', '**');
          break;
        case 'i':
          e.preventDefault();
          wrapSelection('*', '*');
          break;
        case 'k':
          e.preventDefault();
          insertLink();
          break;
        case '`':
          e.preventDefault();
          wrapSelection('`', '`');
          break;
      }
    }
  }

  // ========== 格式操作（操作纯文本 value 后重新渲染）==========
  function wrapSelection(before, after) {
    // 先同步当前编辑器状态
    syncValueFromEditor();

    const sel = window.getSelection();
    const selectedText = sel.toString();

    if (selectedText) {
      // 有选中文本：包裹
      const caretStart = getCaretOffset();
      const caretEnd = caretStart + selectedText.length;
      value = value.substring(0, caretStart) + before + selectedText + after + value.substring(caretEnd);

      suppressRender = true;
      editorEl.innerHTML = renderContent(value);
      setCaretOffset(caretStart + before.length + selectedText.length + after.length);
      suppressRender = false;
      dispatch('input');
    } else {
      // 无选中：插入语法标记，光标在中间
      const caretPos = getCaretOffset();
      const placeholder = before === '**' ? '粗体' : before === '*' ? '斜体' : before === '`' ? 'code' : '';
      value = value.substring(0, caretPos) + before + placeholder + after + value.substring(caretPos);

      suppressRender = true;
      editorEl.innerHTML = renderContent(value);
      setCaretOffset(caretPos + before.length + placeholder.length);
      suppressRender = false;
      dispatch('input');
    }
  }

  function insertLink() {
    syncValueFromEditor();
    const sel = window.getSelection();
    const selectedText = sel.toString();
    const caretPos = getCaretOffset();

    if (selectedText) {
      const caretEnd = caretPos + selectedText.length;
      value = value.substring(0, caretPos) + '[' + selectedText + '](url)' + value.substring(caretEnd);
    } else {
      value = value.substring(0, caretPos) + '[链接文本](url)' + value.substring(caretPos);
    }

    suppressRender = true;
    editorEl.innerHTML = renderContent(value);
    setCaretOffset(caretPos + (selectedText ? selectedText.length + 3 : 6));
    suppressRender = false;
    dispatch('input');
  }

  function insertAtOffset(text) {
    syncValueFromEditor();
    const pos = getCaretOffset();
    value = value.substring(0, pos) + text + value.substring(pos);
    suppressRender = true;
    editorEl.innerHTML = renderContent(value);
    setCaretOffset(pos + text.length);
    suppressRender = false;
    dispatch('input');
  }

  // 工具栏格式按钮
  function formatBold() { wrapSelection('**', '**'); editorEl?.focus(); }
  function formatItalic() { wrapSelection('*', '*'); editorEl?.focus(); }
  function formatCode() { wrapSelection('`', '`'); editorEl?.focus(); }
  function formatLinkBtn() { insertLink(); editorEl?.focus(); }

  function formatList() {
    syncValueFromEditor();
    const pos = getCaretOffset();
    const before = value.substring(0, pos);
    const needNl = before.length > 0 && !before.endsWith('\n');
    insertAtOffset((needNl ? '\n' : '') + '- ');
    editorEl?.focus();
  }

  function formatOrderedList() {
    syncValueFromEditor();
    const pos = getCaretOffset();
    const before = value.substring(0, pos);
    const needNl = before.length > 0 && !before.endsWith('\n');
    insertAtOffset((needNl ? '\n' : '') + '1. ');
    editorEl?.focus();
  }

  function formatQuote() {
    syncValueFromEditor();
    const pos = getCaretOffset();
    const before = value.substring(0, pos);
    const needNl = before.length > 0 && !before.endsWith('\n');
    insertAtOffset((needNl ? '\n' : '') + '> ');
    editorEl?.focus();
  }

  function formatTask() {
    syncValueFromEditor();
    const pos = getCaretOffset();
    const before = value.substring(0, pos);
    const needNl = before.length > 0 && !before.endsWith('\n');
    insertAtOffset((needNl ? '\n' : '') + '- [ ] ');
    editorEl?.focus();
  }

  function addTag() {
    insertAtOffset('#');
    editorEl?.focus();
  }

  // ========== 弹窗回调 ==========
  function handleVoiceResult(e) {
    const text = e.detail;
    if (text) {
      value = (value || '') + text;
      suppressRender = true;
      if (editorEl) editorEl.innerHTML = renderContent(value);
      suppressRender = false;
    }
    showVoice = false;
    editorEl?.focus();
  }

  function handleEmojiSelect(e) {
    const emoji = e.detail;
    if (emoji) {
      insertAtOffset(emoji);
    }
    showEmoji = false;
    editorEl?.focus();
  }

  function handleAIResult(e) {
    const { action, result } = e.detail;
    if (result) {
      if (action === 'continue') {
        value = (value || '') + result;
      } else {
        value = result;
      }
      suppressRender = true;
      if (editorEl) editorEl.innerHTML = renderContent(value);
      suppressRender = false;
    }
    showAI = false;
    editorEl?.focus();
  }

  // ========== 生命周期 ==========
  onMount(() => {
    if (value && editorEl) {
      editorEl.innerHTML = renderContent(value);
    }
  });

  // 外部 value 变化时同步到编辑器（如清空）
  $: if (editorEl && !suppressRender && !isComposing) {
    const currentText = getPlainText().replace(/\n$/, '');
    if (value !== currentText) {
      const html = renderContent(value);
      editorEl.innerHTML = html || '';
      // 光标放末尾
      if (value) {
        const sel = window.getSelection();
        const range = document.createRange();
        range.selectNodeContents(editorEl);
        range.collapse(false);
        sel.removeAllRanges();
        sel.addRange(range);
      }
    }
  }

  // ========== 暴露给父组件 ==========
  export function focus() {
    editorEl?.focus();
  }

  export function insertAtCursor(text) {
    insertAtOffset(text);
  }
</script>

<div class="enhanced-editor">
  <!-- 工具栏 -->
  <div class="toolbar">
    <div class="toolbar-left">
      <button class="tool-btn" on:click={formatBold} title="粗体 (Ctrl+B)" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M6 4h8a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"></path>
          <path d="M6 12h9a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"></path>
        </svg>
      </button>
      <button class="tool-btn" on:click={formatItalic} title="斜体 (Ctrl+I)" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="4" x2="10" y2="4"></line>
          <line x1="14" y1="20" x2="5" y2="20"></line>
          <line x1="15" y1="4" x2="9" y2="20"></line>
        </svg>
      </button>
      <button class="tool-btn" on:click={formatCode} title="代码 (Ctrl+`)" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="16 18 22 12 16 6"></polyline>
          <polyline points="8 6 2 12 8 18"></polyline>
        </svg>
      </button>
      <button class="tool-btn" on:click={formatLinkBtn} title="链接 (Ctrl+K)" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
          <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
        </svg>
      </button>

      <span class="divider"></span>

      <button class="tool-btn" on:click={formatList} title="无序列表" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line>
          <line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line>
        </svg>
      </button>
      <button class="tool-btn" on:click={formatOrderedList} title="有序列表" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="10" y1="6" x2="21" y2="6"></line><line x1="10" y1="12" x2="21" y2="12"></line><line x1="10" y1="18" x2="21" y2="18"></line>
          <path d="M4 6h1v4"></path><path d="M4 10h2"></path>
          <path d="M6 18H4c0-1 2-2 2-3s-1-1.5-2-1"></path>
        </svg>
      </button>
      <button class="tool-btn" on:click={formatQuote} title="引用" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 21c3 0 7-1 7-8V5c0-1.25-.756-2.017-2-2H4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2 1 0 1 0 1 1v1c0 1-1 2-2 2s-1 .008-1 1.031V21z"></path>
          <path d="M15 21c3 0 7-1 7-8V5c0-1.25-.757-2.017-2-2h-4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2h.75c0 2.25.25 4-2.75 4v3z"></path>
        </svg>
      </button>
      <button class="tool-btn" on:click={formatTask} title="任务列表" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="5" width="6" height="6" rx="1"></rect>
          <path d="m3 17 2 2 4-4"></path>
          <path d="M13 6h8"></path><path d="M13 12h8"></path><path d="M13 18h8"></path>
        </svg>
      </button>
      <button class="tool-btn" on:click={addTag} title="添加标签" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="4" y1="9" x2="20" y2="9"></line><line x1="4" y1="15" x2="20" y2="15"></line>
          <line x1="10" y1="3" x2="8" y2="21"></line><line x1="16" y1="3" x2="14" y2="21"></line>
        </svg>
      </button>
    </div>

    <div class="toolbar-right">
      <button class="tool-btn" on:click={() => showVoice = true} title="语音输入" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3Z"></path>
          <path d="M19 10v2a7 7 0 0 1-14 0v-2"></path>
          <line x1="12" y1="19" x2="12" y2="22"></line>
        </svg>
      </button>
      <button class="tool-btn ai-btn" on:click={() => showAI = true} title="AI 辅助" disabled={disabled || !value?.trim()}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2l2.4 7.4h7.6l-6 4.6 2.3 7-6.3-4.6-6.3 4.6 2.3-7-6-4.6h7.6z"></path>
        </svg>
      </button>
      <button class="tool-btn" on:click={() => showEmoji = true} title="表情符号" disabled={disabled}>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
          <line x1="9" y1="9" x2="9.01" y2="9"></line>
          <line x1="15" y1="9" x2="15.01" y2="9"></line>
        </svg>
      </button>
    </div>
  </div>

  <!-- contenteditable 编辑区域（输入即预览） -->
  <div class="editor-body">
    <div
      bind:this={editorEl}
      class="editor-input"
      contenteditable={!disabled}
      role="textbox"
      aria-multiline="true"
      aria-placeholder={placeholder}
      data-placeholder={placeholder}
      on:input={handleInput}
      on:keydown={handleKeydown}
      on:compositionstart={handleCompositionStart}
      on:compositionend={handleCompositionEnd}
      on:paste
      on:drop
      on:dragover
    ></div>

    {#if uploadLoading}
      <div class="upload-indicator">图片上传中...</div>
    {/if}
  </div>
</div>

<!-- 弹窗组件 -->
{#if showVoice}
  <VoiceInput on:result={handleVoiceResult} on:close={() => showVoice = false} />
{/if}

{#if showEmoji}
  <EmojiPicker on:select={handleEmojiSelect} on:close={() => showEmoji = false} />
{/if}

{#if showAI}
  <AIAssist content={value} on:result={handleAIResult} on:close={() => showAI = false} />
{/if}

<style>
  .enhanced-editor {
    border: 1px solid var(--border-2);
    background: var(--panel);
    border-radius: 16px;
    overflow: hidden;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }

  .enhanced-editor:focus-within {
    border-color: rgba(34, 197, 94, 0.4);
    box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.08);
  }

  /* 工具栏 */
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    border-bottom: 1px solid var(--border);
    background: var(--panel-2);
    gap: 8px;
    flex-wrap: wrap;
  }

  .toolbar-left, .toolbar-right {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .tool-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    color: var(--muted);
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .tool-btn:hover:not(:disabled) {
    background: rgba(148, 163, 184, 0.12);
    color: var(--text);
  }

  .tool-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .tool-btn.ai-btn {
    color: #f59e0b;
  }

  .tool-btn.ai-btn:hover:not(:disabled) {
    background: rgba(245, 158, 11, 0.12);
  }

  .divider {
    width: 1px;
    height: 20px;
    background: var(--border);
    margin: 0 6px;
  }

  /* 编辑器主体 */
  .editor-body {
    position: relative;
    padding: 12px 16px;
  }

  .editor-input {
    min-height: 100px;
    outline: none;
    color: var(--text);
    font-size: 15px;
    line-height: 1.7;
    font-family: inherit;
    word-break: break-word;
    white-space: pre-wrap;
  }

  /* placeholder */
  .editor-input:empty::before {
    content: attr(data-placeholder);
    color: var(--muted);
    opacity: 0.6;
    pointer-events: none;
  }

  /* contenteditable 内的格式化元素样式 */
  .editor-input :global(strong) {
    font-weight: 700;
  }

  .editor-input :global(em) {
    font-style: italic;
  }

  .editor-input :global(code) {
    background: rgba(148, 163, 184, 0.15);
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 0.9em;
  }

  .editor-input :global(a) {
    color: var(--accent);
    text-decoration: underline;
  }

  .editor-input :global(ul) {
    margin: 4px 0;
    padding-left: 20px;
  }

  .editor-input :global(li) {
    margin: 2px 0;
  }

  .editor-input :global(hr) {
    border: none;
    border-top: 1px solid var(--border);
    margin: 8px 0;
  }

  .editor-input :global(img) {
    max-width: 100%;
    border-radius: 8px;
    margin: 4px 0;
  }

  .upload-indicator {
    position: absolute;
    bottom: 16px;
    right: 16px;
    font-size: 12px;
    color: var(--accent);
    background: var(--panel);
    padding: 4px 8px;
    border-radius: 4px;
  }

  /* 响应式 */
  @media (max-width: 600px) {
    .toolbar {
      padding: 6px 8px;
    }

    .tool-btn {
      width: 28px;
      height: 28px;
    }

    .tool-btn :global(svg) {
      width: 14px;
      height: 14px;
    }

    .divider {
      margin: 0 4px;
    }
  }
</style>

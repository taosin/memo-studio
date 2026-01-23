<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import { api } from '../utils/api.js';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import Card from '$lib/components/ui/card/card.svelte';

  export let value = '';
  export let placeholder = '开始记录你的想法...';

  const dispatch = createEventDispatcher();
  let editorElement;
  let allTags = [];
  let allNotes = [];
  let showTagSuggestions = false;
  let showNoteSuggestions = false;
  let suggestions = [];
  let suggestionIndex = -1;
  let triggerPosition = { start: 0, end: 0 };
  let currentTrigger = ''; // '#' or '@'

  onMount(async () => {
    await loadData();
    if (editorElement) {
      editorElement.innerHTML = value || '';
      setupEditor();
    }
  });

  async function loadData() {
    try {
      [allTags, allNotes] = await Promise.all([
        api.getTags(),
        api.getNotes()
      ]);
    } catch (err) {
      console.error('加载数据失败:', err);
    }
  }

  function setupEditor() {
    if (!editorElement) return;

    editorElement.addEventListener('input', handleInput);
    editorElement.addEventListener('keydown', handleKeyDown);
    editorElement.addEventListener('click', hideSuggestions);
  }

  function handleInput(e) {
    const text = editorElement.innerText || '';
    const cursorPos = getCaretPosition();
    
    // 检查 # 标签（支持中文和英文）
    const hashMatch = text.substring(0, cursorPos).match(/#([\w\u4e00-\u9fa5]*)$/);
    if (hashMatch) {
      currentTrigger = '#';
      triggerPosition = {
        start: cursorPos - hashMatch[0].length,
        end: cursorPos
      };
      const query = hashMatch[1].toLowerCase();
      suggestions = allTags.filter(tag => 
        tag.name.toLowerCase().includes(query) || tag.name.includes(query)
      );
      showTagSuggestions = suggestions.length > 0;
      suggestionIndex = -1;
      updateValue();
      return;
    }

    // 检查 @ 笔记引用（支持中文和英文）
    const atMatch = text.substring(0, cursorPos).match(/@([\w\u4e00-\u9fa5\s]*)$/);
    if (atMatch) {
      currentTrigger = '@';
      triggerPosition = {
        start: cursorPos - atMatch[0].length,
        end: cursorPos
      };
      const query = atMatch[1].toLowerCase();
      suggestions = allNotes.filter(note => {
        const title = (note.title || '').toLowerCase();
        const content = (note.content || '').toLowerCase();
        return title.includes(query) || content.includes(query);
      });
      showNoteSuggestions = suggestions.length > 0;
      suggestionIndex = -1;
      updateValue();
      return;
    }

    hideSuggestions();
    updateValue();
  }

  function handleKeyDown(e) {
    if (showTagSuggestions || showNoteSuggestions) {
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        suggestionIndex = Math.min(suggestionIndex + 1, suggestions.length - 1);
        return;
      }
      if (e.key === 'ArrowUp') {
        e.preventDefault();
        suggestionIndex = Math.max(suggestionIndex - 1, -1);
        return;
      }
      if (e.key === 'Enter' || e.key === 'Tab') {
        e.preventDefault();
        if (suggestionIndex >= 0 && suggestions[suggestionIndex]) {
          insertSuggestion(suggestions[suggestionIndex]);
        }
        return;
      }
      if (e.key === 'Escape') {
        hideSuggestions();
        return;
      }
    }
  }

  function insertSuggestion(item) {
    if (!editorElement) return;

    const selection = window.getSelection();
    if (selection.rangeCount === 0) return;
    
    const range = selection.getRangeAt(0);
    const textNode = editorElement.childNodes[0] || editorElement;
    
    // 获取当前文本
    const text = editorElement.innerText || '';
    const before = text.substring(0, triggerPosition.start);
    const after = text.substring(triggerPosition.end);
    
    let insertText = '';
    if (currentTrigger === '#') {
      insertText = `#${item.name} `;
    } else if (currentTrigger === '@') {
      insertText = `@${item.title || '无标题'} `;
    }

    // 创建新的文本节点
    const newText = before + insertText + after;
    editorElement.innerHTML = '';
    editorElement.appendChild(document.createTextNode(newText));
    
    // 设置光标位置
    setTimeout(() => {
      setCaretPosition(before.length + insertText.length);
    }, 0);
    
    hideSuggestions();
    updateValue();
  }

  function getCaretPosition() {
    const selection = window.getSelection();
    if (selection.rangeCount === 0) return 0;
    const range = selection.getRangeAt(0);
    const preCaretRange = range.cloneRange();
    preCaretRange.selectNodeContents(editorElement);
    preCaretRange.setEnd(range.endContainer, range.endOffset);
    return preCaretRange.toString().length;
  }

  function setCaretPosition(position) {
    const range = document.createRange();
    const selection = window.getSelection();
    
    let charCount = 0;
    const nodeStack = [editorElement];
    let node;
    let foundStart = false;

    while (!foundStart && (node = nodeStack.pop())) {
      if (node.nodeType === 3) {
        const nextCharCount = charCount + node.length;
        if (position >= charCount && position <= nextCharCount) {
          range.setStart(node, position - charCount);
          range.setEnd(node, position - charCount);
          foundStart = true;
        }
        charCount = nextCharCount;
      } else {
        let i = node.childNodes.length;
        while (i--) {
          nodeStack.push(node.childNodes[i]);
        }
      }
    }

    selection.removeAllRanges();
    selection.addRange(range);
  }

  function hideSuggestions() {
    showTagSuggestions = false;
    showNoteSuggestions = false;
    suggestions = [];
    suggestionIndex = -1;
  }

  function updateValue() {
    value = editorElement.innerHTML;
    dispatch('input', { detail: value });
  }

  $: if (editorElement && value && value !== editorElement.innerHTML) {
    // 只在编辑器不处于焦点状态时更新
    if (!editorElement.contains(document.activeElement)) {
      editorElement.innerHTML = value || '';
    }
  }
</script>

<div class="relative">
  <div
    bind:this={editorElement}
    contenteditable="true"
    class="min-h-[400px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 prose prose-sm dark:prose-invert max-w-none"
    data-placeholder={placeholder}
    role="textbox"
  ></div>

  <!-- 标签建议 -->
  {#if showTagSuggestions && suggestions.length > 0}
    <Card class="absolute z-50 mt-1 max-h-60 overflow-y-auto w-64 shadow-lg">
      <div class="p-2">
        {#each suggestions as suggestion, index}
          <div
            class="p-2 rounded-md cursor-pointer transition-colors {index === suggestionIndex ? 'bg-accent' : 'hover:bg-accent'}"
            on:click={() => insertSuggestion(suggestion)}
            on:mouseenter={() => suggestionIndex = index}
          >
            <Badge
              variant="outline"
              style="border-color: {suggestion.color || '#4ECDC4'}; color: {suggestion.color || '#4ECDC4'}"
            >
              {suggestion.name}
            </Badge>
          </div>
        {/each}
      </div>
    </Card>
  {/if}

  <!-- 笔记建议 -->
  {#if showNoteSuggestions && suggestions.length > 0}
    <Card class="absolute z-50 mt-1 max-h-60 overflow-y-auto w-64 shadow-lg">
      <div class="p-2">
        {#each suggestions as suggestion, index}
          <div
            class="p-2 rounded-md cursor-pointer transition-colors {index === suggestionIndex ? 'bg-accent' : 'hover:bg-accent'}"
            on:click={() => insertSuggestion(suggestion)}
            on:mouseenter={() => suggestionIndex = index}
          >
            <div class="font-medium text-sm">{suggestion.title || '无标题'}</div>
            <div class="text-xs text-muted-foreground line-clamp-1 mt-1">
              {suggestion.content?.substring(0, 50)}
            </div>
          </div>
        {/each}
      </div>
    </Card>
  {/if}
</div>

<style>
  [contenteditable][data-placeholder]:empty:before {
    content: attr(data-placeholder);
    color: hsl(var(--muted-foreground));
    pointer-events: none;
  }

  [contenteditable] {
    outline: none;
  }

  [contenteditable] p {
    margin: 0.5em 0;
  }

  [contenteditable] p:first-child {
    margin-top: 0;
  }

  [contenteditable] p:last-child {
    margin-bottom: 0;
  }
</style>

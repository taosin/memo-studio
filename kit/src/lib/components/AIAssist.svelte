<script>
  import { createEventDispatcher } from 'svelte';
  import { api } from '$lib/api.js';

  export let content = '';

  const dispatch = createEventDispatcher();

  // AI 功能选项
  const actions = [
    { id: 'polish', name: '润色文本', icon: '✨', desc: '优化表达，使文字更流畅' },
    { id: 'continue', name: '继续写作', icon: '📝', desc: '根据上下文续写内容' },
    { id: 'summarize', name: '生成摘要', icon: '📋', desc: '提取关键信息生成摘要' },
    { id: 'translate', name: '翻译内容', icon: '🌍', desc: '翻译为其他语言' },
    { id: 'expand', name: '扩展内容', icon: '📖', desc: '丰富和扩展现有内容' },
    { id: 'simplify', name: '简化表达', icon: '🎯', desc: '使内容更简洁易懂' }
  ];

  let selectedAction = null;
  let loading = false;
  let error = '';
  let result = '';

  async function executeAction(action) {
    selectedAction = action;
    loading = true;
    error = '';
    result = '';

    try {
      const response = await api.aiAssist(action.id, content);
      result = response.result || response.content || '';
    } catch (err) {
      error = err.message || 'AI 处理失败，请稍后重试';
    } finally {
      loading = false;
    }
  }

  function confirm() {
    if (result) {
      dispatch('result', { action: selectedAction.id, result });
    }
  }

  function back() {
    selectedAction = null;
    result = '';
    error = '';
  }

  function close() {
    dispatch('close');
  }
</script>

<div class="modal-backdrop" on:click|self={close}>
  <div class="modal">
    <div class="modal-header">
      <div class="header-left">
        {#if selectedAction}
          <button class="back-btn" on:click={back}>
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="15 18 9 12 15 6"></polyline>
            </svg>
          </button>
        {/if}
        <h3>{selectedAction ? selectedAction.name : 'AI 助手'}</h3>
      </div>
      <button class="close-btn" on:click={close}>
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <div class="modal-body">
      {#if !selectedAction}
        <!-- 功能选择列表 -->
        <div class="actions-list">
          {#each actions as action}
            <button class="action-item" on:click={() => executeAction(action)}>
              <span class="action-icon">{action.icon}</span>
              <div class="action-info">
                <span class="action-name">{action.name}</span>
                <span class="action-desc">{action.desc}</span>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="action-arrow">
                <polyline points="9 18 15 12 9 6"></polyline>
              </svg>
            </button>
          {/each}
        </div>
      {:else}
        <!-- 结果显示 -->
        <div class="result-area">
          <div class="original-section">
            <div class="section-label">原文</div>
            <div class="original-content">{content}</div>
          </div>

          <div class="result-section">
            <div class="section-label">
              {loading ? '处理中...' : '结果'}
            </div>
            
            {#if loading}
              <div class="loading">
                <div class="spinner"></div>
                <span>AI 正在处理...</span>
              </div>
            {:else if error}
              <div class="error">{error}</div>
            {:else if result}
              <div class="result-content">{result}</div>
            {/if}
          </div>
        </div>
      {/if}
    </div>

    {#if selectedAction && result}
      <div class="modal-footer">
        <button class="btn ghost" on:click={back}>
          重新选择
        </button>
        <button class="btn primary" on:click={confirm}>
          使用此结果
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 16px;
  }

  .modal {
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 16px;
    width: 100%;
    max-width: 500px;
    max-height: 80vh;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border);
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
  }

  .back-btn, .close-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    color: var(--muted);
    border-radius: 8px;
    cursor: pointer;
  }

  .back-btn:hover, .close-btn:hover {
    background: rgba(148, 163, 184, 0.1);
    color: var(--text);
  }

  .modal-body {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }

  /* 功能选择列表 */
  .actions-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .action-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    background: var(--panel-2);
    border: 1px solid var(--border);
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.15s ease;
    text-align: left;
  }

  .action-item:hover {
    background: rgba(34, 197, 94, 0.08);
    border-color: rgba(34, 197, 94, 0.3);
  }

  .action-icon {
    font-size: 24px;
    flex-shrink: 0;
  }

  .action-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .action-name {
    font-size: 15px;
    font-weight: 500;
    color: var(--text);
  }

  .action-desc {
    font-size: 13px;
    color: var(--muted);
  }

  .action-arrow {
    color: var(--muted);
    flex-shrink: 0;
  }

  /* 结果区域 */
  .result-area {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .original-section, .result-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .section-label {
    font-size: 12px;
    font-weight: 500;
    color: var(--muted);
    text-transform: uppercase;
  }

  .original-content {
    padding: 12px;
    background: var(--panel-2);
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.6;
    color: var(--muted);
    max-height: 100px;
    overflow-y: auto;
  }

  .result-content {
    padding: 12px;
    background: rgba(34, 197, 94, 0.08);
    border: 1px solid rgba(34, 197, 94, 0.2);
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.6;
    color: var(--text);
    max-height: 200px;
    overflow-y: auto;
    white-space: pre-wrap;
  }

  .loading {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 24px;
    color: var(--muted);
  }

  .spinner {
    width: 24px;
    height: 24px;
    border: 2px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .error {
    padding: 12px;
    background: rgba(248, 113, 113, 0.1);
    border: 1px solid rgba(248, 113, 113, 0.3);
    border-radius: 8px;
    color: #f87171;
    font-size: 14px;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 20px;
    border-top: 1px solid var(--border);
  }

  .btn {
    padding: 10px 20px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .btn.ghost {
    background: transparent;
    border: 1px solid var(--border);
    color: var(--text);
  }

  .btn.ghost:hover {
    background: rgba(148, 163, 184, 0.1);
  }

  .btn.primary {
    background: var(--accent);
    border: none;
    color: white;
  }

  .btn.primary:hover {
    filter: brightness(1.1);
  }
</style>

<script>
  import { onMount } from 'svelte';
  import { toastStore } from '../stores.js';

  let toasts = [];
  const unsub = toastStore.subscribe(value => {
    toasts = value;
  });

  onMount(() => {
    return () => unsub();
  });

  function removeToast(id) {
    toastStore.update(items => items.filter(t => t.id !== id));
  }

  function getIcon(type) {
    switch(type) {
      case 'success': return '✓';
      case 'error': return '✕';
      case 'warning': return '⚠';
      default: return 'ⓘ';
    }
  }
</script>

<div class="toastContainer" role="region" aria-label="通知">
  {#each toasts as toast (toast.id)}
    <div 
      class="toast {toast.type}" 
      role="alert"
      aria-live="polite"
    >
      <div class="toastIcon">{getIcon(toast.type)}</div>
      <div class="toastContent">
        <div class="toastMessage">{toast.message}</div>
      </div>
      <button 
        class="toastClose" 
        on:click={() => removeToast(toast.id)}
        aria-label="关闭通知"
      >
        ✕
      </button>
    </div>
  {/each}
</div>

<style>
  .toastContainer {
    position: fixed;
    top: 70px;
    right: 16px;
    z-index: 9999;
    display: flex;
    flex-direction: column;
    gap: 8px;
    pointer-events: none;
    max-width: 400px;
  }

  .toast {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border-radius: 12px;
    border: 1px solid var(--border);
    background: var(--panel);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    pointer-events: auto;
    animation: slideIn 0.3s ease;
    backdrop-filter: blur(10px);
    min-width: 280px;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateX(100%);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .toast.success {
    border-color: rgba(34, 197, 94, 0.35);
    background: var(--accent-soft);
  }

  .toast.error {
    border-color: rgba(239, 68, 68, 0.35);
    background: rgba(239, 68, 68, 0.12);
  }

  .toast.warning {
    border-color: rgba(251, 191, 36, 0.35);
    background: rgba(251, 191, 36, 0.12);
  }

  .toast.info {
    border-color: rgba(59, 130, 246, 0.35);
    background: rgba(59, 130, 246, 0.12);
  }

  .toastIcon {
    font-size: 18px;
    flex-shrink: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
  }

  .toast.success .toastIcon {
    background: rgba(34, 197, 94, 0.2);
    color: rgba(34, 197, 94, 1);
  }

  .toast.error .toastIcon {
    background: rgba(239, 68, 68, 0.2);
    color: rgba(239, 68, 68, 1);
  }

  .toast.warning .toastIcon {
    background: rgba(251, 191, 36, 0.2);
    color: rgba(251, 191, 36, 1);
  }

  .toast.info .toastIcon {
    background: rgba(59, 130, 246, 0.2);
    color: rgba(59, 130, 246, 1);
  }

  .toastContent {
    flex: 1;
    min-width: 0;
  }

  .toastMessage {
    font-size: 14px;
    line-height: 1.5;
    color: var(--text);
  }

  .toastClose {
    flex-shrink: 0;
    border: none;
    background: transparent;
    color: var(--muted);
    cursor: pointer;
    padding: 4px;
    font-size: 16px;
    border-radius: 6px;
    transition: background 0.15s ease, color 0.15s ease;
  }

  .toastClose:hover {
    background: rgba(148, 163, 184, 0.12);
    color: var(--text);
  }

  @media (max-width: 600px) {
    .toastContainer {
      top: 60px;
      right: 12px;
      left: 12px;
      max-width: none;
    }

    .toast {
      min-width: auto;
    }
  }
</style>

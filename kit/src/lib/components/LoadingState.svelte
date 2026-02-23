<script>
  export let type = 'spinner'; // spinner, skeleton, dots, bar
  export let size = 'md'; // sm, md, lg
  export let text = '加载中...';
</script>

{#if type === 'spinner'}
  <div class="loadingState">
    <div class="spinner {size}"></div>
    {#if text}
      <p class="loadingText">{text}</p>
    {/if}
  </div>
{:else if type === 'dots'}
  <div class="loadingState">
    <div class="loadingDots">
      <span></span>
      <span></span>
      <span></span>
    </div>
    {#if text}
      <p class="loadingText">{text}</p>
    {/if}
  </div>
{:else if type === 'bar'}
  <div class="loadingBar"></div>
{:else if type === 'skeleton'}
  <div class="skeleton">
    <div class="skeletonLine"></div>
    <div class="skeletonLine short"></div>
    <div class="skeletonLine"></div>
  </div>
{/if}

<style>
  .loadingState {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px;
    gap: 16px;
  }

  .spinner {
    border: 3px solid rgba(148, 163, 184, 0.2);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .spinner.sm {
    width: 24px;
    height: 24px;
    border-width: 2px;
  }

  .spinner.md {
    width: 40px;
    height: 40px;
    border-width: 3px;
  }

  .spinner.lg {
    width: 56px;
    height: 56px;
    border-width: 4px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .loadingDots {
    display: flex;
    gap: 8px;
  }

  .loadingDots span {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--accent);
    opacity: 0.6;
    animation: loadingDot 1.4s ease-in-out infinite both;
  }

  .loadingDots span:nth-child(1) { animation-delay: 0s; }
  .loadingDots span:nth-child(2) { animation-delay: 0.2s; }
  .loadingDots span:nth-child(3) { animation-delay: 0.4s; }

  @keyframes loadingDot {
    0%, 80%, 100% { 
      transform: scale(0.8); 
      opacity: 0.5; 
    }
    40% { 
      transform: scale(1.2); 
      opacity: 1; 
    }
  }

  .loadingText {
    margin: 0;
    font-size: 14px;
    color: var(--muted);
  }

  .loadingBar {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(34, 197, 94, 0.5),
      transparent
    );
    background-size: 200% 100%;
    animation: loadingBar 1.5s ease-in-out infinite;
  }

  @keyframes loadingBar {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }

  .skeleton {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    width: 100%;
  }

  .skeletonLine {
    height: 16px;
    background: linear-gradient(
      90deg,
      rgba(148, 163, 184, 0.1) 25%,
      rgba(148, 163, 184, 0.2) 50%,
      rgba(148, 163, 184, 0.1) 75%
    );
    background-size: 200% 100%;
    border-radius: 8px;
    animation: skeleton 1.5s ease-in-out infinite;
  }

  .skeletonLine.short {
    width: 60%;
  }

  @keyframes skeleton {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }
</style>

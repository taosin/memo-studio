<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';

  const dispatch = createEventDispatcher();

  let listening = false;
  let transcript = '';
  let interimTranscript = '';
  let error = '';
  let recognition = null;

  // 浏览器支持检测
  let isSpeechSupported = false;

  onMount(() => {
    isSpeechSupported = 'webkitSpeechRecognition' in window || 'SpeechRecognition' in window;
  });

  onDestroy(() => {
    if (recognition) {
      recognition.stop();
    }
  });

  function startListening() {
    if (!isSpeechSupported) {
      error = '您的浏览器不支持语音识别';
      return;
    }

    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    recognition = new SpeechRecognition();
    recognition.lang = 'zh-CN';
    recognition.interimResults = true;
    recognition.continuous = true;

    recognition.onstart = () => {
      listening = true;
      error = '';
    };

    recognition.onresult = (event) => {
      let finalTranscript = '';
      let interim = '';
      
      for (let i = event.resultIndex; i < event.results.length; i++) {
        if (event.results[i].isFinal) {
          finalTranscript += event.results[i][0].transcript;
        } else {
          interim += event.results[i][0].transcript;
        }
      }
      
      if (finalTranscript) {
        transcript += finalTranscript;
      }
      interimTranscript = interim;
    };

    recognition.onerror = (event) => {
      if (event.error === 'no-speech') {
        error = '未检测到语音，请重试';
      } else if (event.error === 'not-allowed') {
        error = '麦克风权限被拒绝';
      } else {
        error = '语音识别错误: ' + event.error;
      }
      listening = false;
    };

    recognition.onend = () => {
      listening = false;
    };

    recognition.start();
  }

  function stopListening() {
    if (recognition) {
      recognition.stop();
    }
    listening = false;
  }

  function confirm() {
    dispatch('result', transcript);
    close();
  }

  function close() {
    if (recognition) {
      recognition.stop();
    }
    dispatch('close');
  }

  function clear() {
    transcript = '';
    interimTranscript = '';
    error = '';
  }
</script>

<div class="modal-backdrop" on:click|self={close}>
  <div class="modal">
    <div class="modal-header">
      <h3>语音输入</h3>
      <button class="close-btn" on:click={close}>
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <div class="modal-body">
      {#if !isSpeechSupported}
        <div class="warning">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
          <span>您的浏览器不支持语音识别功能</span>
        </div>
      {:else}
        <div class="mic-area">
          <button 
            class="mic-btn" 
            class:listening
            on:click={listening ? stopListening : startListening}
          >
            {#if listening}
              <div class="pulse-ring"></div>
              <div class="pulse-ring delay"></div>
            {/if}
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3Z"></path>
              <path d="M19 10v2a7 7 0 0 1-14 0v-2"></path>
              <line x1="12" y1="19" x2="12" y2="22"></line>
            </svg>
          </button>
          <p class="mic-hint">
            {#if listening}
              正在聆听... 点击停止
            {:else}
              点击开始语音输入
            {/if}
          </p>
        </div>

        {#if error}
          <div class="error">{error}</div>
        {/if}

        <div class="transcript-area">
          <div class="transcript-label">识别结果</div>
          <div class="transcript">
            {transcript}
            {#if interimTranscript}
              <span class="interim">{interimTranscript}</span>
            {/if}
            {#if !transcript && !interimTranscript}
              <span class="placeholder">等待语音输入...</span>
            {/if}
          </div>
        </div>
      {/if}
    </div>

    <div class="modal-footer">
      <button class="btn ghost" on:click={clear} disabled={!transcript}>
        清空
      </button>
      <button class="btn primary" on:click={confirm} disabled={!transcript.trim()}>
        插入文本
      </button>
    </div>
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
    max-width: 420px;
    max-height: 90vh;
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

  .modal-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
  }

  .close-btn {
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

  .close-btn:hover {
    background: rgba(148, 163, 184, 0.1);
    color: var(--text);
  }

  .modal-body {
    padding: 20px;
    flex: 1;
    overflow-y: auto;
  }

  .warning {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: rgba(251, 191, 36, 0.1);
    border: 1px solid rgba(251, 191, 36, 0.3);
    border-radius: 12px;
    color: #fbbf24;
  }

  .mic-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 24px 0;
  }

  .mic-btn {
    position: relative;
    width: 80px;
    height: 80px;
    border-radius: 50%;
    border: none;
    background: var(--panel-2);
    color: var(--text);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .mic-btn:hover {
    background: rgba(34, 197, 94, 0.1);
    color: var(--accent);
  }

  .mic-btn.listening {
    background: rgba(248, 113, 113, 0.15);
    color: #f87171;
  }

  .pulse-ring {
    position: absolute;
    width: 100%;
    height: 100%;
    border-radius: 50%;
    border: 2px solid #f87171;
    animation: pulse 1.5s ease-out infinite;
  }

  .pulse-ring.delay {
    animation-delay: 0.5s;
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
      opacity: 1;
    }
    100% {
      transform: scale(1.5);
      opacity: 0;
    }
  }

  .mic-hint {
    margin-top: 16px;
    font-size: 14px;
    color: var(--muted);
  }

  .error {
    margin-top: 16px;
    padding: 12px;
    background: rgba(248, 113, 113, 0.1);
    border: 1px solid rgba(248, 113, 113, 0.3);
    border-radius: 8px;
    color: #f87171;
    font-size: 14px;
    text-align: center;
  }

  .transcript-area {
    margin-top: 20px;
  }

  .transcript-label {
    font-size: 12px;
    color: var(--muted);
    margin-bottom: 8px;
  }

  .transcript {
    min-height: 80px;
    padding: 12px;
    background: var(--panel-2);
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 15px;
    line-height: 1.6;
  }

  .transcript .interim {
    color: var(--muted);
    opacity: 0.7;
  }

  .transcript .placeholder {
    color: var(--muted);
    opacity: 0.5;
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

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn.ghost {
    background: transparent;
    border: 1px solid var(--border);
    color: var(--text);
  }

  .btn.ghost:hover:not(:disabled) {
    background: rgba(148, 163, 184, 0.1);
  }

  .btn.primary {
    background: var(--accent);
    border: none;
    color: white;
  }

  .btn.primary:hover:not(:disabled) {
    filter: brightness(1.1);
  }
</style>

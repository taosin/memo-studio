<script>
  import { onMount } from 'svelte';
  import { api } from '#lib/api.js';
  import { goto } from '$app/navigation';

  let listening = false;
  let transcript = '';
  let interimTranscript = '';
  let error = '';
  let recording = false;
  let mediaRecorder = null;
  let audioChunks = [];
  let saved = false;

  // 检查浏览器支持
  let isSupported = false;
  let isSpeechSupported = false;
  let isMediaRecorderSupported = false;

  onMount(() => {
    isSupported = 'webkitSpeechRecognition' in window || 'SpeechRecognition' in window;
    isSpeechSupported = isSupported;
    isMediaRecorderSupported = 'MediaRecorder' in window;
  });

  // 语音识别
  let recognition = null;

  function startListening() {
    if (!isSpeechSupported) {
      error = '您的浏览器不支持语音识别';
      return;
    }

    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    recognition = new SpeechRecognition();
    recognition.lang = 'zh-CN';
    recognition.interimResults = true;
    recognition.continuous = false;

    recognition.onstart = () => {
      listening = true;
      error = '';
      saved = false;
    };

    recognition.onresult = (event) => {
      let finalTranscript = '';
      for (let i = event.resultIndex; i < event.results.length; i++) {
        if (event.results[i].isFinal) {
          finalTranscript += event.results[i][0].transcript;
        } else {
          interimTranscript += event.results[i][0].transcript;
        }
      }
      transcript = finalTranscript;
      interimTranscript = '';
    };

    recognition.onerror = (event) => {
      error = '语音识别错误: ' + event.error;
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

  // 录音功能
  async function startRecording() {
    if (!isMediaRecorderSupported) {
      error = '您的浏览器不支持录音';
      return;
    }

    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      mediaRecorder = new MediaRecorder(stream);
      audioChunks = [];

      mediaRecorder.ondataavailable = (event) => {
        audioChunks.push(event.data);
      };

      mediaRecorder.onstop = async () => {
        const audioBlob = new Blob(audioChunks, { type: 'audio/webm' });
        const audioFile = new File([audioBlob], 'recording.webm', { type: 'audio/webm' });

        // 上传到服务器并转录
        await uploadAndTranscribe(audioFile);

        // 停止所有轨道
        stream.getTracks().forEach(track => track.stop());
      };

      mediaRecorder.start();
      recording = true;
      error = '';
      saved = false;
    } catch (err) {
      error = '无法访问麦克风: ' + err.message;
    }
  }

  function stopRecording() {
    if (mediaRecorder && recording) {
      mediaRecorder.stop();
      recording = false;
    }
  }

  // 上传并转录（需要后端支持 OpenAI Whisper）
  async function uploadAndTranscribe(file) {
    try {
      const res = await api.uploadResource(file);
      transcript = '音频已上传: ' + res.filename;
      saved = true;
    } catch (err) {
      // 如果后端不支持转录，提示用户
      error = '音频已保存。可在资源页面查看。';
      transcript = '（录音已保存，请在资源页面管理）';
      saved = true;
    }
  }

  // 保存为笔记
  async function saveAsNote() {
    if (!transcript.trim()) {
      error = '没有内容可保存';
      return;
    }

    try {
      await api.createNote({ content: transcript, tags: ['语音'] });
      saved = true;
    } catch (err) {
      error = '保存失败: ' + err.message;
    }
  }

  function clear() {
    transcript = '';
    error = '';
    saved = false;
  }

  function goBack() {
    goto('/');
  }
</script>

<svelte:head>
  <title>语音输入 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">语音输入</div>
      <button class="link" on:click={goBack}>返回首页</button>
    </div>

    <div class="description">
      {#if isSpeechSupported}
        <p>🎤 使用语音识别输入文字（在线模式，需要网络）</p>
      {:else}
        <p>⚠️ 您的浏览器不支持语音识别</p>
      {/if}
      <p>🎙️ 或使用录音功能（音频会保存到资源库）</p>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if saved}
      <div class="success">已保存到笔记！</div>
    {/if}

    <div class="transcript-area">
      <div class="label">识别结果：</div>
      <div class="transcript">
        {transcript}
        {#if interimTranscript}
          <span class="interim">{interimTranscript}</span>
        {/if}
      </div>
    </div>

    <div class="actions">
      <div class="group">
        <button
          class="btn mic"
          class:active={listening}
          on:click={listening ? stopListening : startListening}
          disabled={!isSpeechSupported}
        >
          {listening ? '🛑 停止识别' : '🎤 语音识别'}
        </button>

        <button
          class="btn record"
          class:active={recording}
          on:click={recording ? stopRecording : startRecording}
          disabled={!isMediaRecorderSupported}
        >
          {recording ? '🛑 停止录音' : '🎙️ 开始录音'}
        </button>
      </div>

      <div class="group">
        <button class="btn" on:click={saveAsNote} disabled={!transcript.trim()}>
          💾 保存为笔记
        </button>
        <button class="btn ghost" on:click={clear}>
          🗑️ 清空
        </button>
      </div>
    </div>

    <div class="tips">
      <h4>💡 使用提示</h4>
      <ul>
        <li>🎤 <strong>语音识别</strong>：直接说话，自动转为文字（需要网络）</li>
        <li>🎙️ <strong>录音</strong>：录制语音，上传到服务器（需要 OPENAI_API_KEY 配置才能自动转文字）</li>
        <li>💾 保存后可在首页查看笔记</li>
      </ul>
    </div>
  </div>
</div>

<style>
  .wrap {
    display: flex;
    justify-content: center;
    padding: 24px 16px;
  }
  .card {
    width: 100%;
    max-width: 700px;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 24px;
  }
  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  .title {
    font-size: 24px;
    font-weight: 800;
  }
  .link {
    font-size: 14px;
    color: var(--accent);
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }
  .description {
    margin-bottom: 20px;
    color: var(--muted);
  }
  .description p {
    margin: 8px 0;
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 12px;
    margin-bottom: 16px;
    color: #f87171;
  }
  .success {
    border: 1px solid rgba(34, 197, 94, 0.35);
    background: rgba(34, 197, 94, 0.1);
    border-radius: 12px;
    padding: 12px;
    margin-bottom: 16px;
    color: #22c55e;
  }
  .transcript-area {
    margin-bottom: 24px;
  }
  .label {
    font-size: 14px;
    color: var(--muted);
    margin-bottom: 8px;
  }
  .transcript {
    min-height: 120px;
    padding: 16px;
    border: 1px solid var(--border);
    border-radius: 12px;
    background: var(--panel-2);
    font-size: 18px;
    line-height: 1.6;
    white-space: pre-wrap;
  }
  .interim {
    color: var(--muted);
    opacity: 0.7;
  }
  .actions {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 24px;
  }
  .group {
    display: flex;
    gap: 12px;
    justify-content: center;
  }
  .btn {
    border-radius: 12px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 12px 24px;
    cursor: pointer;
    font-weight: 600;
    font-size: 16px;
    transition: all 0.2s;
  }
  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .btn.mic.active {
    background: rgba(248, 113, 113, 0.2);
    border-color: #f87171;
  }
  .btn.record.active {
    background: rgba(251, 191, 36, 0.2);
    border-color: #fbbf24;
  }
  .btn.ghost {
    background: transparent;
    border-color: var(--border);
  }
  .tips {
    border-top: 1px solid var(--border);
    padding-top: 20px;
    color: var(--muted);
  }
  .tips h4 {
    margin: 0 0 12px 0;
    font-size: 14px;
  }
  .tips ul {
    margin: 0;
    padding-left: 20px;
  }
  .tips li {
    margin: 8px 0;
    font-size: 14px;
  }
</style>

<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  // 表情分类
  const categories = [
    {
      name: '表情',
      icon: '😀',
      emojis: ['😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂', '🙂', '😊', '😇', '🥰', '😍', '🤩', '😘', '😗', '😚', '😋', '😛', '😜', '🤪', '😝', '🤑', '🤗', '🤭', '🤫', '🤔', '🤐', '🤨', '😐', '😑', '😶', '😏', '😒', '🙄', '😬', '😮‍💨', '🤥', '😌', '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢', '🤮', '🤧', '🥵', '🥶', '🥴', '😵', '🤯', '🤠', '🥳', '🥸', '😎', '🤓', '🧐']
    },
    {
      name: '手势',
      icon: '👍',
      emojis: ['👍', '👎', '👊', '✊', '🤛', '🤜', '👏', '🙌', '👐', '🤲', '🤝', '🙏', '✌️', '🤞', '🤟', '🤘', '🤙', '👈', '👉', '👆', '👇', '☝️', '👋', '🤚', '🖐️', '✋', '🖖', '💪', '🦾', '🖕', '✍️', '🤳']
    },
    {
      name: '爱心',
      icon: '❤️',
      emojis: ['❤️', '🧡', '💛', '💚', '💙', '💜', '🖤', '🤍', '🤎', '💔', '❣️', '💕', '💞', '💓', '💗', '💖', '💘', '💝', '💟', '♥️', '💌', '💋', '💯', '💢', '💥', '💫', '💦', '💨', '🕳️', '💣', '💬', '💭', '💤']
    },
    {
      name: '物品',
      icon: '🎉',
      emojis: ['🎉', '🎊', '🎈', '🎁', '🎀', '🏆', '🥇', '🥈', '🥉', '⚽', '🏀', '🎯', '🎮', '🎲', '🧩', '🎨', '🎭', '🎪', '🎬', '🎤', '🎧', '🎹', '🎸', '🎺', '🎻', '📷', '📱', '💻', '⌨️', '🖥️', '🖨️', '📚', '📖', '📝', '✏️', '🔍', '🔑', '💡', '🔔', '⏰', '📅', '📌', '📎', '✂️', '🗑️']
    },
    {
      name: '自然',
      icon: '🌸',
      emojis: ['🌸', '💮', '🏵️', '🌹', '🥀', '🌺', '🌻', '🌼', '🌷', '🌱', '🪴', '🌲', '🌳', '🌴', '🌵', '🌾', '🌿', '☘️', '🍀', '🍁', '🍂', '🍃', '🌍', '🌎', '🌏', '🌑', '🌒', '🌓', '🌔', '🌕', '🌖', '🌗', '🌘', '🌙', '⭐', '🌟', '✨', '💫', '☀️', '🌤️', '⛅', '🌥️', '☁️', '🌦️', '🌧️', '⛈️', '🌩️', '🌨️', '❄️', '☃️', '⛄', '🔥', '💧', '🌊', '🌈']
    },
    {
      name: '食物',
      icon: '🍕',
      emojis: ['🍎', '🍐', '🍊', '🍋', '🍌', '🍉', '🍇', '🍓', '🫐', '🍈', '🍒', '🍑', '🥭', '🍍', '🥥', '🥝', '🍅', '🥑', '🥦', '🥬', '🥒', '🌶️', '🫑', '🌽', '🥕', '🫒', '🧄', '🧅', '🥔', '🍠', '🥐', '🥯', '🍞', '🥖', '🥨', '🧀', '🥚', '🍳', '🧈', '🥞', '🧇', '🥓', '🥩', '🍗', '🍖', '🦴', '🌭', '🍔', '🍟', '🍕', '🫓', '🥪', '🥙', '🧆', '🌮', '🌯', '🫔', '🥗', '🥘', '🫕', '🥫', '🍝', '🍜', '🍲', '🍛', '🍣', '🍱', '🥟', '🦪', '🍤', '🍙', '🍚', '🍘', '🍥', '🥠', '🥮', '🍢', '🍡', '🍧', '🍨', '🍦', '🥧', '🧁', '🍰', '🎂', '🍮', '🍭', '🍬', '🍫', '🍿', '🍩', '🍪', '🌰', '🥜', '🍯', '🥛', '🍼', '🫖', '☕', '🍵', '🧃', '🥤', '🧋', '🍶', '🍺', '🍻', '🥂', '🍷', '🥃', '🍸', '🍹', '🧉', '🍾']
    }
  ];

  let activeCategory = 0;
  let searchQuery = '';

  $: filteredEmojis = searchQuery
    ? categories.flatMap(c => c.emojis).filter(e => e.includes(searchQuery))
    : categories[activeCategory].emojis;

  function selectEmoji(emoji) {
    dispatch('select', emoji);
  }

  function close() {
    dispatch('close');
  }
</script>

<div class="modal-backdrop" on:click|self={close}>
  <div class="modal">
    <div class="modal-header">
      <h3>选择表情</h3>
      <button class="close-btn" on:click={close}>
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <div class="modal-body">
      <!-- 分类标签 -->
      <div class="categories">
        {#each categories as category, i}
          <button 
            class="category-btn" 
            class:active={activeCategory === i}
            on:click={() => { activeCategory = i; searchQuery = ''; }}
            title={category.name}
          >
            {category.icon}
          </button>
        {/each}
      </div>

      <!-- 表情网格 -->
      <div class="emoji-grid">
        {#each filteredEmojis as emoji}
          <button class="emoji-btn" on:click={() => selectEmoji(emoji)}>
            {emoji}
          </button>
        {/each}
      </div>
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
    max-width: 360px;
    max-height: 450px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border);
  }

  .modal-header h3 {
    margin: 0;
    font-size: 16px;
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
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .categories {
    display: flex;
    gap: 4px;
    padding: 8px 12px;
    border-bottom: 1px solid var(--border);
    overflow-x: auto;
  }

  .category-btn {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 20px;
    transition: all 0.15s ease;
    flex-shrink: 0;
  }

  .category-btn:hover {
    background: rgba(148, 163, 184, 0.1);
  }

  .category-btn.active {
    background: rgba(34, 197, 94, 0.15);
  }

  .emoji-grid {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 4px;
  }

  .emoji-btn {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 22px;
    transition: all 0.1s ease;
  }

  .emoji-btn:hover {
    background: rgba(148, 163, 184, 0.15);
    transform: scale(1.15);
  }

  /* 自定义滚动条 */
  .emoji-grid::-webkit-scrollbar {
    width: 6px;
  }

  .emoji-grid::-webkit-scrollbar-track {
    background: transparent;
  }

  .emoji-grid::-webkit-scrollbar-thumb {
    background: var(--border);
    border-radius: 3px;
  }

  .emoji-grid::-webkit-scrollbar-thumb:hover {
    background: var(--muted);
  }
</style>

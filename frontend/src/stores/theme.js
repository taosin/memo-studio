// 简单的 store 实现
let theme = 'light';

// 从 localStorage 读取主题设置
if (typeof window !== 'undefined') {
  theme = localStorage.getItem('theme') || 'light';
  // 应用主题类
  if (theme === 'dark') {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
}

const subscribers = new Set();

export const themeStore = {
  subscribe(fn) {
    fn(theme);
    subscribers.add(fn);
    return () => subscribers.delete(fn);
  },
  set(value) {
    theme = value;
    if (typeof window !== 'undefined') {
      localStorage.setItem('theme', value);
      // 更新 DOM 类
      if (value === 'dark') {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    }
    subscribers.forEach(fn => fn(theme));
  },
  get() {
    return theme;
  }
};

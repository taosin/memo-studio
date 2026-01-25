import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const STORAGE_KEY = 'memo-theme'; // 'light' | 'dark'

function getInitialTheme() {
  if (!browser) return 'dark';
  const saved = localStorage.getItem(STORAGE_KEY);
  if (saved === 'light' || saved === 'dark') return saved;
  return window.matchMedia?.('(prefers-color-scheme: light)')?.matches ? 'light' : 'dark';
}

export const theme = writable(getInitialTheme());

export function applyTheme(t) {
  if (!browser) return;
  document.documentElement.dataset.theme = t;
  localStorage.setItem(STORAGE_KEY, t);
}

export function toggleTheme() {
  theme.update((t) => (t === 'dark' ? 'light' : 'dark'));
}


import { writable } from 'svelte/store';

export const notesStore = writable([]);
export const tagsStore = writable([]);
export const toastStore = writable([]);

let toastId = 0;

export function showToast(message, type = 'info', duration = 3000) {
  const id = ++toastId;
  const toast = { id, message, type };
  
  toastStore.update(toasts => [...toasts, toast]);
  
  if (duration > 0) {
    setTimeout(() => {
      toastStore.update(toasts => toasts.filter(t => t.id !== id));
    }, duration);
  }
  
  return id;
}

export function removeToast(id) {
  toastStore.update(toasts => toasts.filter(t => t.id !== id));
}

export function clearToasts() {
  toastStore.set([]);
}


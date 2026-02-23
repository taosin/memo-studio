import { writable, derived } from 'svelte/store';
import { goto } from '$app/navigation';

export const notesStore = writable([]);
export const tagsStore = writable([]);
export const toastStore = writable([]);

// Auth store
export const authStore = writable({
  token: null,
  user: null,
  isAuthenticated: false
});

// Check if user is authenticated
export function checkAuth() {
  try {
    if (typeof localStorage === 'undefined') return false;
    const token = localStorage.getItem('token');
    const user = JSON.parse(localStorage.getItem('user') || 'null');
    
    if (token && user) {
      authStore.set({
        token,
        user,
        isAuthenticated: true
      });
      return true;
    }
  } catch (e) {
    console.error('Auth check failed:', e);
  }
  
  authStore.set({
    token: null,
    user: null,
    isAuthenticated: false
  });
  return false;
}

// Require authentication, redirect to login if not authenticated
export function requireAuth(currentPath = '/') {
  const isAuth = checkAuth();
  if (!isAuth) {
    // Save the current path to redirect back after login
    try {
      localStorage.setItem('redirectAfterLogin', currentPath);
    } catch {}
    goto('/login');
    return false;
  }
  return true;
}

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


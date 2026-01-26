import test from 'node:test';
import assert from 'node:assert/strict';

function mockWindow() {
  global.window = {};
  const store = new Map();
  global.localStorage = {
    getItem(k) {
      return store.has(k) ? store.get(k) : null;
    },
    setItem(k, v) {
      store.set(k, String(v));
    },
    removeItem(k) {
      store.delete(k);
    }
  };
  return store;
}

test('authStore login/logout updates subscribers and localStorage', async () => {
  const store = mockWindow();
  const { authStore } = await import('./auth.js');

  const seen = [];
  const unsub = authStore.subscribe((v) => seen.push(v));

  authStore.login('tok', { id: 1, username: 'u' });
  assert.equal(store.get('token'), 'tok');
  assert.ok(store.get('user'));

  authStore.logout();
  assert.equal(store.get('token'), undefined);
  assert.equal(store.get('user'), undefined);

  unsub();
// The 'seen' array should have at least 3 values: initial, after login, and after logout
  assert.ok(seen.length >= 3);
  assert.equal(seen[seen.length - 1].isAuthenticated, false);
});


import test from 'node:test';
import assert from 'node:assert/strict';

import { api } from './api.js';

function mockLocalStorage(token = '') {
  global.localStorage = {
    getItem(k) {
      if (k === 'token') return token;
      if (k === 'user') return '{}';
      return null;
    },
    setItem() {},
    removeItem() {}
  };
}

test('api.listNotes throws when not logged in', async () => {
  mockLocalStorage('');
  await assert.rejects(() => api.listNotes(), /未登录/);
});

test('api.listNotes uses /api/memos when logged in', async () => {
  mockLocalStorage('tkn');
  const calls = [];
  global.fetch = async (url, opts) => {
    calls.push({ url: String(url), opts });
    return {
      ok: true,
      json: async () => []
    };
  };
  const res = await api.listNotes();
  assert.deepEqual(res, []);
  assert.ok(calls[0].url.endsWith('/api/memos'));
  assert.equal(calls[0].opts.headers.Authorization, 'Bearer tkn');
});


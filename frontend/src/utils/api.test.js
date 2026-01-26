import test from 'node:test';
import assert from 'node:assert/strict';

// Note: This file is an ES Module (type=module in package.json)
// We use mock window/localStorage/fetch for minimal unit tests

function mockWindow(token = '') {
  global.window = {};
  global.localStorage = {
    getItem(k) {
      if (k === 'token') return token;
      return null;
    },
    setItem() {},
    removeItem() {}
  };
}

test('frontend api attaches Authorization when token exists', async () => {
  mockWindow('tkn');
  const calls = [];
  global.fetch = async (url, opts) => {
    calls.push({ url: String(url), opts });
    return {
      ok: true,
      json: async () => []
    };
  };

  const { api } = await import('./api.js');
  await api.getTags();

  assert.ok(calls.length === 1);
  assert.ok(calls[0].url.endsWith('/api/tags'));
  assert.equal(calls[0].opts.headers.Authorization, 'Bearer tkn');
});


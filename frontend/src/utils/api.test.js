import test from 'node:test';
import assert from 'node:assert/strict';

// Note: This file is an ES Module (type=module in package.json)
// We use mock window/localStorage/fetch for minimal unit tests

function mockWindow(token = '') {
  global.window = {
    dispatchEvent() {}
  };
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
  mockWindow('test-token');
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

  assert.ok(calls.length === 1, 'Should make exactly one API call');
  assert.ok(calls[0].url.endsWith('/api/v1/tags'), `Expected URL to end with /api/v1/tags, got ${calls[0].url}`);
  assert.equal(calls[0].opts.headers.Authorization, 'Bearer test-token', 'Should include Bearer token in Authorization header');
});

test('frontend api works without token', async () => {
  mockWindow('');
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

  assert.ok(calls.length === 1, 'Should make exactly one API call');
  assert.ok(!calls[0].opts.headers.Authorization, 'Should not include Authorization header when no token');
});

test('frontend api constructs correct URLs', async () => {
  mockWindow('token');
  const calls = [];
  global.fetch = async (url, opts) => {
    calls.push({ url: String(url), opts });
    return {
      ok: true,
      json: async () => ({})
    };
  };

  const { api } = await import('./api.js');
  
  // Test various API endpoints
  await api.getTags();
  assert.ok(calls[calls.length - 1].url.includes('/api/v1/tags'), 'getTags should use /api/v1/tags');
  
  await api.getNotes();
  assert.ok(calls[calls.length - 1].url.includes('/api/v1/notes'), 'getNotes should use /api/v1/notes');
});


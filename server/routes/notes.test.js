const test = require('node:test');
const assert = require('node:assert/strict');

const router = require('./notes');

function findHandler(method, path) {
  const layer = router.stack.find((l) => l.route && l.route.path === path && l.route.methods[method]);
  if (!layer) throw new Error(`route not found: ${method.toUpperCase()} ${path}`);
  return layer.route.stack[0].handle;
}

function mockRes() {
  return {
    statusCode: 200,
    body: null,
    status(code) {
      this.statusCode = code;
      return this;
    },
    json(obj) {
      this.body = obj;
      return this;
    }
  };
}

test('POST / returns 400 when content missing', async () => {
  const handler = findHandler('post', '/');
  const req = { body: { title: 't', content: '' } };
  const res = mockRes();
  handler(req, res);
  assert.equal(res.statusCode, 400);
});

test('GET /search returns 400 when keyword missing', async () => {
  const handler = findHandler('get', '/search');
  const req = { query: { keyword: '' } };
  const res = mockRes();
  handler(req, res);
  assert.equal(res.statusCode, 400);
});


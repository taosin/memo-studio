const test = require('node:test');
const assert = require('node:assert/strict');

const router = require('./auth');

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

test('POST /register returns 400 when missing fields', async () => {
  const handler = findHandler('post', '/register');
  const req = { body: { username: '', password: '' } };
  const res = mockRes();
  handler(req, res);
  assert.equal(res.statusCode, 400);
});

test('POST /login returns 400 when missing fields', async () => {
  const handler = findHandler('post', '/login');
  const req = { body: { username: '', password: '' } };
  const res = mockRes();
  handler(req, res);
  assert.equal(res.statusCode, 400);
});


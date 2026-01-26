import test from 'node:test';
import assert from 'node:assert/strict';

function mockDom() {
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
  const classes = new Set();
  global.document = {
    documentElement: {
      classList: {
        add(c) {
          classes.add(c);
        },
        remove(c) {
          classes.delete(c);
        },
        contains(c) {
          return classes.has(c);
        }
      }
    }
  };
  return { store, classes };
}

test('themeStore set() updates localStorage and toggles dark class', async () => {
  const { store, classes } = mockDom();
  const { themeStore } = await import('./theme.js');

  themeStore.set('dark');
  assert.equal(store.get('theme'), 'dark');
  assert.equal(classes.has('dark'), true);

  themeStore.set('light');
  assert.equal(store.get('theme'), 'light');
  assert.equal(classes.has('dark'), false);
});


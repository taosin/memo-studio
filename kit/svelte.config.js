import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  compilerOptions: {
    css: "injected",
    warningFilter: (warning) => {
      if (warning.code && warning.code.startsWith('a11y_')) return false;
      return true;
    }
  },
  kit: {
    adapter: adapter({
      // 产出纯静态站点，并提供 SPA fallback，方便 Go 侧 NoRoute 回退到 index.html
      fallback: 'index.html'
    })
  }
};

export default config;


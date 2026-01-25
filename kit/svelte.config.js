import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      // 产出纯静态站点，并提供 SPA fallback，方便 Go 侧 NoRoute 回退到 index.html
      fallback: 'index.html'
    })
  }
};

export default config;


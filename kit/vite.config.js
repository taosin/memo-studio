import adapter from '@sveltejs/adapter-static';
import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
  plugins: [sveltekit({
      compilerOptions: {
        css: "injected",
        warningFilter: (warning) => {
          if (warning.code && warning.code.startsWith('a11y_')) return false;
          return true;
        }
      },
      adapter: adapter({ fallback: 'index.html' })
  })],
  server: {
    port: 9001,
    proxy: {
      '/api': {
        target: 'http://localhost:9000',
        changeOrigin: true
      }
    }
  }
});

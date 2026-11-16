# memo-studio (SvelteKit)

## 开发

```bash
cd kit
npm install
npm run dev
```

默认端口：`9001`，并通过代理把 `/api` 转发到 `http://localhost:9000`。

## 构建（静态）

```bash
cd kit
npm run build
```

构建产物在 `kit/build/`，用于被 Go 服务托管/内嵌。


构建环境：使用本目录 .nvmrc 指定的 Node 版本，在 kit/ 内运行 npm ci、npm test 和 npm run build。Vite 8 与 vite-plugin-svelte 7 应配套更新。

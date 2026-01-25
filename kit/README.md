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


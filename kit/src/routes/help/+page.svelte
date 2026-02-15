<script>
  const APP_VERSION = "1.0.0";

  const shortcuts = [
    { keys: "Ctrl / Cmd + Enter", desc: "保存当前笔记" },
    { keys: "顶栏主题按钮", desc: "切换深色/浅色主题" },
    { keys: "#标签名", desc: "在笔记内容中使用，可自动关联标签" },
    { keys: "导出页", desc: "支持 JSON 与 Markdown 格式导出备份" },
    { keys: "导入页", desc: "粘贴 JSON 数组批量导入笔记" },
    { keys: "笔记本", desc: "将笔记归类到不同笔记本，便于按主题管理" },
    { keys: "统计页", desc: "查看笔记、标签、资源、笔记本数量" },
    { keys: "随机回顾", desc: "在首页或标签下可进行随机回顾" },
  ];
</script>

<svelte:head>
  <title>帮助 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">帮助</div>
      <a class="link" href="/">返回首页</a>
    </div>
    <p class="hint">常用操作与快捷键说明。</p>

    <section class="section">
      <h3 class="sectionTitle">快捷键与操作</h3>
      <table class="table">
        <thead>
          <tr>
            <th>操作</th>
            <th>说明</th>
          </tr>
        </thead>
        <tbody>
          {#each shortcuts as s}
            <tr>
              <td class="keys">{s.keys}</td>
              <td class="desc">{s.desc}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </section>

    <section class="section">
      <h3 class="sectionTitle">功能概览</h3>
      <ul class="list">
        <li><strong>首页</strong>：写笔记、搜索、按标签筛选、随机回顾。</li>
        <li><strong>笔记本</strong>：创建笔记本，将笔记归类；同一笔记可属于多个笔记本。</li>
        <li><strong>标签库</strong>：管理标签、合并标签。</li>
        <li><strong>资源库</strong>：上传附件，在笔记中引用。</li>
        <li><strong>统计</strong>：查看笔记、标签、资源、笔记本、置顶数量。</li>
        <li><strong>导出</strong>：JSON 或 Markdown 格式备份。</li>
        <li><strong>导入</strong>：JSON 数组批量导入笔记。</li>
        <li><strong>设置</strong>：主题、每页条数、默认视图（本机存储）。</li>
        <li><strong>个人信息</strong>：修改用户名、邮箱、密码。</li>
      </ul>
    </section>

    <section class="section">
      <h3 class="sectionTitle">常见问题</h3>
      <ul class="list">
        <li><strong>如何备份数据？</strong> 使用「导出」页面，选择 JSON 或 Markdown 格式下载。JSON 含完整数据，Markdown 便于阅读与迁移。</li>
        <li><strong>如何从别处迁移笔记？</strong> 在「导入」页面粘贴 JSON 数组，格式为 <code>[{'{'}"title":"标题","content":"内容","tags":["标签"]{'}'}]</code>，点击「开始导入」即可。</li>
        <li><strong>笔记本和标签有什么区别？</strong> 标签用于分类与筛选（如 #工作 #学习）；笔记本用于按主题归类，一篇笔记可以同时属于多个笔记本。</li>
        <li><strong>置顶笔记在哪里？</strong> 首页列表会优先显示置顶笔记；统计页可查看置顶数量。在笔记详情中可设置置顶。</li>
        <li><strong>随机回顾怎么用？</strong> 在首页或标签页可使用随机回顾功能，系统会随机抽取笔记供你复习，支持按标签、按天数筛选。</li>
      </ul>
    </section>

    <section class="section">
      <h3 class="sectionTitle">关于</h3>
      <p class="muted">
        Memo Studio 版本：<strong>{APP_VERSION}</strong>。极简记录，专注书写与回顾。
      </p>
    </section>
  </div>
</div>

<style>
  .wrap {
    padding: 24px 16px;
    max-width: 640px;
    margin: 0 auto;
  }
  .card {
    width: 100%;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 16px;
  }
  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    margin-bottom: 12px;
  }
  .title {
    font-weight: 800;
  }
  .link {
    font-size: 12px;
    text-decoration: underline;
    color: inherit;
  }
  .hint {
    font-size: 13px;
    color: var(--muted);
    margin-bottom: 20px;
  }
  .section {
    margin-bottom: 24px;
  }
  .section:last-child {
    margin-bottom: 0;
  }
  .sectionTitle {
    font-size: 14px;
    font-weight: 700;
    margin: 0 0 12px;
    color: var(--text);
  }
  .table {
    width: 100%;
    border-collapse: collapse;
    font-size: 14px;
  }
  .table th,
  .table td {
    border: 1px solid var(--border);
    padding: 10px 12px;
    text-align: left;
  }
  .table th {
    background: var(--panel-2);
    font-weight: 600;
  }
  .keys {
    font-family: ui-monospace, monospace;
    white-space: nowrap;
  }
  .desc {
    color: var(--muted);
  }
  .list {
    margin: 0;
    padding-left: 20px;
    color: var(--text);
    line-height: 1.6;
  }
  .list li {
    margin-bottom: 8px;
  }
  .muted {
    font-size: 13px;
    color: var(--muted);
    margin: 0;
  }
</style>

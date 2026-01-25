const GITHUB_API = 'https://api.github.com';

function required(name) {
  const v = process.env[name];
  if (!v || String(v).trim() === '') throw new Error(`Missing env: ${name}`);
  return String(v);
}

function optional(name, fallback = '') {
  const v = process.env[name];
  if (!v || String(v).trim() === '') return fallback;
  return String(v);
}

function truncate(s, maxChars) {
  s = String(s ?? '');
  if (s.length <= maxChars) return s;
  return s.slice(0, maxChars) + `\n\n[... 已截断，原始长度 ${s.length} chars ...]\n`;
}

async function ghRequest(method, path, { headers = {}, body } = {}) {
  const token = required('GITHUB_TOKEN');
  const res = await fetch(`${GITHUB_API}${path}`, {
    method,
    headers: {
      Authorization: `Bearer ${token}`,
      'X-GitHub-Api-Version': '2022-11-28',
      ...headers
    },
    body: body ? JSON.stringify(body) : undefined
  });

  if (!res.ok) {
    const text = await res.text().catch(() => '');
    throw new Error(`GitHub API ${method} ${path} failed: ${res.status} ${res.statusText}\n${text}`);
  }
  return res;
}

async function getPrDiff(owner, repo, prNumber) {
  const res = await ghRequest('GET', `/repos/${owner}/${repo}/pulls/${prNumber}`, {
    headers: { Accept: 'application/vnd.github.v3.diff' }
  });
  return await res.text();
}

async function listIssueComments(owner, repo, prNumber) {
  const res = await ghRequest('GET', `/repos/${owner}/${repo}/issues/${prNumber}/comments?per_page=100`);
  return await res.json();
}

async function createComment(owner, repo, prNumber, body) {
  await ghRequest('POST', `/repos/${owner}/${repo}/issues/${prNumber}/comments`, { body: { body } });
}

async function updateComment(owner, repo, commentId, body) {
  await ghRequest('PATCH', `/repos/${owner}/${repo}/issues/comments/${commentId}`, { body: { body } });
}

async function callModel(prompt) {
  const apiKey = required('AI_REVIEW_API_KEY');
  const model = required('AI_REVIEW_MODEL');
  const baseUrl = optional('AI_REVIEW_BASE_URL', 'https://api.openai.com/v1').replace(/\/+$/, '');

  const res = await fetch(`${baseUrl}/chat/completions`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${apiKey}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      model,
      temperature: 0.2,
      messages: [
        {
          role: 'system',
          content:
            '你是一个严谨的代码审查员（CR）。请用中文输出，关注：正确性、边界条件、安全、可维护性、性能、可测试性。' +
            '输出必须是 Markdown，并包含：## 总结、## 主要风险（按严重程度排序）、## 具体建议（可定位到文件/函数/行范围）、## 建议测试计划。' +
            '如果 PR 很大，请先给出“高优先级改动清单”。不要输出与代码无关的泛泛而谈。'
        },
        { role: 'user', content: prompt }
      ]
    })
  });

  if (!res.ok) {
    const text = await res.text().catch(() => '');
    throw new Error(`AI API failed: ${res.status} ${res.statusText}\n${text}`);
  }
  const data = await res.json();
  const msg = data?.choices?.[0]?.message?.content;
  if (!msg) throw new Error('AI API returned empty content');
  return String(msg);
}

function buildPrompt({ repoFull, prNumber, diff }) {
  const cappedDiff = truncate(diff, 180_000);
  return [
    `仓库：${repoFull}`,
    `PR：#${prNumber}`,
    '',
    '下面是该 PR 的 unified diff（可能被截断）：',
    '```diff',
    cappedDiff,
    '```',
    '',
    '要求：请基于 diff 进行代码审查；如发现潜在漏洞/破坏性改动/迁移风险，请明确指出并给出可执行的修改建议。'
  ].join('\n');
}

function wrapComment(content) {
  const marker = '<!-- AI-PR-REVIEW -->';
  const footer = [
    '',
    '---',
    '_由 AI 自动生成的 CR 建议：仅供参考；如与项目约定冲突，以人工审查为准。_'
  ].join('\n');
  return `${marker}\n${content}\n${footer}`;
}

async function main() {
  const repoFull = required('REPO'); // owner/repo
  const prNumber = Number(required('PR_NUMBER'));
  const [owner, repo] = repoFull.split('/');
  if (!owner || !repo) throw new Error(`Invalid REPO: ${repoFull}`);

  // 如果未配置 AI Key/Model，则直接跳过（不让工作流失败，方便外部贡献者 PR）
  const hasKey = !!optional('AI_REVIEW_API_KEY');
  const hasModel = !!optional('AI_REVIEW_MODEL');
  if (!hasKey || !hasModel) {
    console.log('AI_REVIEW_API_KEY / AI_REVIEW_MODEL 未配置，跳过 AI Review。');
    return;
  }

  const diff = await getPrDiff(owner, repo, prNumber);
  const prompt = buildPrompt({ repoFull, prNumber, diff });
  const review = await callModel(prompt);
  const body = wrapComment(review);

  const comments = await listIssueComments(owner, repo, prNumber);
  const existing = comments
    .filter((c) => typeof c?.body === 'string' && c.body.includes('<!-- AI-PR-REVIEW -->'))
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))[0];

  if (existing?.id) {
    await updateComment(owner, repo, existing.id, body);
    console.log(`Updated AI review comment: ${existing.id}`);
  } else {
    await createComment(owner, repo, prNumber, body);
    console.log('Created AI review comment');
  }
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});


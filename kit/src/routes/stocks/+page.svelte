<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let stockCode = '';
  let stockName = '';
  let loading = false;
  let error = '';
  let stockInfo = null;
  let stockHistory = [];
  let searchResults = [];
  let hotStocks = [];

  const API_BASE = 'http://localhost:9000/api';

  onMount(async () => {
    await loadHotStocks();
  });

  // 获取热门股票
  async function loadHotStocks() {
    try {
      const res = await fetch(`${API_BASE}/stocks/hot`);
      const data = await res.json();
      hotStocks = data.stocks || [];
    } catch (e) {
      // 使用默认
      hotStocks = [
        { code: '000001', name: '平安银行', market: '深圳' },
        { code: '600519', name: '贵州茅台', market: '上海' },
        { code: '600036', name: '招商银行', market: '上海' },
        { code: '000002', name: '万 科Ａ', market: '深圳' },
        { code: '601398', name: '工商银行', market: '上海' },
      ];
    }
  }

  async function searchStock() {
    if (!stockCode && !stockName) return;

    loading = true;
    error = '';
    searchResults = [];

    try {
      const res = await fetch(`${API_BASE}/stocks/search?q=${encodeURIComponent(stockCode || stockName)}`);
      const data = await res.json();
      searchResults = data.results || [];
    } catch (e) {
      error = '搜索失败: ' + e.message;
    } finally {
      loading = false;
    }
  }

  async function getStockDetail(code) {
    loading = true;
    error = '';
    stockInfo = null;

    try {
      const res = await fetch(`${API_BASE}/stocks/${code}`);
      if (!res.ok) throw new Error('获取失败');
      const data = await res.json();
      stockInfo = data.stock;

      // 如果有股票信息，自动分析
      if (stockInfo) {
        await analyzeStock(code);
      }
    } catch (e) {
      // 使用模拟数据
      stockInfo = {
        code: code,
        name: '未知股票',
        market: '深圳',
        price: 100.00,
        change: 1.50,
        changePercent: 1.52,
        volume: 5000000,
        pe: 20.0,
        high: 102.00,
        low: 98.50,
        open: 99.00,
        preClose: 98.50,
      };
      error = '无法获取实时数据，显示模拟数据';
    } finally {
      loading = false;
    }
  }

  async function analyzeStock(code) {
    try {
      const res = await fetch(`${API_BASE}/stocks/analyze`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code: code })
      });
      if (res.ok) {
        const data = await res.json();
        if (data.analysis) {
          stockInfo.analysis = data.analysis;
        }
      }
    } catch (e) {
      // 本地生成分析
      stockInfo.analysis = generateLocalAnalysis(stockInfo);
    }
  }

  function generateLocalAnalysis(stock) {
    const isUp = stock.change > 0;
    const trend = isUp ? '上涨' : '下跌';

    const signals = [];
    signals.push({ icon: isUp ? '📈' : '📉', text: isUp ? '价格上涨，技术面偏强' : '价格波动，关注支撑位' });

    if (stock.volume > 10000000) {
      signals.push({ icon: '📊', text: `成交量 ${(stock.volume/10000).toFixed(0)} 万手` });
    }

    if (stock.pe > 0) {
      if (stock.pe > 50) {
        signals.push({ icon: '💰', text: `市盈率 ${stock.pe.toFixed(1)} 倍，估值偏高` });
      } else if (stock.pe < 20) {
        signals.push({ icon: '💵', text: `市盈率 ${stock.pe.toFixed(1)} 倍，估值合理` });
      }
    }

    let suggestion = '建议观望';
    if (isUp && stock.change < 3) {
      suggestion = '可持有，关注上方压力位';
    } else if (!isUp && stock.change > -3) {
      suggestion = '可适当补仓，设置止损';
    } else if (isUp) {
      suggestion = '涨幅较大，建议减仓';
    }

    return {
      summary: `${stock.name}（${stock.code}）今日${trend}${Math.abs(stock.change).toFixed(2)}元（${Math.abs(stock.changePercent).toFixed(2)}%），当前价格¥${stock.price.toFixed(2)}`,
      signals,
      suggestion,
      risks: ['市场整体回调风险', '行业政策变化影响', '公司业绩不及预期'],
      tips: ['分散投资，不要满仓', '设置止损位', '关注基本面', '保持长期心态']
    };
  }

  function saveToNotes() {
    if (!stockInfo) return;

    const analysis = stockInfo.analysis || {};
    const content = `📊 股票分析 - ${stockInfo.name}（${stockInfo.code}）

💰 当前价格: ¥${stockInfo.price?.toFixed(2) || 'N/A'}
📈 涨跌幅: ${stockInfo.change > 0 ? '+' : ''}${stockInfo.change?.toFixed(2) || 0} (${stockInfo.changePercent?.toFixed(2) || 0}%)

📋 基本信息:
- 市盈率: ${stockInfo.pe || 'N/A'}
- 成交量: ${stockInfo.volume ? (stockInfo.volume/10000).toFixed(0) + '万手' : 'N/A'}
- 最高/最低: ¥${stockInfo.high?.toFixed(2) || 'N/A'} / ¥${stockInfo.low?.toFixed(2) || 'N/A'}

${analysis.summary ? `🧠 分析结论:
${analysis.summary}

💡 建议: ${analysis.suggestion || '观望'}` : ''}

---
分析时间: ${new Date().toLocaleString('zh-CN')}
`;

    localStorage.setItem('stock_analysis_draft', content);
    goto('/?new=true');
  }

  function goBack() {
    goto('/');
  }
</script>

<svelte:head>
  <title>股票分析 - Memo Studio</title>
</svelte:head>

<div class="wrap">
  <div class="card">
    <div class="row">
      <div class="title">📈 股票分析</div>
      <button class="link" on:click={goBack}>返回首页</button>
    </div>

    <div class="description">
      查询股票信息，进行技术分析和风险管理。
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <div class="search-box">
      <input
        type="text"
        bind:value={stockCode}
        placeholder="输入股票代码（如 600519）"
        on:keydown={(e) => e.key === 'Enter' && searchStock()}
      />
      <button class="btn" on:click={searchStock} disabled={loading}>
        {loading ? '搜索中...' : '🔍 搜索'}
      </button>
    </div>

    {#if searchResults.length > 0}
      <div class="section">
        <h3>🔎 搜索结果</h3>
        <div class="stock-list">
          {#each searchResults as stock}
            <button class="stock-item" on:click={() => getStockDetail(stock.code)}>
              <div class="stock-info">
                <span class="code">{stock.code}</span>
                <span class="name">{stock.name}</span>
                <span class="market">{stock.market}</span>
              </div>
              <span class="view-btn">查看 →</span>
            </button>
          {/each}
        </div>
      </div>
    {/if}

    {#if !stockInfo}
      <div class="section">
        <h3>🔥 热门股票</h3>
        <div class="stock-list">
          {#each hotStocks as stock}
            <button class="stock-item" on:click={() => getStockDetail(stock.code)}>
              <div class="stock-info">
                <span class="code">{stock.code}</span>
                <span class="name">{stock.name}</span>
                <span class="market">{stock.market}</span>
              </div>
              <span class="view-btn">查看 →</span>
            </button>
          {/each}
        </div>
      </div>
    {/if}

    {#if stockInfo}
      <div class="section stock-detail">
        <div class="detail-header">
          <div>
            <h3>{stockInfo.name}（{stockInfo.code}）</h3>
            <span class="update-time">市场: {stockInfo.market || '未知'}</span>
          </div>
          <div class="current-price" class:up={stockInfo.change > 0} class:down={stockInfo.change < 0}>
            <span class="price">¥{stockInfo.price?.toFixed(2) || 'N/A'}</span>
            <span class="change">
              {stockInfo.change > 0 ? '+' : ''}{stockInfo.change?.toFixed(2) || 0}
              ({stockInfo.changePercent?.toFixed(2) || 0}%)
            </span>
          </div>
        </div>

        <div class="price-info">
          <div class="price-item">
            <span class="label">今开</span>
            <span class="value">{stockInfo.open?.toFixed(2) || '-'}</span>
          </div>
          <div class="price-item">
            <span class="label">昨收</span>
            <span class="value">{stockInfo.preClose?.toFixed(2) || '-'}</span>
          </div>
          <div class="price-item">
            <span class="label">最高</span>
            <span class="value">{stockInfo.high?.toFixed(2) || '-'}</span>
          </div>
          <div class="price-item">
            <span class="label">最低</span>
            <span class="value">{stockInfo.low?.toFixed(2) || '-'}</span>
          </div>
          <div class="price-item">
            <span class="label">成交量</span>
            <span class="value">{stockInfo.volume ? (stockInfo.volume/10000).toFixed(0) + '万' : '-'}</span>
          </div>
          <div class="price-item">
            <span class="label">市盈率</span>
            <span class="value">{stockInfo.pe?.toFixed(1) || '-'}</span>
          </div>
        </div>

        <div class="actions">
          <button class="btn" on:click={() => analyzeStock(stockInfo.code)} disabled={loading}>
            🧠 AI 分析
          </button>
          <button class="btn ghost" on:click={saveToNotes}>
            💾 保存分析
          </button>
          <button class="btn ghost" on:click={() => { stockInfo = null; searchResults = []; }}>
            🔙 返回
          </button>
        </div>

        {#if stockInfo.analysis}
          <div class="analysis-result">
            <h4>🧠 分析结论</h4>
            <div class="summary">{stockInfo.analysis.summary}</div>

            {#if stockInfo.analysis.signals}
              <div class="signals">
                {#each stockInfo.analysis.signals as signal}
                  <div class="signal-item">
                    <span class="icon">{signal.icon}</span>
                    <span class="text">{signal.text}</span>
                  </div>
                {/each}
              </div>
            {/if}

            <div class="suggestion">
              <span class="label">💡 建议</span>
              <span class="text">{stockInfo.analysis.suggestion || '观望'}</span>
            </div>

            {#if stockInfo.analysis.risks}
              <div class="risks">
                <h5>⚠️ 风险提示</h5>
                <ul>
                  {#each stockInfo.analysis.risks as risk}
                    <li>{risk}</li>
                  {/each}
                </ul>
              </div>
            {/if}

            {#if stockInfo.analysis.tips}
              <div class="tips">
                <h5>📚 投资小贴士</h5>
                <ul>
                  {#each stockInfo.analysis.tips as tip}
                    <li>{tip}</li>
                  {/each}
                </ul>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    {/if}

    <div class="disclaimer">
      <h4>⚠️ 风险提示</h4>
      <ul>
        <li>股市有风险，投资需谨慎</li>
        <li>本页面仅供学习参考，不构成投资建议</li>
        <li>请根据自身风险承受能力理性投资</li>
      </ul>
    </div>
  </div>
</div>

<style>
  .wrap {
    display: flex;
    justify-content: center;
    padding: 24px 16px;
  }
  .card {
    width: 100%;
    max-width: 800px;
    border: 1px solid var(--border);
    background: var(--panel);
    border-radius: 14px;
    padding: 24px;
  }
  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }
  .title {
    font-size: 24px;
    font-weight: 800;
  }
  .link {
    font-size: 14px;
    color: var(--accent);
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }
  .description {
    color: var(--muted);
    margin-bottom: 20px;
  }
  .error {
    border: 1px solid rgba(248, 113, 113, 0.35);
    background: rgba(248, 113, 113, 0.1);
    border-radius: 12px;
    padding: 12px;
    margin-bottom: 16px;
    color: #f87171;
  }
  .search-box {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
  }
  .search-box input {
    flex: 1;
    padding: 12px 16px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--panel-2);
    color: inherit;
    font-size: 14px;
  }
  .btn {
    border-radius: 10px;
    border: 1px solid rgba(34, 197, 94, 0.55);
    background: var(--accent-soft);
    color: inherit;
    padding: 12px 24px;
    cursor: pointer;
    font-weight: 600;
  }
  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .btn.ghost {
    background: transparent;
  }
  .section {
    margin-bottom: 24px;
  }
  .section h3 {
    margin: 0 0 12px 0;
    font-size: 16px;
    font-weight: 600;
  }
  .stock-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .stock-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 14px 16px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--panel-2);
    cursor: pointer;
    text-align: left;
    transition: all 0.2s;
  }
  .stock-item:hover {
    border-color: var(--accent);
  }
  .stock-info .code {
    font-weight: 600;
    margin-right: 10px;
  }
  .stock-info .name {
    margin-right: 10px;
  }
  .stock-info .market {
    font-size: 12px;
    color: var(--muted);
  }
  .view-btn {
    font-size: 13px;
    color: var(--accent);
  }
  .stock-detail {
    background: var(--panel-2);
    border-radius: 12px;
    padding: 20px;
  }
  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;
  }
  .detail-header h3 {
    margin: 0 0 8px 0;
    font-size: 20px;
  }
  .update-time {
    font-size: 12px;
    color: var(--muted);
  }
  .current-price {
    text-align: right;
  }
  .current-price .price {
    font-size: 28px;
    font-weight: 800;
    display: block;
  }
  .current-price .change {
    font-size: 14px;
  }
  .current-price.up .price { color: #22c55e; }
  .current-price.up .change { color: #22c55e; }
  .current-price.down .price { color: #f87171; }
  .current-price.down .change { color: #f87171; }
  .price-info {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin-bottom: 20px;
    padding: 16px;
    background: var(--panel);
    border-radius: 10px;
  }
  .price-item {
    text-align: center;
  }
  .price-item .label {
    font-size: 12px;
    color: var(--muted);
    display: block;
  }
  .price-item .value {
    font-weight: 600;
    font-size: 14px;
  }
  .actions {
    display: flex;
    gap: 12px;
    justify-content: center;
    margin-top: 20px;
  }
  .analysis-result {
    margin-top: 24px;
    padding-top: 20px;
    border-top: 1px solid var(--border);
  }
  .analysis-result h4 {
    margin: 0 0 12px 0;
  }
  .summary {
    font-size: 16px;
    line-height: 1.6;
    margin-bottom: 16px;
    padding: 12px;
    background: var(--panel);
    border-radius: 8px;
  }
  .signals {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
  }
  .signal-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    background: var(--panel);
    border-radius: 8px;
  }
  .signal-item .icon {
    font-size: 18px;
  }
  .suggestion {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px;
    background: rgba(34, 197, 94, 0.1);
    border-radius: 8px;
    margin-bottom: 16px;
  }
  .suggestion .label {
    font-weight: 600;
  }
  .risks, .tips {
    margin-bottom: 16px;
  }
  .risks h5, .tips h5 {
    margin: 0 0 8px 0;
    font-size: 14px;
  }
  .risks ul, .tips ul {
    margin: 0;
    padding-left: 20px;
  }
  .risks li {
    color: #f87171;
    margin: 4px 0;
  }
  .tips li {
    color: var(--muted);
    margin: 4px 0;
  }
  .disclaimer {
    margin-top: 24px;
    padding-top: 20px;
    border-top: 1px solid var(--border);
    font-size: 13px;
    color: var(--muted);
  }
  .disclaimer h4 {
    margin: 0 0 8px 0;
    font-size: 14px;
  }
  .disclaimer ul {
    margin: 0;
    padding-left: 20px;
  }
  .disclaimer li {
    margin: 4px 0;
  }
</style>

export function buildHeatmap(notes, days = 90) {
  const now = new Date();
  const start = new Date(now);
  start.setDate(start.getDate() - (days - 1));
  start.setHours(0, 0, 0, 0);

  const counts = new Map();
  for (const n of Array.isArray(notes) ? notes : []) {
    const d = new Date(n.created_at || n.createdAt || n.updated_at || n.updatedAt || Date.now());
    d.setHours(0, 0, 0, 0);
    const key = d.toISOString().slice(0, 10);
    counts.set(key, (counts.get(key) || 0) + 1);
  }

  const cells = [];
  for (let i = 0; i < days; i++) {
    const d = new Date(start);
    d.setDate(start.getDate() + i);
    const key = d.toISOString().slice(0, 10);
    cells.push({
      date: key,
      count: counts.get(key) || 0
    });
  }

  const max = Math.max(0, ...cells.map((c) => c.count));
  return { cells, max };
}

export function heatColor(count, max) {
  if (!max || count <= 0) return 'rgba(148, 163, 184, 0.12)';
  const t = Math.min(1, count / max);
  // 绿色渐变（与暗色背景更搭）
  const a = 0.18 + 0.72 * t;
  return `rgba(34, 197, 94, ${a.toFixed(3)})`;
}


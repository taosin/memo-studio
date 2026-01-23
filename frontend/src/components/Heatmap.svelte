<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';

  let notes = [];
  let heatmapData = {};

  onMount(async () => {
    await loadNotes();
    generateHeatmapData();
  });

  async function loadNotes() {
    try {
      notes = await api.getNotes();
      generateHeatmapData();
    } catch (err) {
      console.error('加载笔记失败:', err);
    }
  }

  function generateHeatmapData() {
    const data = {};
    const today = new Date();
    const oneYearAgo = new Date(today);
    oneYearAgo.setFullYear(today.getFullYear() - 1);

    // 初始化过去一年的日期
    for (let d = new Date(oneYearAgo); d <= today; d.setDate(d.getDate() + 1)) {
      const dateStr = formatDate(d);
      data[dateStr] = 0;
    }

    // 统计每天的笔记数量
    notes.forEach(note => {
      const dateStr = formatDate(new Date(note.created_at));
      if (data[dateStr] !== undefined) {
        data[dateStr]++;
      }
    });

    heatmapData = data;
  }

  function formatDate(date) {
    return date.toISOString().split('T')[0];
  }

  function getIntensity(count) {
    if (count === 0) return 'bg-muted';
    if (count === 1) return 'bg-primary/20';
    if (count <= 3) return 'bg-primary/40';
    if (count <= 5) return 'bg-primary/60';
    return 'bg-primary';
  }

  function getWeekData() {
    const weeks = [];
    const today = new Date();
    const oneYearAgo = new Date(today);
    oneYearAgo.setFullYear(today.getFullYear() - 1);
    
    // 对齐到一年前的周一
    const startDate = new Date(oneYearAgo);
    const dayOfWeek = startDate.getDay();
    const diff = startDate.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1);
    startDate.setDate(diff);

    let currentDate = new Date(startDate);
    let week = [];
    
    while (currentDate <= today) {
      const dateStr = formatDate(currentDate);
      week.push({
        date: dateStr,
        count: heatmapData[dateStr] || 0,
        isToday: dateStr === formatDate(today)
      });

      if (week.length === 7) {
        weeks.push(week);
        week = [];
      }

      currentDate.setDate(currentDate.getDate() + 1);
    }

    if (week.length > 0) {
      weeks.push(week);
    }

    return weeks;
  }

  $: weeks = getWeekData();
  $: maxCount = Math.max(...Object.values(heatmapData), 1);
</script>

<div class="w-full">
  <div class="mb-4">
    <h3 class="text-lg font-semibold mb-2">笔记热力图</h3>
    <p class="text-sm text-muted-foreground">过去一年的笔记记录</p>
  </div>

  <div class="overflow-x-auto -mx-4 px-4">
    <div class="flex gap-1 min-w-max">
      <!-- 星期标签 -->
      <div class="flex flex-col gap-1 mr-2">
        <div class="h-4"></div>
        {#each ['一', '三', '五', '日'] as day}
          <div class="h-3 text-xs text-muted-foreground text-center">{day}</div>
        {/each}
      </div>

      <!-- 热力图 -->
      <div class="flex gap-1">
        {#each weeks as week}
          <div class="flex flex-col gap-1">
            {#each week as day}
              <div
                class="w-3 h-3 rounded-sm {getIntensity(day.count)} {day.isToday ? 'ring-2 ring-primary ring-offset-1' : ''}"
                title="{day.date}: {day.count} 条笔记"
              ></div>
            {/each}
          </div>
        {/each}
      </div>
    </div>
  </div>

  <!-- 图例 -->
  <div class="flex items-center gap-4 mt-4 text-xs text-muted-foreground">
    <span>较少</span>
    <div class="flex gap-1">
      <div class="w-3 h-3 rounded-sm bg-muted"></div>
      <div class="w-3 h-3 rounded-sm bg-primary/20"></div>
      <div class="w-3 h-3 rounded-sm bg-primary/40"></div>
      <div class="w-3 h-3 rounded-sm bg-primary/60"></div>
      <div class="w-3 h-3 rounded-sm bg-primary"></div>
    </div>
    <span>较多</span>
  </div>
</div>

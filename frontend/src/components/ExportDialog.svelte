<script>
  import { createEventDispatcher } from 'svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/ui/card/card-title.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import { exportNotes } from '../utils/export.js';
  import { api } from '../utils/api.js';

  export let visible = false;
  export let selectedNotes = [];

  const dispatch = createEventDispatcher();

  let loading = false;
  let exportFormat = 'markdown';
  let exportAll = true;

  async function handleExport() {
    loading = true;
    try {
      let notesToExport = selectedNotes;
      
      if (exportAll) {
        notesToExport = await api.getNotes();
      }

      if (notesToExport.length === 0) {
        alert('没有可导出的笔记');
        return;
      }

      await exportNotes(notesToExport, exportFormat);
      dispatch('exported');
      dispatch('close');
    } catch (err) {
      alert('导出失败: ' + err.message);
    } finally {
      loading = false;
    }
  }
</script>

{#if visible}
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" on:click={() => dispatch('close')}>
    <Card class="w-full max-w-md" on:click|stopPropagation>
      <CardHeader>
        <CardTitle>导出笔记</CardTitle>
      </CardHeader>
      <CardContent class="p-3 space-y-4">
        <div>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              type="radio"
              bind:group={exportAll}
              value={true}
            />
            <span>导出所有笔记</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer mt-2">
            <input
              type="radio"
              bind:group={exportAll}
              value={false}
            />
            <span>导出选中的笔记 ({selectedNotes.length})</span>
          </label>
        </div>

        <div>
          <label class="block text-sm font-medium mb-2">导出格式</label>
          <div class="space-y-2">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" bind:group={exportFormat} value="markdown" />
              <span>Markdown (.md)</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" bind:group={exportFormat} value="json" />
              <span>JSON (.json)</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" bind:group={exportFormat} value="csv" />
              <span>CSV (.csv)</span>
            </label>
          </div>
        </div>

        <div class="flex gap-2">
          <Button on:click={handleExport} disabled={loading} class="flex-1">
            {loading ? '导出中...' : '导出'}
          </Button>
          <Button variant="outline" on:click={() => dispatch('close')}>
            取消
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
{/if}

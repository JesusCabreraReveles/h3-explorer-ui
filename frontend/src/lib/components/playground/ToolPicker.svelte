<script lang="ts">
	import { selectTool, tool, type Tool } from '$lib/stores/playground';

	const tools: { id: Tool; label: string; hint: string }[] = [
		{ id: 'gridDisk', label: 'Grid disk', hint: 'gridDisk(k) / gridRing(k)' },
		{ id: 'gridPath', label: 'Grid path', hint: 'gridPathCells' },
		{ id: 'polygon', label: 'Polygon', hint: 'polygonToCells' }
	];
</script>

<div class="grid grid-cols-3 gap-1.5">
	{#each tools as t (t.id)}
		<button
			type="button"
			onclick={() => selectTool(t.id)}
			title={t.hint}
			class="rounded-md border px-2 py-2 text-center text-xs font-medium transition
				{$tool === t.id
				? 'border-accent/50 bg-accent/10 text-slate-100'
				: 'border-edge bg-surface-800 text-slate-400 hover:border-edge-strong'}"
		>
			{t.label}
		</button>
	{/each}
</div>

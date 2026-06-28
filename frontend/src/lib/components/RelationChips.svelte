<script lang="ts">
	import { selectCell } from '$lib/stores/explorer';

	interface Props {
		label: string;
		indexes: string[];
		/** Optional cap; extra items are summarized as "+N more". */
		limit?: number;
	}

	let { label, indexes, limit = 12 }: Props = $props();

	const shown = $derived(indexes.slice(0, limit));
	const overflow = $derived(Math.max(0, indexes.length - limit));
</script>

{#if indexes.length > 0}
	<div class="space-y-1.5">
		<div class="flex items-center justify-between">
			<span class="text-[11px] font-medium text-slate-400">{label}</span>
			<span class="font-mono text-[11px] text-slate-600">{indexes.length}</span>
		</div>
		<div class="flex flex-wrap gap-1">
			{#each shown as index (index)}
				<button
					type="button"
					onclick={() => selectCell(index)}
					title={`Inspect ${index}`}
					class="rounded border border-edge bg-surface-800 px-1.5 py-0.5 font-mono text-[11px] text-slate-300 transition hover:border-accent hover:text-accent"
				>
					{index}
				</button>
			{/each}
			{#if overflow > 0}
				<span class="px-1.5 py-0.5 text-[11px] text-slate-600">+{overflow} more</span>
			{/if}
		</div>
	</div>
{/if}

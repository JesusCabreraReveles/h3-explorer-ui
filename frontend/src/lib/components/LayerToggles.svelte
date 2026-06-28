<script lang="ts">
	import { layers, selectedCell, setLayer, type LayerState } from '$lib/stores/explorer';

	const items: { key: keyof LayerState; label: string; swatch: string }[] = [
		{ key: 'neighbors', label: 'Neighbors', swatch: 'bg-emerald-400' },
		{ key: 'children', label: 'Children', swatch: 'bg-violet-400' },
		{ key: 'parent', label: 'Parent', swatch: 'bg-amber-400' }
	];

	const disabled = $derived($selectedCell === null);
</script>

<div class="flex flex-wrap gap-1.5">
	{#each items as item (item.key)}
		<button
			type="button"
			{disabled}
			aria-pressed={$layers[item.key]}
			onclick={() => setLayer(item.key, !$layers[item.key])}
			class="flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-[11px] transition disabled:opacity-40
				{$layers[item.key]
				? 'border-accent/50 bg-accent/10 text-slate-100'
				: 'border-edge bg-surface-800 text-slate-400 hover:border-edge-strong'}"
		>
			<span class="h-2 w-2 rounded-full {item.swatch} {$layers[item.key] ? '' : 'opacity-40'}"
			></span>
			{item.label}
		</button>
	{/each}
</div>

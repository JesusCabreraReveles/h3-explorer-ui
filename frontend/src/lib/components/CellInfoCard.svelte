<script lang="ts">
	import { selectedCell } from '$lib/stores/explorer';
	import { formatArea, formatLatLng, formatLength } from '$lib/utils/format';
	import RelationChips from './RelationChips.svelte';

	let copied = $state(false);
	let showBoundary = $state(false);

	async function copyIndex(index: string): Promise<void> {
		try {
			await navigator.clipboard.writeText(index);
			copied = true;
			setTimeout(() => (copied = false), 1200);
		} catch {
			// Clipboard may be unavailable (insecure context); ignore silently.
		}
	}
</script>

{#if $selectedCell}
	{@const cell = $selectedCell}
	<div class="space-y-4">
		<div>
			<div class="mb-1 flex items-center justify-between">
				<span class="text-xs font-medium text-slate-400">H3 Index</span>
				<button
					type="button"
					onclick={() => copyIndex(cell.index)}
					class="text-[11px] text-slate-500 transition hover:text-accent"
				>
					{copied ? 'copied ✓' : 'copy'}
				</button>
			</div>
			<code
				class="block break-all rounded-md border border-edge bg-surface-950 px-2.5 py-1.5 font-mono text-sm text-accent"
			>
				{cell.index}
			</code>
		</div>

		<dl
			class="grid grid-cols-2 gap-px overflow-hidden rounded-md border border-edge bg-edge text-sm"
		>
			{#snippet stat(label: string, value: string)}
				<div class="bg-surface-850 px-2.5 py-2">
					<dt class="text-[11px] text-slate-500">{label}</dt>
					<dd class="mt-0.5 font-mono text-slate-200">{value}</dd>
				</div>
			{/snippet}

			{@render stat('Resolution', String(cell.resolution))}
			{@render stat('Base cell', String(cell.baseCell))}
			{@render stat('Area', formatArea(cell.areaKm2))}
			{@render stat('Edge length', formatLength(cell.edgeLengthKm))}
			{@render stat('Center', formatLatLng(cell.center.lat, cell.center.lng))}
			{@render stat('Icosahedron face', cell.icosahedronFaces.join(', '))}
		</dl>

		<div class="flex flex-wrap gap-1.5">
			{#if cell.isPentagon}
				<span
					class="rounded-full border border-amber-500/40 bg-amber-500/10 px-2 py-0.5 text-[11px] text-amber-300"
				>
					pentagon
				</span>
			{/if}
			<span
				class="rounded-full border border-edge bg-surface-800 px-2 py-0.5 text-[11px] text-slate-400"
			>
				{cell.isClassIII ? 'Class III' : 'Class II'}
			</span>
		</div>

		<!-- Topology navigation: click any related index to inspect it. -->
		{#if cell.parent}
			<RelationChips label="Parent" indexes={[cell.parent]} />
		{/if}
		<RelationChips label="Neighbors" indexes={cell.neighbors} />
		<RelationChips label={`Children (res ${cell.resolution + 1})`} indexes={cell.children} />

		<!-- Boundary vertices. -->
		<div>
			<button
				type="button"
				onclick={() => (showBoundary = !showBoundary)}
				class="flex w-full items-center justify-between text-[11px] font-medium text-slate-400 hover:text-slate-200"
			>
				<span>Boundary coordinates</span>
				<span class="font-mono text-slate-600">
					{cell.boundary.length} · {showBoundary ? '▾' : '▸'}
				</span>
			</button>
			{#if showBoundary}
				<ol
					class="mt-1.5 space-y-px overflow-hidden rounded-md border border-edge font-mono text-[11px]"
				>
					{#each cell.boundary as vertex, i (i)}
						<li class="flex gap-2 bg-surface-850 px-2.5 py-1 text-slate-300">
							<span class="w-4 text-slate-600">{i}</span>
							<span>{formatLatLng(vertex.lat, vertex.lng, 6)}</span>
						</li>
					{/each}
				</ol>
			{/if}
		</div>
	</div>
{:else}
	<p class="text-sm text-slate-500">
		Click anywhere on the map, or enter coordinates above, to index a point into its H3 cell.
	</p>
{/if}

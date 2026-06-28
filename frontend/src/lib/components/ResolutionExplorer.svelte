<script lang="ts">
	import { onMount } from 'svelte';
	import { changeResolution, resolution } from '$lib/stores/explorer';
	import { loadResolutions, resolutions, resolutionsError } from '$lib/stores/resolutions';
	import { formatArea, formatCount, formatLength } from '$lib/utils/format';

	onMount(() => {
		void loadResolutions();
	});
</script>

<div class="space-y-2">
	<p class="text-[11px] text-slate-500">
		Every H3 resolution, from continent-scale (0) to ~1&nbsp;m² (15). Select one to re-index the
		active point instantly.
	</p>

	{#if $resolutionsError}
		<p class="rounded-md border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-xs text-rose-300">
			{$resolutionsError}
		</p>
	{/if}

	<div class="overflow-hidden rounded-md border border-edge">
		<table class="w-full border-collapse text-left text-xs">
			<thead class="bg-surface-850 text-[10px] uppercase tracking-wider text-slate-500">
				<tr>
					<th class="px-2 py-1.5 font-medium">Res</th>
					<th class="px-2 py-1.5 font-medium">Avg area</th>
					<th class="px-2 py-1.5 font-medium">Edge</th>
					<th class="px-2 py-1.5 text-right font-medium">Total cells</th>
				</tr>
			</thead>
			<tbody>
				{#each $resolutions as row (row.resolution)}
					<tr
						onclick={() => changeResolution(row.resolution)}
						class="cursor-pointer border-t border-edge/60 transition
							{row.resolution === $resolution
							? 'bg-accent/10 text-slate-100'
							: 'text-slate-300 hover:bg-surface-800'}"
					>
						<td class="px-2 py-1.5 font-mono font-semibold">
							{#if row.resolution === $resolution}
								<span class="text-accent">{row.resolution}</span>
							{:else}
								{row.resolution}
							{/if}
						</td>
						<td class="px-2 py-1.5 font-mono">{formatArea(row.avgAreaKm2)}</td>
						<td class="px-2 py-1.5 font-mono">{formatLength(row.avgEdgeLengthKm)}</td>
						<td class="px-2 py-1.5 text-right font-mono text-slate-400">
							{formatCount(row.totalCells)}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>

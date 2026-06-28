<script lang="ts">
	import ToolPicker from './ToolPicker.svelte';
	import { resolve } from '$app/paths';
	import { MAX_RESOLUTION, MIN_RESOLUTION } from '$lib/config';
	import {
		animateDisk,
		cells,
		closePolygon,
		count,
		draftPolygon,
		errorMessage,
		k,
		loading,
		PLAYGROUND_MAX_K,
		polygonClosed,
		reset,
		resolution,
		ringOnly,
		setK,
		setResolution,
		showBoundaries,
		showLabels,
		toggleRing,
		tool
	} from '$lib/stores/playground';

	const canClose = $derived($draftPolygon.length >= 3 && !$polygonClosed);
</script>

<aside
	class="flex h-full w-full flex-col gap-5 overflow-y-auto border-r border-edge bg-surface-900 p-5 md:w-[380px]"
>
	<header class="flex items-center justify-between">
		<div class="flex items-center gap-2.5">
			<div class="grid h-9 w-9 place-items-center rounded-lg bg-accent/15 text-lg">🧪</div>
			<div>
				<h1 class="text-sm font-semibold tracking-tight text-slate-100">H3 Playground</h1>
				<p class="text-[11px] text-slate-500">Experiment with H3 algorithms</p>
			</div>
		</div>
		<a
			href={resolve('/')}
			class="rounded-md border border-edge bg-surface-800 px-2.5 py-1 text-[11px] text-slate-400 transition hover:border-accent hover:text-accent"
		>
			← Explorer
		</a>
	</header>

	<section class="space-y-2">
		<h2 class="text-[11px] font-semibold uppercase tracking-wider text-slate-500">Tool</h2>
		<ToolPicker />
	</section>

	<!-- Resolution (applies to every tool). -->
	<section class="space-y-2">
		<div class="flex items-center justify-between">
			<span class="text-xs font-medium text-slate-400">Resolution</span>
			<span class="font-mono text-sm font-semibold text-accent">{$resolution}</span>
		</div>
		<input
			type="range"
			min={MIN_RESOLUTION}
			max={MAX_RESOLUTION}
			value={$resolution}
			oninput={(e) => setResolution(Number(e.currentTarget.value))}
			class="h-1.5 w-full cursor-pointer appearance-none rounded-full bg-surface-700 accent-accent"
		/>
	</section>

	<!-- Tool-specific controls. -->
	{#if $tool === 'gridDisk'}
		<section class="space-y-3">
			<div class="flex items-center justify-between">
				<span class="text-xs font-medium text-slate-400">k = {$k}</span>
				<label class="flex items-center gap-1.5 text-[11px] text-slate-400">
					<input
						type="checkbox"
						checked={$ringOnly}
						onchange={(e) => toggleRing(e.currentTarget.checked)}
						class="accent-accent"
					/>
					hollow ring
				</label>
			</div>
			<input
				type="range"
				min="0"
				max={PLAYGROUND_MAX_K}
				value={$k}
				oninput={(e) => setK(Number(e.currentTarget.value))}
				class="h-1.5 w-full cursor-pointer appearance-none rounded-full bg-surface-700 accent-accent"
			/>
			<button
				type="button"
				onclick={() => animateDisk()}
				class="w-full rounded-md border border-edge bg-surface-800 px-3 py-1.5 text-xs text-slate-300 transition hover:border-accent hover:text-accent"
			>
				▶ Animate growth
			</button>
			<p class="text-[11px] text-slate-500">Click the map to set the origin cell.</p>
		</section>
	{:else if $tool === 'gridPath'}
		<section class="space-y-2">
			<p class="text-[11px] text-slate-500">
				Click an origin cell, then a destination at the same resolution. Far-apart or
				cross-resolution cells have no path.
			</p>
		</section>
	{:else}
		<section class="space-y-3">
			<p class="text-[11px] text-slate-500">
				Click the map to drop polygon vertices ({$draftPolygon.length}), then close it to fill with
				cells.
			</p>
			<button
				type="button"
				disabled={!canClose}
				onclick={() => closePolygon()}
				class="w-full rounded-md bg-accent px-3 py-1.5 text-xs font-semibold text-surface-950 transition hover:brightness-110 disabled:opacity-40"
			>
				Close polygon → cells
			</button>
		</section>
	{/if}

	<!-- Display options. -->
	<section class="space-y-2">
		<h2 class="text-[11px] font-semibold uppercase tracking-wider text-slate-500">Display</h2>
		<div class="flex flex-wrap gap-3 text-xs text-slate-400">
			<label class="flex items-center gap-1.5">
				<input type="checkbox" bind:checked={$showBoundaries} class="accent-accent" />
				Boundaries
			</label>
			<label class="flex items-center gap-1.5">
				<input type="checkbox" bind:checked={$showLabels} class="accent-accent" />
				H3 index labels
			</label>
		</div>
	</section>

	<!-- Status / result. -->
	<section class="space-y-2">
		<div
			class="flex items-center justify-between rounded-md border border-edge bg-surface-850 px-3 py-2"
		>
			<span class="text-xs text-slate-400">Cells</span>
			<span class="font-mono text-sm font-semibold text-slate-100">
				{$loading ? '…' : $count.toLocaleString('en-US')}
			</span>
		</div>
		{#if $errorMessage}
			<p
				class="rounded-md border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-xs text-rose-300"
			>
				{$errorMessage}
			</p>
		{/if}
		{#if $cells.length > 0 || $draftPolygon.length > 0}
			<button
				type="button"
				onclick={() => reset()}
				class="w-full rounded-md border border-edge bg-surface-800 px-3 py-1.5 text-xs text-slate-400 transition hover:border-rose-500/50 hover:text-rose-300"
			>
				Clear
			</button>
		{/if}
	</section>

	<footer class="mt-auto border-t border-edge pt-3 text-[11px] text-slate-600">
		Phase 5 · H3 Playground
	</footer>
</aside>

<script lang="ts">
	import CoordinateSearch from './CoordinateSearch.svelte';
	import IndexSearch from './IndexSearch.svelte';
	import ResolutionSlider from './ResolutionSlider.svelte';
	import LayerToggles from './LayerToggles.svelte';
	import CellInfoCard from './CellInfoCard.svelte';
	import ResolutionExplorer from './ResolutionExplorer.svelte';
	import KeyboardHints from './KeyboardHints.svelte';
	import { resolve } from '$app/paths';
	import { errorMessage } from '$lib/stores/explorer';

	type Tab = 'inspector' | 'resolutions';
	let tab = $state<Tab>('inspector');

	const tabs: { id: Tab; label: string }[] = [
		{ id: 'inspector', label: 'Inspector' },
		{ id: 'resolutions', label: 'Resolutions' }
	];
</script>

<aside
	class="flex h-full w-full flex-col gap-5 overflow-y-auto border-r border-edge bg-surface-900 p-5 md:w-[380px]"
>
	<header class="flex items-center justify-between">
		<div class="flex items-center gap-2.5">
			<div class="grid h-9 w-9 place-items-center rounded-lg bg-accent/15 text-lg">🗺️</div>
			<div>
				<h1 class="text-sm font-semibold tracking-tight text-slate-100">H3 Explorer UI</h1>
				<p class="text-[11px] text-slate-500">Uber H3 interactive playground</p>
			</div>
		</div>
		<a
			href={resolve('/playground')}
			class="rounded-md border border-edge bg-surface-800 px-2.5 py-1 text-[11px] text-slate-400 transition hover:border-accent hover:text-accent"
		>
			Playground →
		</a>
	</header>

	<section class="space-y-3">
		<h2 class="text-[11px] font-semibold uppercase tracking-wider text-slate-500">Search</h2>
		<CoordinateSearch />
		<IndexSearch />
	</section>

	<section class="space-y-3">
		<ResolutionSlider />
	</section>

	{#if $errorMessage}
		<p class="rounded-md border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-xs text-rose-300">
			{$errorMessage}
		</p>
	{/if}

	<!-- Tabbed panel: cell inspector vs. resolution explorer. -->
	<div class="flex flex-col gap-3">
		<div class="flex gap-1 rounded-lg bg-surface-850 p-1">
			{#each tabs as t (t.id)}
				<button
					type="button"
					onclick={() => (tab = t.id)}
					class="flex-1 rounded-md px-3 py-1.5 text-xs font-medium transition
						{tab === t.id ? 'bg-surface-700 text-slate-100' : 'text-slate-500 hover:text-slate-300'}"
				>
					{t.label}
				</button>
			{/each}
		</div>

		{#if tab === 'inspector'}
			<LayerToggles />
			<CellInfoCard />
		{:else}
			<ResolutionExplorer />
		{/if}
	</div>

	<KeyboardHints
		hints={[
			{ keys: ['+', '−'], label: 'Resolution' },
			{ keys: ['N'], label: 'Toggle neighbors' },
			{ keys: ['C'], label: 'Toggle children' },
			{ keys: ['P'], label: 'Toggle parent' },
			{ keys: ['Esc'], label: 'Clear selection' }
		]}
	/>

	<footer class="mt-auto border-t border-edge pt-3 text-[11px] text-slate-600">
		Inspector & Resolution Explorer
	</footer>
</aside>

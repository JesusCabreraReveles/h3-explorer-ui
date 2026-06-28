<script lang="ts">
	import { DEFAULT_CENTER } from '$lib/config';
	import { indexPoint, loading } from '$lib/stores/explorer';

	// Seed the form with the default map center ([lng, lat]).
	let lat = $state(DEFAULT_CENTER[1]);
	let lng = $state(DEFAULT_CENTER[0]);
	let geoError = $state<string | null>(null);

	function submit(event: SubmitEvent): void {
		event.preventDefault();
		void indexPoint({ lat, lng });
	}

	function useMyLocation(): void {
		geoError = null;
		if (!navigator.geolocation) {
			geoError = 'Geolocation is not available in this browser.';
			return;
		}
		navigator.geolocation.getCurrentPosition(
			(pos) => {
				lat = pos.coords.latitude;
				lng = pos.coords.longitude;
				void indexPoint({ lat, lng });
			},
			() => {
				geoError = 'Could not get your location.';
			}
		);
	}
</script>

<form class="space-y-3" onsubmit={submit}>
	<div class="grid grid-cols-2 gap-3">
		<label class="block">
			<span class="mb-1 block text-xs font-medium text-slate-400">Latitude</span>
			<input
				type="number"
				step="any"
				min="-90"
				max="90"
				bind:value={lat}
				class="w-full rounded-md border border-edge bg-surface-800 px-2.5 py-1.5 font-mono text-sm text-slate-100 outline-none focus:border-accent"
			/>
		</label>
		<label class="block">
			<span class="mb-1 block text-xs font-medium text-slate-400">Longitude</span>
			<input
				type="number"
				step="any"
				min="-180"
				max="180"
				bind:value={lng}
				class="w-full rounded-md border border-edge bg-surface-800 px-2.5 py-1.5 font-mono text-sm text-slate-100 outline-none focus:border-accent"
			/>
		</label>
	</div>

	<div class="flex gap-2">
		<button
			type="submit"
			disabled={$loading}
			class="flex-1 rounded-md bg-accent px-3 py-1.5 text-sm font-semibold text-surface-950 transition hover:brightness-110 disabled:opacity-50"
		>
			{$loading ? 'Indexing…' : 'Index coordinate'}
		</button>
		<button
			type="button"
			onclick={useMyLocation}
			title="Use my location"
			class="rounded-md border border-edge bg-surface-800 px-3 py-1.5 text-sm text-slate-300 transition hover:border-accent hover:text-accent"
		>
			◎
		</button>
	</div>

	{#if geoError}
		<p class="text-xs text-rose-400">{geoError}</p>
	{/if}
</form>

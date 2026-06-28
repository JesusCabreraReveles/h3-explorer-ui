<script lang="ts">
	import { MAX_RESOLUTION, MIN_RESOLUTION } from '$lib/config';
	import { changeResolution, resolution } from '$lib/stores/explorer';

	function onSlide(event: Event): void {
		const value = Number((event.currentTarget as HTMLInputElement).value);
		void changeResolution(value);
	}

	function step(delta: number): void {
		const next = Math.min(MAX_RESOLUTION, Math.max(MIN_RESOLUTION, $resolution + delta));
		if (next !== $resolution) {
			void changeResolution(next);
		}
	}
</script>

<div class="space-y-2">
	<div class="flex items-center justify-between">
		<span class="text-xs font-medium text-slate-400">Resolution</span>
		<span class="font-mono text-sm font-semibold text-accent">{$resolution}</span>
	</div>

	<div class="flex items-center gap-2">
		<button
			type="button"
			onclick={() => step(-1)}
			aria-label="Decrease resolution"
			class="h-7 w-7 shrink-0 rounded-md border border-edge bg-surface-800 text-slate-300 transition hover:border-accent hover:text-accent"
		>
			−
		</button>
		<input
			type="range"
			min={MIN_RESOLUTION}
			max={MAX_RESOLUTION}
			value={$resolution}
			oninput={onSlide}
			class="h-1.5 w-full cursor-pointer appearance-none rounded-full bg-surface-700 accent-accent"
		/>
		<button
			type="button"
			onclick={() => step(1)}
			aria-label="Increase resolution"
			class="h-7 w-7 shrink-0 rounded-md border border-edge bg-surface-800 text-slate-300 transition hover:border-accent hover:text-accent"
		>
			+
		</button>
	</div>

	<div class="flex justify-between px-0.5 text-[10px] text-slate-600">
		<span>{MIN_RESOLUTION}</span>
		<span>{MAX_RESOLUTION}</span>
	</div>
</div>

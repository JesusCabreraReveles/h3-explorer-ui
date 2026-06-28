<script lang="ts">
	import { selectCell, loading } from '$lib/stores/explorer';
	import { isH3IndexFormat } from '$lib/utils/h3';

	let value = $state('');
	let touched = $state(false);

	const valid = $derived(isH3IndexFormat(value));

	function submit(event: SubmitEvent): void {
		event.preventDefault();
		touched = true;
		if (valid) {
			void selectCell(value.trim().toLowerCase());
		}
	}
</script>

<form class="space-y-1.5" onsubmit={submit}>
	<span class="block text-xs font-medium text-slate-400">H3 index</span>
	<div class="flex gap-2">
		<input
			type="text"
			placeholder="8928308280fffff"
			spellcheck="false"
			autocomplete="off"
			bind:value
			oninput={() => (touched = false)}
			class="w-full rounded-md border border-edge bg-surface-800 px-2.5 py-1.5 font-mono text-sm text-slate-100 outline-none focus:border-accent"
		/>
		<button
			type="submit"
			disabled={$loading || !valid}
			class="rounded-md border border-edge bg-surface-800 px-3 py-1.5 text-sm text-slate-300 transition hover:border-accent hover:text-accent disabled:opacity-40 disabled:hover:border-edge disabled:hover:text-slate-300"
		>
			Inspect
		</button>
	</div>
	{#if touched && value && !valid}
		<p class="text-[11px] text-rose-400">Expected 15–16 hexadecimal characters.</p>
	{/if}
</form>

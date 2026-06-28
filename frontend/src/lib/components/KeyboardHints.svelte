<script lang="ts">
	interface Hint {
		keys: string[];
		label: string;
	}

	let { hints }: { hints: Hint[] } = $props();
</script>

<div class="space-y-1.5">
	<span class="text-[11px] font-semibold uppercase tracking-wider text-slate-500">Shortcuts</span>
	<ul class="space-y-1">
		{#each hints as hint (hint.label)}
			<li class="flex items-center justify-between gap-2 text-[11px] text-slate-500">
				<span>{hint.label}</span>
				<span class="flex gap-1">
					{#each hint.keys as key (key)}
						<kbd
							class="rounded border border-edge bg-surface-800 px-1.5 py-0.5 font-mono text-[10px] text-slate-400"
						>
							{key}
						</kbd>
					{/each}
				</span>
			</li>
		{/each}
	</ul>
</div>

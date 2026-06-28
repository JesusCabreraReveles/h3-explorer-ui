<script lang="ts">
	import { copyText, downloadText } from '$lib/utils/download';

	export interface ExportItem {
		label: string;
		content: string;
		filename: string;
		mime: string;
	}

	let { items }: { items: ExportItem[] } = $props();

	let copiedLabel = $state<string | null>(null);

	async function copy(item: ExportItem): Promise<void> {
		if (await copyText(item.content)) {
			copiedLabel = item.label;
			setTimeout(() => (copiedLabel = null), 1200);
		}
	}
</script>

{#if items.length > 0}
	<div class="space-y-1.5">
		{#each items as item (item.label)}
			<div class="flex items-center gap-2">
				<span class="flex-1 text-xs text-slate-400">{item.label}</span>
				<button
					type="button"
					onclick={() => copy(item)}
					aria-label={`Copy ${item.label}`}
					class="rounded border border-edge bg-surface-800 px-2 py-0.5 text-[11px] text-slate-300 transition hover:border-accent hover:text-accent"
				>
					{copiedLabel === item.label ? 'copied ✓' : 'copy'}
				</button>
				<button
					type="button"
					onclick={() => downloadText(item.filename, item.content, item.mime)}
					aria-label={`Download ${item.label}`}
					class="rounded border border-edge bg-surface-800 px-2 py-0.5 text-[11px] text-slate-300 transition hover:border-accent hover:text-accent"
				>
					download
				</button>
			</div>
		{/each}
	</div>
{/if}

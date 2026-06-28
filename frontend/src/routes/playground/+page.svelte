<script lang="ts">
	import PlaygroundMap from '$lib/components/playground/PlaygroundMap.svelte';
	import PlaygroundPanel from '$lib/components/playground/PlaygroundPanel.svelte';
</script>

<svelte:head>
	<title>H3 Playground — H3 Explorer UI</title>
</svelte:head>

<div class="flex h-full w-full flex-col-reverse md:flex-row">
	<PlaygroundPanel />
	<main class="relative h-full min-h-0 flex-1">
		<PlaygroundMap />
	</main>
</div>

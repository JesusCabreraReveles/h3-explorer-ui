<script lang="ts">
	import MapView from '$lib/components/MapView.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
</script>

<svelte:head>
	<title>H3 Explorer UI — interactive playground for Uber's H3</title>
</svelte:head>

<div class="flex h-full w-full flex-col-reverse md:flex-row">
	<Sidebar />
	<main class="relative h-full min-h-0 flex-1">
		<MapView />
	</main>
</div>

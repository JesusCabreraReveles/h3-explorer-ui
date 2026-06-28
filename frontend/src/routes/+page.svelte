<script lang="ts">
	import MapView from '$lib/components/MapView.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import { MAX_RESOLUTION, MIN_RESOLUTION } from '$lib/config';
	import {
		changeResolution,
		clearSelection,
		layers,
		resolution,
		setLayer
	} from '$lib/stores/explorer';
	import { isEditableTarget } from '$lib/utils/keyboard';

	function onKeydown(event: KeyboardEvent): void {
		if (isEditableTarget(event.target)) return;

		switch (event.key) {
			case '+':
			case '=':
				changeResolution(Math.min(MAX_RESOLUTION, $resolution + 1));
				break;
			case '-':
			case '_':
				changeResolution(Math.max(MIN_RESOLUTION, $resolution - 1));
				break;
			case 'n':
			case 'N':
				setLayer('neighbors', !$layers.neighbors);
				break;
			case 'c':
			case 'C':
				setLayer('children', !$layers.children);
				break;
			case 'p':
			case 'P':
				setLayer('parent', !$layers.parent);
				break;
			case 'Escape':
				clearSelection();
				break;
			default:
				return;
		}
		event.preventDefault();
	}
</script>

<svelte:head>
	<title>H3 Explorer UI — interactive playground for Uber's H3</title>
</svelte:head>

<svelte:window onkeydown={onKeydown} />

<div class="flex h-full w-full flex-col-reverse md:flex-row">
	<Sidebar />
	<main class="relative h-full min-h-0 flex-1">
		<MapView />
	</main>
</div>

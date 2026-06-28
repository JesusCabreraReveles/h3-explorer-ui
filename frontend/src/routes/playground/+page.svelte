<script lang="ts">
	import PlaygroundMap from '$lib/components/playground/PlaygroundMap.svelte';
	import PlaygroundPanel from '$lib/components/playground/PlaygroundPanel.svelte';
	import { MAX_RESOLUTION, MIN_RESOLUTION } from '$lib/config';
	import {
		animateDisk,
		closePolygon,
		reset,
		resolution,
		selectTool,
		setResolution,
		showBoundaries,
		showLabels
	} from '$lib/stores/playground';
	import { isEditableTarget } from '$lib/utils/keyboard';

	function onKeydown(event: KeyboardEvent): void {
		if (isEditableTarget(event.target)) return;

		switch (event.key) {
			case '1':
				selectTool('gridDisk');
				break;
			case '2':
				selectTool('gridPath');
				break;
			case '3':
				selectTool('polygon');
				break;
			case '+':
			case '=':
				setResolution(Math.min(MAX_RESOLUTION, $resolution + 1));
				break;
			case '-':
			case '_':
				setResolution(Math.max(MIN_RESOLUTION, $resolution - 1));
				break;
			case 'l':
			case 'L':
				showLabels.update((v) => !v);
				break;
			case 'b':
			case 'B':
				showBoundaries.update((v) => !v);
				break;
			case 'a':
			case 'A':
				animateDisk();
				break;
			case 'Enter':
				closePolygon();
				break;
			case 'Escape':
				reset();
				break;
			default:
				return;
		}
		event.preventDefault();
	}
</script>

<svelte:head>
	<title>H3 Playground — H3 Explorer UI</title>
</svelte:head>

<svelte:window onkeydown={onKeydown} />

<div class="flex h-full w-full flex-col-reverse md:flex-row">
	<PlaygroundPanel />
	<main class="relative h-full min-h-0 flex-1">
		<PlaygroundMap />
	</main>
</div>

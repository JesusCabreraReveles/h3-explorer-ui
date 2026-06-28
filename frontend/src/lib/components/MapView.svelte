<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import maplibregl, { type Map as MlMap, type GeoJSONSource } from 'maplibre-gl';
	import type { FeatureCollection } from 'geojson';
	import { createMap } from '$lib/map/createMap';
	import { indexPoint, overlays, selectCell, selectedCell } from '$lib/stores/explorer';
	import { boundaryToPolygon, cellsToFeatureCollection } from '$lib/utils/geojson';
	import type { CellGeometry, CellInfo } from '$lib/types/h3';

	const CELL_SOURCE = 'selected-cell';
	const NEIGHBORS_SOURCE = 'overlay-neighbors';
	const CHILDREN_SOURCE = 'overlay-children';
	const PARENT_SOURCE = 'overlay-parent';

	const ACCENT = '#38bdf8';
	const NEIGHBOR_COLOR = '#34d399';
	const CHILD_COLOR = '#a78bfa';
	const PARENT_COLOR = '#fbbf24';

	// Layers whose polygons select-on-click, topmost first.
	const CLICKABLE_LAYERS = ['cell-fill', 'children-fill', 'neighbors-fill', 'parent-fill'];

	let container: HTMLDivElement;
	let map: MlMap | undefined;
	let marker: maplibregl.Marker | undefined;
	let ready = $state(false);

	onMount(() => {
		map = createMap(container);

		map.on('load', () => {
			if (!map) return;

			// Overlay sources/layers are added first so the selected cell draws on top.
			addCellSource(PARENT_SOURCE);
			map.addLayer({
				id: 'parent-fill',
				type: 'fill',
				source: PARENT_SOURCE,
				paint: { 'fill-color': PARENT_COLOR, 'fill-opacity': 0.04 }
			});
			map.addLayer({
				id: 'parent-outline',
				type: 'line',
				source: PARENT_SOURCE,
				paint: { 'line-color': PARENT_COLOR, 'line-width': 1.5, 'line-dasharray': [2, 1] }
			});

			addCellSource(CHILDREN_SOURCE);
			map.addLayer({
				id: 'children-fill',
				type: 'fill',
				source: CHILDREN_SOURCE,
				paint: { 'fill-color': CHILD_COLOR, 'fill-opacity': 0.12 }
			});
			map.addLayer({
				id: 'children-outline',
				type: 'line',
				source: CHILDREN_SOURCE,
				paint: { 'line-color': CHILD_COLOR, 'line-width': 0.75 }
			});

			addCellSource(NEIGHBORS_SOURCE);
			map.addLayer({
				id: 'neighbors-fill',
				type: 'fill',
				source: NEIGHBORS_SOURCE,
				paint: { 'fill-color': NEIGHBOR_COLOR, 'fill-opacity': 0.1 }
			});
			map.addLayer({
				id: 'neighbors-outline',
				type: 'line',
				source: NEIGHBORS_SOURCE,
				paint: { 'line-color': NEIGHBOR_COLOR, 'line-width': 1 }
			});

			addCellSource(CELL_SOURCE);
			map.addLayer({
				id: 'cell-fill',
				type: 'fill',
				source: CELL_SOURCE,
				paint: { 'fill-color': ACCENT, 'fill-opacity': 0.2 }
			});
			map.addLayer({
				id: 'cell-outline',
				type: 'line',
				source: CELL_SOURCE,
				paint: { 'line-color': ACCENT, 'line-width': 2.5 }
			});

			// Pointer cursor over clickable cells.
			for (const layer of CLICKABLE_LAYERS) {
				map.on('mouseenter', layer, () => setCursor('pointer'));
				map.on('mouseleave', layer, () => setCursor(''));
			}

			ready = true;
		});

		// A click selects a rendered cell if one is under the cursor, otherwise it
		// indexes the raw coordinate. Centralizing the hit-test here avoids
		// per-layer handlers double-firing.
		map.on('click', (e) => {
			const hit = map?.queryRenderedFeatures(e.point, { layers: presentLayers() });
			const index = hit?.[0]?.properties?.index as string | undefined;
			if (index) {
				void selectCell(index);
			} else {
				void indexPoint({ lat: e.lngLat.lat, lng: e.lngLat.lng });
			}
		});
	});

	onDestroy(() => {
		marker?.remove();
		map?.remove();
	});

	// Redraw the selected cell (and recenter) when the selection changes.
	$effect(() => {
		const cell = $selectedCell;
		if (ready && map) {
			drawSelected(cell);
		}
	});

	// Redraw overlays whenever their geometry changes.
	$effect(() => {
		const data = $overlays;
		if (ready && map) {
			setData(NEIGHBORS_SOURCE, data.neighbors);
			setData(CHILDREN_SOURCE, data.children);
			setData(PARENT_SOURCE, data.parent);
		}
	});

	function addCellSource(id: string): void {
		map?.addSource(id, { type: 'geojson', data: emptyData() });
	}

	function emptyData(): FeatureCollection {
		return { type: 'FeatureCollection', features: [] };
	}

	function setCursor(value: string): void {
		if (map) map.getCanvas().style.cursor = value;
	}

	/** Only query layers that currently exist (post-load). */
	function presentLayers(): string[] {
		return CLICKABLE_LAYERS.filter((id) => map?.getLayer(id));
	}

	function setData(sourceId: string, cells: CellGeometry[]): void {
		const source = map?.getSource(sourceId) as GeoJSONSource | undefined;
		source?.setData(cellsToFeatureCollection(cells));
	}

	function drawSelected(cell: CellInfo | null): void {
		const source = map?.getSource(CELL_SOURCE) as GeoJSONSource | undefined;
		if (!map || !source) return;

		if (!cell) {
			source.setData(emptyData());
			marker?.remove();
			marker = undefined;
			return;
		}

		source.setData({
			type: 'Feature',
			geometry: boundaryToPolygon(cell.boundary),
			properties: { index: cell.index }
		});

		const center: [number, number] = [cell.center.lng, cell.center.lat];
		marker ??= new maplibregl.Marker({ color: ACCENT });
		marker.setLngLat(center).addTo(map);

		const targetZoom = Math.min(16, Math.max(2, cell.resolution + 2));
		map.flyTo({ center, zoom: Math.max(map.getZoom(), targetZoom), speed: 0.8 });
	}
</script>

<div bind:this={container} class="h-full w-full" data-testid="map"></div>

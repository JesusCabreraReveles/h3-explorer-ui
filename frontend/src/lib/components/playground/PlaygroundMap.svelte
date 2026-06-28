<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { type Map as MlMap, type GeoJSONSource } from 'maplibre-gl';
	import type { Feature, FeatureCollection } from 'geojson';
	import { createMap } from '$lib/map/createMap';
	import { boundaryToPolygon, cellsToFeatureCollection } from '$lib/utils/geojson';
	import {
		cells,
		draftPolygon,
		handleClick,
		polygonClosed,
		showBoundaries,
		showLabels
	} from '$lib/stores/playground';
	import type { LatLng } from '$lib/types/h3';

	const RESULT_SOURCE = 'pg-result';
	const DRAFT_SOURCE = 'pg-draft';
	const DRAFT_POINTS_SOURCE = 'pg-draft-points';
	const ACCENT = '#38bdf8';

	let container: HTMLDivElement;
	let map: MlMap | undefined;
	let ready = $state(false);

	onMount(() => {
		map = createMap(container);

		map.on('load', () => {
			if (!map) return;

			map.addSource(RESULT_SOURCE, { type: 'geojson', data: empty() });
			map.addLayer({
				id: 'pg-result-fill',
				type: 'fill',
				source: RESULT_SOURCE,
				paint: { 'fill-color': ACCENT, 'fill-opacity': 0.18 }
			});
			map.addLayer({
				id: 'pg-result-outline',
				type: 'line',
				source: RESULT_SOURCE,
				paint: { 'line-color': ACCENT, 'line-width': 1.25 }
			});
			map.addLayer({
				id: 'pg-result-labels',
				type: 'symbol',
				source: RESULT_SOURCE,
				layout: {
					'text-field': ['get', 'index'],
					'text-size': 9,
					'text-font': ['Open Sans Regular', 'Noto Sans Regular'],
					'text-allow-overlap': false,
					visibility: 'none'
				},
				paint: {
					'text-color': '#e2e8f0',
					'text-halo-color': '#0d1118',
					'text-halo-width': 1.2
				}
			});

			// Draft polygon (vertices + connecting ring) while drawing.
			map.addSource(DRAFT_SOURCE, { type: 'geojson', data: empty() });
			map.addLayer({
				id: 'pg-draft-fill',
				type: 'fill',
				source: DRAFT_SOURCE,
				paint: { 'fill-color': ACCENT, 'fill-opacity': 0.08 }
			});
			map.addLayer({
				id: 'pg-draft-line',
				type: 'line',
				source: DRAFT_SOURCE,
				paint: { 'line-color': ACCENT, 'line-width': 1.5, 'line-dasharray': [2, 1] }
			});

			map.addSource(DRAFT_POINTS_SOURCE, { type: 'geojson', data: empty() });
			map.addLayer({
				id: 'pg-draft-points',
				type: 'circle',
				source: DRAFT_POINTS_SOURCE,
				paint: {
					'circle-radius': 4,
					'circle-color': ACCENT,
					'circle-stroke-color': '#0d1118',
					'circle-stroke-width': 1.5
				}
			});

			ready = true;
		});

		map.on('click', (e) => {
			void handleClick({ lat: e.lngLat.lat, lng: e.lngLat.lng });
		});
	});

	onDestroy(() => map?.remove());

	// Result cells.
	$effect(() => {
		const data = $cells;
		if (ready && map) {
			source(RESULT_SOURCE)?.setData(cellsToFeatureCollection(data));
		}
	});

	// Draft polygon (line/fill + vertices).
	$effect(() => {
		const vertices = $draftPolygon;
		const closed = $polygonClosed;
		if (ready && map) {
			source(DRAFT_SOURCE)?.setData(draftGeometry(vertices, closed));
			source(DRAFT_POINTS_SOURCE)?.setData(pointsGeometry(vertices));
		}
	});

	// Display toggles.
	$effect(() => {
		const labels = $showLabels;
		if (ready && map) {
			map.setLayoutProperty('pg-result-labels', 'visibility', labels ? 'visible' : 'none');
		}
	});
	$effect(() => {
		const boundaries = $showBoundaries;
		if (ready && map) {
			map.setLayoutProperty('pg-result-outline', 'visibility', boundaries ? 'visible' : 'none');
		}
	});

	function source(id: string): GeoJSONSource | undefined {
		return map?.getSource(id) as GeoJSONSource | undefined;
	}

	function empty(): FeatureCollection {
		return { type: 'FeatureCollection', features: [] };
	}

	function draftGeometry(vertices: LatLng[], closed: boolean): Feature | FeatureCollection {
		if (vertices.length === 0) return empty();
		if (closed && vertices.length >= 3) {
			return { type: 'Feature', geometry: boundaryToPolygon(vertices), properties: {} };
		}
		return {
			type: 'Feature',
			geometry: { type: 'LineString', coordinates: vertices.map((v) => [v.lng, v.lat]) },
			properties: {}
		};
	}

	function pointsGeometry(vertices: LatLng[]): FeatureCollection {
		return {
			type: 'FeatureCollection',
			features: vertices.map((v) => ({
				type: 'Feature',
				geometry: { type: 'Point', coordinates: [v.lng, v.lat] },
				properties: {}
			}))
		};
	}
</script>

<div bind:this={container} class="h-full w-full" data-testid="playground-map"></div>

// Shared MapLibre bootstrap so the explorer and the playground configure their
// maps identically (same basemap, controls, defaults) without duplication.

import maplibregl, { type Map as MlMap } from 'maplibre-gl';
import { DEFAULT_CENTER, DEFAULT_ZOOM, MAP_STYLE_URL } from '$lib/config';

export function createMap(container: HTMLElement): MlMap {
	const map = new maplibregl.Map({
		container,
		style: MAP_STYLE_URL,
		center: DEFAULT_CENTER,
		zoom: DEFAULT_ZOOM,
		attributionControl: { compact: true }
	});

	map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'top-right');
	map.addControl(
		new maplibregl.GeolocateControl({ positionOptions: { enableHighAccuracy: true } }),
		'top-right'
	);

	return map;
}

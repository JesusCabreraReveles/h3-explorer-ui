// Conversion helpers between the API's H3 geometry (arrays of {lat,lng}) and
// GeoJSON (arrays of [lng,lat]) for rendering with MapLibre. Kept dependency-free
// and pure so they are easy to unit test.

import type { Boundary, CellGeometry, CellInfo } from '$lib/types/h3';

export interface GeoJSONPolygon {
	type: 'Polygon';
	coordinates: number[][][];
}

export interface GeoJSONFeature<P = Record<string, unknown>> {
	type: 'Feature';
	geometry: GeoJSONPolygon;
	properties: P;
}

export interface GeoJSONFeatureCollection<P = Record<string, unknown>> {
	type: 'FeatureCollection';
	features: GeoJSONFeature<P>[];
}

/** A renderable cell is either a full CellInfo or a lightweight CellGeometry. */
type RenderableCell = Pick<CellGeometry, 'index' | 'boundary'> & Partial<CellInfo>;

/**
 * Converts an H3 boundary ring into a closed GeoJSON Polygon. GeoJSON requires
 * the first and last positions of a linear ring to be identical, which H3
 * boundaries are not, so we close the ring explicitly.
 */
export function boundaryToPolygon(boundary: Boundary): GeoJSONPolygon {
	const ring = boundary.map(({ lat, lng }) => [lng, lat]);
	if (ring.length > 0) {
		const [firstLng, firstLat] = ring[0];
		const [lastLng, lastLat] = ring[ring.length - 1];
		if (firstLng !== lastLng || firstLat !== lastLat) {
			ring.push([firstLng, firstLat]);
		}
	}
	return { type: 'Polygon', coordinates: [ring] };
}

/** Wraps a single cell as a GeoJSON Feature carrying its H3 index. */
export function cellToFeature(cell: RenderableCell): GeoJSONFeature<{ index: string }> {
	return {
		type: 'Feature',
		geometry: boundaryToPolygon(cell.boundary),
		properties: { index: cell.index }
	};
}

/** Builds a FeatureCollection from many cells (grid-disk, children, …). */
export function cellsToFeatureCollection(
	cells: RenderableCell[]
): GeoJSONFeatureCollection<{ index: string }> {
	return {
		type: 'FeatureCollection',
		features: cells.map(cellToFeature)
	};
}

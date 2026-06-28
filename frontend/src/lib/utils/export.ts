// Pure serializers for exporting H3 selections. Kept side-effect-free (no DOM,
// no clipboard) so they are easy to unit test; the browser-only download/copy
// helpers live in utils/download.ts.

import type { Boundary, CellGeometry, CellInfo } from '$lib/types/h3';
import {
	boundaryToPolygon,
	cellsToFeatureCollection,
	type GeoJSONFeatureCollection
} from './geojson';

/** A single renderable cell (full info or lightweight geometry). */
type AnyCell = Pick<CellGeometry, 'index' | 'boundary'>;

/** Serializes many cells as a pretty-printed GeoJSON FeatureCollection. */
export function cellsToGeoJSON(cells: AnyCell[]): string {
	return JSON.stringify(cellsToFeatureCollection(cells), null, 2);
}

/**
 * Serializes a single inspected cell as a GeoJSON FeatureCollection whose one
 * feature carries the cell's key metrics as properties.
 */
export function cellInfoToGeoJSON(cell: CellInfo): string {
	const fc: GeoJSONFeatureCollection = {
		type: 'FeatureCollection',
		features: [
			{
				type: 'Feature',
				geometry: boundaryToPolygon(cell.boundary),
				properties: {
					index: cell.index,
					resolution: cell.resolution,
					areaKm2: cell.areaKm2,
					edgeLengthKm: cell.edgeLengthKm,
					baseCell: cell.baseCell,
					isPentagon: cell.isPentagon
				}
			}
		]
	};
	return JSON.stringify(fc, null, 2);
}

/** One H3 index per line. */
export function indexesToText(cells: { index: string }[]): string {
	return cells.map((c) => c.index).join('\n');
}

/** Boundary vertices as `lat,lng` CSV with a header row. */
export function boundaryToCSV(boundary: Boundary): string {
	const rows = boundary.map((v) => `${v.lat},${v.lng}`);
	return ['lat,lng', ...rows].join('\n');
}

import { describe, expect, it } from 'vitest';
import { boundaryToCSV, cellInfoToGeoJSON, cellsToGeoJSON, indexesToText } from './export';
import type { Boundary, CellInfo } from '$lib/types/h3';

const triangle: Boundary = [
	{ lat: 1, lng: 2 },
	{ lat: 3, lng: 4 },
	{ lat: 5, lng: 6 }
];

const cell: CellInfo = {
	index: '8928308280fffff',
	resolution: 9,
	center: { lat: 1, lng: 2 },
	boundary: triangle,
	areaKm2: 0.1,
	areaM2: 100000,
	edgeLengthKm: 0.2,
	edgeLengthM: 200,
	baseCell: 20,
	icosahedronFaces: [7],
	isPentagon: false,
	isClassIII: true,
	children: [],
	numChildren: 7,
	neighbors: []
};

describe('cellsToGeoJSON', () => {
	it('produces a valid FeatureCollection with one feature per cell', () => {
		const json = cellsToGeoJSON([
			{ index: 'a', boundary: triangle },
			{ index: 'b', boundary: triangle }
		]);
		const parsed = JSON.parse(json);
		expect(parsed.type).toBe('FeatureCollection');
		expect(parsed.features).toHaveLength(2);
		expect(parsed.features[0].properties.index).toBe('a');
	});

	it('is pretty-printed', () => {
		expect(cellsToGeoJSON([{ index: 'a', boundary: triangle }])).toContain('\n');
	});
});

describe('cellInfoToGeoJSON', () => {
	it('carries cell metrics as feature properties', () => {
		const parsed = JSON.parse(cellInfoToGeoJSON(cell));
		const props = parsed.features[0].properties;
		expect(props.index).toBe('8928308280fffff');
		expect(props.resolution).toBe(9);
		expect(props.baseCell).toBe(20);
		expect(props.isPentagon).toBe(false);
	});

	it('closes the polygon ring', () => {
		const ring = JSON.parse(cellInfoToGeoJSON(cell)).features[0].geometry.coordinates[0];
		expect(ring[0]).toEqual(ring[ring.length - 1]);
	});
});

describe('indexesToText', () => {
	it('writes one index per line', () => {
		expect(indexesToText([{ index: 'a' }, { index: 'b' }, { index: 'c' }])).toBe('a\nb\nc');
	});

	it('handles an empty selection', () => {
		expect(indexesToText([])).toBe('');
	});
});

describe('boundaryToCSV', () => {
	it('emits a header and lat,lng rows', () => {
		expect(boundaryToCSV(triangle)).toBe('lat,lng\n1,2\n3,4\n5,6');
	});
});

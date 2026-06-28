import { describe, expect, it } from 'vitest';
import { boundaryToPolygon, cellsToFeatureCollection, cellToFeature } from './geojson';
import type { Boundary } from '$lib/types/h3';

const triangle: Boundary = [
	{ lat: 1, lng: 2 },
	{ lat: 3, lng: 4 },
	{ lat: 5, lng: 6 }
];

describe('boundaryToPolygon', () => {
	it('swaps to [lng, lat] order', () => {
		const poly = boundaryToPolygon(triangle);
		expect(poly.coordinates[0][0]).toEqual([2, 1]);
	});

	it('closes the ring by repeating the first vertex', () => {
		const poly = boundaryToPolygon(triangle);
		const ring = poly.coordinates[0];
		expect(ring).toHaveLength(triangle.length + 1);
		expect(ring[0]).toEqual(ring[ring.length - 1]);
	});

	it('does not double-close an already closed ring', () => {
		const closed: Boundary = [...triangle, triangle[0]];
		const poly = boundaryToPolygon(closed);
		expect(poly.coordinates[0]).toHaveLength(closed.length);
	});

	it('handles an empty boundary', () => {
		expect(boundaryToPolygon([]).coordinates[0]).toEqual([]);
	});
});

describe('cellToFeature / cellsToFeatureCollection', () => {
	it('carries the H3 index in properties', () => {
		const feature = cellToFeature({ index: '8928308280fffff', boundary: triangle });
		expect(feature.type).toBe('Feature');
		expect(feature.properties.index).toBe('8928308280fffff');
	});

	it('builds a FeatureCollection from many cells', () => {
		const fc = cellsToFeatureCollection([
			{ index: 'a', boundary: triangle },
			{ index: 'b', boundary: triangle }
		]);
		expect(fc.type).toBe('FeatureCollection');
		expect(fc.features).toHaveLength(2);
		expect(fc.features.map((f) => f.properties.index)).toEqual(['a', 'b']);
	});
});

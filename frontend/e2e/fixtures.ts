import type { Page } from '@playwright/test';

/** A deterministic CellInfo used to mock inspector responses. */
export const SF_CELL = {
	index: '8928308280fffff',
	resolution: 9,
	center: { lat: 37.7767, lng: -122.4184 },
	boundary: [
		{ lat: 37.7752, lng: -122.4172 },
		{ lat: 37.7769, lng: -122.4161 },
		{ lat: 37.7784, lng: -122.4174 },
		{ lat: 37.7782, lng: -122.4197 },
		{ lat: 37.7765, lng: -122.4208 },
		{ lat: 37.775, lng: -122.4195 }
	],
	areaKm2: 0.109,
	areaM2: 109398,
	edgeLengthKm: 0.2,
	edgeLengthM: 200.8,
	baseCell: 20,
	icosahedronFaces: [7],
	isPentagon: false,
	isClassIII: true,
	parent: '8828308281fffff',
	children: ['8a28308280c7fff', '8a28308280cffff'],
	numChildren: 7,
	neighbors: ['89283082803ffff', '89283082807ffff']
};

/** Mocks the H3 endpoints the explorer/playground call, so e2e needs no backend. */
export async function mockApi(page: Page): Promise<void> {
	await page.route('**/api/h3/inspect', (route) => route.fulfill({ json: SF_CELL }));
	await page.route('**/api/h3/from-coordinates', (route) => route.fulfill({ json: SF_CELL }));
	await page.route('**/api/h3/resolutions', (route) =>
		route.fulfill({
			json: {
				resolutions: Array.from({ length: 16 }, (_, r) => ({
					resolution: r,
					avgAreaKm2: 4357449 / 7 ** r,
					avgEdgeLengthKm: 1281 / 2.6 ** r,
					totalCells: 122 * 7 ** r
				}))
			}
		})
	);
}

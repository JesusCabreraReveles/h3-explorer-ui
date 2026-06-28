import { describe, expect, it, vi } from 'vitest';
import { ApiError, H3Api } from './api';
import type { CellInfo } from '$lib/types/h3';

function jsonResponse(body: unknown, status = 200): Response {
	return new Response(JSON.stringify(body), {
		status,
		headers: { 'content-type': 'application/json' }
	});
}

const sampleCell: Partial<CellInfo> = {
	index: '8928308280fffff',
	resolution: 9,
	baseCell: 20
};

describe('H3Api', () => {
	it('posts to from-coordinates and returns the cell', async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(sampleCell));
		const api = new H3Api('', fetchFn);

		const cell = await api.fromCoordinates(37.77, -122.41, 9);

		expect(cell.index).toBe('8928308280fffff');
		const [url, init] = fetchFn.mock.calls[0];
		expect(url).toBe('/api/h3/from-coordinates');
		expect(init.method).toBe('POST');
		expect(JSON.parse(init.body)).toEqual({ lat: 37.77, lng: -122.41, resolution: 9 });
	});

	it('unwraps the resolutions envelope', async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ resolutions: [{ resolution: 0 }, { resolution: 1 }] }));
		const api = new H3Api('', fetchFn);

		const res = await api.resolutions();

		expect(res).toHaveLength(2);
		expect(res[1].resolution).toBe(1);
	});

	it('throws a typed ApiError carrying the backend code', async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(
				jsonResponse({ error: { code: 'invalid_cell', message: 'bad index' } }, 400)
			);
		const api = new H3Api('', fetchFn);

		await expect(api.inspect('nope')).rejects.toMatchObject({
			name: 'ApiError',
			status: 400,
			code: 'invalid_cell',
			message: 'bad index'
		});
	});

	it('falls back to a generic ApiError on non-JSON failures', async () => {
		const fetchFn = vi.fn().mockResolvedValue(new Response('boom', { status: 500 }));
		const api = new H3Api('', fetchFn);

		const err = await api.gridDisk('x', 1).catch((e) => e);
		expect(err).toBeInstanceOf(ApiError);
		expect(err.status).toBe(500);
		expect(err.code).toBe('http_error');
	});

	it('applies the base URL prefix', async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(sampleCell));
		const api = new H3Api('http://localhost:8080', fetchFn);

		await api.inspect('8928308280fffff');

		expect(fetchFn.mock.calls[0][0]).toBe('http://localhost:8080/api/h3/inspect');
	});
});

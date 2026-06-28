// Typed client for the H3 Explorer backend REST API. This is the single place
// the frontend talks to the network; components and stores depend on these
// methods, never on fetch directly. The fetch implementation is injectable so
// the client is trivially unit-testable and usable from SvelteKit load
// functions.

import { API_BASE } from '$lib/config';
import type {
	BoundaryResponse,
	CellCollection,
	CellGeometry,
	CellInfo,
	MultiPolygon,
	ResolutionInfo,
	ApiErrorBody,
	Boundary
} from '$lib/types/h3';

type FetchFn = typeof fetch;

/** Error thrown for any non-2xx API response, carrying the backend's code. */
export class ApiError extends Error {
	constructor(
		readonly status: number,
		readonly code: string,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

export class H3Api {
	constructor(
		private readonly base: string = API_BASE,
		private readonly fetchFn: FetchFn = fetch
	) {}

	// --- Indexing & inspection ---

	fromCoordinates(lat: number, lng: number, resolution: number): Promise<CellInfo> {
		return this.post<CellInfo>('/api/h3/from-coordinates', { lat, lng, resolution });
	}

	inspect(index: string): Promise<CellInfo> {
		return this.post<CellInfo>('/api/h3/inspect', { index });
	}

	toBoundary(index: string): Promise<BoundaryResponse> {
		return this.post<BoundaryResponse>('/api/h3/to-boundary', { index });
	}

	async resolutions(): Promise<ResolutionInfo[]> {
		const { resolutions } = await this.get<{ resolutions: ResolutionInfo[] }>(
			'/api/h3/resolutions'
		);
		return resolutions;
	}

	// --- Traversal & hierarchy ---

	gridDisk(index: string, k: number): Promise<CellCollection> {
		return this.post<CellCollection>('/api/h3/grid-disk', { index, k });
	}

	gridRing(index: string, k: number): Promise<CellCollection> {
		return this.post<CellCollection>('/api/h3/grid-ring', { index, k });
	}

	gridPath(origin: string, destination: string): Promise<CellCollection> {
		return this.post<CellCollection>('/api/h3/grid-path', { origin, destination });
	}

	parent(index: string, resolution: number): Promise<CellGeometry> {
		return this.post<CellGeometry>('/api/h3/parent', { index, resolution });
	}

	children(index: string, resolution: number): Promise<CellCollection> {
		return this.post<CellCollection>('/api/h3/children', { index, resolution });
	}

	neighbors(index: string): Promise<CellCollection> {
		return this.post<CellCollection>('/api/h3/neighbors', { index });
	}

	// --- Polygon operations ---

	polygonToCells(
		polygon: Boundary,
		resolution: number,
		holes: Boundary[] = []
	): Promise<CellCollection> {
		return this.post<CellCollection>('/api/h3/polygon-to-cells', { polygon, holes, resolution });
	}

	cellsToMultiPolygon(indexes: string[]): Promise<MultiPolygon> {
		return this.post<MultiPolygon>('/api/h3/cells-to-multi-polygon', { indexes });
	}

	// --- Internals ---

	private get<T>(path: string): Promise<T> {
		return this.request<T>(path, { method: 'GET' });
	}

	private post<T>(path: string, body: unknown): Promise<T> {
		return this.request<T>(path, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify(body)
		});
	}

	private async request<T>(path: string, init: RequestInit): Promise<T> {
		const res = await this.fetchFn(`${this.base}${path}`, init);
		if (!res.ok) {
			throw await toApiError(res);
		}
		return (await res.json()) as T;
	}
}

/** Builds an ApiError from a failed response, tolerating non-JSON bodies. */
async function toApiError(res: Response): Promise<ApiError> {
	try {
		const body = (await res.json()) as ApiErrorBody;
		if (body?.error?.code) {
			return new ApiError(res.status, body.error.code, body.error.message);
		}
	} catch {
		// fall through to a generic error below
	}
	return new ApiError(res.status, 'http_error', `request failed with status ${res.status}`);
}

/** Default singleton used throughout the app. */
export const api = new H3Api();

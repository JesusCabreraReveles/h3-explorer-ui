// The explorer feature store: the single source of truth for "which point /
// resolution / cell is selected, and which related cells are overlaid". UI
// components render these stores and call the exported actions; they never call
// the API directly. This keeps orchestration out of the UI layer.

import { get, writable } from 'svelte/store';
import { api, ApiError } from '$lib/services/api';
import { DEFAULT_RESOLUTION, MAX_RESOLUTION, MIN_RESOLUTION } from '$lib/config';
import type { CellGeometry, CellInfo, LatLng } from '$lib/types/h3';

/** Toggleable map overlays relative to the selected cell. */
export interface LayerState {
	neighbors: boolean;
	children: boolean;
	parent: boolean;
}

/** Resolved geometry for each enabled overlay. */
export interface OverlayData {
	neighbors: CellGeometry[];
	children: CellGeometry[];
	parent: CellGeometry[];
}

const NO_OVERLAYS: OverlayData = { neighbors: [], children: [], parent: [] };

/** Currently selected H3 resolution (0–15). */
export const resolution = writable<number>(DEFAULT_RESOLUTION);

/** The cell resulting from the last successful lookup, or null. */
export const selectedCell = writable<CellInfo | null>(null);

/** The coordinate that produced the current selection (drives re-indexing). */
export const activePoint = writable<LatLng | null>(null);

/** True while a primary lookup is in flight. */
export const loading = writable<boolean>(false);

/** Human-readable error from the last failed request, or null. */
export const errorMessage = writable<string | null>(null);

/** Which related-cell overlays are enabled. */
export const layers = writable<LayerState>({ neighbors: false, children: false, parent: false });

/** Geometry for the currently enabled overlays. */
export const overlays = writable<OverlayData>(NO_OVERLAYS);

/**
 * Indexes a coordinate at the current resolution and selects the resulting
 * cell. Triggered by the search form and by clicking empty map.
 */
export async function indexPoint(point: LatLng): Promise<void> {
	activePoint.set(point);
	await runLookup(point, get(resolution));
}

/**
 * Selects an existing cell by its H3 index (inspector navigation, index search,
 * clicking an overlay polygon). The active point and resolution follow the cell
 * so subsequent resolution changes re-index from its center.
 */
export async function selectCell(index: string): Promise<void> {
	loading.set(true);
	errorMessage.set(null);
	try {
		const cell = await api.inspect(index);
		resolution.set(cell.resolution);
		activePoint.set(cell.center);
		selectedCell.set(cell);
		await refreshOverlays(cell);
	} catch (err) {
		errorMessage.set(toMessage(err));
	} finally {
		loading.set(false);
	}
}

/**
 * Changes the active resolution and, if a point is selected, re-indexes it so
 * the map updates immediately.
 */
export async function changeResolution(next: number): Promise<void> {
	resolution.set(next);
	const point = get(activePoint);
	if (point) {
		await runLookup(point, next);
	}
}

/** Enables or disables an overlay layer and refreshes its geometry. */
export async function setLayer(name: keyof LayerState, on: boolean): Promise<void> {
	layers.update((current) => ({ ...current, [name]: on }));
	const cell = get(selectedCell);
	if (cell) {
		await refreshOverlays(cell);
	} else {
		overlays.set(NO_OVERLAYS);
	}
}

/** Clears the current selection and overlays. */
export function clearSelection(): void {
	selectedCell.set(null);
	activePoint.set(null);
	overlays.set(NO_OVERLAYS);
	errorMessage.set(null);
}

async function runLookup(point: LatLng, res: number): Promise<void> {
	loading.set(true);
	errorMessage.set(null);
	try {
		const cell = await api.fromCoordinates(point.lat, point.lng, res);
		selectedCell.set(cell);
		await refreshOverlays(cell);
	} catch (err) {
		errorMessage.set(toMessage(err));
	} finally {
		loading.set(false);
	}
}

/**
 * Fetches geometry for each enabled overlay in parallel. Overlay failures are
 * swallowed so they never break the primary selection flow.
 */
async function refreshOverlays(cell: CellInfo): Promise<void> {
	const want = get(layers);
	const next: OverlayData = { neighbors: [], children: [], parent: [] };

	const tasks: Promise<void>[] = [];
	if (want.neighbors) {
		tasks.push(api.neighbors(cell.index).then((r) => void (next.neighbors = r.cells)));
	}
	if (want.children && cell.resolution < MAX_RESOLUTION) {
		tasks.push(
			api.children(cell.index, cell.resolution + 1).then((r) => void (next.children = r.cells))
		);
	}
	if (want.parent && cell.resolution > MIN_RESOLUTION) {
		tasks.push(api.parent(cell.index, cell.resolution - 1).then((g) => void (next.parent = [g])));
	}

	try {
		await Promise.all(tasks);
	} catch {
		// Keep whatever resolved; overlays are best-effort.
	}
	overlays.set(next);
}

function toMessage(err: unknown): string {
	return err instanceof ApiError ? err.message : 'Request failed. Is the backend running?';
}

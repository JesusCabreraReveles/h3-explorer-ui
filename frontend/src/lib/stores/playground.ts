// State and orchestration for the H3 Playground — a visual sandbox for the
// grid and polygon algorithms. As with the explorer, UI components render these
// stores and call the actions; all H3 work goes through the API client.

import { derived, get, writable } from 'svelte/store';
import { api, ApiError } from '$lib/services/api';
import { DEFAULT_RESOLUTION } from '$lib/config';
import type { CellGeometry, LatLng } from '$lib/types/h3';

export type Tool = 'gridDisk' | 'gridPath' | 'polygon';

/** Upper bound on k in the playground (keeps animation + labels manageable). */
export const PLAYGROUND_MAX_K = 12;

export const tool = writable<Tool>('gridDisk');
export const resolution = writable<number>(DEFAULT_RESOLUTION);

// --- Tool inputs ---
export const k = writable<number>(2);
export const ringOnly = writable<boolean>(false);
/** Vertices of the polygon currently being drawn (polygon tool). */
export const draftPolygon = writable<LatLng[]>([]);
export const polygonClosed = writable<boolean>(false);

// --- Display options ---
export const showLabels = writable<boolean>(false);
export const showBoundaries = writable<boolean>(true);

// --- Results / status ---
export const cells = writable<CellGeometry[]>([]);
export const loading = writable<boolean>(false);
export const errorMessage = writable<string | null>(null);
export const count = derived(cells, ($cells) => $cells.length);

// --- Internal interaction state ---
const originCell = writable<string | null>(null); // gridDisk origin
const pathFrom = writable<string | null>(null); // gridPath origin
let animationToken = 0;

/** Switches tool and resets all transient state. */
export function selectTool(next: Tool): void {
	tool.set(next);
	resetState();
}

/** Handles a map click according to the active tool. */
export async function handleClick(point: LatLng): Promise<void> {
	switch (get(tool)) {
		case 'gridDisk':
			return startGridDisk(point);
		case 'gridPath':
			return advanceGridPath(point);
		case 'polygon':
			return addPolygonVertex(point);
	}
}

/** Sets k and recomputes the disk/ring if an origin is set. */
export async function setK(value: number): Promise<void> {
	cancelAnimation();
	k.set(clampK(value));
	if (get(tool) === 'gridDisk' && get(originCell)) {
		await computeDisk(get(originCell)!, get(k));
	}
}

/** Toggles filled disk vs. hollow ring and recomputes. */
export async function toggleRing(value: boolean): Promise<void> {
	ringOnly.set(value);
	if (get(tool) === 'gridDisk' && get(originCell)) {
		await computeDisk(get(originCell)!, get(k));
	}
}

/** Changes resolution and recomputes whatever the active tool last produced. */
export async function setResolution(res: number): Promise<void> {
	resolution.set(res);
	await recompute();
}

/** Animates the disk growing from k=0 to the current k. */
export async function animateDisk(): Promise<void> {
	const origin = get(originCell);
	if (get(tool) !== 'gridDisk' || !origin) return;

	const target = get(k);
	const token = ++animationToken;
	for (let i = 0; i <= target; i++) {
		if (token !== animationToken) return; // superseded/cancelled
		k.set(i);
		await computeDisk(origin, i);
		await delay(160);
	}
}

/** Closes the drafted polygon and computes its covering cells. */
export async function closePolygon(): Promise<void> {
	const vertices = get(draftPolygon);
	if (get(tool) !== 'polygon' || vertices.length < 3) return;
	polygonClosed.set(true);
	await withStatus(async () => {
		const result = await api.polygonToCells(vertices, get(resolution));
		cells.set(result.cells);
	});
}

/** Clears the current result and interaction state. */
export function reset(): void {
	resetState();
}

// --- Tool implementations -------------------------------------------------

async function startGridDisk(point: LatLng): Promise<void> {
	cancelAnimation();
	await withStatus(async () => {
		const cell = await api.fromCoordinates(point.lat, point.lng, get(resolution));
		originCell.set(cell.index);
		await computeDisk(cell.index, get(k));
	});
}

async function computeDisk(origin: string, kValue: number): Promise<void> {
	const result = get(ringOnly)
		? await api.gridRing(origin, kValue)
		: await api.gridDisk(origin, kValue);
	cells.set(result.cells);
}

async function advanceGridPath(point: LatLng): Promise<void> {
	const from = get(pathFrom);
	// No origin yet, or a completed path → start a fresh one.
	if (!from || get(cells).length > 0) {
		await withStatus(async () => {
			const cell = await api.fromCoordinates(point.lat, point.lng, get(resolution));
			pathFrom.set(cell.index);
			cells.set([]);
		});
		return;
	}
	await withStatus(async () => {
		const cell = await api.fromCoordinates(point.lat, point.lng, get(resolution));
		const result = await api.gridPath(from, cell.index);
		cells.set(result.cells);
	});
}

async function addPolygonVertex(point: LatLng): Promise<void> {
	if (get(polygonClosed)) {
		// Starting a new polygon after a previous one.
		resetState();
	}
	draftPolygon.update((vertices) => [...vertices, point]);
}

/** Re-runs the active tool's computation (after a resolution change). */
async function recompute(): Promise<void> {
	const currentTool = get(tool);
	if (currentTool === 'gridDisk' && get(originCell)) {
		await withStatus(() => computeDisk(get(originCell)!, get(k)));
	} else if (currentTool === 'polygon' && get(polygonClosed)) {
		await closePolygon();
	}
	// gridPath uses fixed-resolution endpoints; a resolution change starts over.
	else if (currentTool === 'gridPath') {
		resetState();
	}
}

// --- Helpers --------------------------------------------------------------

function resetState(): void {
	cancelAnimation();
	cells.set([]);
	originCell.set(null);
	pathFrom.set(null);
	draftPolygon.set([]);
	polygonClosed.set(false);
	errorMessage.set(null);
}

function cancelAnimation(): void {
	animationToken++;
}

function clampK(value: number): number {
	return Math.min(PLAYGROUND_MAX_K, Math.max(0, Math.round(value)));
}

async function withStatus(fn: () => Promise<void>): Promise<void> {
	loading.set(true);
	errorMessage.set(null);
	try {
		await fn();
	} catch (err) {
		errorMessage.set(
			err instanceof ApiError ? err.message : 'Request failed. Is the backend running?'
		);
	} finally {
		loading.set(false);
	}
}

function delay(ms: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

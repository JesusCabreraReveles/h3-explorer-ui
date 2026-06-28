// Lazily-loaded metadata for every H3 resolution (powers the resolution
// explorer). Loaded once and cached for the session.

import { writable } from 'svelte/store';
import { api } from '$lib/services/api';
import type { ResolutionInfo } from '$lib/types/h3';

export const resolutions = writable<ResolutionInfo[]>([]);
export const resolutionsError = writable<string | null>(null);

let loaded = false;
let inflight: Promise<void> | null = null;

/** Fetches the resolution table once; subsequent calls are no-ops. */
export function loadResolutions(): Promise<void> {
	if (loaded) return Promise.resolve();
	if (inflight) return inflight;

	inflight = api
		.resolutions()
		.then((rows) => {
			resolutions.set(rows);
			loaded = true;
		})
		.catch(() => {
			resolutionsError.set('Could not load resolution metadata.');
		})
		.finally(() => {
			inflight = null;
		});

	return inflight;
}

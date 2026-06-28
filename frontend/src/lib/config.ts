// Application-wide configuration constants. The API is always called
// same-origin (relative paths) — dev via the Vite proxy, prod via the
// hooks.server.ts proxy — so no API URL needs to be injected into the bundle.

/** Base path for the backend API. Empty string == same origin. */
export const API_BASE = '';

/**
 * Free, key-less dark basemap from CARTO. Suits the dark, developer-oriented
 * theme and keeps the project zero-config to run.
 */
export const MAP_STYLE_URL = 'https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json';

/** Initial map view: downtown San Francisco. Note MapLibre uses [lng, lat]. */
export const DEFAULT_CENTER: [number, number] = [-122.41795, 37.775938];
export const DEFAULT_ZOOM = 11;

/** H3 resolution bounds (mirrors the backend domain constants). */
export const MIN_RESOLUTION = 0;
export const MAX_RESOLUTION = 15;
export const DEFAULT_RESOLUTION = 9;

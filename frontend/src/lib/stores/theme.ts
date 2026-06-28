import { writable } from 'svelte/store';

export type Theme = 'dark' | 'light';

/**
 * The app ships dark-first (it is a developer tool). The store exists so a
 * light theme can be added later without touching components — they already
 * read the theme rather than assuming it.
 */
export const theme = writable<Theme>('dark');

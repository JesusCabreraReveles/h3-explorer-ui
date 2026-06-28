import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';

// In development, Vite proxies /api (and /health) to the Go backend so the
// browser always talks to a same-origin endpoint. This mirrors the production
// hooks.server.ts proxy. Override the target with API_PROXY_TARGET when the
// backend lives elsewhere (e.g. inside docker-compose).
const apiProxyTarget = process.env.API_PROXY_TARGET ?? 'http://localhost:8080';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit({
			compilerOptions: {
				// Force runes mode across the project (libraries excepted).
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter()
		})
	],
	server: {
		proxy: {
			'/api': { target: apiProxyTarget, changeOrigin: true },
			'/health': { target: apiProxyTarget, changeOrigin: true }
		}
	},
	test: {
		// Unit tests cover pure modules (services, utils) and need no DOM.
		environment: 'node',
		include: ['src/**/*.{test,spec}.ts']
	}
});

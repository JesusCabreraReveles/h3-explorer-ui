import { defineConfig, devices } from '@playwright/test';

// E2E runs against the Vite dev server. Tests mock /api at the network boundary
// (page.route), so no backend is required and runs are deterministic.
export default defineConfig({
	testDir: './e2e',
	fullyParallel: true,
	forbidOnly: !!process.env.CI,
	retries: process.env.CI ? 1 : 0,
	workers: process.env.CI ? 1 : undefined,
	reporter: 'list',
	use: {
		baseURL: 'http://localhost:5173',
		trace: 'on-first-retry',
		// SwiftShader gives headless Chromium a WebGL context for MapLibre.
		launchOptions: { args: ['--enable-unsafe-swiftshader'] }
	},
	projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
	webServer: {
		command: 'npm run dev',
		url: 'http://localhost:5173',
		reuseExistingServer: !process.env.CI,
		timeout: 120_000
	}
});

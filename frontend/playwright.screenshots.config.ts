import { defineConfig, devices } from '@playwright/test';

// Captures real screenshots against a running backend (drives the live API, not
// mocks) for the README. Run via `npm run screenshots` with the backend up.
const viewport = { width: 1440, height: 900 };

export default defineConfig({
	testDir: './e2e',
	testMatch: 'screenshots.ts',
	reporter: 'list',
	use: {
		baseURL: 'http://localhost:5173',
		viewport,
		launchOptions: { args: ['--enable-unsafe-swiftshader'] }
	},
	projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'], viewport } }],
	webServer: {
		command: 'npm run dev',
		url: 'http://localhost:5173',
		reuseExistingServer: true,
		timeout: 120_000
	}
});

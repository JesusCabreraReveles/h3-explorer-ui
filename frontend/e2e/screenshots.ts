import { test } from '@playwright/test';

// These run against the live backend (no mocking) so the captured hexagons are
// real H3 geometry. Output goes to ../docs for the README.

test('explorer — inspector with neighbors overlay', async ({ page }) => {
	await page.goto('/');
	await page.getByPlaceholder('8928308280fffff').fill('8928308280fffff');
	await page.getByRole('button', { name: 'Inspect', exact: true }).click();
	await page.getByText('Base cell').waitFor();
	await page.getByRole('button', { name: 'Neighbors' }).click();
	// Let the basemap tiles, flyTo, and overlays settle.
	await page.waitForTimeout(3000);
	await page.screenshot({ path: '../docs/screenshot-inspector.png' });
});

test('playground — grid disk', async ({ page }) => {
	await page.goto('/playground');
	const map = page.getByTestId('playground-map');
	await map.click({ position: { x: 720, y: 460 } });
	await page.waitForTimeout(3000);
	await page.screenshot({ path: '../docs/screenshot-playground.png' });
});

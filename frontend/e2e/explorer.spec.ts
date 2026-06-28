import { expect, test } from '@playwright/test';
import { mockApi } from './fixtures';

test.beforeEach(async ({ page }) => {
	await mockApi(page);
});

test('renders the explorer shell', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByRole('heading', { name: 'H3 Explorer UI' })).toBeVisible();
	await expect(page.getByTestId('map')).toBeVisible();
});

test('validates an H3 index before inspecting', async ({ page }) => {
	await page.goto('/');
	const input = page.getByPlaceholder('8928308280fffff');
	await input.fill('not-a-valid-index');
	await input.press('Enter');
	await expect(page.getByText('Expected 15–16 hexadecimal characters.')).toBeVisible();
});

test('inspects a cell by index and shows its details', async ({ page }) => {
	await page.goto('/');
	await page.getByPlaceholder('8928308280fffff').fill('8928308280fffff');
	await page.getByRole('button', { name: 'Inspect', exact: true }).click();

	// The inspector renders the cell's index and metrics.
	await expect(page.getByText('Base cell')).toBeVisible();
	await expect(page.getByText('20', { exact: true })).toBeVisible();
	await expect(page.getByRole('button', { name: 'Copy GeoJSON' })).toBeVisible();
});

test('navigates to the playground', async ({ page }) => {
	await page.goto('/');
	await page.getByRole('link', { name: 'Playground →' }).click();
	await expect(page).toHaveURL(/\/playground$/);
	await expect(page.getByRole('heading', { name: 'H3 Playground' })).toBeVisible();
});

import { describe, expect, it } from 'vitest';
import { isH3IndexFormat } from './h3';

describe('isH3IndexFormat', () => {
	it('accepts a canonical 15-char index', () => {
		expect(isH3IndexFormat('8928308280fffff')).toBe(true);
	});

	it('accepts upper case and surrounding whitespace', () => {
		expect(isH3IndexFormat('  8928308280FFFFF  ')).toBe(true);
	});

	it('rejects too-short or too-long strings', () => {
		expect(isH3IndexFormat('8928308')).toBe(false);
		expect(isH3IndexFormat('8928308280fffff00')).toBe(false);
	});

	it('rejects non-hex characters', () => {
		expect(isH3IndexFormat('8928308280ggggg')).toBe(false);
		expect(isH3IndexFormat('37.77,-122.41')).toBe(false);
	});
});

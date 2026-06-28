import { describe, expect, it } from 'vitest';
import { isEditableTarget } from './keyboard';

describe('isEditableTarget', () => {
	it('detects inputs, textareas and selects (case-insensitive)', () => {
		expect(isEditableTarget({ tagName: 'INPUT' } as unknown as EventTarget)).toBe(true);
		expect(isEditableTarget({ tagName: 'textarea' } as unknown as EventTarget)).toBe(true);
		expect(isEditableTarget({ tagName: 'SELECT' } as unknown as EventTarget)).toBe(true);
	});

	it('detects contenteditable elements', () => {
		expect(
			isEditableTarget({ tagName: 'DIV', isContentEditable: true } as unknown as EventTarget)
		).toBe(true);
	});

	it('returns false for non-editable elements and null', () => {
		expect(isEditableTarget({ tagName: 'DIV' } as unknown as EventTarget)).toBe(false);
		expect(isEditableTarget({ tagName: 'BUTTON' } as unknown as EventTarget)).toBe(false);
		expect(isEditableTarget(null)).toBe(false);
	});
});

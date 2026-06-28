// Helpers for global keyboard shortcuts.

/**
 * Returns true when the event target is a text-editing element, so shortcut
 * handlers can bail out and not hijack typing in inputs/textareas.
 */
export function isEditableTarget(target: EventTarget | null): boolean {
	const el = target as (HTMLElement & { tagName?: string }) | null;
	if (!el || typeof el.tagName !== 'string') return false;
	const tag = el.tagName.toUpperCase();
	return tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT' || el.isContentEditable === true;
}

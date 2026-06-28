// Lightweight client-side helpers for H3 index strings. Authoritative
// validation still happens on the backend; this is only to drive UI affordances
// (enable a button, show a hint) without a round-trip.

/**
 * Returns true if the string is shaped like an H3 cell index: 15–16 hexadecimal
 * characters. This is a format check, not a validity check.
 */
export function isH3IndexFormat(value: string): boolean {
	return /^[0-9a-f]{15,16}$/i.test(value.trim());
}

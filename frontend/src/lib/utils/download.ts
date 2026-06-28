// Browser-only side effects for export: triggering a file download and copying
// to the clipboard. Separated from the pure serializers (utils/export.ts) so
// those stay unit-testable.

/** Triggers a client-side download of text content as a named file. */
export function downloadText(filename: string, content: string, mime: string): void {
	const blob = new Blob([content], { type: mime });
	const url = URL.createObjectURL(blob);
	const anchor = document.createElement('a');
	anchor.href = url;
	anchor.download = filename;
	anchor.click();
	URL.revokeObjectURL(url);
}

/** Copies text to the clipboard, returning whether it succeeded. */
export async function copyText(text: string): Promise<boolean> {
	try {
		await navigator.clipboard.writeText(text);
		return true;
	} catch {
		return false;
	}
}

// Human-friendly formatting for the metrics shown in the inspector. Pure
// functions, unit-tested, and reused across components.

/** Formats an area given in km², switching to m² for small cells. */
export function formatArea(km2: number): string {
	if (km2 < 1) {
		return `${formatNumber(km2 * 1_000_000)} m²`;
	}
	return `${formatNumber(km2)} km²`;
}

/** Formats a length given in km, switching to m below 1 km. */
export function formatLength(km: number): string {
	if (km < 1) {
		return `${formatNumber(km * 1000)} m`;
	}
	return `${formatNumber(km)} km`;
}

/** Formats a number with up to 3 significant decimals and thousands separators. */
export function formatNumber(value: number): string {
	const decimals = value >= 100 ? 0 : value >= 1 ? 2 : 4;
	return value.toLocaleString('en-US', {
		maximumFractionDigits: decimals
	});
}

/** Formats a coordinate pair as "lat, lng" with fixed precision. */
export function formatLatLng(lat: number, lng: number, precision = 5): string {
	return `${lat.toFixed(precision)}, ${lng.toFixed(precision)}`;
}

/** Compact integer formatting (e.g. 33,897,029,882). */
export function formatCount(value: number): string {
	return value.toLocaleString('en-US');
}

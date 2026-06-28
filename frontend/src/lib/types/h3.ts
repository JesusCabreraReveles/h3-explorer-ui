// TypeScript mirror of the backend domain DTOs (see backend/internal/domain).
// Keeping these in sync is what lets the whole frontend stay strictly typed
// against the API contract.

export interface LatLng {
	lat: number;
	lng: number;
}

/** Ordered ring of boundary vertices (6 for hexagons, 5 for pentagons). */
export type Boundary = LatLng[];

/** Full inspector description of a single H3 cell. */
export interface CellInfo {
	index: string;
	resolution: number;
	center: LatLng;
	boundary: Boundary;
	areaKm2: number;
	areaM2: number;
	edgeLengthKm: number;
	edgeLengthM: number;
	baseCell: number;
	icosahedronFaces: number[];
	isPentagon: boolean;
	isClassIII: boolean;
	parent?: string;
	children: string[];
	numChildren: number;
	neighbors: string[];
}

/** Minimal renderable description of a cell (grid/hierarchy endpoints). */
export interface CellGeometry {
	index: string;
	center: LatLng;
	boundary: Boundary;
}

/** Aggregate metadata for one H3 resolution. */
export interface ResolutionInfo {
	resolution: number;
	avgAreaKm2: number;
	avgEdgeLengthKm: number;
	totalCells: number;
}

export interface Polygon {
	outer: Boundary;
	holes: Boundary[];
}

export interface MultiPolygon {
	polygons: Polygon[];
}

/** Common envelope for endpoints returning many cells. */
export interface CellCollection {
	count: number;
	cells: CellGeometry[];
}

export interface BoundaryResponse {
	index: string;
	boundary: Boundary;
}

/** Error envelope returned by every endpoint on failure. */
export interface ApiErrorBody {
	error: {
		code: string;
		message: string;
	};
}

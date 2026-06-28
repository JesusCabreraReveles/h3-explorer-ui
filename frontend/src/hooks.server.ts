import type { Handle } from '@sveltejs/kit';

// In production (adapter-node) the browser hits the SvelteKit server, which
// transparently proxies API calls to the Go backend so everything is
// same-origin — no CORS, no public API URL baked into the client bundle. In
// development the Vite proxy plays the same role (see vite.config.ts).
const API_PROXY_TARGET = process.env.API_PROXY_TARGET ?? 'http://localhost:8080';

// Hop-by-hop headers that must not be forwarded verbatim. We also drop the
// upstream content-encoding/length because fetch has already decoded the body.
const STRIP_RESPONSE_HEADERS = ['content-encoding', 'content-length', 'transfer-encoding'];

function isProxied(pathname: string): boolean {
	return pathname.startsWith('/api/') || pathname === '/health';
}

export const handle: Handle = async ({ event, resolve }) => {
	const { pathname, search } = event.url;

	if (!isProxied(pathname)) {
		return resolve(event);
	}

	const target = `${API_PROXY_TARGET}${pathname}${search}`;
	const headers = new Headers(event.request.headers);
	headers.delete('host');

	const hasBody = event.request.method !== 'GET' && event.request.method !== 'HEAD';
	const upstream = await fetch(target, {
		method: event.request.method,
		headers,
		body: hasBody ? await event.request.arrayBuffer() : undefined
	});

	const responseHeaders = new Headers(upstream.headers);
	for (const h of STRIP_RESPONSE_HEADERS) {
		responseHeaders.delete(h);
	}

	return new Response(upstream.body, {
		status: upstream.status,
		statusText: upstream.statusText,
		headers: responseHeaders
	});
};

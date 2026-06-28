// The app is a client-rendered interactive map tool: MapLibre needs the DOM and
// there is no SEO benefit to SSR here. Disabling SSR keeps hydration simple and
// avoids guarding every browser-only API. The adapter-node server still runs
// (and proxies /api via hooks.server.ts).
export const ssr = false;
export const prerender = false;

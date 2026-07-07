// All API calls are prefixed with /api. This keeps them from colliding with the
// SPA's own client routes (e.g. the /credit-analyses page vs. the
// /credit-analyses API endpoint). The prefix is stripped before reaching the
// backend, by the dev-server proxy in development (proxy.conf.json) and by nginx
// in the Docker image (nginx.conf) — so the Go backend still serves its routes
// at the root and needs no change.
export const environment = {
  apiBaseUrl: '/api',
};

// apiBaseUrl is empty in development: requests are relative and forwarded to the
// Go backend by the Angular dev-server proxy (proxy.conf.json), which sidesteps
// browser CORS. In production, set this to the API origin (or serve both behind
// the same host) and add CORS on the backend.
export const environment = {
  apiBaseUrl: '',
};

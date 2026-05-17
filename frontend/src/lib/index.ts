// Re-export common lib entry points. Prefer direct $lib/... imports in app code.
export type { ApiEnvelope, ApiStatus, AppConfig } from './types/api';
export { api, apiEnvelope, ApiError, fetchAppConfig } from './api';

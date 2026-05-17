import { apiEnvelope } from '$lib/api/client';
import type { AppConfig } from '$lib/types/api';

/** Loads runtime config from GET /api/config (served by Go Fiber). */
export function fetchAppConfig(): Promise<AppConfig> {
	return apiEnvelope<AppConfig>('/api/config');
}

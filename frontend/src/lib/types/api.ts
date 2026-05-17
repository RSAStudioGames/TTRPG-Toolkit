/** Mirrors backend JSON envelope (internal/api/response.go). */
export type ApiStatus = 'success' | 'error';

export interface ApiEnvelope<T = unknown> {
	status: ApiStatus;
	data?: T;
	message?: string;
	errors?: string[];
}

/** GET /api/config — runtime values from Go (no frontend .env). */
export interface AppConfig {
	api_base_url: string;
	ws_url: string;
}

import type { ApiEnvelope } from '$lib/types/api';

/** Thrown when the API returns status "error" or a non-OK HTTP status. */
export class ApiError extends Error {
	readonly status: number;
	readonly errors: string[];

	constructor(message: string, status: number, errors: string[] = []) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
		this.errors = errors;
	}
}

/**
 * Typed fetch for same-origin API routes (session cookies when auth exists).
 */
export async function api<T>(path: string, init?: RequestInit): Promise<T> {
	const response = await fetch(path, {
		...init,
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json',
			...init?.headers
		}
	});

	if (!response.ok) {
		const text = await response.text();
		throw new ApiError(text || response.statusText, response.status);
	}

	if (response.status === 204) {
		return undefined as T;
	}

	const contentType = response.headers.get('content-type');
	if (contentType?.includes('application/json')) {
		return response.json() as Promise<T>;
	}

	return (await response.text()) as T;
}

/** Unwraps the rulebook success envelope; throws ApiError on failure. */
export async function apiEnvelope<T>(path: string, init?: RequestInit): Promise<T> {
	const body = await api<ApiEnvelope<T>>(path, init);

	if (body.status === 'error') {
		throw new ApiError(body.message ?? 'Request failed', 400, body.errors ?? []);
	}

	if (body.data === undefined) {
		throw new ApiError('Response missing data', 500);
	}

	return body.data;
}

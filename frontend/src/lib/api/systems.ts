import { apiEnvelope, ApiError } from './client';
import type { ApiEnvelope } from '$lib/types/api';
import type { DeletePreview, GameSystem, SystemListResponse } from '$lib/types/system';

export async function listSystems(params?: {
	page?: number;
	per_page?: number;
	status?: string;
	is_active?: boolean;
}): Promise<SystemListResponse> {
	const q = new URLSearchParams();
	if (params?.page) q.set('page', String(params.page));
	if (params?.per_page) q.set('per_page', String(params.per_page));
	if (params?.status) q.set('status', params.status);
	if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
	const qs = q.toString();
	return apiEnvelope<SystemListResponse>(`/api/systems${qs ? `?${qs}` : ''}`);
}

export async function getSystem(id: string): Promise<GameSystem> {
	return apiEnvelope<GameSystem>(`/api/systems/${id}`);
}

export async function createSystem(body: Record<string, unknown>): Promise<GameSystem> {
	return apiEnvelope<GameSystem>('/api/systems', {
		method: 'POST',
		body: JSON.stringify(body)
	});
}

export async function updateSystem(id: string, body: Record<string, unknown>): Promise<GameSystem> {
	return apiEnvelope<GameSystem>(`/api/systems/${id}`, {
		method: 'PUT',
		body: JSON.stringify(body)
	});
}

export async function deleteSystem(id: string): Promise<void> {
	const response = await fetch(`/api/systems/${id}`, { method: 'DELETE', credentials: 'include' });
	if (!response.ok) {
		const text = await response.text();
		throw new ApiError(text || response.statusText, response.status);
	}
	if (response.status !== 204) {
		const body = (await response.json()) as ApiEnvelope<unknown>;
		if (body.status === 'error') {
			throw new ApiError(body.message ?? 'Request failed', 400, body.errors ?? []);
		}
	}
}

export async function getDeletePreview(id: string): Promise<DeletePreview> {
	return apiEnvelope<DeletePreview>(`/api/systems/${id}/delete-preview`);
}

async function systemAction(id: string, action: string): Promise<GameSystem> {
	return apiEnvelope<GameSystem>(`/api/systems/${id}/${action}`, { method: 'POST' });
}

export const publishSystem = (id: string) => systemAction(id, 'publish');
export const lockSystem = (id: string) => systemAction(id, 'lock');
export const unlockSystem = (id: string) => systemAction(id, 'unlock');
export const archiveSystem = (id: string) => systemAction(id, 'archive');
export const restoreSystem = (id: string) => systemAction(id, 'restore');
export const cloneSystem = (id: string) => systemAction(id, 'clone');
export const forkSystem = (id: string) => systemAction(id, 'fork');

export async function uploadSystemImage(
	id: string,
	kind: 'icon' | 'cover',
	file: File
): Promise<GameSystem> {
	const form = new FormData();
	form.append('file', file);
	const path = `/api/systems/${id}/${kind}`;
	const response = await fetch(path, { method: 'POST', body: form, credentials: 'include' });
	const body = (await response.json()) as ApiEnvelope<GameSystem>;
	if (!response.ok || body.status === 'error') {
		throw new ApiError(body.message ?? 'Upload failed', response.status, body.errors ?? []);
	}
	if (!body.data) throw new ApiError('Response missing data', 500);
	return body.data;
}

export async function saveSystemTemplate(
	id: string,
	template_name: string,
	template_description: string
): Promise<void> {
	await apiEnvelope<{ saved: boolean }>(`/api/systems/${id}/save-template`, {
		method: 'POST',
		body: JSON.stringify({ template_name, template_description })
	});
}

export async function exportSystem(id: string): Promise<Blob> {
	const response = await fetch(`/api/systems/${id}/export`, { credentials: 'include' });
	if (!response.ok) throw new ApiError('Failed to export system', response.status);
	return response.blob();
}

export async function importSystem(file: File): Promise<GameSystem> {
	const form = new FormData();
	form.append('file', file);
	const response = await fetch('/api/systems/import', { method: 'POST', body: form, credentials: 'include' });
	const body = (await response.json()) as ApiEnvelope<GameSystem>;
	if (!response.ok || body.status === 'error') {
		throw new ApiError(body.message ?? 'Failed to import Template', response.status, body.errors ?? []);
	}
	if (!body.data) throw new ApiError('Response missing data', 500);
	return body.data;
}

import { listSystems } from '$lib/api/systems';
import type { GameSystem } from '$lib/types/system';

/** Shared reactive store — read fields directly in components for Svelte 5 reactivity. */
export const systemsState = $state({
	systems: [] as GameSystem[],
	loading: false,
	error: null as string | null,
	page: 1,
	totalPages: 1
});

export async function loadSystems(opts?: { page?: number; is_active?: boolean }) {
	systemsState.loading = true;
	systemsState.error = null;
	try {
		const res = await listSystems({
			page: opts?.page ?? systemsState.page,
			per_page: 20,
			is_active: opts?.is_active
		});
		systemsState.systems = res.items ?? [];
		systemsState.page = res.page;
		systemsState.totalPages = res.total_pages;
	} catch (e) {
		systemsState.error = e instanceof Error ? e.message : 'Failed to load Game Systems';
		systemsState.systems = [];
	} finally {
		systemsState.loading = false;
	}
}

export function upsertSystem(system: GameSystem) {
	const i = systemsState.systems.findIndex((s) => s.id === system.id);
	if (i >= 0) {
		systemsState.systems = systemsState.systems.map((s, idx) => (idx === i ? system : s));
	} else {
		systemsState.systems = [system, ...systemsState.systems];
	}
}

export function removeSystem(id: string) {
	systemsState.systems = systemsState.systems.filter((s) => s.id !== id);
}

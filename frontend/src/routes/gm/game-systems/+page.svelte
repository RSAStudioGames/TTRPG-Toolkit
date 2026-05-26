<script lang="ts">
	import { onMount } from 'svelte';
	import BackButton from '$lib/components/BackButton.svelte';
	import DashboardPage from '$lib/components/DashboardPage.svelte';
	import LogoutButton from '$lib/components/LogoutButton.svelte';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import SystemEditorModal from '$lib/components/SystemEditorModal.svelte';
	import { loadSystems, systemsState } from '$lib/stores/systemsStore.svelte';
	import type { GameSystem } from '$lib/types/system';

	let editorOpen = $state(false);
	let editorTarget = $state<GameSystem | null>(null);

	onMount(() => {
		loadSystems();
	});

	function openCreate() {
		editorTarget = null;
		editorOpen = true;
	}

	function openEdit(sys: GameSystem) {
		editorTarget = sys;
		editorOpen = true;
	}

	function closeEditor() {
		editorOpen = false;
		editorTarget = null;
	}
</script>

<DashboardPage title="Game Systems" role="gm" contentLayout="stack">
	{#snippet headerActions()}
		<BackButton href="/gm/system-management" parentTitle="System Management" />
		<LogoutButton />
	{/snippet}

	{#if systemsState.loading}
		<p class="status-msg">Loading game systems…</p>
	{:else if systemsState.error}
		<p class="status-msg error">{systemsState.error}</p>
	{:else if systemsState.systems.length === 0}
		<div class="empty-state">
			<p class="status-msg">No game systems yet. Create one to get started.</p>
			<button type="button" class="btn-create" onclick={openCreate}>Create</button>
		</div>
	{:else}
		<div class="list-toolbar">
			<button type="button" class="btn-create" onclick={openCreate}>Create</button>
		</div>
		<ul class="systems-list">
			{#each systemsState.systems as sys (sys.id)}
				<li>
					<button type="button" class="system-card" onclick={() => openEdit(sys)}>
						<div class="card-top">
							<strong>{sys.name}</strong>
							<StatusBadge status={sys.status} />
						</div>
						<span class="meta">{sys.slug}</span>
						{#if !sys.is_active}
							<span class="inactive">Inactive</span>
						{/if}
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</DashboardPage>

<SystemEditorModal
	open={editorOpen}
	initialSystem={editorTarget}
	onclose={closeEditor}
	onsaved={() => loadSystems()}
/>

<style>
	.btn-create {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 8px;
		background: var(--accent-gm-muted);
		font-weight: 600;
		font-family: inherit;
		cursor: pointer;
		transition: filter 0.15s ease;
	}

	.btn-create:hover {
		filter: brightness(0.97);
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		margin-top: 2rem;
	}

	.empty-state .status-msg {
		margin-top: 0;
	}

	.list-toolbar {
		display: flex;
		justify-content: flex-end;
		width: 100%;
		margin-bottom: 0.75rem;
	}

	.status-msg {
		text-align: center;
		color: var(--text-muted);
		margin-top: 2rem;
	}

	.status-msg.error {
		color: #b91c1c;
	}

	.systems-list {
		list-style: none;
		margin: 0;
		padding: 0;
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.system-card {
		width: 100%;
		text-align: left;
		padding: 1rem 1.25rem;
		border: 2px solid var(--card-border-gm);
		border-radius: 12px;
		background: var(--card-tint-gm);
		cursor: pointer;
		font: inherit;
		transition: var(--transition-card);
	}

	.system-card:hover {
		border-color: var(--card-border-hover-gm);
		box-shadow: var(--shadow-card-hover-gm);
	}

	.card-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
	}

	.meta {
		font-size: 0.85rem;
		color: var(--text-muted);
	}

	.inactive {
		display: inline-block;
		margin-top: 0.35rem;
		font-size: 0.75rem;
		color: #b45309;
		font-weight: 600;
	}
</style>

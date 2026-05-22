<script lang="ts">
	import { onMount } from 'svelte';
	import TE_ResourceModal from './TE_ResourceModal.svelte';
	import { deleteResource, listResources } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import { resourceFormatLabel, resourceTypeLabel } from '$lib/utils/resourceDefaults';
	import type { ResourceResponse } from '$lib/types/mechanics';

	interface Props {
		systemId: string;
		disabled?: boolean;
	}

	let { systemId, disabled = false }: Props = $props();

	let resources = $state<ResourceResponse[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let modalOpen = $state(false);
	let editingResource = $state<ResourceResponse | null>(null);

	const sortedResources = $derived(
		[...resources].sort((a, b) => a.sort_order - b.sort_order || a.name.localeCompare(b.name))
	);

	async function load() {
		loading = true;
		error = null;
		try {
			const data = await listResources(systemId);
			resources = data.items ?? [];
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to load resources';
		} finally {
			loading = false;
		}
	}

	function openCreate() {
		editingResource = null;
		modalOpen = true;
	}

	function openEdit(resource: ResourceResponse) {
		editingResource = resource;
		modalOpen = true;
	}

	function closeModal() {
		modalOpen = false;
		editingResource = null;
	}

	async function handleDelete(resource: ResourceResponse) {
		if (!confirm(`Delete resource "${resource.name}"?`)) return;
		error = null;
		try {
			await deleteResource(systemId, resource.id);
			await load();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to delete resource';
		}
	}

	onMount(() => {
		load();
	});
</script>

<div class="resources-tab">
	{#if loading}
		<p class="status-msg">Loading resources…</p>
	{:else}
		{#if error}
			<p class="form-error">{error}</p>
		{/if}

		<div class="mechanics-table-wrap">
			<table class="mechanics-table">
				<thead>
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Format</th>
						<th class="col-actions">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each sortedResources as resource (resource.id)}
						<tr>
							<td>{resource.name}</td>
							<td>{resourceTypeLabel(resource.type)}</td>
							<td>{resourceFormatLabel(resource.config.current_max_format)}</td>
							<td class="col-actions">
								<button
									type="button"
									class="btn-sm"
									{disabled}
									onclick={() => openEdit(resource)}
								>
									Edit
								</button>
								<button
									type="button"
									class="btn-sm btn-danger-outline"
									{disabled}
									onclick={() => handleDelete(resource)}
								>
									Delete
								</button>
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="4" class="empty-row">No resources defined yet.</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div class="tab-footer">
			<button type="button" class="btn-primary" {disabled} onclick={openCreate}>
				Create Resource
			</button>
		</div>
	{/if}
</div>

<TE_ResourceModal
	open={modalOpen}
	{systemId}
	resource={editingResource}
	siblingResourceCount={resources.length}
	{disabled}
	onclose={closeModal}
	onsaved={load}
/>

<style>
	.resources-tab {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.status-msg {
		margin: 0;
		color: var(--text-muted, #6b7280);
		font-size: 0.9rem;
	}

	.form-error {
		color: #b91c1c;
		margin: 0;
	}

	.mechanics-table-wrap {
		overflow-x: auto;
	}

	.mechanics-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.9rem;
	}

	.mechanics-table th,
	.mechanics-table td {
		padding: 0.5rem;
		border-bottom: 1px solid #e5e7eb;
		text-align: left;
		vertical-align: middle;
	}

	.col-actions {
		white-space: nowrap;
	}

	.empty-row {
		color: var(--text-muted, #6b7280);
		font-style: italic;
	}

	.tab-footer {
		margin-top: 0.25rem;
	}

	.btn-primary {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		background: var(--accent-gm, #c9a227);
		color: #1a1a1a;
		font-weight: 600;
		cursor: pointer;
		font: inherit;
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-sm {
		padding: 0.25rem 0.55rem;
		font-size: 0.8rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
		margin-right: 0.35rem;
	}

	.btn-danger-outline {
		border-color: #fca5a5;
		color: #b91c1c;
	}
</style>

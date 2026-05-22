<script lang="ts">
	import BaseOverlay from './BaseOverlay.svelte';
	import ImageUploadZone from './ImageUploadZone.svelte';
	import StatusBadge from './StatusBadge.svelte';
	import SystemDeleteModal from './SystemDeleteModal.svelte';
	import SystemMetadataForm from './SystemMetadataForm.svelte';
	import SystemSaveTemplateModal from './SystemSaveTemplateModal.svelte';
	import TE_AttributesTab from './TE_AttributesTab.svelte';
	import TE_ResolutionTab from './TE_ResolutionTab.svelte';
	import TE_SkillsTab from './TE_SkillsTab.svelte';
	import {
		archiveSystem,
		cloneSystem,
		exportSystem,
		forkSystem,
		lockSystem,
		publishSystem,
		restoreSystem,
		unlockSystem,
		updateSystem,
		uploadSystemImage
	} from '$lib/api/systems';
	import { systemsState, upsertSystem } from '$lib/stores/systemsStore.svelte';
	import {
		formValuesToPayload,
		systemToFormValues,
		type GameSystem,
		type SystemFormValues
	} from '$lib/types/system';
	import { slugValidationMessage } from '$lib/utils/slug';
	import { ApiError } from '$lib/api/client';

	interface Props {
		system: GameSystem | null;
		onclose: () => void;
	}

	let { system = $bindable(null), onclose }: Props = $props();

	let form = $state<SystemFormValues | null>(null);
	let editing = $state(false);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let deleteOpen = $state(false);
	let templateOpen = $state(false);

	const OVERLAY_TABS = [
		{ id: 'identity', label: 'Identity' },
		{ id: 'resolution', label: 'Resolution' },
		{ id: 'attributes', label: 'Attributes' },
		{ id: 'skills', label: 'Skills' },
		{ id: 'progression', label: 'Progression' },
		{ id: 'resources', label: 'Resources' },
		{ id: 'action_economy', label: 'Action Economy' }
	] as const;

	type OverlayTabId = (typeof OVERLAY_TABS)[number]['id'];

	const TAB_PLACEHOLDERS: Partial<Record<OverlayTabId, string>> = {
		progression: 'Coming in Step 9.',
		resources: 'Coming in Step 9.',
		action_economy: 'Coming in Step 10.'
	};

	let activeTab = $state<OverlayTabId>('identity');

	const readOnly = $derived(
		system != null && (system.status === 'locked' || system.status === 'archived')
	);

	const deleteDisabled = $derived(system?.is_protected === true);

	const parentOptions = $derived(
		systemsState.systems
			.filter((s) => s.status === 'published' && s.id !== system?.id)
			.map((s) => ({ id: s.id, name: s.name }))
	);

	$effect(() => {
		if (system) {
			form = systemToFormValues(system);
			editing = false;
			error = null;
			activeTab = 'identity';
		}
	});

	const activePlaceholder = $derived(TAB_PLACEHOLDERS[activeTab] ?? null);

	function close() {
		system = null;
		onclose();
	}

	async function handleSave() {
		if (!system || !form) return;
		const slugErr = slugValidationMessage(form.slug);
		if (slugErr) {
			error = slugErr;
			return;
		}
		saving = true;
		error = null;
		try {
			const updated = await updateSystem(system.id, formValuesToPayload(form));
			upsertSystem(updated);
			system = updated;
			editing = false;
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save system';
		} finally {
			saving = false;
		}
	}

	async function runAction(fn: (id: string) => Promise<GameSystem>) {
		if (!system) return;
		try {
			const updated = await fn(system.id);
			upsertSystem(updated);
			system = updated;
			form = systemToFormValues(updated);
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Action failed';
		}
	}

	async function handleExport() {
		if (!system) return;
		try {
			const blob = await exportSystem(system.id);
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `${system.slug}-export.json`;
			a.click();
			URL.revokeObjectURL(url);
		} catch {
			error = 'Failed to export system';
		}
	}

	async function onIcon(file: File) {
		if (!system) return;
		const updated = await uploadSystemImage(system.id, 'icon', file);
		upsertSystem(updated);
		system = updated;
	}

	async function onCover(file: File) {
		if (!system) return;
		const updated = await uploadSystemImage(system.id, 'cover', file);
		upsertSystem(updated);
		system = updated;
	}
</script>

{#if system && form}
	<BaseOverlay open={true} onclose={close} title={system.name}>
		{#snippet headerExtra()}
			<StatusBadge status={system!.status} />
		{/snippet}

		<div class="actions-bar">
			{#if system.status === 'draft'}
				<button type="button" class="btn-sm" onclick={() => runAction(publishSystem)}>Publish</button>
			{/if}
			{#if system.status === 'published'}
				<button type="button" class="btn-sm" onclick={() => runAction(lockSystem)}>Lock</button>
			{/if}
			{#if system.status === 'locked'}
				<button type="button" class="btn-sm" onclick={() => runAction(unlockSystem)}>Unlock</button>
			{/if}
			{#if system.status === 'draft' || system.status === 'published'}
				<button type="button" class="btn-sm" onclick={() => runAction(archiveSystem)}>Archive</button>
			{/if}
			{#if system.status === 'archived'}
				<button type="button" class="btn-sm" onclick={() => runAction(restoreSystem)}>Restore</button>
			{/if}
			<button type="button" class="btn-sm" onclick={() => runAction(cloneSystem)}>Clone</button>
			<button type="button" class="btn-sm" onclick={() => runAction(forkSystem)}>Fork</button>
			<button type="button" class="btn-sm" onclick={() => (templateOpen = true)}>Save as Template</button>
			<button type="button" class="btn-sm" onclick={handleExport}>Export</button>
		</div>

		<div class="overlay-tabs" role="tablist" aria-label="System editor sections">
			{#each OVERLAY_TABS as tab}
				<button
					type="button"
					role="tab"
					class="overlay-tab"
					class:active={activeTab === tab.id}
					aria-selected={activeTab === tab.id}
					onclick={() => (activeTab = tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</div>

		<div class="tab-panel" role="tabpanel">
			{#if activeTab === 'identity'}
				{#if !readOnly && !editing}
					<button type="button" class="btn-edit" onclick={() => (editing = true)}>Edit</button>
				{/if}

				{#if editing && form}
					<SystemMetadataForm
						bind:form
						showStatusExtras={true}
						showTagsRulebooks={true}
						{parentOptions}
					/>
					<ImageUploadZone variant="icon" previewUrl={system.icon_url} onfile={onIcon} />
					<ImageUploadZone variant="cover" previewUrl={system.cover_url} onfile={onCover} />
					{#if error}
						<p class="form-error">{error}</p>
					{/if}
					<div class="edit-actions">
						<button
							type="button"
							class="btn-secondary"
							onclick={() => {
								editing = false;
								if (system) form = systemToFormValues(system);
							}}
						>
							Cancel
						</button>
						<button type="button" class="btn-primary" disabled={saving} onclick={handleSave}>Save</button>
					</div>
				{:else if form}
					<SystemMetadataForm
						bind:form
						disabled={true}
						showStatusExtras={true}
						showTagsRulebooks={true}
						{parentOptions}
					/>
					{#if system.icon_url}
						<ImageUploadZone
							variant="icon"
							previewUrl={system.icon_url}
							disabled={readOnly}
							onfile={onIcon}
						/>
					{/if}
					{#if system.cover_url}
						<ImageUploadZone
							variant="cover"
							previewUrl={system.cover_url}
							disabled={readOnly}
							onfile={onCover}
						/>
					{/if}
				{/if}
			{:else if activeTab === 'resolution'}
				<TE_ResolutionTab systemId={system.id} disabled={readOnly} />
			{:else if activeTab === 'attributes'}
				<TE_AttributesTab systemId={system.id} disabled={readOnly} />
			{:else if activeTab === 'skills'}
				<TE_SkillsTab systemId={system.id} disabled={readOnly} />
			{:else if activePlaceholder}
				<p class="tab-placeholder">{activePlaceholder}</p>
			{/if}
		</div>

		{#snippet footer()}
			<button
				type="button"
				class="btn-danger"
				disabled={deleteDisabled}
				title={deleteDisabled ? 'Disable Deletion Protection to delete' : undefined}
				onclick={() => {
					if (!deleteDisabled) deleteOpen = true;
				}}
			>
				Delete
			</button>
		{/snippet}
	</BaseOverlay>

	<SystemDeleteModal
		open={deleteOpen}
		{system}
		onclose={() => (deleteOpen = false)}
		ondeleted={close}
	/>
	<SystemSaveTemplateModal open={templateOpen} {system} onclose={() => (templateOpen = false)} />
{/if}

<style>
	.actions-bar {
		display: flex;
		flex-wrap: wrap;
		gap: 0.35rem;
		margin-bottom: 1rem;
	}

	.btn-sm {
		padding: 0.3rem 0.6rem;
		font-size: 0.8rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.btn-edit {
		margin-bottom: 1rem;
		padding: 0.4rem 0.85rem;
		border: 1px solid var(--accent-gm);
		border-radius: 6px;
		background: var(--accent-gm-muted);
		font-weight: 600;
		cursor: pointer;
	}

	.edit-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 1rem;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		font: inherit;
	}

	.btn-primary {
		background: var(--accent-gm-muted);
		border: none;
	}

	.btn-secondary {
		background: #fff;
		border: 1px solid #d1d5db;
	}

	.btn-danger {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		background: #dc2626;
		color: #fff;
		font-weight: 600;
		cursor: pointer;
	}

	.btn-danger:disabled {
		opacity: 0.45;
		cursor: not-allowed;
	}

	.form-error {
		color: #b91c1c;
		font-size: 0.85rem;
	}

	.overlay-tabs {
		display: flex;
		flex-wrap: wrap;
		gap: 0.25rem;
		margin-bottom: 1rem;
		border-bottom: 1px solid #e5e7eb;
		padding-bottom: 0;
	}

	.overlay-tab {
		padding: 0.5rem 0.85rem;
		border: none;
		border-bottom: 2px solid transparent;
		margin-bottom: -1px;
		background: transparent;
		font: inherit;
		font-size: 0.875rem;
		font-weight: 600;
		color: #6b7280;
		cursor: pointer;
	}

	.overlay-tab:hover {
		color: var(--text-ink, #1a1a1a);
	}

	.overlay-tab.active {
		color: var(--text-ink, #1a1a1a);
		border-bottom-color: var(--accent-gm, #c9a227);
	}

	.tab-panel {
		min-height: 8rem;
	}

	.tab-placeholder {
		color: #6b7280;
		font-size: 0.9rem;
		margin: 1rem 0;
	}
</style>

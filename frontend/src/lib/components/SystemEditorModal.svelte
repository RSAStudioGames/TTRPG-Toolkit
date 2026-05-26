<script lang="ts">
	import { untrack } from 'svelte';
	import BaseModal from './BaseModal.svelte';
	import ImageUploadZone from './ImageUploadZone.svelte';
	import StatusBadge from './StatusBadge.svelte';
	import SystemDeleteModal from './SystemDeleteModal.svelte';
	import SystemMetadataForm from './SystemMetadataForm.svelte';
	import SystemSaveTemplateModal from './SystemSaveTemplateModal.svelte';
	import TE_ActionEconomyTab from './TE_ActionEconomyTab.svelte';
	import TE_AttributesTab from './TE_AttributesTab.svelte';
	import TE_ProgressionTab from './TE_ProgressionTab.svelte';
	import TE_ResolutionTab from './TE_ResolutionTab.svelte';
	import TE_ResourcesTab from './TE_ResourcesTab.svelte';
	import TE_SkillsTab from './TE_SkillsTab.svelte';
	import {
		archiveSystem,
		cloneSystem,
		createSystem,
		exportSystem,
		forkSystem,
		importSystem,
		lockSystem,
		publishSystem,
		restoreSystem,
		unlockSystem,
		updateSystem,
		uploadSystemImage
	} from '$lib/api/systems';
	import { ApiError } from '$lib/api/client';
	import { systemsState, upsertSystem } from '$lib/stores/systemsStore.svelte';
	import {
		defaultFormValues,
		formValuesToPayload,
		systemToFormValues,
		type GameSystem,
		type SystemFormValues
	} from '$lib/types/system';
	import { slugValidationMessage } from '$lib/utils/slug';

	interface Props {
		open: boolean;
		/** Pass a system to open in edit mode; pass null to open in create-wizard mode. */
		initialSystem: GameSystem | null;
		onclose: () => void;
		onsaved?: () => void;
	}

	let { open, initialSystem, onclose, onsaved }: Props = $props();

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

	let system = $state<GameSystem | null>(null);
	let form = $state<SystemFormValues>(defaultFormValues());
	let iconFile = $state<File | null>(null);
	let coverFile = $state<File | null>(null);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let importError = $state<string | null>(null);
	let activeTab = $state<OverlayTabId>('identity');
	let tabDirty = $state<Partial<Record<OverlayTabId, boolean>>>({});
	let modalOpenByTab = $state({
		attributes: false,
		skills: false,
		resources: false
	});
	let discardConfirmOpen = $state(false);
	let deleteOpen = $state(false);
	let templateOpen = $state(false);

	let attributesTab = $state<{ closeModalDiscard: () => void } | undefined>();
	let skillsTab = $state<{ closeModalDiscard: () => void } | undefined>();
	let resourcesTab = $state<{ closeModalDiscard: () => void } | undefined>();
	let resolutionTab = $state<{ save: () => Promise<boolean> } | undefined>();
	let progressionTab = $state<{ save: () => Promise<boolean> } | undefined>();
	let actionEconomyTab = $state<{ save: () => Promise<boolean> } | undefined>();

	let footerBusy = $state(false);

	const SETTINGS_TABS: OverlayTabId[] = ['resolution', 'progression', 'action_economy'];
	const LIST_TABS: OverlayTabId[] = ['attributes', 'skills', 'resources'];

	const mode = $derived<'create' | 'edit'>(system == null ? 'create' : 'edit');

	const readOnly = $derived(
		system != null && (system.status === 'locked' || system.status === 'archived')
	);

	const deleteDisabled = $derived(system?.is_protected === true);

	const parentOptions = $derived(
		systemsState.systems
			.filter((s) => s.status === 'published' && s.id !== system?.id)
			.map((s) => ({ id: s.id, name: s.name }))
	);

	const modalTitle = $derived.by(() => {
		if (mode === 'create') return 'Create TTRPG Game System';
		return `Edit TTRPG Game System — ${system?.name ?? ''}`;
	});

	const lastTabId: OverlayTabId = OVERLAY_TABS[OVERLAY_TABS.length - 1].id;

	function isTabLocked(id: OverlayTabId): boolean {
		return mode === 'create' && id !== 'identity';
	}

	function identityFormDirty(): boolean {
		if (mode === 'create') {
			if (iconFile || coverFile) return true;
			return JSON.stringify(form) !== JSON.stringify(defaultFormValues());
		}
		if (!system) return false;
		return JSON.stringify(form) !== JSON.stringify(systemToFormValues(system));
	}

	function anyTabDirty(): boolean {
		if (identityFormDirty()) return true;
		for (const t of OVERLAY_TABS) {
			if (t.id !== 'identity' && tabDirty[t.id]) return true;
		}
		return false;
	}

	function setTabDirty(tab: OverlayTabId, dirty: boolean) {
		tabDirty = { ...tabDirty, [tab]: dirty };
	}

	function modalOpenOnTab(tab: OverlayTabId): boolean {
		if (tab === 'attributes') return modalOpenByTab.attributes;
		if (tab === 'skills') return modalOpenByTab.skills;
		if (tab === 'resources') return modalOpenByTab.resources;
		return false;
	}

	function closeModalOnTab(tab: OverlayTabId) {
		if (tab === 'attributes') attributesTab?.closeModalDiscard();
		else if (tab === 'skills') skillsTab?.closeModalDiscard();
		else if (tab === 'resources') resourcesTab?.closeModalDiscard();
	}

	function tabHasUnsavedDirty(tab: OverlayTabId): boolean {
		if (tab === 'identity') return identityFormDirty();
		return tabDirty[tab] === true;
	}

	function confirmLeaveCurrentTab(): boolean {
		if (modalOpenOnTab(activeTab)) {
			return confirm('You have unsaved changes in the open form. Leave without saving?');
		}
		if (tabHasUnsavedDirty(activeTab)) {
			return confirm('You have unsaved changes. Leave without saving?');
		}
		return true;
	}

	function trySwitchTab(next: OverlayTabId) {
		if (next === activeTab) return;
		if (isTabLocked(next)) return;
		if (!confirmLeaveCurrentTab()) return;
		if (modalOpenOnTab(activeTab)) closeModalOnTab(activeTab);
		activeTab = next;
	}

	function advanceTab() {
		const idx = OVERLAY_TABS.findIndex((t) => t.id === activeTab);
		if (idx < OVERLAY_TABS.length - 1) {
			activeTab = OVERLAY_TABS[idx + 1].id;
		} else {
			confirmClose();
		}
	}

	const settingsSaveLabel = $derived.by(() => {
		if (mode === 'edit') return 'Save';
		return activeTab === lastTabId ? 'Save & Finish' : 'Save & Next';
	});

	const listNextLabel = $derived(activeTab === lastTabId ? 'Finish' : 'Next');

	async function handleSettingsSave() {
		footerBusy = true;
		try {
			let ok = false;
			if (activeTab === 'resolution') ok = (await resolutionTab?.save()) ?? false;
			else if (activeTab === 'progression') ok = (await progressionTab?.save()) ?? false;
			else if (activeTab === 'action_economy') ok = (await actionEconomyTab?.save()) ?? false;
			if (ok && mode === 'create') advanceTab();
		} finally {
			footerBusy = false;
		}
	}

	function initializeForOpen() {
		const target = initialSystem;
		if (target) {
			system = target;
			form = systemToFormValues(target);
		} else {
			system = null;
			form = defaultFormValues();
		}
		iconFile = null;
		coverFile = null;
		activeTab = 'identity';
		tabDirty = {};
		modalOpenByTab = { attributes: false, skills: false, resources: false };
		error = null;
		importError = null;
		saving = false;
		discardConfirmOpen = false;
	}

	$effect(() => {
		if (!open) return;
		untrack(initializeForOpen);
	});

	function requestClose() {
		if (anyTabDirty() || modalOpenOnTab(activeTab)) {
			discardConfirmOpen = true;
			return;
		}
		confirmClose();
	}

	function confirmClose() {
		discardConfirmOpen = false;
		if (modalOpenOnTab(activeTab)) closeModalOnTab(activeTab);
		onclose();
	}

	async function handleNextIdentity() {
		const trimmedName = form.name.trim();
		if (!trimmedName || trimmedName.length < 3) {
			error = 'System Name is required (at least 3 characters).';
			return;
		}
		const slugErr = slugValidationMessage(form.slug);
		if (slugErr) {
			error = slugErr;
			return;
		}
		saving = true;
		error = null;
		try {
			let created = await createSystem(formValuesToPayload(form));
			if (iconFile) created = await uploadSystemImage(created.id, 'icon', iconFile);
			if (coverFile) created = await uploadSystemImage(created.id, 'cover', coverFile);
			upsertSystem(created);
			system = created;
			form = systemToFormValues(created);
			iconFile = null;
			coverFile = null;
			tabDirty = {};
			activeTab = 'resolution';
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to create system';
		} finally {
			saving = false;
		}
	}

	async function handleSaveIdentity() {
		if (!system) return;
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
			form = systemToFormValues(updated);
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save system';
		} finally {
			saving = false;
		}
	}

	function handleCancelIdentity() {
		const dirty = identityFormDirty();
		if (dirty && !confirm('Discard changes on this tab?')) return;
		if (mode === 'create') {
			form = defaultFormValues();
			iconFile = null;
			coverFile = null;
		} else if (system) {
			form = systemToFormValues(system);
			iconFile = null;
			coverFile = null;
		}
		error = null;
	}

	async function handleImport(e: Event) {
		const input = e.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		importError = null;
		try {
			const s = await importSystem(file);
			upsertSystem(s);
			system = s;
			form = systemToFormValues(s);
			iconFile = null;
			coverFile = null;
			tabDirty = {};
			activeTab = 'identity';
			onsaved?.();
		} catch (err) {
			importError = err instanceof ApiError ? err.message : 'Failed to import Template';
		}
		input.value = '';
	}

	async function runAction(fn: (id: string) => Promise<GameSystem>) {
		if (!system) return;
		try {
			const updated = await fn(system.id);
			upsertSystem(updated);
			system = updated;
			form = systemToFormValues(updated);
			onsaved?.();
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
		if (!system) {
			iconFile = file;
			return;
		}
		const updated = await uploadSystemImage(system.id, 'icon', file);
		upsertSystem(updated);
		system = updated;
		form = systemToFormValues(updated);
		onsaved?.();
	}

	async function onCover(file: File) {
		if (!system) {
			coverFile = file;
			return;
		}
		const updated = await uploadSystemImage(system.id, 'cover', file);
		upsertSystem(updated);
		system = updated;
		form = systemToFormValues(updated);
		onsaved?.();
	}
</script>

<BaseModal
	{open}
	title={modalTitle}
	wide
	closeOnBackdrop={false}
	closeOnEscape={!discardConfirmOpen && !deleteOpen && !templateOpen}
	onclose={requestClose}
>
	{#snippet headerExtra()}
		<div class="editor-tabs" role="tablist" aria-label="System editor steps">
			{#each OVERLAY_TABS as tab, idx (tab.id)}
				<button
					type="button"
					role="tab"
					class="editor-tab"
					class:active={activeTab === tab.id}
					class:locked={isTabLocked(tab.id)}
					aria-selected={activeTab === tab.id}
					disabled={isTabLocked(tab.id)}
					title={isTabLocked(tab.id) ? 'Complete the Identity step to unlock' : undefined}
					onclick={() => trySwitchTab(tab.id)}
				>
					<span class="step-number">{idx + 1}</span>
					<span class="step-label">{tab.label}</span>
				</button>
			{/each}
		</div>
	{/snippet}

	{#if mode === 'edit' && system}
		<div class="status-row">
			<StatusBadge status={system.status} />
		</div>
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
	{/if}

	<div class="tab-panel" role="tabpanel">
		{#if activeTab === 'identity'}
			{#if mode === 'create'}
				<div class="import-zone">
					<span class="import-label">Import from JSON</span>
					<label class="file-picker">
						<input
							type="file"
							class="file-picker-input"
							accept=".json,application/json"
							onchange={handleImport}
						/>
						<span class="file-picker-button">
							<i class="fas fa-file-import" aria-hidden="true"></i>
							Choose JSON file
						</span>
						<span class="file-picker-hint">Import a validated system export</span>
					</label>
					{#if importError}
						<p class="form-error">{importError}</p>
					{/if}
				</div>

				<SystemMetadataForm bind:form />

				<ImageUploadZone variant="icon" onfile={onIcon} disabled={saving} />
				<ImageUploadZone variant="cover" onfile={onCover} disabled={saving} />

				{#if error}
					<p class="form-error">{error}</p>
				{/if}
			{:else if system}
				<SystemMetadataForm
					bind:form
					disabled={readOnly}
					showStatusExtras={true}
					showTagsRulebooks={true}
					{parentOptions}
				/>
				<ImageUploadZone
					variant="icon"
					previewUrl={system.icon_url}
					disabled={readOnly}
					onfile={onIcon}
				/>
				<ImageUploadZone
					variant="cover"
					previewUrl={system.cover_url}
					disabled={readOnly}
					onfile={onCover}
				/>

				{#if error}
					<p class="form-error">{error}</p>
				{/if}
			{/if}
		{:else if system}
			{#if activeTab === 'resolution'}
				<TE_ResolutionTab
					bind:this={resolutionTab}
					systemId={system.id}
					disabled={readOnly}
					onDirtyChange={(d) => setTabDirty('resolution', d)}
				/>
			{:else if activeTab === 'attributes'}
				<TE_AttributesTab
					bind:this={attributesTab}
					systemId={system.id}
					disabled={readOnly}
					onModalOpenChange={(o) => {
						modalOpenByTab = { ...modalOpenByTab, attributes: o };
					}}
				/>
			{:else if activeTab === 'skills'}
				<TE_SkillsTab
					bind:this={skillsTab}
					systemId={system.id}
					disabled={readOnly}
					onModalOpenChange={(o) => {
						modalOpenByTab = { ...modalOpenByTab, skills: o };
					}}
				/>
			{:else if activeTab === 'progression'}
				<TE_ProgressionTab
					bind:this={progressionTab}
					systemId={system.id}
					disabled={readOnly}
					onDirtyChange={(d) => setTabDirty('progression', d)}
				/>
			{:else if activeTab === 'resources'}
				<TE_ResourcesTab
					bind:this={resourcesTab}
					systemId={system.id}
					disabled={readOnly}
					onModalOpenChange={(o) => {
						modalOpenByTab = { ...modalOpenByTab, resources: o };
					}}
				/>
			{:else if activeTab === 'action_economy'}
				<TE_ActionEconomyTab
					bind:this={actionEconomyTab}
					systemId={system.id}
					disabled={readOnly}
					onDirtyChange={(d) => setTabDirty('action_economy', d)}
				/>
			{/if}
		{/if}
	</div>

	{#snippet footer()}
		<div class="footer-left">
			{#if mode === 'edit'}
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
			{/if}
		</div>
		<div class="footer-right">
			<button type="button" class="btn-secondary" disabled={footerBusy || saving} onclick={requestClose}>
				Cancel
			</button>
			{#if activeTab === 'identity'}
				{#if mode === 'create'}
					<button
						type="button"
						class="btn-primary"
						disabled={saving}
						onclick={handleNextIdentity}
					>
						{saving ? 'Creating…' : 'Next'}
					</button>
				{:else if !readOnly}
					<button
						type="button"
						class="btn-primary"
						disabled={saving}
						onclick={handleSaveIdentity}
					>
						{saving ? 'Saving…' : 'Save'}
					</button>
				{/if}
			{:else if SETTINGS_TABS.includes(activeTab)}
				{#if !readOnly}
					<button
						type="button"
						class="btn-primary"
						disabled={footerBusy}
						onclick={handleSettingsSave}
					>
						{footerBusy ? 'Saving…' : settingsSaveLabel}
					</button>
				{/if}
			{:else if LIST_TABS.includes(activeTab) && mode === 'create'}
				<button type="button" class="btn-primary" onclick={advanceTab}>
					{listNextLabel}
				</button>
			{/if}
		</div>
	{/snippet}
</BaseModal>

<BaseModal
	open={discardConfirmOpen}
	title="Discard changes?"
	stacked
	closeOnBackdrop={false}
	onclose={() => (discardConfirmOpen = false)}
>
	<p class="discard-msg">Are you sure? Any information you entered will be lost.</p>
	{#snippet footer()}
		<button type="button" class="btn-secondary" onclick={() => (discardConfirmOpen = false)}>
			Keep editing
		</button>
		<button type="button" class="btn-danger" onclick={confirmClose}>Discard</button>
	{/snippet}
</BaseModal>

{#if system}
	<SystemDeleteModal
		open={deleteOpen}
		{system}
		onclose={() => (deleteOpen = false)}
		ondeleted={() => {
			deleteOpen = false;
			onsaved?.();
			onclose();
		}}
	/>
	<SystemSaveTemplateModal open={templateOpen} {system} onclose={() => (templateOpen = false)} />
{/if}

<style>
	.editor-tabs {
		display: flex;
		flex-wrap: wrap;
		gap: 0.25rem;
		padding-top: 0.25rem;
	}

	.editor-tab {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.4rem 0.75rem;
		border: none;
		border-bottom: 2px solid transparent;
		background: transparent;
		font: inherit;
		font-size: 0.85rem;
		font-weight: 600;
		color: #6b7280;
		cursor: pointer;
	}

	.editor-tab:hover:not(:disabled) {
		color: var(--text-ink, #1a1a1a);
	}

	.editor-tab.active {
		color: var(--text-ink, #1a1a1a);
		border-bottom-color: var(--accent-gm, #c9a227);
	}

	.editor-tab:disabled,
	.editor-tab.locked {
		opacity: 0.45;
		cursor: not-allowed;
	}

	.step-number {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 1.4rem;
		height: 1.4rem;
		border-radius: 50%;
		background: var(--accent-gm-muted, #f5ecd4);
		color: var(--text-ink, #1a1a1a);
		font-size: 0.75rem;
	}

	.editor-tab.active .step-number {
		background: var(--accent-gm, #c9a227);
	}

	.step-label {
		white-space: nowrap;
	}

	.status-row {
		display: flex;
		align-items: center;
		margin-bottom: 0.75rem;
	}

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

	.tab-panel {
		min-height: 8rem;
	}

	.import-zone {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		margin-bottom: 1rem;
		border: 1px dashed var(--accent-gm-muted, #e7d59a);
		border-radius: 8px;
		background: rgba(255, 255, 255, 0.5);
	}

	.import-label {
		font-weight: 600;
		font-size: 0.85rem;
	}

	.file-picker {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.file-picker-input {
		display: none;
	}

	.file-picker-button {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.4rem 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		font-weight: 600;
		font-size: 0.85rem;
		cursor: pointer;
		width: fit-content;
	}

	.file-picker-hint {
		font-size: 0.8rem;
		color: var(--text-muted, #6b7280);
	}

	.footer-left {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-right: auto;
	}

	.footer-right {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font: inherit;
		font-weight: 600;
		cursor: pointer;
	}

	.btn-primary {
		background: var(--accent-gm-muted, #f5ecd4);
		border: none;
		color: var(--text-ink, #1a1a1a);
	}

	.btn-secondary {
		background: #fff;
		border: 1px solid #d1d5db;
	}

	.btn-primary:disabled,
	.btn-secondary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
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
		margin: 0;
	}

	.discard-msg {
		margin: 0;
		font-size: 0.95rem;
		line-height: 1.5;
	}
</style>

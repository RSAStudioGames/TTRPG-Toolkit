<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import ImageUploadZone from './ImageUploadZone.svelte';
	import SystemMetadataForm from './SystemMetadataForm.svelte';
	import { createSystem, importSystem, uploadSystemImage } from '$lib/api/systems';
	import { upsertSystem } from '$lib/stores/systemsStore.svelte';
	import { defaultFormValues, formValuesToPayload } from '$lib/types/system';
	import { slugValidationMessage } from '$lib/utils/slug';
	import { ApiError } from '$lib/api/client';

	interface Props {
		open: boolean;
		onclose: () => void;
		onsaved?: () => void;
	}

	let { open, onclose, onsaved }: Props = $props();

	let form = $state(defaultFormValues());
	let saving = $state(false);
	let error = $state<string | null>(null);
	let iconFile = $state<File | null>(null);
	let coverFile = $state<File | null>(null);
	let importError = $state<string | null>(null);
	let discardConfirmOpen = $state(false);

	function isFormDirty(): boolean {
		if (iconFile || coverFile) return true;
		return JSON.stringify(form) !== JSON.stringify(defaultFormValues());
	}

	function requestClose() {
		if (isFormDirty()) {
			discardConfirmOpen = true;
			return;
		}
		confirmClose();
	}

	function confirmClose() {
		discardConfirmOpen = false;
		reset();
		onclose();
	}

	function reset() {
		form = defaultFormValues();
		iconFile = null;
		coverFile = null;
		error = null;
		importError = null;
	}

	async function handleSave() {
		const slugErr = slugValidationMessage(form.slug);
		if (!form.name.trim() || form.name.trim().length < 3) {
			error = 'System Name is required (at least 3 characters).';
			return;
		}
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
			confirmClose();
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to create system';
		} finally {
			saving = false;
		}
	}

	async function handleImport(e: Event) {
		const input = e.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		importError = null;
		try {
			const s = await importSystem(file);
			upsertSystem(s);
			confirmClose();
			onsaved?.();
		} catch (err) {
			importError = err instanceof ApiError ? err.message : 'Failed to import Template';
		}
		input.value = '';
	}
</script>

<BaseModal
	{open}
	title="Create Game System"
	closeOnBackdrop={false}
	closeOnEscape={!discardConfirmOpen}
	onclose={requestClose}
>
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

	<ImageUploadZone
		variant="icon"
		onfile={(f) => (iconFile = f)}
		disabled={saving}
	/>
	<ImageUploadZone
		variant="cover"
		onfile={(f) => (coverFile = f)}
		disabled={saving}
	/>

	{#if error}
		<p class="form-error">{error}</p>
	{/if}

	{#snippet footer()}
		<button type="button" class="btn-secondary" onclick={requestClose}>Cancel</button>
		<button type="button" class="btn-primary" disabled={saving} onclick={handleSave}>
			{saving ? 'Creating…' : 'Create'}
		</button>
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

<style>
	.discard-msg {
		margin: 0;
		font-size: 0.95rem;
		line-height: 1.5;
	}

	.btn-danger {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font: inherit;
		font-weight: 600;
		cursor: pointer;
		border: 1px solid #fca5a5;
		background: #fef2f2;
		color: #b91c1c;
	}
	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font: inherit;
		font-weight: 600;
		cursor: pointer;
		border: none;
	}

	.btn-primary {
		background: var(--accent-gm-muted);
		color: var(--text-ink);
	}

	.btn-secondary {
		background: #fff;
		border: 1px solid #d1d5db;
	}
</style>

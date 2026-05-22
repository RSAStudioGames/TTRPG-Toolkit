<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import { deleteSystem, getDeletePreview } from '$lib/api/systems';
	import { removeSystem } from '$lib/stores/systemsStore.svelte';
	import type { GameSystem } from '$lib/types/system';
	import { ApiError } from '$lib/api/client';

	interface Props {
		open: boolean;
		system: GameSystem | null;
		onclose: () => void;
		ondeleted?: () => void;
	}

	let { open, system, onclose, ondeleted }: Props = $props();

	let confirmName = $state('');
	let previewText = $state('');
	let deleting = $state(false);
	let error = $state<string | null>(null);

	$effect(() => {
		if (open && system) {
			confirmName = '';
			error = null;
			getDeletePreview(system.id)
				.then((p) => {
					previewText = `This will permanently delete this system and ${p.total_associated} associated items.`;
				})
				.catch(() => {
					previewText = 'This will permanently delete this system and its associated items.';
				});
		}
	});

	const canDelete = $derived(system != null && confirmName === system.name);

	async function handleDelete() {
		if (!system || !canDelete) return;
		deleting = true;
		error = null;
		try {
			await deleteSystem(system.id);
			removeSystem(system.id);
			onclose();
			ondeleted?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to delete system';
		} finally {
			deleting = false;
		}
	}
</script>

<BaseModal {open} title="Delete Game System" onclose={onclose}>
	{#if system}
		<p>{previewText}</p>
		<p>Type <strong>{system.name}</strong> to confirm:</p>
		<input type="text" bind:value={confirmName} placeholder={system.name} />
		{#if error}
			<p class="form-error">{error}</p>
		{/if}
	{/if}

	{#snippet footer()}
		<button type="button" class="btn-secondary" onclick={onclose}>Cancel</button>
		<button type="button" class="btn-danger" disabled={!canDelete || deleting} onclick={handleDelete}>
			{deleting ? 'Deleting…' : 'Delete'}
		</button>
	{/snippet}
</BaseModal>

<style>
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
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		padding: 0.5rem 1rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		font: inherit;
		cursor: pointer;
	}

	input {
		width: 100%;
		padding: 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
		margin-top: 0.5rem;
	}

	.form-error {
		color: #b91c1c;
		font-size: 0.85rem;
	}
</style>

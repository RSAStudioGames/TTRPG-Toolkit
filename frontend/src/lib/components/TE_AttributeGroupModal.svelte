<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import {
		createAttributeGroup,
		updateAttributeGroup
	} from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import type { AttributeGroupResponse } from '$lib/types/mechanics';

	interface Props {
		open: boolean;
		systemId: string;
		group?: AttributeGroupResponse | null;
		siblingGroupCount: number;
		disabled?: boolean;
		onclose: () => void;
		onsaved?: () => void;
	}

	let {
		open,
		systemId,
		group = null,
		siblingGroupCount,
		disabled = false,
		onclose,
		onsaved
	}: Props = $props();

	let name = $state('');
	let saving = $state(false);
	let error = $state<string | null>(null);

	const isEdit = $derived(group != null);
	const modalTitle = $derived(isEdit ? 'Edit Group' : 'Create Group');

	function resetForm() {
		name = group?.name ?? '';
		error = null;
	}

	$effect(() => {
		if (open) resetForm();
	});

	async function handleSave() {
		const trimmed = name.trim();
		if (!trimmed) {
			error = 'Group name is required.';
			return;
		}
		saving = true;
		error = null;
		try {
			if (isEdit && group) {
				await updateAttributeGroup(systemId, group.id, { name: trimmed });
			} else {
				await createAttributeGroup(systemId, {
					name: trimmed,
					sort_order: siblingGroupCount
				});
			}
			onclose();
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save group';
		} finally {
			saving = false;
		}
	}
</script>

<BaseModal {open} title={modalTitle} onclose={onclose}>
	<div class="group-form">
		<div class="form-field">
			<label for="group-name">Group Name</label>
			<input id="group-name" type="text" bind:value={name} {disabled} required />
		</div>
		{#if error}
			<p class="form-error">{error}</p>
		{/if}
	</div>

	{#snippet footer()}
		<button type="button" class="btn-secondary" onclick={onclose}>Cancel</button>
		<button type="button" class="btn-primary" disabled={disabled || saving} onclick={handleSave}>
			{saving ? 'Saving…' : 'Save'}
		</button>
	{/snippet}
</BaseModal>

<style>
	.form-field label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		margin-bottom: 0.35rem;
	}

	.form-field input {
		width: 100%;
		padding: 0.45rem 0.6rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.form-error {
		color: #b91c1c;
		margin: 0;
		font-size: 0.875rem;
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
		border: none;
		background: var(--accent-gm, #c9a227);
		color: #1a1a1a;
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		border: 1px solid #d1d5db;
		background: #fff;
	}
</style>

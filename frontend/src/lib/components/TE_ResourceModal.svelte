<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import { createResource, updateResource } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		ADD_RECOVERY_EVENT_LABEL,
		MIN_VAL_HELPER,
		RECOVERY_TRIGGER_OPTIONS,
		RESOURCE_FORMAT_OPTIONS,
		RESOURCE_TYPE_OPTIONS
	} from '$lib/constants/resourceOptions';
	import {
		defaultResourceForm,
		formToCreatePayload,
		formToUpdatePayload,
		resourceToForm,
		validateResourceForm,
		type ResourceFormState
	} from '$lib/utils/resourceDefaults';
	import type { ResourceResponse } from '$lib/types/mechanics';

	interface Props {
		open: boolean;
		systemId: string;
		resource?: ResourceResponse | null;
		siblingResourceCount: number;
		disabled?: boolean;
		onclose: () => void;
		onsaved?: () => void;
	}

	let {
		open,
		systemId,
		resource = null,
		siblingResourceCount,
		disabled = false,
		onclose,
		onsaved
	}: Props = $props();

	let form = $state<ResourceFormState>(defaultResourceForm());
	let saving = $state(false);
	let error = $state<string | null>(null);

	const isEdit = $derived(resource != null);
	const modalTitle = $derived(isEdit ? 'Edit Resource' : 'Create Resource');

	function resetForm() {
		form = resource ? resourceToForm(resource) : defaultResourceForm();
		error = null;
	}

	$effect(() => {
		if (open) resetForm();
	});

	function addRecovery() {
		form.config.recovery_schedules = [
			...(form.config.recovery_schedules ?? []),
			{ trigger: 'short_rest', amount: '', conditions: '' }
		];
	}

	function removeRecovery(index: number) {
		form.config.recovery_schedules = (form.config.recovery_schedules ?? []).filter(
			(_, i) => i !== index
		);
	}

	async function handleSave() {
		const validation = validateResourceForm(form);
		if (validation) {
			error = validation;
			return;
		}
		saving = true;
		error = null;
		try {
			if (isEdit && resource) {
				await updateResource(systemId, resource.id, formToUpdatePayload(form));
			} else {
				await createResource(systemId, formToCreatePayload(form, siblingResourceCount));
			}
			onclose();
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save resource';
		} finally {
			saving = false;
		}
	}
</script>

<BaseModal {open} title={modalTitle} wide onclose={onclose}>
	<div class="resource-form">
		<div class="form-field">
			<label for="res-name">Name</label>
			<input id="res-name" type="text" bind:value={form.name} {disabled} required />
		</div>

		<div class="form-field">
			<label for="res-type">Type</label>
			<select id="res-type" bind:value={form.type} {disabled}>
				{#each RESOURCE_TYPE_OPTIONS as opt}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>
		</div>

		<fieldset class="format-group">
			<legend>Value Bounds</legend>
			{#each RESOURCE_FORMAT_OPTIONS as opt}
				<label class="radio-option">
					<input
						type="radio"
						name="resource-format"
						value={opt.value}
						checked={form.config.current_max_format === opt.value}
						{disabled}
						onchange={() => {
							form.config.current_max_format = opt.value;
						}}
					/>
					<span>{opt.label}</span>
				</label>
			{/each}
		</fieldset>

		<div class="form-field">
			<label for="res-min">Minimum</label>
			<input id="res-min" type="number" bind:value={form.config.min_val} {disabled} />
			<p class="hint">{MIN_VAL_HELPER}</p>
		</div>

		<div class="form-field">
			<label for="res-max">Maximum Value</label>
			<input
				id="res-max"
				type="text"
				bind:value={form.config.max_val_formula}
				{disabled}
				placeholder="Number or formula"
			/>
		</div>

		<div class="dynamic-section">
			<h4>Recovery Schedule</h4>
			{#each form.config.recovery_schedules ?? [] as row, i}
				<div class="recovery-row">
					<div class="form-field">
						<label for="res-trigger-{i}">Trigger</label>
						<select
							id="res-trigger-{i}"
							bind:value={form.config.recovery_schedules![i].trigger}
							{disabled}
						>
							{#each RECOVERY_TRIGGER_OPTIONS as opt}
								<option value={opt.value}>{opt.label}</option>
							{/each}
						</select>
					</div>
					<div class="form-field">
						<label for="res-amount-{i}">Recovery Amount</label>
						<input
							id="res-amount-{i}"
							type="text"
							bind:value={form.config.recovery_schedules![i].amount}
							{disabled}
							placeholder="Amount or formula"
						/>
					</div>
					<div class="form-field">
						<label for="res-cond-{i}">Conditions</label>
						<input
							id="res-cond-{i}"
							type="text"
							bind:value={form.config.recovery_schedules![i].conditions}
							{disabled}
						/>
					</div>
					<button
						type="button"
						class="btn-sm btn-danger-outline"
						{disabled}
						aria-label="Delete row"
						onclick={() => removeRecovery(i)}
					>
						Delete
					</button>
				</div>
			{/each}
			<button type="button" class="btn-secondary-sm" {disabled} onclick={addRecovery}>
				{ADD_RECOVERY_EVENT_LABEL}
			</button>
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
	.resource-form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.form-field label,
	.format-group legend {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		margin-bottom: 0.35rem;
	}

	.form-field input,
	.form-field select {
		width: 100%;
		padding: 0.45rem 0.6rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.format-group {
		border: none;
		margin: 0;
		padding: 0;
	}

	.radio-option {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		margin-bottom: 0.35rem;
		cursor: pointer;
	}

	.hint {
		margin: 0.25rem 0 0;
		font-size: 0.8rem;
		color: var(--text-muted, #6b7280);
	}

	.dynamic-section h4 {
		margin: 0 0 0.5rem;
		font-size: 0.9rem;
	}

	.recovery-row {
		display: grid;
		gap: 0.5rem;
		padding: 0.75rem;
		margin-bottom: 0.5rem;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		background: #fafafa;
	}

	.btn-sm {
		justify-self: start;
		padding: 0.25rem 0.55rem;
		font-size: 0.8rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.btn-danger-outline {
		border-color: #fca5a5;
		color: #b91c1c;
	}

	.btn-secondary-sm {
		padding: 0.35rem 0.75rem;
		font-size: 0.85rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.form-error {
		color: #b91c1c;
		margin: 0;
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

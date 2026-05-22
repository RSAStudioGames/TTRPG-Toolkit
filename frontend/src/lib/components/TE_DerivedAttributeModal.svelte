<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import TE_FormulaBuilder from './TE_FormulaBuilder.svelte';
	import { createAttribute, updateAttribute } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		DERIVATION_FORMULA_HELPER,
		RECALCULATE_TRIGGER_OPTIONS
	} from '$lib/constants/attributeOptions';
	import {
		defaultDerivedAttributeForm,
		derivedAttributeToForm,
		derivationFormulaVariables,
		formToCreatePayload,
		formToUpdatePayload,
		validateDerivedAttributeForm,
		type AttributeFormState
	} from '$lib/utils/attributeDefaults';
	import type { AttributeResponse } from '$lib/types/mechanics';

	interface Props {
		open: boolean;
		systemId: string;
		attribute?: AttributeResponse | null;
		siblingAttributes: AttributeResponse[];
		disabled?: boolean;
		onclose: () => void;
		onsaved?: () => void;
	}

	let {
		open,
		systemId,
		attribute = null,
		siblingAttributes,
		disabled = false,
		onclose,
		onsaved
	}: Props = $props();

	let form = $state<AttributeFormState>(defaultDerivedAttributeForm());
	let saving = $state(false);
	let error = $state<string | null>(null);

	const isEdit = $derived(attribute != null);
	const modalTitle = $derived(isEdit ? 'Edit Derived Attribute' : 'Create Derived Attribute');

	const derivationVariables = $derived(
		derivationFormulaVariables(siblingAttributes, attribute?.id)
	);

	function resetForm() {
		form = attribute ? derivedAttributeToForm(attribute) : defaultDerivedAttributeForm();
		error = null;
	}

	$effect(() => {
		if (open) resetForm();
	});

	function toggleTrigger(value: string, checked: boolean) {
		const current = form.config.recalculate_triggers ?? [];
		if (checked) {
			form.config.recalculate_triggers = [...current, value];
		} else {
			form.config.recalculate_triggers = current.filter((t) => t !== value);
		}
	}

	async function handleSave() {
		form.config.is_derived = true;
		form.config.caching_rule = 'on_trigger';
		const validation = validateDerivedAttributeForm(form);
		if (validation) {
			error = validation;
			return;
		}
		saving = true;
		error = null;
		try {
			if (isEdit && attribute) {
				await updateAttribute(systemId, attribute.id, formToUpdatePayload(form));
			} else {
				const derivedCount = siblingAttributes.filter((a) => a.config.is_derived).length;
				await createAttribute(systemId, formToCreatePayload(form, derivedCount));
			}
			onclose();
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save derived attribute';
		} finally {
			saving = false;
		}
	}
</script>

<BaseModal {open} title={modalTitle} onclose={onclose}>
	<div class="derived-form">
		<div class="form-field">
			<label for="derived-name">Name</label>
			<input id="derived-name" type="text" bind:value={form.name} {disabled} required />
		</div>

		<TE_FormulaBuilder
			bind:value={form.config.derivation_formula}
			label="Derivation Formula"
			{systemId}
			variables={derivationVariables}
			{disabled}
		/>
		<p class="hint">{DERIVATION_FORMULA_HELPER}</p>

		<fieldset class="checkbox-group">
			<legend>Recalculate Triggers</legend>
			{#each RECALCULATE_TRIGGER_OPTIONS as opt}
				<label class="checkbox-row">
					<input
						type="checkbox"
						checked={(form.config.recalculate_triggers ?? []).includes(opt.value)}
						{disabled}
						onchange={(e) =>
							toggleTrigger(opt.value, (e.currentTarget as HTMLInputElement).checked)}
					/>
					<span>{opt.label}</span>
				</label>
			{/each}
		</fieldset>

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
	.derived-form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.form-field label,
	.checkbox-group legend {
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

	.checkbox-group {
		border: none;
		margin: 0;
		padding: 0;
	}

	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		font-weight: 500;
		margin-bottom: 0.35rem;
		cursor: pointer;
	}

	.hint {
		margin: 0;
		font-size: 0.8rem;
		color: var(--text-muted, #6b7280);
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

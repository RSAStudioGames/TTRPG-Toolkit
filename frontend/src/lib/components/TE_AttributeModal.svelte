<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import TE_FormulaBuilder from './TE_FormulaBuilder.svelte';
	import { createAttribute, updateAttribute } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		ATTRIBUTE_TYPE_CUSTOM,
		ATTRIBUTE_TYPE_DESCRIPTIVE,
		ATTRIBUTE_TYPE_NUMERIC,
		ATTRIBUTE_TYPE_RANK_TIER,
		ATTRIBUTE_TYPE_STEP_DIE,
		ATTRIBUTE_TYPE_OPTIONS,
		CACHING_RULE_LIVE,
		CACHING_RULE_ON_TRIGGER,
		CACHING_RULE_OPTIONS,
		MODIFIER_DISPLAY_OPTIONS,
		MODIFIER_FORMULA_HELPER,
		NUMERIC_FORMAT_OPTIONS,
		RECALCULATE_TRIGGER_OPTIONS,
		STEP_DIE_OPTIONS
	} from '$lib/constants/attributeOptions';
	import {
		attributeToForm,
		defaultAttributeConfig,
		defaultAttributeForm,
		derivationFormulaVariables,
		descendantAttributeIds,
		formToCreatePayload,
		formToUpdatePayload,
		validateAttributeForm,
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

	let form = $state<AttributeFormState>(defaultAttributeForm());
	let saving = $state(false);
	let error = $state<string | null>(null);

	const isEdit = $derived(attribute != null);
	const modalTitle = $derived(isEdit ? 'Edit Attribute' : 'Create Attribute');

	const derivationVariables = $derived(
		derivationFormulaVariables(siblingAttributes, attribute?.id)
	);

	const excludedParentIds = $derived.by(() => {
		if (!attribute) return new Set<string>();
		const d = descendantAttributeIds(siblingAttributes, attribute.id);
		d.add(attribute.id);
		return d;
	});

	const parentOptions = $derived(
		siblingAttributes.filter((a) => !excludedParentIds.has(a.id))
	);

	const showNumeric = $derived(form.type === ATTRIBUTE_TYPE_NUMERIC);
	const showStepDie = $derived(form.type === ATTRIBUTE_TYPE_STEP_DIE);
	const showDescriptive = $derived(form.type === ATTRIBUTE_TYPE_DESCRIPTIVE);
	const showRankTier = $derived(form.type === ATTRIBUTE_TYPE_RANK_TIER);
	const showCustom = $derived(form.type === ATTRIBUTE_TYPE_CUSTOM);
	const showDerived = $derived(form.config.is_derived);
	const showTriggers = $derived(
		showDerived && form.config.caching_rule === CACHING_RULE_ON_TRIGGER
	);

	const minRowsForType = $derived(
		showStepDie || showDescriptive || showRankTier ? 1 : 0
	);

	function resetForm() {
		form = attribute ? attributeToForm(attribute) : defaultAttributeForm();
		error = null;
	}

	$effect(() => {
		if (open) resetForm();
	});

	$effect(() => {
		if (form.config.is_derived && !form.config.caching_rule) {
			form.config.caching_rule = CACHING_RULE_LIVE;
		}
	});

	function onTypeChange(next: string) {
		form.type = next;
		form.config = defaultAttributeConfig(next);
	}

	function toggleTrigger(value: string, checked: boolean) {
		const current = form.config.recalculate_triggers ?? [];
		if (checked) {
			form.config.recalculate_triggers = [...current, value];
		} else {
			form.config.recalculate_triggers = current.filter((t) => t !== value);
		}
	}

	function addStepDie() {
		form.config.step_dice = [...(form.config.step_dice ?? []), 'd6'];
	}

	function removeStepDie(index: number) {
		const list = form.config.step_dice ?? [];
		if (list.length <= minRowsForType) return;
		form.config.step_dice = list.filter((_, i) => i !== index);
	}

	function addDescriptive() {
		form.config.descriptive_map = [
			...(form.config.descriptive_map ?? []),
			{ label: '', value: 0 }
		];
	}

	function removeDescriptive(index: number) {
		const list = form.config.descriptive_map ?? [];
		if (list.length <= minRowsForType) return;
		form.config.descriptive_map = list.filter((_, i) => i !== index);
	}

	function addRank() {
		form.config.rank_map = [
			...(form.config.rank_map ?? []),
			{ rank_name: '', numeric_backing: 0 }
		];
	}

	function removeRank(index: number) {
		const list = form.config.rank_map ?? [];
		if (list.length <= minRowsForType) return;
		form.config.rank_map = list.filter((_, i) => i !== index);
	}

	async function handleSave() {
		const validation = validateAttributeForm(form);
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
				await createAttribute(
					systemId,
					formToCreatePayload(form, siblingAttributes.length)
				);
			}
			onclose();
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save attribute';
		} finally {
			saving = false;
		}
	}
</script>

<BaseModal {open} title={modalTitle} wide onclose={onclose}>
	<div class="attr-form">
		<div class="form-field">
			<label for="attr-name">Name</label>
			<input id="attr-name" type="text" bind:value={form.name} {disabled} required minlength="3" />
		</div>

		<div class="form-field">
			<label for="attr-group">Group Name</label>
			<input id="attr-group" type="text" bind:value={form.group_name} {disabled} />
		</div>

		<div class="form-field">
			<label for="attr-type">Type</label>
			<select
				id="attr-type"
				value={form.type}
				{disabled}
				onchange={(e) => onTypeChange((e.currentTarget as HTMLSelectElement).value)}
			>
				{#each ATTRIBUTE_TYPE_OPTIONS as opt}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>
		</div>

		{#if showNumeric}
			<div class="form-row">
				<div class="form-field">
					<label for="attr-min">Default Minimum</label>
					<input
						id="attr-min"
						type="number"
						bind:value={form.config.min}
						{disabled}
					/>
				</div>
				<div class="form-field">
					<label for="attr-max">Default Maximum</label>
					<input
						id="attr-max"
						type="number"
						bind:value={form.config.max}
						{disabled}
					/>
				</div>
			</div>
			<div class="form-field">
				<label for="attr-numeric-format">Numeric format</label>
				<select id="attr-numeric-format" bind:value={form.config.numeric_format} {disabled}>
					{#each NUMERIC_FORMAT_OPTIONS as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>
		{/if}

		{#if showStepDie}
			<div class="dynamic-section">
				<h4>Step dice</h4>
				{#each form.config.step_dice ?? [] as die, i}
					<div class="dynamic-row">
						<select bind:value={form.config.step_dice![i]} {disabled}>
							{#each STEP_DIE_OPTIONS as opt}
								<option value={opt.value}>{opt.label}</option>
							{/each}
						</select>
						<button
							type="button"
							class="btn-sm btn-danger-outline"
							disabled={disabled || (form.config.step_dice?.length ?? 0) <= minRowsForType}
							title={(form.config.step_dice?.length ?? 0) <= minRowsForType
								? 'At least one die is required'
								: undefined}
							aria-label="Delete row"
							onclick={() => removeStepDie(i)}
						>
							Delete
						</button>
					</div>
				{/each}
				<button type="button" class="btn-secondary-sm" {disabled} onclick={addStepDie}>
					Create Die
				</button>
			</div>
		{/if}

		{#if showDescriptive}
			<div class="dynamic-section">
				<h4>Descriptive mappings</h4>
				{#each form.config.descriptive_map ?? [] as entry, i}
					<div class="dynamic-row dynamic-row--wide">
						<input
							type="text"
							placeholder="Label"
							bind:value={form.config.descriptive_map![i].label}
							{disabled}
						/>
						<input
							type="number"
							placeholder="Numeric Value"
							bind:value={form.config.descriptive_map![i].value}
							{disabled}
						/>
						<button
							type="button"
							class="btn-sm btn-danger-outline"
							disabled={disabled || (form.config.descriptive_map?.length ?? 0) <= minRowsForType}
							title={(form.config.descriptive_map?.length ?? 0) <= minRowsForType
								? 'At least one mapping is required'
								: undefined}
							aria-label="Delete row"
							onclick={() => removeDescriptive(i)}
						>
							Delete
						</button>
					</div>
				{/each}
				<button type="button" class="btn-secondary-sm" {disabled} onclick={addDescriptive}>
					Create Mapping
				</button>
			</div>
		{/if}

		{#if showRankTier}
			<div class="dynamic-section">
				<h4>Rank map</h4>
				{#each form.config.rank_map ?? [] as entry, i}
					<div class="dynamic-row dynamic-row--wide">
						<input
							type="text"
							placeholder="Rank name"
							bind:value={form.config.rank_map![i].rank_name}
							{disabled}
						/>
						<input
							type="number"
							placeholder="Numeric backing"
							bind:value={form.config.rank_map![i].numeric_backing}
							{disabled}
						/>
						<button
							type="button"
							class="btn-sm btn-danger-outline"
							disabled={disabled || (form.config.rank_map?.length ?? 0) <= minRowsForType}
							title={(form.config.rank_map?.length ?? 0) <= minRowsForType
								? 'At least one rank is required'
								: undefined}
							aria-label="Delete row"
							onclick={() => removeRank(i)}
						>
							Delete
						</button>
					</div>
				{/each}
				<button type="button" class="btn-secondary-sm" {disabled} onclick={addRank}>
					Create Rank
				</button>
			</div>
		{/if}

		{#if showCustom}
			<p class="hint">Custom attributes use minimal configuration. Add modifier or derived rules below if needed.</p>
		{/if}

		<TE_FormulaBuilder
			bind:value={form.config.modifier_formula}
			label="Modifier Formula"
			{systemId}
			variables={['score']}
			{disabled}
		/>
		<p class="hint">{MODIFIER_FORMULA_HELPER}</p>

		<fieldset class="radio-group">
			<legend>Modifier Display</legend>
			{#each MODIFIER_DISPLAY_OPTIONS as opt}
				<label class="radio-option">
					<input
						type="radio"
						name="modifier-display"
						value={opt.value}
						checked={form.config.modifier_display === opt.value}
						{disabled}
						onchange={() => {
							form.config.modifier_display = opt.value;
						}}
					/>
					<span>{opt.label}</span>
				</label>
			{/each}
		</fieldset>

		<label class="checkbox-row">
			<input type="checkbox" bind:checked={form.config.is_derived} {disabled} />
			<span>Enable Derived Attributes</span>
		</label>

		{#if showDerived}
			<TE_FormulaBuilder
				bind:value={form.config.derivation_formula}
				label="Derivation Formula"
				{systemId}
				variables={derivationVariables}
				{disabled}
			/>

			<fieldset class="radio-group">
				<legend>Caching</legend>
				{#each CACHING_RULE_OPTIONS as opt}
					<label class="radio-option">
						<input
							type="radio"
							name="caching-rule"
							value={opt.value}
							checked={form.config.caching_rule === opt.value}
							{disabled}
							onchange={() => {
								form.config.caching_rule = opt.value;
							}}
						/>
						<span>{opt.label}</span>
					</label>
				{/each}
			</fieldset>

			{#if showTriggers}
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
			{/if}
		{/if}

		<div class="form-field">
			<label for="attr-parent">Parent attribute</label>
			<select id="attr-parent" bind:value={form.parent_attribute_id} {disabled}>
				<option value="">None</option>
				{#each parentOptions as p}
					<option value={p.id}>{p.name}</option>
				{/each}
			</select>
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
	.attr-form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.form-field label,
	.radio-group legend,
	.checkbox-group legend {
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

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.75rem;
	}

	.dynamic-section h4 {
		margin: 0 0 0.5rem;
		font-size: 0.9rem;
	}

	.dynamic-row {
		display: flex;
		gap: 0.5rem;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.dynamic-row select,
	.dynamic-row input {
		flex: 1;
		padding: 0.4rem 0.55rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.dynamic-row--wide {
		flex-wrap: wrap;
	}

	.hint {
		margin: 0;
		font-size: 0.8rem;
		color: var(--text-muted, #6b7280);
	}

	.radio-group,
	.checkbox-group {
		border: none;
		margin: 0;
		padding: 0;
	}

	.radio-option,
	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		margin-bottom: 0.35rem;
		cursor: pointer;
	}

	.checkbox-row {
		font-weight: 600;
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

	.btn-secondary-sm {
		padding: 0.35rem 0.75rem;
		font-size: 0.85rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.btn-sm {
		padding: 0.25rem 0.55rem;
		font-size: 0.8rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
		flex-shrink: 0;
	}

	.btn-danger-outline {
		border-color: #fca5a5;
		color: #b91c1c;
	}
</style>

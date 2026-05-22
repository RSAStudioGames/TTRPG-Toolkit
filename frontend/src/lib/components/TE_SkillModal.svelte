<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import { createSkill, updateSkill } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		SKILL_CATEGORY_CUSTOM,
		SKILL_CATEGORY_OPTIONS,
		SKILL_TYPE_MULTI_TIER,
		SKILL_TYPE_NUMERIC,
		SKILL_TYPE_OPTIONS,
		SPECIALIZATION_EXAMPLES_HELPER
	} from '$lib/constants/skillOptions';
	import {
		defaultSkillConfig,
		defaultSkillForm,
		formToCreatePayload,
		formToUpdatePayload,
		skillToForm,
		validateSkillForm,
		type SkillFormState
	} from '$lib/utils/skillDefaults';
	import type { AttributeResponse, SkillResponse } from '$lib/types/mechanics';

	interface Props {
		open: boolean;
		systemId: string;
		skill?: SkillResponse | null;
		attributes: AttributeResponse[];
		siblingSkillCount: number;
		disabled?: boolean;
		onclose: () => void;
		onsaved?: () => void;
	}

	let {
		open,
		systemId,
		skill = null,
		attributes,
		siblingSkillCount,
		disabled = false,
		onclose,
		onsaved
	}: Props = $props();

	let form = $state<SkillFormState>(defaultSkillForm());
	let saving = $state(false);
	let error = $state<string | null>(null);

	const isEdit = $derived(skill != null);
	const modalTitle = $derived(isEdit ? 'Edit Skill' : 'Create Skill');
	const showMultiTier = $derived(form.type === SKILL_TYPE_MULTI_TIER);
	const showNumeric = $derived(form.type === SKILL_TYPE_NUMERIC);
	const showCustomCategory = $derived(form.category_preset === SKILL_CATEGORY_CUSTOM);
	const showSpecializationBonus = $derived(form.config.allow_specializations);
	const minTierRows = $derived(showMultiTier ? 1 : 0);

	function resetForm() {
		form = skill ? skillToForm(skill) : defaultSkillForm();
		error = null;
	}

	$effect(() => {
		if (open) resetForm();
	});

	function onTypeChange(next: string) {
		form.type = next;
		form.config = defaultSkillConfig(next);
	}

	function addTier() {
		form.config.tiers = [...(form.config.tiers ?? []), { tier_name: '', numeric_backing: 0 }];
	}

	function removeTier(index: number) {
		const list = form.config.tiers ?? [];
		if (list.length <= minTierRows) return;
		form.config.tiers = list.filter((_, i) => i !== index);
	}

	async function handleSave() {
		const validation = validateSkillForm(form);
		if (validation) {
			error = validation;
			return;
		}
		saving = true;
		error = null;
		try {
			if (isEdit && skill) {
				await updateSkill(systemId, skill.id, formToUpdatePayload(form));
			} else {
				await createSkill(systemId, formToCreatePayload(form, siblingSkillCount));
			}
			onclose();
			onsaved?.();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save skill';
		} finally {
			saving = false;
		}
	}
</script>

<BaseModal {open} title={modalTitle} onclose={onclose}>
	<div class="skill-form">
		<div class="form-field">
			<label for="skill-name">Name</label>
			<input id="skill-name" type="text" bind:value={form.name} {disabled} required />
		</div>

		<div class="form-field">
			<label for="skill-linked-attr">Linked Attribute</label>
			<select id="skill-linked-attr" bind:value={form.linked_attribute_id} {disabled}>
				<option value="">None</option>
				{#each attributes as attr}
					<option value={attr.id}>{attr.name}</option>
				{/each}
			</select>
		</div>

		<div class="form-field">
			<label for="skill-type">Rating Type</label>
			<select
				id="skill-type"
				value={form.type}
				{disabled}
				onchange={(e) => onTypeChange((e.currentTarget as HTMLSelectElement).value)}
			>
				{#each SKILL_TYPE_OPTIONS as opt}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>
		</div>

		{#if showMultiTier}
			<div class="dynamic-section">
				<h4>Tiers</h4>
				{#each form.config.tiers ?? [] as tier, i}
					<div class="dynamic-row">
						<input
							type="text"
							placeholder="Tier Name"
							bind:value={form.config.tiers![i].tier_name}
							{disabled}
						/>
						<input
							type="number"
							placeholder="Numeric Backing"
							bind:value={form.config.tiers![i].numeric_backing}
							{disabled}
						/>
						<button
							type="button"
							class="btn-sm btn-danger-outline"
							disabled={disabled || (form.config.tiers?.length ?? 0) <= minTierRows}
							title={(form.config.tiers?.length ?? 0) <= minTierRows
								? 'At least one tier is required'
								: undefined}
							aria-label="Delete row"
							onclick={() => removeTier(i)}
						>
							Delete
						</button>
					</div>
				{/each}
				<button type="button" class="btn-secondary-sm" {disabled} onclick={addTier}>
					Create Tier
				</button>
			</div>
		{/if}

		{#if showNumeric}
			<div class="form-row">
				<div class="form-field">
					<label for="skill-min">Min Bound</label>
					<input id="skill-min" type="number" bind:value={form.config.min} {disabled} />
				</div>
				<div class="form-field">
					<label for="skill-max">Max Bound</label>
					<input id="skill-max" type="number" bind:value={form.config.max} {disabled} />
				</div>
			</div>
		{/if}

		<div class="form-field">
			<label for="skill-category">Category</label>
			<select id="skill-category" bind:value={form.category_preset} {disabled}>
				{#each SKILL_CATEGORY_OPTIONS as opt}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>
		</div>

		{#if showCustomCategory}
			<div class="form-field">
				<label for="skill-category-custom">Custom Category</label>
				<input
					id="skill-category-custom"
					type="text"
					bind:value={form.category_custom}
					{disabled}
					placeholder="Enter category name"
				/>
			</div>
		{/if}

		<label class="checkbox-row">
			<input type="checkbox" bind:checked={form.config.allow_specializations} {disabled} />
			<span>Allow Specializations</span>
		</label>

		{#if showSpecializationBonus}
			<div class="form-field">
				<label for="skill-spec-bonus">Specialization Bonus</label>
				<input
					id="skill-spec-bonus"
					type="number"
					bind:value={form.config.specialization_bonus}
					{disabled}
				/>
			</div>
			<div class="form-field">
				<label for="skill-spec-examples">Specialization Examples</label>
				<textarea
					id="skill-spec-examples"
					rows="3"
					{disabled}
					readonly
					placeholder={SPECIALIZATION_EXAMPLES_HELPER}
				></textarea>
				<p class="hint">{SPECIALIZATION_EXAMPLES_HELPER}</p>
			</div>
		{/if}

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
	.skill-form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.form-field label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		margin-bottom: 0.35rem;
	}

	.form-field input,
	.form-field select,
	.form-field textarea {
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

	.dynamic-row input {
		flex: 1;
		padding: 0.4rem 0.55rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
	}

	.hint {
		margin: 0.35rem 0 0;
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

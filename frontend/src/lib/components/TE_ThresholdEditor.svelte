<script lang="ts">
	import { LADDER_OPERATOR_OPTIONS } from '$lib/constants/resolutionOptions';
	import type { LadderTier } from '$lib/types/mechanics';

	interface Props {
		tiers: LadderTier[];
		disabled?: boolean;
	}

	let { tiers = $bindable([]), disabled = false }: Props = $props();

	function addTier() {
		tiers = [...tiers, { label: '', operator: '>=', value: 0 }];
	}

	function removeTier(index: number) {
		tiers = tiers.filter((_, i) => i !== index);
	}
</script>

<div class="threshold-editor">
	<h4>Threshold Ladder</h4>
	{#each tiers as tier, i (i)}
		<div class="ladder-row">
			<div class="form-field">
				<label for="tier-label-{i}">Threshold Label</label>
				<input id="tier-label-{i}" type="text" bind:value={tier.label} {disabled} />
			</div>
			<div class="form-field">
				<label for="tier-op-{i}">Operator</label>
				<select id="tier-op-{i}" bind:value={tier.operator} {disabled}>
					{#each LADDER_OPERATOR_OPTIONS as op}
						<option value={op.value}>{op.label}</option>
					{/each}
				</select>
			</div>
			<div class="form-field">
				<label for="tier-val-{i}">Value</label>
				<input id="tier-val-{i}" type="number" bind:value={tier.value} {disabled} />
			</div>
			{#if !disabled}
				<button type="button" class="btn-row-remove" onclick={() => removeTier(i)}>Remove</button>
			{/if}
		</div>
	{/each}
	{#if !disabled}
		<button type="button" class="btn-add" onclick={addTier}>Add Tier</button>
	{/if}
</div>

<style>
	.threshold-editor h4 {
		margin: 0 0 0.65rem;
		font-size: 0.95rem;
		font-weight: 600;
	}

	.form-field {
		margin-bottom: 0;
	}

	.form-field label {
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

	.ladder-row {
		display: grid;
		grid-template-columns: 1fr auto auto auto;
		gap: 0.75rem;
		align-items: end;
		margin-bottom: 0.75rem;
		padding: 0.75rem;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		background: #fafafa;
	}

	@media (max-width: 640px) {
		.ladder-row {
			grid-template-columns: 1fr;
		}
	}

	.btn-add,
	.btn-row-remove {
		padding: 0.35rem 0.75rem;
		font-size: 0.85rem;
		border-radius: 6px;
		cursor: pointer;
		font: inherit;
	}

	.btn-add {
		border: 1px dashed #9ca3af;
		background: #fff;
		color: #374151;
	}

	.btn-row-remove {
		border: 1px solid #e5e7eb;
		background: #fff;
		color: #6b7280;
		margin-bottom: 0.35rem;
	}
</style>

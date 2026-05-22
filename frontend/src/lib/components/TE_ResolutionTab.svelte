<script lang="ts">
	import { onMount } from 'svelte';
	import TE_FormulaBuilder from './TE_FormulaBuilder.svelte';
	import { getMechanics, saveResolutionConfig } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		CRIT_FAILURE_TRIGGER_OPTIONS,
		CRIT_SUCCESS_TRIGGER_OPTIONS,
		CRIT_TRIGGER_EXCEED_DC,
		LADDER_OPERATOR_OPTIONS,
		MECHANIC_TYPE_OPTIONS,
		RESOLUTION_TYPE_CUSTOM,
		RESOLUTION_TYPE_OPTIONS,
		SUCCESS_METHOD_OPTIONS,
		SUCCESS_METHOD_TARGET_NUMBER,
		SUCCESS_METHOD_THRESHOLD_LADDER
	} from '$lib/constants/resolutionOptions';
	import {
		defaultResolutionConfig,
		normalizeResolutionConfig
	} from '$lib/utils/resolutionDefaults';
	import type { ResolutionConfig } from '$lib/types/mechanics';

	interface Props {
		systemId: string;
		disabled?: boolean;
	}

	let { systemId, disabled = false }: Props = $props();

	let config = $state<ResolutionConfig>(defaultResolutionConfig());
	let loading = $state(true);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let saveMessage = $state<string | null>(null);
	let validationErrors = $state<string[]>([]);

	const formulaVariables = $derived.by(() => {
		const names = new Set<string>(['modifier', 'target']);
		const dv = config.success_determination.default_target_variable?.trim();
		if (dv) {
			const cleaned = dv.replace(/^\{|\}$/g, '').trim();
			if (cleaned) names.add(cleaned);
		}
		const matches = config.roll_expression.matchAll(/\{([^}]+)\}/g);
		for (const m of matches) {
			const n = m[1]?.trim();
			if (n) names.add(n);
		}
		return [...names];
	});

	const showCustomName = $derived(config.resolution_type === RESOLUTION_TYPE_CUSTOM);
	const showLadder = $derived(
		config.success_determination.method === SUCCESS_METHOD_THRESHOLD_LADDER
	);
	const showTargetVariable = $derived(
		config.success_determination.method === SUCCESS_METHOD_TARGET_NUMBER
	);
	const showCritSuccessExceed = $derived(
		config.critical_mechanics.crit_success_trigger === CRIT_TRIGGER_EXCEED_DC
	);

	async function load() {
		loading = true;
		error = null;
		saveMessage = null;
		validationErrors = [];
		try {
			const data = await getMechanics(systemId);
			config = normalizeResolutionConfig(data.resolution_config);
		} catch (e) {
			if (e instanceof ApiError && e.status === 404) {
				config = defaultResolutionConfig();
			} else {
				error = e instanceof ApiError ? e.message : 'Failed to load resolution settings';
			}
		} finally {
			loading = false;
		}
	}

	async function handleSave() {
		saving = true;
		error = null;
		saveMessage = null;
		validationErrors = [];
		try {
			const data = await saveResolutionConfig(systemId, config);
			config = normalizeResolutionConfig(data.resolution_config);
			saveMessage = 'Resolution settings saved.';
		} catch (e) {
			if (e instanceof ApiError) {
				error = e.message;
				validationErrors = e.errors;
			} else {
				error = 'Failed to save resolution settings';
			}
		} finally {
			saving = false;
		}
	}

	function addLadderTier() {
		config.success_determination.threshold_ladder = [
			...config.success_determination.threshold_ladder,
			{ label: '', operator: '>=', value: 0 }
		];
	}

	function removeLadderTier(index: number) {
		config.success_determination.threshold_ladder =
			config.success_determination.threshold_ladder.filter((_, i) => i !== index);
	}

	function addMechanic() {
		config.advantage_disadvantage = [
			...config.advantage_disadvantage,
			{ name: '', mechanic_type: 'add_modifier' }
		];
	}

	function removeMechanic(index: number) {
		config.advantage_disadvantage = config.advantage_disadvantage.filter((_, i) => i !== index);
	}

	onMount(() => {
		load();
	});
</script>

<div class="resolution-tab">
	{#if loading}
		<p class="status-msg">Loading resolution settings…</p>
	{:else}
		<section class="form-section">
			<h3>Resolution Paradigm</h3>
			<fieldset class="radio-paradigm" {disabled}>
				{#each RESOLUTION_TYPE_OPTIONS as opt}
					<label class="radio-option">
						<input
							type="radio"
							name="resolution-type"
							value={opt.value}
							checked={config.resolution_type === opt.value}
							{disabled}
							onchange={() => {
								config.resolution_type = opt.value;
							}}
						/>
						<span>{opt.label}</span>
					</label>
				{/each}
			</fieldset>
			{#if showCustomName}
				<div class="form-field">
					<label for="custom-paradigm-name">Custom Paradigm Name</label>
					<input
						id="custom-paradigm-name"
						type="text"
						bind:value={config.custom_paradigm_name}
						{disabled}
						required
					/>
				</div>
			{/if}
		</section>

		<section class="form-section">
			<h3>Roll Syntax Engine</h3>
			<TE_FormulaBuilder
				bind:value={config.roll_expression}
				label="Roll Expression"
				{systemId}
				variables={formulaVariables}
				{disabled}
			/>
		</section>

		<section class="form-section">
			<h3>Success Determination</h3>
			<div class="form-field">
				<label for="success-method">Method</label>
				<select id="success-method" bind:value={config.success_determination.method} {disabled}>
					{#each SUCCESS_METHOD_OPTIONS as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>

			{#if showTargetVariable}
				<div class="form-field">
					<label for="default-target-var">Default Target Number Variable</label>
					<input
						id="default-target-var"
						type="text"
						placeholder="e.g., target"
						bind:value={config.success_determination.default_target_variable}
						{disabled}
					/>
				</div>
			{/if}

			{#if showLadder}
				<div class="ladder-list">
					<h4>Threshold Ladder</h4>
					{#each config.success_determination.threshold_ladder as tier, i (i)}
						<div class="ladder-row">
							<div class="form-field">
								<label for="tier-label-{i}">Threshold Label</label>
								<input
									id="tier-label-{i}"
									type="text"
									bind:value={tier.label}
									{disabled}
								/>
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
								<input
									id="tier-val-{i}"
									type="number"
									bind:value={tier.value}
									{disabled}
								/>
							</div>
							{#if !disabled}
								<button type="button" class="btn-row-remove" onclick={() => removeLadderTier(i)}>
									Remove
								</button>
							{/if}
						</div>
					{/each}
					{#if !disabled}
						<button type="button" class="btn-add" onclick={addLadderTier}>Create Tier</button>
					{/if}
				</div>
			{/if}
		</section>

		<section class="form-section">
			<details class="crit-details">
				<summary>Critical Mechanics</summary>
				<div class="crit-block">
					<label class="checkbox-row">
						<input
							type="checkbox"
							bind:checked={config.critical_mechanics.enable_crit_success}
							{disabled}
						/>
						Enable Critical Success
					</label>
					{#if config.critical_mechanics.enable_crit_success}
						<div class="form-field">
							<label for="crit-success-trigger">Trigger</label>
							<select
								id="crit-success-trigger"
								bind:value={config.critical_mechanics.crit_success_trigger}
								{disabled}
							>
								<option value="">Select…</option>
								{#each CRIT_SUCCESS_TRIGGER_OPTIONS as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</select>
						</div>
						{#if showCritSuccessExceed}
							<div class="form-field">
								<label for="crit-exceed-amount">Exceed Amount</label>
								<input
									id="crit-exceed-amount"
									type="number"
									min="0"
									bind:value={config.critical_mechanics.crit_success_exceed_amount}
									{disabled}
								/>
							</div>
						{/if}
					{/if}

					<label class="checkbox-row">
						<input
							type="checkbox"
							bind:checked={config.critical_mechanics.enable_crit_failure}
							{disabled}
						/>
						Enable Critical Failure
					</label>
					{#if config.critical_mechanics.enable_crit_failure}
						<div class="form-field">
							<label for="crit-failure-trigger">Trigger</label>
							<select
								id="crit-failure-trigger"
								bind:value={config.critical_mechanics.crit_failure_trigger}
								{disabled}
							>
								<option value="">Select…</option>
								{#each CRIT_FAILURE_TRIGGER_OPTIONS as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</select>
						</div>
					{/if}
				</div>
			</details>
		</section>

		<section class="form-section">
			<h3>Advantage / Disadvantage</h3>
			<div class="mechanics-table-wrap">
				<table class="mechanics-table">
					<thead>
						<tr>
							<th scope="col">Name</th>
							<th scope="col">Mechanic Type</th>
							{#if !disabled}
								<th scope="col"><span class="sr-only">Actions</span></th>
							{/if}
						</tr>
					</thead>
					<tbody>
						{#each config.advantage_disadvantage as entry, i (i)}
							<tr>
								<td>
									<input
										type="text"
										bind:value={entry.name}
										placeholder="Name"
										{disabled}
										aria-label="Mechanic name"
									/>
								</td>
								<td>
									<select bind:value={entry.mechanic_type} {disabled} aria-label="Mechanic type">
										{#each MECHANIC_TYPE_OPTIONS as opt}
											<option value={opt.value}>{opt.label}</option>
										{/each}
									</select>
								</td>
								{#if !disabled}
									<td>
										<button
											type="button"
											class="btn-row-remove"
											onclick={() => removeMechanic(i)}
										>
											Remove
										</button>
									</td>
								{/if}
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			{#if !disabled}
				<button type="button" class="btn-add" onclick={addMechanic}>Create Mechanic</button>
			{/if}
		</section>

		{#if error}
			<p class="form-error">{error}</p>
			{#if validationErrors.length > 0}
				<ul class="error-list">
					{#each validationErrors as msg}
						<li>{msg}</li>
					{/each}
				</ul>
			{/if}
		{/if}
		{#if saveMessage}
			<p class="save-success">{saveMessage}</p>
		{/if}

		<div class="save-row">
			<button
				type="button"
				class="btn-primary"
				disabled={disabled || saving}
				onclick={handleSave}
			>
				{saving ? 'Saving…' : 'Save Resolution'}
			</button>
		</div>
	{/if}
</div>

<style>
	.resolution-tab {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-section {
		margin-bottom: 1.25rem;
	}

	.form-section h3,
	.form-section h4 {
		margin: 0 0 0.65rem;
		font-size: 0.95rem;
		font-weight: 600;
	}

	.form-field {
		margin-bottom: 0.75rem;
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
		max-width: 28rem;
		padding: 0.45rem 0.6rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.radio-paradigm {
		border: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
	}

	.radio-option {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		cursor: pointer;
	}

	.radio-option input {
		accent-color: var(--accent-gm, #c9a227);
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

	.crit-details {
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		padding: 0.5rem 0.75rem;
	}

	.crit-details summary {
		font-weight: 600;
		cursor: pointer;
		padding: 0.35rem 0;
	}

	.crit-block {
		padding-top: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		font-weight: 600;
	}

	.mechanics-table-wrap {
		overflow-x: auto;
		margin-bottom: 0.75rem;
	}

	.mechanics-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.9rem;
	}

	.mechanics-table th,
	.mechanics-table td {
		padding: 0.5rem;
		border-bottom: 1px solid #e5e7eb;
		text-align: left;
		vertical-align: middle;
	}

	.mechanics-table input,
	.mechanics-table select {
		width: 100%;
		min-width: 8rem;
		padding: 0.4rem 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.btn-add,
	.btn-row-remove {
		padding: 0.35rem 0.75rem;
		border-radius: 6px;
		font-size: 0.85rem;
		font: inherit;
		cursor: pointer;
	}

	.btn-add {
		border: 1px solid var(--accent-gm, #c9a227);
		background: var(--accent-gm-muted, #f5ecd4);
		font-weight: 600;
	}

	.btn-row-remove {
		border: 1px solid #d1d5db;
		background: #fff;
	}

	.btn-primary {
		padding: 0.5rem 1.25rem;
		border: none;
		border-radius: 6px;
		background: var(--accent-gm-muted, #f5ecd4);
		font-weight: 600;
		font: inherit;
		cursor: pointer;
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.save-row {
		margin-top: 0.5rem;
		padding-top: 0.75rem;
		border-top: 1px solid #e5e7eb;
	}

	.status-msg {
		color: #6b7280;
	}

	.form-error {
		color: #b91c1c;
		font-size: 0.85rem;
	}

	.error-list {
		margin: 0.25rem 0 0;
		padding-left: 1.25rem;
		color: #b91c1c;
		font-size: 0.85rem;
	}

	.save-success {
		color: #15803d;
		font-size: 0.85rem;
	}

	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}
</style>

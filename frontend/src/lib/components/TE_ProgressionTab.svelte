<script lang="ts">
	import { onMount } from 'svelte';
	import TE_TableEditor from './TE_TableEditor.svelte';
	import { getMechanics, saveProgressionConfig } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		COST_TABLE_COLUMNS,
		LABEL_ALLOW_MILESTONE,
		LABEL_ALLOW_UNDO,
		LABEL_GM_APPROVAL,
		PROGRESSION_PARADIGM_OPTIONS,
		XP_TABLE_COLUMNS
	} from '$lib/constants/progressionOptions';
	import {
		defaultProgressionConfig,
		normalizeProgressionConfig,
		showLevelBasedBlock,
		showPointBuyBlock
	} from '$lib/utils/progressionDefaults';
	import type { ProgressionConfig } from '$lib/types/mechanics';

	interface Props {
		systemId: string;
		disabled?: boolean;
		onDirtyChange?: (dirty: boolean) => void;
	}

	let { systemId, disabled = false, onDirtyChange }: Props = $props();

	let config = $state<ProgressionConfig>(defaultProgressionConfig());
	let loading = $state(true);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let saveMessage = $state<string | null>(null);
	let savedSnapshot = $state('');

	function snapshot(): string {
		syncConfigFromTables();
		return JSON.stringify(config);
	}

	function markDirty() {
		onDirtyChange?.(snapshot() !== savedSnapshot);
	}

	const showLevel = $derived(showLevelBasedBlock(config.paradigm));
	const showPointBuy = $derived(showPointBuyBlock(config.paradigm));

	let xpRows = $state<Record<string, string | number>[]>([]);
	let costRows = $state<Record<string, string | number>[]>([]);

	function syncTablesFromConfig() {
		xpRows = (config.level_based?.xp_table ?? []).map((r) => ({
			level: r.level,
			xp_required: r.xp_required
		}));
		costRows = (config.point_buy?.cost_table ?? []).map((r) => ({
			rating: r.rating,
			cost: r.cost
		}));
	}

	function syncConfigFromTables() {
		if (!config.level_based) config.level_based = defaultProgressionConfig().level_based!;
		if (!config.point_buy) config.point_buy = defaultProgressionConfig().point_buy!;
		config.level_based.xp_table = xpRows.map((r) => ({
			level: Number(r.level) || 0,
			xp_required: Number(r.xp_required) || 0
		}));
		config.point_buy.cost_table = costRows.map((r) => ({
			rating: Number(r.rating) || 0,
			cost: Number(r.cost) || 0
		}));
	}

	async function load() {
		loading = true;
		error = null;
		saveMessage = null;
		try {
			const data = await getMechanics(systemId);
			config = normalizeProgressionConfig(data.progression_config);
			syncTablesFromConfig();
			savedSnapshot = snapshot();
			onDirtyChange?.(false);
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to load progression settings';
		} finally {
			loading = false;
		}
	}

	async function handleSave() {
		syncConfigFromTables();
		saving = true;
		error = null;
		saveMessage = null;
		try {
			const data = await saveProgressionConfig(systemId, config);
			config = normalizeProgressionConfig(data.progression_config);
			syncTablesFromConfig();
			savedSnapshot = snapshot();
			onDirtyChange?.(false);
			saveMessage = 'Progression settings saved.';
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save progression settings';
		} finally {
			saving = false;
		}
	}

	function addXpRow() {
		xpRows = [...xpRows, { level: xpRows.length + 1, xp_required: 0 }];
	}

	function addCostRow() {
		costRows = [...costRows, { rating: 0, cost: 0 }];
	}

	$effect(() => {
		if (loading) return;
		markDirty();
	});

	onMount(() => {
		load();
	});
</script>

<div class="progression-tab">
	{#if loading}
		<p class="status-msg">Loading progression settings…</p>
	{:else}
		<section class="form-section">
			<h3>Paradigm</h3>
			<fieldset class="radio-paradigm" {disabled}>
				{#each PROGRESSION_PARADIGM_OPTIONS as opt}
					<label class="radio-option">
						<input
							type="radio"
							name="progression-paradigm"
							value={opt.value}
							checked={config.paradigm === opt.value}
							{disabled}
							onchange={() => {
								config.paradigm = opt.value;
								const d = defaultProgressionConfig();
								if (!config.level_based) config.level_based = { ...d.level_based! };
								if (!config.point_buy) config.point_buy = { ...d.point_buy! };
							}}
						/>
						<span>{opt.label}</span>
					</label>
				{/each}
			</fieldset>
		</section>

		{#if showLevel}
			<section class="form-section">
				<h3>Level-Based</h3>
				<div class="form-row">
					<div class="form-field">
						<label for="prog-min-level">Min Level</label>
						<input
							id="prog-min-level"
							type="number"
							bind:value={config.level_based!.min_level}
							{disabled}
						/>
					</div>
					<div class="form-field">
						<label for="prog-max-level">Max Level</label>
						<input
							id="prog-max-level"
							type="number"
							bind:value={config.level_based!.max_level}
							{disabled}
						/>
					</div>
				</div>
				<h4>XP Thresholds</h4>
				<TE_TableEditor
					columns={XP_TABLE_COLUMNS}
					bind:rows={xpRows}
					{disabled}
					addLabel="Add Row"
					minRows={1}
					onadd={addXpRow}
					onremove={(i) => {
						xpRows = xpRows.filter((_, idx) => idx !== i);
					}}
				/>
				<label class="checkbox-row">
					<input type="checkbox" bind:checked={config.level_based!.allow_milestone} {disabled} />
					<span>{LABEL_ALLOW_MILESTONE}</span>
				</label>
			</section>
		{/if}

		{#if showPointBuy}
			<section class="form-section">
				<h3>Point-Buy</h3>
				<div class="form-field">
					<label for="prog-starting-pool">Starting Pool</label>
					<input
						id="prog-starting-pool"
						type="number"
						bind:value={config.point_buy!.starting_pool}
						{disabled}
					/>
				</div>
				<h4>Cost Table</h4>
				<TE_TableEditor
					columns={COST_TABLE_COLUMNS}
					bind:rows={costRows}
					{disabled}
					addLabel="Add Row"
					minRows={1}
					onadd={addCostRow}
					onremove={(i) => {
						costRows = costRows.filter((_, idx) => idx !== i);
					}}
				/>
			</section>
		{/if}

		<section class="form-section">
			<h3>Universal Advancement Features</h3>
			<label class="checkbox-row">
				<input type="checkbox" bind:checked={config.gm_approval} {disabled} />
				<span>{LABEL_GM_APPROVAL}</span>
			</label>
			<label class="checkbox-row">
				<input type="checkbox" bind:checked={config.allow_undo} {disabled} />
				<span>{LABEL_ALLOW_UNDO}</span>
			</label>
		</section>

		{#if error}
			<p class="form-error">{error}</p>
		{/if}
		{#if saveMessage}
			<p class="status-msg success">{saveMessage}</p>
		{/if}

		<div class="tab-actions">
			<button type="button" class="btn-primary" disabled={disabled || saving} onclick={handleSave}>
				{saving ? 'Saving…' : 'Save Progression'}
			</button>
		</div>
	{/if}
</div>

<style>
	.progression-tab {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-section h3,
	.form-section h4 {
		margin: 0 0 0.65rem;
		font-size: 0.95rem;
		font-weight: 600;
	}

	.form-field label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		margin-bottom: 0.35rem;
	}

	.form-field input {
		width: 100%;
		max-width: 12rem;
		padding: 0.45rem 0.6rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.form-row {
		display: flex;
		gap: 1rem;
		flex-wrap: wrap;
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

	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		font-weight: 600;
		margin-top: 0.5rem;
		cursor: pointer;
	}

	.status-msg {
		margin: 0;
		color: var(--text-muted, #6b7280);
		font-size: 0.9rem;
	}

	.status-msg.success {
		color: #166534;
	}

	.form-error {
		color: #b91c1c;
		margin: 0;
	}

	.tab-actions {
		margin-top: 0.5rem;
	}

	.btn-primary {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		background: var(--accent-gm, #c9a227);
		color: #1a1a1a;
		font-weight: 600;
		cursor: pointer;
		font: inherit;
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>

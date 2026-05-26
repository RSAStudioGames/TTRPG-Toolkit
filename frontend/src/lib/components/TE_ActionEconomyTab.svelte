<script lang="ts">
	import TE_FormulaBuilder from './TE_FormulaBuilder.svelte';
	import TE_TableEditor from './TE_TableEditor.svelte';
	import { getMechanics, saveActionEconomyConfig } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		ACTION_COST_TABLE_COLUMNS,
		ACTION_SLOT_TABLE_COLUMNS,
		ALLOWANCE_UNLIMITED,
		COMBAT_TIME_TRACKING_OPTIONS,
		INITIATIVE_PERSISTENCE_OPTIONS,
		INITIATIVE_SYSTEM_OPTIONS,
		initiativeNeedsExpression,
		initiativeNeedsStaticValue,
		REFRESH_SCOPE_OPTIONS,
		ROUND_TIME_OPTIONS,
		SYSTEM_TYPE_OPTIONS,
		TIE_BREAKING_OPTIONS,
		TOKEN_REFRESH_OPTIONS,
		TURN_STRUCTURE_OPTIONS
	} from '$lib/constants/actionEconomyOptions';
	import {
		costRowsFromPool,
		defaultActionEconomyConfig,
		defaultActionSlotEntry,
		defaultComboEntry,
		normalizeActionEconomyConfig,
		normalizeActionSlot,
		poolCostsFromRows,
		slotToTableRow,
		tableRowToSlot
	} from '$lib/utils/actionEconomyDefaults';
	import type { ActionEconomyConfig, ActionSlotEntry } from '$lib/types/mechanics';

	interface Props {
		systemId: string;
		disabled?: boolean;
		onDirtyChange?: (dirty: boolean) => void;
	}

	let { systemId, disabled = false, onDirtyChange }: Props = $props();

	export async function save(): Promise<boolean> {
		return await handleSave();
	}

	export function cancel(): void {
		handleCancel();
	}

	let config = $state<ActionEconomyConfig>(defaultActionEconomyConfig());
	let actionSlots = $state<ActionSlotEntry[]>([]);
	let costRows = $state<Record<string, string | number>[]>([]);
	let slotRows = $state<Record<string, string | number>[]>([]);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let saveMessage = $state<string | null>(null);
	let savedSnapshot = $state('');

	const showCustomTurn = $derived(config.turn_structure === 'custom');
	const showTokenTurn = $derived(config.turn_structure === 'token_based');
	const showPointPool = $derived(config.system_type === 'point_pool');
	const showTypeSlots = $derived(config.system_type === 'type_slots');
	const showCustomRound = $derived(config.round_time_definition === 'custom');
	const showInitiativeExpr = $derived(initiativeNeedsExpression(config.initiative_system));
	const showStaticInit = $derived(initiativeNeedsStaticValue(config.initiative_system));

	const actionNameOptions = $derived(
		actionSlots.map((s) => s.name.trim()).filter((n) => n.length > 0)
	);

	function snapshot(): string {
		syncFromUi();
		return JSON.stringify({ config, actionSlots });
	}

	function markDirty() {
		onDirtyChange?.(snapshot() !== savedSnapshot);
	}

	function syncFromUi() {
		if (!config.point_pool) config.point_pool = defaultActionEconomyConfig().point_pool!;
		config.point_pool.action_cost_table = poolCostsFromRows(costRows);
		const slots: ActionSlotEntry[] = [];
		for (let i = 0; i < slotRows.length; i++) {
			slots.push(tableRowToSlot(slotRows[i], actionSlots[i]));
		}
		actionSlots = slots;
		config.action_slots = slots;
	}

	function syncToUi() {
		if (!config.point_pool) config.point_pool = defaultActionEconomyConfig().point_pool!;
		costRows = costRowsFromPool(config.point_pool);
		const slots = config.action_slots ?? [];
		actionSlots = slots.length > 0 ? slots.map((s) => normalizeActionSlot(s)) : [];
		slotRows = actionSlots.map(slotToTableRow);
	}

	function slotAt(index: number): ActionSlotEntry {
		if (!actionSlots[index]) {
			actionSlots[index] = defaultActionSlotEntry();
		}
		return actionSlots[index];
	}

	function setUnlimited(index: number, unlimited: boolean) {
		const slot = slotAt(index);
		if (unlimited) {
			slot.allowance = ALLOWANCE_UNLIMITED;
			slot.allowance_scope = '';
			slotRows[index] = { ...slotRows[index], _unlimited: 1, allowance: 1 };
		} else {
			slot.allowance = slot.allowance === ALLOWANCE_UNLIMITED ? 1 : Math.max(1, slot.allowance);
			slot.allowance_scope = slot.allowance_scope || 'per_turn';
			slotRows[index] = { ...slotRows[index], _unlimited: 0, allowance: slot.allowance };
		}
		actionSlots = [...actionSlots];
		markDirty();
	}

	function isUnlimited(index: number): boolean {
		return Number(slotRows[index]?._unlimited) === 1 || slotAt(index).allowance === ALLOWANCE_UNLIMITED;
	}

	async function load() {
		loading = true;
		error = null;
		saveMessage = null;
		try {
			const data = await getMechanics(systemId);
			config = normalizeActionEconomyConfig(data.action_economy_config);
			syncToUi();
			savedSnapshot = snapshot();
			onDirtyChange?.(false);
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to load action economy settings';
		} finally {
			loading = false;
		}
	}

	async function handleSave(): Promise<boolean> {
		syncFromUi();
		saving = true;
		error = null;
		saveMessage = null;
		const body: ActionEconomyConfig = {
			...config,
			action_slots: actionSlots,
			point_pool: config.point_pool
		};
		try {
			const data = await saveActionEconomyConfig(systemId, body);
			config = normalizeActionEconomyConfig(data.action_economy_config);
			syncToUi();
			savedSnapshot = snapshot();
			onDirtyChange?.(false);
			saveMessage = 'Action economy settings saved.';
			return true;
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save action economy settings';
			return false;
		} finally {
			saving = false;
		}
	}

	function handleCancel() {
		const isDirty = snapshot() !== savedSnapshot;
		if (isDirty && !confirm('Discard changes on this tab?')) return;
		if (savedSnapshot) {
			try {
				const parsed = JSON.parse(savedSnapshot) as {
					config: ActionEconomyConfig;
					actionSlots: ActionSlotEntry[];
				};
				config = normalizeActionEconomyConfig(parsed.config);
				actionSlots = (parsed.actionSlots ?? []).map((s) => normalizeActionSlot(s));
				syncToUi();
			} catch {
				config = defaultActionEconomyConfig();
				actionSlots = [];
				syncToUi();
			}
		}
		error = null;
		saveMessage = null;
		onDirtyChange?.(false);
	}

	function addCostRow() {
		costRows = [...costRows, { name: '', cost: 0 }];
		markDirty();
	}

	function addSlotRow() {
		const slot = defaultActionSlotEntry();
		actionSlots = [...actionSlots, slot];
		slotRows = [...slotRows, slotToTableRow(slot)];
		markDirty();
	}

	function addCombo(slotIndex: number) {
		const slot = slotAt(slotIndex);
		slot.combos = [...(slot.combos ?? []), defaultComboEntry()];
		actionSlots = [...actionSlots];
		markDirty();
	}

	function removeCombo(slotIndex: number, comboIndex: number) {
		const slot = slotAt(slotIndex);
		slot.combos = (slot.combos ?? []).filter((_, i) => i !== comboIndex);
		actionSlots = [...actionSlots];
		markDirty();
	}

	function toggleComboComponent(slotIndex: number, comboIndex: number, name: string) {
		const slot = slotAt(slotIndex);
		const combos = slot.combos ?? [];
		const combo = combos[comboIndex];
		if (!combo) return;
		const set = new Set(combo.component_names);
		if (set.has(name)) set.delete(name);
		else set.add(name);
		combo.component_names = [...set];
		actionSlots = [...actionSlots];
		markDirty();
	}

	$effect(() => {
		if (config.turn_structure === 'token_based' && !config.token_turn) {
			config.token_turn = { tokens_per_round: 1, refresh_on: 'per_round' };
		}
	});

	$effect(() => {
		if (loading) return;
		void slotRows;
		void costRows;
		markDirty();
	});

	$effect(() => {
		if (systemId) load();
	});
</script>

<div class="action-economy-tab">
	{#if loading}
		<p class="status-msg">Loading action economy settings…</p>
	{:else}
		<section class="form-section">
			<h3>Turn Structure</h3>
			<div class="form-field">
				<label for="ae-turn-structure">Turn Structure</label>
				<select
					id="ae-turn-structure"
					bind:value={config.turn_structure}
					{disabled}
					onchange={markDirty}
				>
					{#each TURN_STRUCTURE_OPTIONS as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>
			{#if showCustomTurn}
				<div class="form-field">
					<label for="ae-custom-turn">Custom Turn Structure Name</label>
					<input
						id="ae-custom-turn"
						type="text"
						bind:value={config.custom_turn_structure_name}
						{disabled}
						oninput={markDirty}
					/>
				</div>
			{/if}
			{#if showTokenTurn}
				<div class="form-subsection">
					<h4>Token-based Turn</h4>
					<div class="form-row">
						<div class="form-field">
							<label for="ae-tokens-per-round">Tokens Per Round</label>
							<input
								id="ae-tokens-per-round"
								type="number"
								min="1"
								bind:value={config.token_turn!.tokens_per_round}
								{disabled}
								oninput={markDirty}
							/>
						</div>
						<div class="form-field">
							<label for="ae-token-refresh">Refresh On</label>
							<select
								id="ae-token-refresh"
								bind:value={config.token_turn!.refresh_on}
								{disabled}
								onchange={markDirty}
							>
								{#each TOKEN_REFRESH_OPTIONS as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</select>
						</div>
					</div>
				</div>
			{/if}
		</section>

		<section class="form-section">
			<h3>Action Point System</h3>
			<fieldset class="radio-group" {disabled}>
				{#each SYSTEM_TYPE_OPTIONS as opt}
					<label class="radio-option">
						<input
							type="radio"
							name="ae-system-type"
							value={opt.value}
							checked={config.system_type === opt.value}
							{disabled}
							onchange={() => {
								config.system_type = opt.value;
								if (opt.value === 'point_pool' && !config.point_pool) {
									config.point_pool = defaultActionEconomyConfig().point_pool!;
									costRows = costRowsFromPool(config.point_pool);
								}
								if (opt.value === 'type_slots' && actionSlots.length === 0) {
									addSlotRow();
								}
								markDirty();
							}}
						/>
						<span>{opt.label}</span>
					</label>
				{/each}
			</fieldset>

			{#if showPointPool}
				<div class="form-subsection">
					<div class="form-row">
						<div class="form-field">
							<label for="ae-pool-size">Point Pool Size</label>
							<input
								id="ae-pool-size"
								type="number"
								min="0"
								bind:value={config.point_pool!.points_per_pool}
								{disabled}
								oninput={markDirty}
							/>
						</div>
						<div class="form-field">
							<label for="ae-pool-refresh">Refresh Scope</label>
							<select
								id="ae-pool-refresh"
								bind:value={config.point_pool!.refresh_scope}
								{disabled}
								onchange={markDirty}
							>
								{#each REFRESH_SCOPE_OPTIONS as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</select>
						</div>
					</div>
					<h4>Action Cost Table</h4>
					<TE_TableEditor
						columns={ACTION_COST_TABLE_COLUMNS}
						bind:rows={costRows}
						{disabled}
						addLabel="Add Row"
						minRows={1}
						onadd={addCostRow}
						onremove={(i) => {
							costRows = costRows.filter((_, idx) => idx !== i);
							markDirty();
						}}
					/>
				</div>
			{/if}

			{#if showTypeSlots}
				<div class="form-subsection">
					<h4>Action Type Definitions</h4>
					<TE_TableEditor
						columns={ACTION_SLOT_TABLE_COLUMNS}
						bind:rows={slotRows}
						{disabled}
						expandable={true}
						addLabel="Add Action Type"
						minRows={1}
						onadd={addSlotRow}
						onremove={(i) => {
							slotRows = slotRows.filter((_, idx) => idx !== i);
							actionSlots = actionSlots.filter((_, idx) => idx !== i);
							markDirty();
						}}
					>
						{#snippet rowDetail({ index })}
							{@const slot = slotAt(index)}
							{@const unlimited = isUnlimited(index)}
							<div class="slot-detail-grid">
								<div class="form-row">
									<label class="checkbox-row">
										<input
											type="checkbox"
											checked={unlimited}
											{disabled}
											onchange={(e) =>
												setUnlimited(index, (e.currentTarget as HTMLInputElement).checked)}
										/>
										<span>Unlimited</span>
									</label>
									<div class="form-field">
										<label for="ae-scope-{index}">Scope</label>
										<select
											id="ae-scope-{index}"
											value={slot.allowance_scope || 'per_turn'}
											disabled={disabled || unlimited}
											onchange={(e) => {
												slot.allowance_scope = (e.currentTarget as HTMLSelectElement).value;
												actionSlots = [...actionSlots];
												markDirty();
											}}
										>
											{#each REFRESH_SCOPE_OPTIONS as opt}
												<option value={opt.value}>{opt.label}</option>
											{/each}
										</select>
									</div>
								</div>
								{#if String(slotRows[index]?.carry_over) === 'convert'}
									<div class="form-field">
										<label for="ae-convert-{index}">Convert Target</label>
										<input
											id="ae-convert-{index}"
											type="text"
											bind:value={slot.convert_target}
											{disabled}
											oninput={() => {
												actionSlots = [...actionSlots];
												markDirty();
											}}
										/>
									</div>
								{/if}
								<label class="checkbox-row">
									<input
										type="checkbox"
										bind:checked={slot.is_reaction}
										{disabled}
										onchange={() => {
											actionSlots = [...actionSlots];
											markDirty();
										}}
									/>
									<span>Reaction / triggered action</span>
								</label>
								{#if slot.is_reaction}
									<div class="form-field">
										<label for="ae-trigger-{index}">Trigger Condition</label>
										<input
											id="ae-trigger-{index}"
											type="text"
											bind:value={slot.reaction_trigger}
											{disabled}
											oninput={() => {
												actionSlots = [...actionSlots];
												markDirty();
											}}
										/>
									</div>
								{/if}
								<label class="checkbox-row">
									<input
										type="checkbox"
										bind:checked={slot.is_free_action}
										{disabled}
										onchange={() => {
											actionSlots = [...actionSlots];
											markDirty();
										}}
									/>
									<span>Free action</span>
								</label>
								{#if slot.is_free_action}
									<div class="form-field">
										<label for="ae-free-limits-{index}">Limits on Use</label>
										<input
											id="ae-free-limits-{index}"
											type="text"
											bind:value={slot.free_action_limits}
											{disabled}
											oninput={() => {
												actionSlots = [...actionSlots];
												markDirty();
											}}
										/>
									</div>
								{/if}
								<div class="form-field">
									<label for="ae-interrupt-{index}">Action Interruption Rules</label>
									<input
										id="ae-interrupt-{index}"
										type="text"
										bind:value={slot.interruption_rules}
										{disabled}
										oninput={() => {
											actionSlots = [...actionSlots];
											markDirty();
										}}
									/>
								</div>
								<div class="form-field">
									<label for="ae-delay-{index}">Delay / Ready Action Rules</label>
									<input
										id="ae-delay-{index}"
										type="text"
										bind:value={slot.delay_ready_rules}
										{disabled}
										oninput={() => {
											actionSlots = [...actionSlots];
											markDirty();
										}}
									/>
								</div>
								<div class="combos-block">
									<h5>Action Combos</h5>
									{#each slot.combos ?? [] as combo, ci (ci)}
										<div class="combo-row">
											<div class="form-field">
												<label for="ae-combo-name-{index}-{ci}">Combo Name</label>
												<input
													id="ae-combo-name-{index}-{ci}"
													type="text"
													bind:value={combo.combo_name}
													{disabled}
													oninput={() => {
														actionSlots = [...actionSlots];
														markDirty();
													}}
												/>
											</div>
											<fieldset class="combo-components">
												<legend>Components</legend>
												{#if actionNameOptions.length === 0}
													<p class="hint">Define action names in the table above first.</p>
												{:else}
													{#each actionNameOptions as name}
														<label class="checkbox-row">
															<input
																type="checkbox"
																checked={combo.component_names.includes(name)}
																{disabled}
																onchange={() => toggleComboComponent(index, ci, name)}
															/>
															<span>{name}</span>
														</label>
													{/each}
												{/if}
											</fieldset>
											{#if !disabled}
												<button
													type="button"
													class="btn-text-danger"
													onclick={() => removeCombo(index, ci)}
												>
													Remove Combo
												</button>
											{/if}
										</div>
									{/each}
									{#if !disabled}
										<button type="button" class="btn-add-small" onclick={() => addCombo(index)}>
											Add Combo
										</button>
									{/if}
								</div>
							</div>
						{/snippet}
					</TE_TableEditor>
				</div>
			{/if}
		</section>

		<details class="round-combat-details form-section">
			<summary>Round & Combat Duration</summary>
			<div class="details-body">
				<div class="form-field">
					<label for="ae-round-time">Round Time Definition</label>
					<select
						id="ae-round-time"
						bind:value={config.round_time_definition}
						{disabled}
						onchange={markDirty}
					>
						{#each ROUND_TIME_OPTIONS as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</select>
				</div>
				{#if showCustomRound}
					<div class="form-field">
						<label for="ae-custom-round">Custom Round Time Definition</label>
						<input
							id="ae-custom-round"
							type="text"
							bind:value={config.custom_round_time_definition}
							{disabled}
							oninput={markDirty}
						/>
					</div>
				{/if}
				<div class="form-field">
					<label for="ae-combat-tracking">Combat Time Tracking Mode</label>
					<select
						id="ae-combat-tracking"
						bind:value={config.combat_time_tracking_mode}
						{disabled}
						onchange={markDirty}
					>
						{#each COMBAT_TIME_TRACKING_OPTIONS as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</select>
				</div>
				<div class="form-field">
					<label for="ae-escalation">Time Escalation Rules</label>
					<input
						id="ae-escalation"
						type="text"
						bind:value={config.time_escalation_rules}
						{disabled}
						oninput={markDirty}
					/>
				</div>
			</div>
		</details>

		<section class="form-section">
			<h3>Initiative</h3>
			<div class="form-field">
				<label for="ae-initiative-system">Initiative System</label>
				<select
					id="ae-initiative-system"
					bind:value={config.initiative_system}
					{disabled}
					onchange={markDirty}
				>
					{#each INITIATIVE_SYSTEM_OPTIONS as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>
			{#if showInitiativeExpr}
				<TE_FormulaBuilder
					bind:value={config.initiative_expression}
					label="Initiative Roll Expression"
					variables={['dexterity_mod']}
					{systemId}
					{disabled}
				/>
			{/if}
			{#if showStaticInit}
				<div class="form-field">
					<label for="ae-static-init">Static Initiative Value</label>
					<input
						id="ae-static-init"
						type="text"
						bind:value={config.static_initiative_value}
						{disabled}
						oninput={markDirty}
					/>
				</div>
			{/if}
			<div class="form-field">
				<label for="ae-init-modifiers">Initiative Modifiers</label>
				<textarea
					id="ae-init-modifiers"
					rows="3"
					bind:value={config.initiative_modifiers}
					{disabled}
					oninput={markDirty}
				></textarea>
			</div>
			<div class="form-field">
				<label for="ae-init-persist">Initiative Persistence</label>
				<select
					id="ae-init-persist"
					bind:value={config.initiative_persistence}
					{disabled}
					onchange={markDirty}
				>
					{#each INITIATIVE_PERSISTENCE_OPTIONS as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>
			<div class="form-field">
				<label for="ae-tie-breaking">Tie-breaking</label>
				<select
					id="ae-tie-breaking"
					bind:value={config.tie_breaking}
					{disabled}
					onchange={markDirty}
				>
					{#each TIE_BREAKING_OPTIONS as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>
		</section>

		{#if error}
			<p class="form-error">{error}</p>
		{/if}
		{#if saveMessage}
			<p class="status-msg success">{saveMessage}</p>
		{/if}

		{#if saving}
			<p class="status-msg">Saving…</p>
		{/if}
	{/if}
</div>

<style>
	.action-economy-tab {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-section h3,
	.form-section h4,
	.form-section h5 {
		margin: 0 0 0.65rem;
		font-size: 0.95rem;
		font-weight: 600;
	}

	.form-subsection {
		margin-top: 1rem;
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
		max-width: 24rem;
		padding: 0.45rem 0.6rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.form-field textarea {
		max-width: 100%;
		resize: vertical;
	}

	.form-row {
		display: flex;
		gap: 1rem;
		flex-wrap: wrap;
		align-items: flex-end;
	}

	.radio-group {
		border: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
	}

	.radio-option,
	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		cursor: pointer;
	}

	.round-combat-details {
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		padding: 0.5rem 0.75rem;
	}

	.round-combat-details summary {
		font-weight: 600;
		cursor: pointer;
		padding: 0.35rem 0;
	}

	.details-body {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		margin-top: 0.75rem;
	}

	.slot-detail-grid {
		display: flex;
		flex-direction: column;
		gap: 0.65rem;
	}

	.combos-block {
		margin-top: 0.5rem;
	}

	.combo-row {
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		padding: 0.65rem;
		margin-bottom: 0.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.combo-components {
		border: none;
		margin: 0;
		padding: 0;
	}

	.combo-components legend {
		font-size: 0.8rem;
		font-weight: 600;
		margin-bottom: 0.35rem;
	}

	.hint {
		margin: 0;
		font-size: 0.85rem;
		color: #6b7280;
	}

	.btn-add-small {
		padding: 0.3rem 0.65rem;
		font-size: 0.85rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.btn-text-danger {
		align-self: flex-start;
		padding: 0;
		border: none;
		background: none;
		color: #b91c1c;
		font-size: 0.85rem;
		cursor: pointer;
		font: inherit;
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

</style>

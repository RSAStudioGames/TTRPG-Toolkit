import { ALLOWANCE_UNLIMITED } from '$lib/constants/actionEconomyOptions';
import type {
	ActionComboEntry,
	ActionCostEntry,
	ActionEconomyConfig,
	ActionSlotEntry,
	PointPoolConfig
} from '$lib/types/mechanics';

export function defaultPointPoolConfig(): PointPoolConfig {
	return {
		points_per_pool: 3,
		refresh_scope: 'per_turn',
		action_cost_table: [{ name: 'Attack', cost: 2 }, { name: 'Move', cost: 1 }]
	};
}

export function defaultActionSlotEntry(): ActionSlotEntry {
	return {
		name: '',
		allowance: 1,
		allowance_scope: 'per_turn',
		carry_over: 'expire',
		is_reaction: false,
		is_free_action: false,
		combos: []
	};
}

export function defaultActionEconomyConfig(): ActionEconomyConfig {
	return {
		turn_structure: 'turn_based',
		system_type: 'point_pool',
		point_pool: defaultPointPoolConfig(),
		action_slots: [],
		round_time_definition: 'one_minute',
		combat_time_tracking_mode: 'round_counting',
		initiative_system: 'standard_rolled',
		initiative_persistence: 'persist_combat',
		initiative_expression: '1d20 + {dexterity_mod}',
		static_initiative_value: '',
		tie_breaking: 'attacker_wins'
	};
}

export function normalizeActionEconomyConfig(
	raw: Partial<ActionEconomyConfig> | Record<string, unknown> | undefined
): ActionEconomyConfig {
	const base = defaultActionEconomyConfig();
	if (!raw || typeof raw !== 'object') return base;
	const r = raw as Partial<ActionEconomyConfig>;
	return {
		turn_structure: r.turn_structure ?? base.turn_structure,
		custom_turn_structure_name: r.custom_turn_structure_name ?? '',
		token_turn: r.token_turn,
		system_type: r.system_type ?? base.system_type,
		point_pool: normalizePointPool(r.point_pool, base.point_pool!),
		action_slots: (r.action_slots ?? []).map(normalizeActionSlot),
		round_time_definition: r.round_time_definition ?? base.round_time_definition,
		custom_round_time_definition: r.custom_round_time_definition ?? '',
		combat_time_tracking_mode: r.combat_time_tracking_mode ?? base.combat_time_tracking_mode,
		time_escalation_rules: r.time_escalation_rules ?? '',
		initiative_system: r.initiative_system ?? base.initiative_system,
		initiative_persistence: r.initiative_persistence ?? base.initiative_persistence,
		initiative_expression: r.initiative_expression ?? base.initiative_expression ?? '',
		static_initiative_value: r.static_initiative_value ?? '',
		initiative_modifiers: r.initiative_modifiers ?? '',
		tie_breaking: r.tie_breaking ?? base.tie_breaking
	};
}

function normalizePointPool(
	raw: Partial<PointPoolConfig> | undefined,
	fallback: PointPoolConfig
): PointPoolConfig {
	if (!raw) return { ...fallback };
	return {
		points_per_pool: raw.points_per_pool ?? fallback.points_per_pool,
		refresh_scope: raw.refresh_scope ?? fallback.refresh_scope,
		action_cost_table: (raw.action_cost_table ?? fallback.action_cost_table ?? []).map((row) => ({
			name: row.name ?? '',
			cost: Number(row.cost) || 0
		}))
	};
}

export function normalizeActionSlot(raw: Partial<ActionSlotEntry>): ActionSlotEntry {
	return {
		name: raw.name ?? '',
		allowance: raw.allowance ?? 1,
		allowance_scope: raw.allowance_scope ?? 'per_turn',
		carry_over: raw.carry_over ?? 'expire',
		convert_target: raw.convert_target ?? '',
		is_reaction: raw.is_reaction ?? false,
		reaction_trigger: raw.reaction_trigger ?? '',
		is_free_action: raw.is_free_action ?? false,
		free_action_limits: raw.free_action_limits ?? '',
		interruption_rules: raw.interruption_rules ?? '',
		delay_ready_rules: raw.delay_ready_rules ?? '',
		combos: (raw.combos ?? []).map((c) => ({
			combo_name: c.combo_name ?? '',
			component_names: [...(c.component_names ?? [])]
		}))
	};
}

export function slotToTableRow(slot: ActionSlotEntry): Record<string, string | number> {
	const unlimited = slot.allowance === ALLOWANCE_UNLIMITED;
	return {
		name: slot.name,
		allowance: unlimited ? 1 : slot.allowance,
		carry_over: slot.carry_over,
		_unlimited: unlimited ? 1 : 0
	};
}

export function tableRowToSlot(
	row: Record<string, string | number>,
	existing?: ActionSlotEntry
): ActionSlotEntry {
	const unlimited = Number(row._unlimited) === 1;
	const base = existing ?? defaultActionSlotEntry();
	return {
		...base,
		name: String(row.name ?? ''),
		allowance: unlimited ? ALLOWANCE_UNLIMITED : Math.max(1, Number(row.allowance) || 1),
		carry_over: String(row.carry_over ?? 'expire')
	};
}

export function costRowsFromPool(pool: PointPoolConfig): Record<string, string | number>[] {
	return (pool.action_cost_table ?? []).map((r: ActionCostEntry) => ({
		name: r.name,
		cost: r.cost
	}));
}

export function poolCostsFromRows(rows: Record<string, string | number>[]): ActionCostEntry[] {
	return rows.map((r) => ({
		name: String(r.name ?? ''),
		cost: Math.max(0, Number(r.cost) || 0)
	}));
}

export function defaultComboEntry(): ActionComboEntry {
	return { combo_name: '', component_names: [] };
}

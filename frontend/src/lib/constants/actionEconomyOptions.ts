import type { TableColumn } from '$lib/components/TE_TableEditor.svelte';

export const ALLOWANCE_UNLIMITED = -1;

export const TURN_STRUCTURE_OPTIONS = [
	{ value: 'turn_based', label: 'Turn-based' },
	{ value: 'popcorn', label: 'Popcorn' },
	{ value: 'side_based', label: 'Side-based' },
	{ value: 'token_based', label: 'Token-based' },
	{ value: 'no_formal', label: 'No formal turn structure' },
	{ value: 'simultaneous', label: 'Simultaneous' },
	{ value: 'custom', label: 'Custom / user-defined' }
] as const;

export const SYSTEM_TYPE_OPTIONS = [
	{ value: 'point_pool', label: 'Action Point Pool' },
	{ value: 'type_slots', label: 'Action Type Slots' }
] as const;

export const REFRESH_SCOPE_OPTIONS = [
	{ value: 'per_turn', label: 'Turn' },
	{ value: 'per_round', label: 'Round' }
] as const;

export const CARRY_OVER_OPTIONS = [
	{ value: 'expire', label: 'Expire' },
	{ value: 'convert', label: 'Convert' }
] as const;

export const INITIATIVE_SYSTEM_OPTIONS = [
	{ value: 'standard_rolled', label: 'Standard rolled initiative' },
	{ value: 'static', label: 'Static initiative (fixed value)' },
	{ value: 'group', label: 'Group initiative (one roll for whole side)' },
	{ value: 'side_based_individual', label: 'Side-based with individual within side' },
	{ value: 'card_based', label: 'Card-based initiative (draw from deck)' },
	{ value: 'popcorn_pass', label: 'Popcorn / initiative pass' }
] as const;

export const INITIATIVE_PERSISTENCE_OPTIONS = [
	{ value: 'reroll_each_round', label: 'Reroll each round' },
	{ value: 'persist_combat', label: 'Persist entire combat' }
] as const;

export const TIE_BREAKING_OPTIONS = [
	{ value: 'attacker_wins', label: 'Attacker wins' },
	{ value: 'defender_wins', label: 'Defender wins' },
	{ value: 'roll_again', label: 'Roll again' },
	{ value: 'partial_success', label: 'Partial success' },
	{ value: 'stipulation', label: 'Stipulation' }
] as const;

export const ROUND_TIME_OPTIONS = [
	{ value: 'six_seconds', label: '6 seconds' },
	{ value: 'ten_seconds', label: '10 seconds' },
	{ value: 'one_minute', label: '1 minute' },
	{ value: 'narrative', label: 'Narrative' },
	{ value: 'custom', label: 'Custom' }
] as const;

export const COMBAT_TIME_TRACKING_OPTIONS = [
	{ value: 'round_counting', label: 'Round counting' },
	{ value: 'real_time_mapping', label: 'Real-time mapping' },
	{ value: 'narrative', label: 'Narrative' }
] as const;

export const TOKEN_REFRESH_OPTIONS = [
	{ value: 'per_round', label: 'Per round' },
	{ value: 'per_turn', label: 'Per turn' }
] as const;

export const ACTION_COST_TABLE_COLUMNS: TableColumn[] = [
	{ key: 'name', label: 'Action Name', type: 'text' },
	{ key: 'cost', label: 'Cost', type: 'number' }
];

export const ACTION_SLOT_TABLE_COLUMNS: TableColumn[] = [
	{ key: 'name', label: 'Action Name', type: 'text' },
	{ key: 'allowance', label: 'Per-Turn Allowance', type: 'number' },
	{
		key: 'carry_over',
		label: 'Carry-over Rules',
		type: 'select',
		options: CARRY_OVER_OPTIONS.map((o) => ({ value: o.value, label: o.label }))
	}
];

export const DEFAULT_INITIATIVE_EXPRESSION = '1d20 + {dexterity_mod}';

export function initiativeNeedsExpression(system: string): boolean {
	return system === 'standard_rolled' || system === 'side_based_individual';
}

export function initiativeNeedsStaticValue(system: string): boolean {
	return system === 'static';
}

import type { SelectOption } from './resolutionOptions';

export const ATTRIBUTE_TYPE_NUMERIC = 'numeric';
export const ATTRIBUTE_TYPE_STEP_DIE = 'step_die';
export const ATTRIBUTE_TYPE_DESCRIPTIVE = 'descriptive';
export const ATTRIBUTE_TYPE_RANK_TIER = 'rank_tier';
export const ATTRIBUTE_TYPE_CUSTOM = 'custom';

export const ATTRIBUTE_TYPE_OPTIONS: SelectOption[] = [
	{ value: ATTRIBUTE_TYPE_NUMERIC, label: 'Numeric' },
	{ value: ATTRIBUTE_TYPE_STEP_DIE, label: 'Step die' },
	{ value: ATTRIBUTE_TYPE_DESCRIPTIVE, label: 'Descriptive/rating' },
	{ value: ATTRIBUTE_TYPE_RANK_TIER, label: 'Rank/tier' },
	{ value: ATTRIBUTE_TYPE_CUSTOM, label: 'Custom' }
];

export const NUMERIC_FORMAT_OPTIONS: SelectOption[] = [
	{ value: 'integer', label: 'Integer' },
	{ value: 'float', label: 'Float' }
];

export const MODIFIER_DISPLAY_OPTIONS: SelectOption[] = [
	{ value: 'signed', label: 'Signed (+3)' },
	{ value: 'absolute', label: 'Absolute (3)' }
];

export const STEP_DIE_OPTIONS: SelectOption[] = [
	{ value: 'd4', label: 'd4' },
	{ value: 'd6', label: 'd6' },
	{ value: 'd8', label: 'd8' },
	{ value: 'd10', label: 'd10' },
	{ value: 'd12', label: 'd12' },
	{ value: 'd20', label: 'd20' }
];

export const CACHING_RULE_LIVE = 'live';
export const CACHING_RULE_ON_TRIGGER = 'on_trigger';

export const CACHING_RULE_OPTIONS: SelectOption[] = [
	{ value: CACHING_RULE_LIVE, label: 'Live (recalculate every time)' },
	{ value: CACHING_RULE_ON_TRIGGER, label: 'On trigger (cache until recalculate)' }
];

export const RECALCULATE_TRIGGER_OPTIONS: SelectOption[] = [
	{ value: 'attribute_change', label: 'On attribute change' },
	{ value: 'level_up', label: 'On level up' },
	{ value: 'rest', label: 'On rest' },
	{ value: 'condition_apply', label: 'On condition apply' }
];

export const MODIFIER_FORMULA_HELPER =
	'Reference the attribute score as {score}. Example: floor(({score} - 10) / 2)';

export const DERIVATION_FORMULA_HELPER =
	'Reference attributes like {strength_mod}. Example: {constitution} + {size_rating}';

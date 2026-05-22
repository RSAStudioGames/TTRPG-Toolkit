import type { SelectOption } from './resolutionOptions';

export const SKILL_TYPE_BINARY = 'binary';
export const SKILL_TYPE_MULTI_TIER = 'multi_tier';
export const SKILL_TYPE_NUMERIC = 'numeric';
export const SKILL_TYPE_STEP_DIE = 'step_die';
export const SKILL_TYPE_RANK_NAME = 'rank_name';

export const SKILL_TYPE_OPTIONS: SelectOption[] = [
	{ value: SKILL_TYPE_BINARY, label: 'Binary' },
	{ value: SKILL_TYPE_MULTI_TIER, label: 'Multi-tier' },
	{ value: SKILL_TYPE_NUMERIC, label: 'Numeric rating' },
	{ value: SKILL_TYPE_STEP_DIE, label: 'Step die' },
	{ value: SKILL_TYPE_RANK_NAME, label: 'Rank name mapping' }
];

export const SKILL_CATEGORY_CUSTOM = 'Custom';

export const SKILL_CATEGORY_OPTIONS: SelectOption[] = [
	{ value: 'Combat', label: 'Combat' },
	{ value: 'Social', label: 'Social' },
	{ value: 'Knowledge', label: 'Knowledge' },
	{ value: 'Physical', label: 'Physical' },
	{ value: 'Crafting', label: 'Crafting' },
	{ value: 'Magic', label: 'Magic' },
	{ value: 'Vehicle', label: 'Vehicle' },
	{ value: SKILL_CATEGORY_CUSTOM, label: 'Custom' }
];

export const SPECIALIZATION_EXAMPLES_HELPER =
	'e.g. For Stealth: Hide, Move Silently. For Knowledge: History, Nature.';

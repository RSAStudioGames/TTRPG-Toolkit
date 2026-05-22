import type { SelectOption } from './resolutionOptions';

export const PROGRESSION_PARADIGM_LEVEL_BASED = 'level_based';
export const PROGRESSION_PARADIGM_POINT_BUY = 'point_buy';
export const PROGRESSION_PARADIGM_MILESTONE = 'milestone';
export const PROGRESSION_PARADIGM_LIFEPATH = 'lifepath';
export const PROGRESSION_PARADIGM_PROGRESS_TRACK = 'progress_track';
export const PROGRESSION_PARADIGM_NO_ADVANCEMENT = 'no_advancement';
export const PROGRESSION_PARADIGM_CUSTOM = 'custom';

export const PROGRESSION_PARADIGM_OPTIONS: SelectOption[] = [
	{ value: PROGRESSION_PARADIGM_LEVEL_BASED, label: 'Level-based' },
	{ value: PROGRESSION_PARADIGM_POINT_BUY, label: 'Point-buy' },
	{ value: PROGRESSION_PARADIGM_MILESTONE, label: 'Milestone' },
	{ value: PROGRESSION_PARADIGM_LIFEPATH, label: 'Lifepath' },
	{ value: PROGRESSION_PARADIGM_PROGRESS_TRACK, label: 'Progress track' },
	{ value: PROGRESSION_PARADIGM_NO_ADVANCEMENT, label: 'No advancement' },
	{ value: PROGRESSION_PARADIGM_CUSTOM, label: 'Custom' }
];

export const LABEL_GM_APPROVAL = 'Require GM Approval for Advancement';
export const LABEL_ALLOW_UNDO = 'Allow Undo of Last Advancement';
export const LABEL_ALLOW_MILESTONE =
	'Allow Milestone Leveling (Ignore XP Thresholds)';

export const XP_TABLE_COLUMNS = [
	{ key: 'level', label: 'Level', type: 'number' as const },
	{ key: 'xp_required', label: 'XP Required', type: 'number' as const }
];

export const COST_TABLE_COLUMNS = [
	{ key: 'rating', label: 'Rating', type: 'number' as const },
	{ key: 'cost', label: 'Cost', type: 'number' as const }
];

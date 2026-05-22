import {
	PROGRESSION_PARADIGM_LEVEL_BASED,
	PROGRESSION_PARADIGM_MILESTONE,
	PROGRESSION_PARADIGM_POINT_BUY
} from '$lib/constants/progressionOptions';
import type { ProgressionConfig } from '$lib/types/mechanics';

export function defaultProgressionConfig(): ProgressionConfig {
	return {
		paradigm: PROGRESSION_PARADIGM_LEVEL_BASED,
		level_based: {
			min_level: 1,
			max_level: 20,
			xp_table: [{ level: 1, xp_required: 0 }],
			allow_milestone: false
		},
		point_buy: {
			starting_pool: 0,
			cost_table: [{ rating: 0, cost: 0 }]
		},
		gm_approval: false,
		allow_undo: false
	};
}

export function normalizeProgressionConfig(raw: ProgressionConfig | undefined): ProgressionConfig {
	const base = defaultProgressionConfig();
	if (!raw || typeof raw !== 'object') return base;
	return {
		paradigm: raw.paradigm || base.paradigm,
		level_based: {
			min_level: raw.level_based?.min_level ?? base.level_based!.min_level,
			max_level: raw.level_based?.max_level ?? base.level_based!.max_level,
			xp_table:
				raw.level_based?.xp_table?.length
					? [...raw.level_based.xp_table]
					: base.level_based!.xp_table,
			allow_milestone: raw.level_based?.allow_milestone ?? false
		},
		point_buy: {
			starting_pool: raw.point_buy?.starting_pool ?? 0,
			cost_table:
				raw.point_buy?.cost_table?.length
					? [...raw.point_buy.cost_table]
					: base.point_buy!.cost_table
		},
		gm_approval: raw.gm_approval ?? false,
		allow_undo: raw.allow_undo ?? false
	};
}

export function showLevelBasedBlock(paradigm: string): boolean {
	return (
		paradigm === PROGRESSION_PARADIGM_LEVEL_BASED ||
		paradigm === PROGRESSION_PARADIGM_MILESTONE
	);
}

export function showPointBuyBlock(paradigm: string): boolean {
	return paradigm === PROGRESSION_PARADIGM_POINT_BUY;
}

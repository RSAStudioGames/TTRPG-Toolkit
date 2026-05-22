export interface SelectOption {
	value: string;
	label: string;
}

export const RESOLUTION_TYPE_OPTIONS: SelectOption[] = [
	{ value: 'single_die_modifier', label: 'Single die + modifier' },
	{ value: 'dice_pool_count_successes', label: 'Dice pool (count successes)' },
	{ value: 'dice_pool_take_highest_lowest', label: 'Dice pool (take highest/lowest)' },
	{ value: 'dice_pool_sum', label: 'Dice pool (sum)' },
	{ value: 'step_dice', label: 'Step die' },
	{ value: 'card_draw', label: 'Card draw' },
	{ value: 'coin_flip', label: 'Coin flip' },
	{ value: 'narrative', label: 'Narrative / no randomizer' },
	{ value: 'custom', label: 'Custom' }
];

export const SUCCESS_METHOD_OPTIONS: SelectOption[] = [
	{ value: 'target_number', label: 'Target number' },
	{ value: 'opposed_roll', label: 'Opposed roll' },
	{ value: 'success_threshold_ladder', label: 'Success threshold ladder' },
	{ value: 'success_counting', label: 'Success counting' },
	{ value: 'binary_pass_fail', label: 'Binary pass/fail' },
	{ value: 'margin_of_success_failure', label: 'Margin of success/failure' },
	{ value: 'shifts_raises', label: 'Shifts / raises' },
	{ value: 'custom_tiers', label: 'Custom tiers' }
];

export const LADDER_OPERATOR_OPTIONS: SelectOption[] = [
	{ value: '<=', label: '<=' },
	{ value: '<', label: '<' },
	{ value: '>=', label: '>=' },
	{ value: '>', label: '>' },
	{ value: '==', label: '==' }
];

export const MECHANIC_TYPE_OPTIONS: SelectOption[] = [
	{ value: 'keep_highest_n', label: 'Keep highest N' },
	{ value: 'keep_lowest_n', label: 'Keep lowest N' },
	{ value: 'add_modifier', label: 'Add modifier' },
	{ value: 'reroll', label: 'Reroll' },
	{ value: 'shift_tier', label: 'Shift tier' }
];

export const CRIT_SUCCESS_TRIGGER_OPTIONS: SelectOption[] = [
	{ value: 'natural_max', label: 'Natural max' },
	{ value: 'exceed_dc_by_n', label: 'Exceed DC by N' }
];

export const CRIT_FAILURE_TRIGGER_OPTIONS: SelectOption[] = [
	{ value: 'natural_1', label: 'Natural 1' },
	{ value: 'botch_in_pool', label: 'Botch in pool' },
	{ value: 'zero_successes', label: 'Zero successes' }
];

export const RESOLUTION_TYPE_CUSTOM = 'custom';
export const SUCCESS_METHOD_THRESHOLD_LADDER = 'success_threshold_ladder';
export const SUCCESS_METHOD_TARGET_NUMBER = 'target_number';
export const CRIT_TRIGGER_EXCEED_DC = 'exceed_dc_by_n';

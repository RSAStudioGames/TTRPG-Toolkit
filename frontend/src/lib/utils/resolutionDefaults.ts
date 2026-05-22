import type { ResolutionConfig } from '$lib/types/mechanics';
import {
	RESOLUTION_TYPE_CUSTOM,
	SUCCESS_METHOD_TARGET_NUMBER
} from '$lib/constants/resolutionOptions';

export function defaultResolutionConfig(): ResolutionConfig {
	return {
		resolution_type: 'single_die_modifier',
		roll_expression: '',
		custom_paradigm_name: '',
		success_determination: {
			method: SUCCESS_METHOD_TARGET_NUMBER,
			threshold_ladder: [],
			default_target_variable: ''
		},
		critical_mechanics: {
			enable_crit_success: false,
			crit_success_trigger: '',
			crit_success_exceed_amount: 0,
			enable_crit_failure: false,
			crit_failure_trigger: ''
		},
		advantage_disadvantage: [
			{ name: 'Advantage', mechanic_type: 'keep_highest_n' },
			{ name: 'Disadvantage', mechanic_type: 'keep_lowest_n' }
		]
	};
}

/** Merge API payload with defaults for missing nested fields. */
export function normalizeResolutionConfig(raw: Partial<ResolutionConfig> | undefined): ResolutionConfig {
	const base = defaultResolutionConfig();
	if (!raw) return base;
	return {
		resolution_type: raw.resolution_type ?? base.resolution_type,
		roll_expression: raw.roll_expression ?? base.roll_expression,
		custom_paradigm_name: raw.custom_paradigm_name ?? base.custom_paradigm_name,
		success_determination: {
			method: raw.success_determination?.method ?? base.success_determination.method,
			threshold_ladder: raw.success_determination?.threshold_ladder ?? [],
			default_target_variable:
				raw.success_determination?.default_target_variable ??
				base.success_determination.default_target_variable
		},
		critical_mechanics: {
			enable_crit_success:
				raw.critical_mechanics?.enable_crit_success ?? base.critical_mechanics.enable_crit_success,
			crit_success_trigger:
				raw.critical_mechanics?.crit_success_trigger ?? base.critical_mechanics.crit_success_trigger,
			crit_success_exceed_amount:
				raw.critical_mechanics?.crit_success_exceed_amount ??
				base.critical_mechanics.crit_success_exceed_amount,
			enable_crit_failure:
				raw.critical_mechanics?.enable_crit_failure ?? base.critical_mechanics.enable_crit_failure,
			crit_failure_trigger:
				raw.critical_mechanics?.crit_failure_trigger ?? base.critical_mechanics.crit_failure_trigger
		},
		advantage_disadvantage:
			raw.advantage_disadvantage !== undefined
				? raw.advantage_disadvantage
				: base.advantage_disadvantage
	};
}

export { RESOLUTION_TYPE_CUSTOM };

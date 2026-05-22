/** POST /api/systems/{id}/validate-formula success data */
export interface ValidateFormulaResponse {
	valid: boolean;
}

export interface LadderTier {
	label: string;
	operator: string;
	value: number;
}

export interface SuccessDetermination {
	method: string;
	threshold_ladder: LadderTier[];
	default_target_variable: string;
}

export interface CriticalMechanics {
	enable_crit_success: boolean;
	crit_success_trigger: string;
	crit_success_exceed_amount: number;
	enable_crit_failure: boolean;
	crit_failure_trigger: string;
}

export interface AdvantageDisadvantageEntry {
	name: string;
	mechanic_type: string;
}

export interface ResolutionConfig {
	resolution_type: string;
	roll_expression: string;
	custom_paradigm_name?: string;
	success_determination: SuccessDetermination;
	critical_mechanics: CriticalMechanics;
	advantage_disadvantage: AdvantageDisadvantageEntry[];
}

/** GET /api/systems/{id}/mechanics success data */
export interface MechanicsResponse {
	id: string;
	system_id: string;
	resolution_config: ResolutionConfig;
	progression_config: Record<string, unknown>;
	action_economy_config: Record<string, unknown>;
}

export interface DescriptiveMapEntry {
	label: string;
	value: number;
}

export interface RankMapEntry {
	rank_name: string;
	numeric_backing: number;
}

export interface AttributeConfig {
	min?: number;
	max?: number;
	numeric_format?: string;
	step_dice?: string[];
	descriptive_map?: DescriptiveMapEntry[];
	rank_map?: RankMapEntry[];
	modifier_formula?: string;
	modifier_display?: string;
	is_derived: boolean;
	derivation_formula?: string;
	caching_rule?: string;
	recalculate_triggers?: string[];
}

export interface AttributeResponse {
	id: string;
	system_id: string;
	group_name?: string | null;
	parent_attribute_id?: string | null;
	name: string;
	type: string;
	config: AttributeConfig;
	sort_order: number;
}

export interface ListAttributesResponse {
	items: AttributeResponse[];
}

export interface CreateAttributePayload {
	group_name?: string | null;
	parent_attribute_id?: string | null;
	name: string;
	type: string;
	config: AttributeConfig;
	sort_order: number;
}

export interface UpdateAttributePayload {
	group_name?: string | null;
	parent_attribute_id?: string | null;
	name?: string;
	type?: string;
	config?: AttributeConfig;
	sort_order?: number;
}

export interface SkillTierEntry {
	tier_name: string;
	numeric_backing: number;
}

export interface SkillConfig {
	tiers?: SkillTierEntry[];
	min?: number;
	max?: number;
	allow_specializations: boolean;
	specialization_bonus?: number;
}

export interface SkillResponse {
	id: string;
	system_id: string;
	name: string;
	linked_attribute_id?: string | null;
	type: string;
	category?: string | null;
	config: SkillConfig;
	sort_order: number;
}

export interface ListSkillsResponse {
	items: SkillResponse[];
}

export interface CreateSkillPayload {
	name: string;
	linked_attribute_id?: string | null;
	type: string;
	category?: string | null;
	config: SkillConfig;
	sort_order: number;
}

export interface UpdateSkillPayload {
	name?: string;
	linked_attribute_id?: string | null;
	type?: string;
	category?: string | null;
	config?: SkillConfig;
	sort_order?: number;
}

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

export interface XPTableEntry {
	level: number;
	xp_required: number;
}

export interface CostTableEntry {
	rating: number;
	cost: number;
}

export interface LevelBasedConfig {
	min_level: number;
	max_level: number;
	xp_table?: XPTableEntry[];
	allow_milestone: boolean;
}

export interface PointBuyConfig {
	starting_pool: number;
	cost_table?: CostTableEntry[];
}

export interface ProgressionConfig {
	paradigm: string;
	level_based?: LevelBasedConfig;
	point_buy?: PointBuyConfig;
	gm_approval: boolean;
	allow_undo: boolean;
}

export interface ActionCostEntry {
	name: string;
	cost: number;
}

export interface ActionComboEntry {
	combo_name: string;
	component_names: string[];
}

export interface PointPoolConfig {
	points_per_pool: number;
	refresh_scope: string;
	action_cost_table?: ActionCostEntry[];
}

export interface ActionSlotEntry {
	name: string;
	allowance: number;
	allowance_scope?: string;
	carry_over: string;
	convert_target?: string;
	is_reaction: boolean;
	reaction_trigger?: string;
	is_free_action: boolean;
	free_action_limits?: string;
	interruption_rules?: string;
	delay_ready_rules?: string;
	combos?: ActionComboEntry[];
}

export interface TokenTurnConfig {
	tokens_per_round: number;
	refresh_on: string;
}

export interface ActionEconomyConfig {
	turn_structure: string;
	custom_turn_structure_name?: string;
	token_turn?: TokenTurnConfig;
	system_type: string;
	point_pool?: PointPoolConfig;
	action_slots?: ActionSlotEntry[];
	round_time_definition: string;
	custom_round_time_definition?: string;
	combat_time_tracking_mode: string;
	time_escalation_rules?: string;
	initiative_system: string;
	initiative_persistence: string;
	initiative_expression: string;
	static_initiative_value: string;
	initiative_modifiers?: string;
	tie_breaking: string;
}

/** GET /api/systems/{id}/mechanics success data */
export interface MechanicsResponse {
	id: string;
	system_id: string;
	resolution_config: ResolutionConfig;
	progression_config: ProgressionConfig;
	action_economy_config: ActionEconomyConfig;
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

export interface RecoveryScheduleEntry {
	trigger: string;
	amount: string;
	conditions?: string;
}

export interface ResourceConfig {
	current_max_format: string;
	min_val: number;
	max_val_formula: string;
	recovery_schedules?: RecoveryScheduleEntry[];
}

export interface ResourceResponse {
	id: string;
	system_id: string;
	name: string;
	type: string;
	config: ResourceConfig;
	sort_order: number;
}

export interface ListResourcesResponse {
	items: ResourceResponse[];
}

export interface CreateResourcePayload {
	name: string;
	type: string;
	config: ResourceConfig;
	sort_order: number;
}

export interface UpdateResourcePayload {
	name?: string;
	type?: string;
	config?: ResourceConfig;
	sort_order?: number;
}

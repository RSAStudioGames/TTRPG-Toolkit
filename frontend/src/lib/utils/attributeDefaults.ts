import {
	ATTRIBUTE_TYPE_CUSTOM,
	ATTRIBUTE_TYPE_DESCRIPTIVE,
	ATTRIBUTE_TYPE_NUMERIC,
	ATTRIBUTE_TYPE_RANK_TIER,
	ATTRIBUTE_TYPE_STEP_DIE,
} from '$lib/constants/attributeOptions';
import type {
	AttributeConfig,
	AttributeResponse,
	CreateAttributePayload,
	UpdateAttributePayload
} from '$lib/types/mechanics';

/** Form-local config with required formula strings for bindable fields. */
export interface AttributeFormConfig extends AttributeConfig {
	modifier_formula: string;
	derivation_formula: string;
}

export interface AttributeFormState {
	group_name: string;
	parent_attribute_id: string;
	name: string;
	type: string;
	config: AttributeFormConfig;
}

function normalizeConfigFields(cfg: AttributeConfig): AttributeFormConfig {
	return {
		...cfg,
		modifier_formula: cfg.modifier_formula ?? '',
		derivation_formula: cfg.derivation_formula ?? ''
	};
}

export function defaultAttributeConfig(type: string): AttributeFormConfig {
	const base: AttributeFormConfig = {
		modifier_formula: '',
		derivation_formula: '',
		modifier_display: 'signed',
		is_derived: false
	};

	switch (type) {
		case ATTRIBUTE_TYPE_NUMERIC:
			return withNormalizedConfig({
				...base,
				min: 0,
				max: 20,
				numeric_format: 'integer'
			});
		case ATTRIBUTE_TYPE_STEP_DIE:
			return withNormalizedConfig({ ...base, step_dice: ['d6'] });
		case ATTRIBUTE_TYPE_DESCRIPTIVE:
			return withNormalizedConfig({
				...base,
				descriptive_map: [{ label: '', value: 0 }]
			});
		case ATTRIBUTE_TYPE_RANK_TIER:
			return withNormalizedConfig({
				...base,
				rank_map: [{ rank_name: '', numeric_backing: 0 }]
			});
		case ATTRIBUTE_TYPE_CUSTOM:
		default:
			return normalizeConfigFields(base);
	}
}

function withNormalizedConfig(cfg: AttributeFormConfig): AttributeFormConfig {
	return normalizeConfigFields(cfg);
}

export function defaultAttributeForm(type = ATTRIBUTE_TYPE_NUMERIC): AttributeFormState {
	return {
		group_name: '',
		parent_attribute_id: '',
		name: '',
		type,
		config: defaultAttributeConfig(type)
	};
}

export function attributeToForm(attr: AttributeResponse): AttributeFormState {
	return {
		group_name: attr.group_name ?? '',
		parent_attribute_id: attr.parent_attribute_id ?? '',
		name: attr.name,
		type: attr.type,
		config: normalizeConfigFields(structuredClone(attr.config))
	};
}

export function derivationFormulaVariables(
	siblings: AttributeResponse[],
	excludeId?: string
): string[] {
	return siblings
		.filter((a) => a.id !== excludeId)
		.map((a) => a.name.trim())
		.filter((n) => n.length > 0);
}

/** Collect descendant ids of attrId for parent-select exclusion. */
export function descendantAttributeIds(
	attributes: AttributeResponse[],
	attrId: string
): Set<string> {
	const childrenByParent = new Map<string, string[]>();
	for (const a of attributes) {
		const pid = a.parent_attribute_id;
		if (!pid) continue;
		const list = childrenByParent.get(pid) ?? [];
		list.push(a.id);
		childrenByParent.set(pid, list);
	}
	const out = new Set<string>();
	const stack = [...(childrenByParent.get(attrId) ?? [])];
	while (stack.length > 0) {
		const id = stack.pop()!;
		if (out.has(id)) continue;
		out.add(id);
		stack.push(...(childrenByParent.get(id) ?? []));
	}
	return out;
}

function sanitizeConfigForSubmit(config: AttributeConfig): AttributeConfig {
	const out = structuredClone(config);
	if (!out.is_derived) {
		out.derivation_formula = undefined;
		out.caching_rule = undefined;
		out.recalculate_triggers = undefined;
	}
	if (out.caching_rule !== 'on_trigger') {
		out.recalculate_triggers = undefined;
	}
	return out;
}

function optionalGroupName(group: string): string | null | undefined {
	const t = group.trim();
	return t === '' ? null : t;
}

function optionalParentId(parent: string): string | null | undefined {
	const t = parent.trim();
	return t === '' ? null : t;
}

export function formToCreatePayload(
	form: AttributeFormState,
	sortOrder: number
): CreateAttributePayload {
	return {
		name: form.name.trim(),
		type: form.type,
		group_name: optionalGroupName(form.group_name),
		parent_attribute_id: optionalParentId(form.parent_attribute_id),
		config: sanitizeConfigForSubmit(form.config),
		sort_order: sortOrder
	};
}

export function formToUpdatePayload(form: AttributeFormState): UpdateAttributePayload {
	return {
		name: form.name.trim(),
		type: form.type,
		group_name: optionalGroupName(form.group_name),
		parent_attribute_id: optionalParentId(form.parent_attribute_id),
		config: sanitizeConfigForSubmit(form.config)
	};
}

export function validateAttributeForm(form: AttributeFormState): string | null {
	const name = form.name.trim();
	if (name.length < 3) {
		return 'Name is required (at least 3 characters).';
	}
	const cfg = form.config;
	if (form.type === ATTRIBUTE_TYPE_NUMERIC) {
		if (!cfg.numeric_format) return 'Numeric format is required.';
		if ((cfg.min ?? 0) > (cfg.max ?? 0)) return 'Default Minimum cannot exceed Default Maximum.';
	}
	if (form.type === ATTRIBUTE_TYPE_STEP_DIE && (cfg.step_dice?.length ?? 0) === 0) {
		return 'At least one die is required for Step die attributes.';
	}
	if (form.type === ATTRIBUTE_TYPE_DESCRIPTIVE) {
		if ((cfg.descriptive_map?.length ?? 0) === 0) {
			return 'At least one mapping is required for Descriptive attributes.';
		}
		for (const e of cfg.descriptive_map ?? []) {
			if (!e.label.trim()) return 'Each mapping must have a label.';
		}
	}
	if (form.type === ATTRIBUTE_TYPE_RANK_TIER) {
		if ((cfg.rank_map?.length ?? 0) === 0) {
			return 'At least one rank is required for Rank/tier attributes.';
		}
		for (const e of cfg.rank_map ?? []) {
			if (!e.rank_name.trim()) return 'Each rank must have a name.';
		}
	}
	if (cfg.is_derived) {
		if (!cfg.derivation_formula?.trim()) {
			return 'Derivation Formula is required when derived attributes are enabled.';
		}
		if (!cfg.caching_rule) return 'Caching rule is required for derived attributes.';
		if (
			cfg.caching_rule === 'on_trigger' &&
			(cfg.recalculate_triggers?.length ?? 0) === 0
		) {
			return 'Select at least one Recalculate Trigger when caching is on trigger.';
		}
	}
	return null;
}

export function attributeTypeLabel(type: string): string {
	const labels: Record<string, string> = {
		numeric: 'Numeric',
		step_die: 'Step die',
		descriptive: 'Descriptive/rating',
		rank_tier: 'Rank/tier',
		custom: 'Custom'
	};
	return labels[type] ?? type;
}

export function modifierDisplayLabel(value: string | undefined): string {
	if (value === 'absolute') return 'Absolute (3)';
	if (value === 'signed') return 'Signed (+3)';
	return value ?? '—';
}

import {
	SKILL_CATEGORY_CUSTOM,
	SKILL_CATEGORY_OPTIONS,
	SKILL_TYPE_BINARY,
	SKILL_TYPE_MULTI_TIER,
	SKILL_TYPE_NUMERIC
} from '$lib/constants/skillOptions';
import type {
	CreateSkillPayload,
	SkillConfig,
	SkillResponse,
	UpdateSkillPayload
} from '$lib/types/mechanics';

const PRESET_CATEGORIES = new Set(SKILL_CATEGORY_OPTIONS.map((o) => o.value));

export interface SkillFormState {
	name: string;
	linked_attribute_id: string;
	type: string;
	category_preset: string;
	category_custom: string;
	config: SkillConfig;
}

export function defaultSkillConfig(type: string): SkillConfig {
	const base: SkillConfig = {
		allow_specializations: false,
		specialization_bonus: 2
	};
	switch (type) {
		case SKILL_TYPE_MULTI_TIER:
			return { ...base, tiers: [{ tier_name: '', numeric_backing: 0 }] };
		case SKILL_TYPE_NUMERIC:
			return { ...base, min: 0, max: 10 };
		default:
			return base;
	}
}

export function defaultSkillForm(type = SKILL_TYPE_BINARY): SkillFormState {
	return {
		name: '',
		linked_attribute_id: '',
		type,
		category_preset: 'Combat',
		category_custom: '',
		config: defaultSkillConfig(type)
	};
}

export function skillToForm(skill: SkillResponse): SkillFormState {
	const cat = skill.category?.trim() ?? '';
	const isPreset = cat !== '' && PRESET_CATEGORIES.has(cat) && cat !== SKILL_CATEGORY_CUSTOM;
	return {
		name: skill.name,
		linked_attribute_id: skill.linked_attribute_id ?? '',
		type: skill.type,
		category_preset: isPreset ? cat : cat === '' ? 'Combat' : SKILL_CATEGORY_CUSTOM,
		category_custom: isPreset ? '' : cat,
		config: structuredClone(skill.config)
	};
}

function resolveCategory(form: SkillFormState): string | null {
	if (form.category_preset === SKILL_CATEGORY_CUSTOM) {
		const custom = form.category_custom.trim();
		return custom === '' ? null : custom;
	}
	return form.category_preset || null;
}

function optionalLinkedId(id: string): string | null {
	const t = id.trim();
	return t === '' ? null : t;
}

function sanitizeConfigForSubmit(skillType: string, config: SkillConfig): SkillConfig {
	const out = structuredClone(config);
	if (!out.allow_specializations) {
		out.specialization_bonus = undefined;
	} else if (!out.specialization_bonus) {
		out.specialization_bonus = 2;
	}
	if (skillType !== SKILL_TYPE_MULTI_TIER) {
		out.tiers = undefined;
	}
	if (skillType !== SKILL_TYPE_NUMERIC) {
		out.min = undefined;
		out.max = undefined;
	}
	return out;
}

export function formToCreatePayload(form: SkillFormState, sortOrder: number): CreateSkillPayload {
	return {
		name: form.name.trim(),
		type: form.type,
		linked_attribute_id: optionalLinkedId(form.linked_attribute_id),
		category: resolveCategory(form),
		config: sanitizeConfigForSubmit(form.type, form.config),
		sort_order: sortOrder
	};
}

export function formToUpdatePayload(form: SkillFormState): UpdateSkillPayload {
	return {
		name: form.name.trim(),
		type: form.type,
		linked_attribute_id: optionalLinkedId(form.linked_attribute_id),
		category: resolveCategory(form),
		config: sanitizeConfigForSubmit(form.type, form.config)
	};
}

export function validateSkillForm(form: SkillFormState): string | null {
	if (!form.name.trim()) {
		return 'Name is required.';
	}
	if (form.type === SKILL_TYPE_MULTI_TIER) {
		if ((form.config.tiers?.length ?? 0) === 0) {
			return 'At least one tier is required for Multi-tier skills.';
		}
		for (const t of form.config.tiers ?? []) {
			if (!t.tier_name.trim()) return 'Each tier must have a name.';
		}
	}
	if (form.type === SKILL_TYPE_NUMERIC && (form.config.min ?? 0) > (form.config.max ?? 0)) {
		return 'Min Bound cannot exceed Max Bound.';
	}
	if (form.category_preset === SKILL_CATEGORY_CUSTOM && !form.category_custom.trim()) {
		return 'Enter a custom category name.';
	}
	return null;
}

export function skillTypeLabel(type: string): string {
	const labels: Record<string, string> = {
		binary: 'Binary',
		multi_tier: 'Multi-tier',
		numeric: 'Numeric rating',
		step_die: 'Step die',
		rank_name: 'Rank name mapping'
	};
	return labels[type] ?? type;
}

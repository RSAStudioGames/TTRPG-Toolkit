import {
	RECOVERY_TRIGGER_CUSTOM,
	RESOURCE_FORMAT_INTEGER
} from '$lib/constants/resourceOptions';
import type {
	CreateResourcePayload,
	ResourceConfig,
	ResourceResponse,
	UpdateResourcePayload
} from '$lib/types/mechanics';

export interface ResourceFormState {
	name: string;
	type: string;
	config: ResourceConfig;
}

export function defaultResourceConfig(): ResourceConfig {
	return {
		current_max_format: RESOURCE_FORMAT_INTEGER,
		min_val: 0,
		max_val_formula: '',
		recovery_schedules: []
	};
}

export function defaultResourceForm(): ResourceFormState {
	return {
		name: '',
		type: 'pool',
		config: defaultResourceConfig()
	};
}

export function resourceToForm(resource: ResourceResponse): ResourceFormState {
	return {
		name: resource.name,
		type: resource.type,
		config: structuredClone(resource.config)
	};
}

export function formToCreatePayload(
	form: ResourceFormState,
	sortOrder: number
): CreateResourcePayload {
	return {
		name: form.name.trim(),
		type: form.type,
		config: sanitizeResourceConfig(form.config),
		sort_order: sortOrder
	};
}

export function formToUpdatePayload(form: ResourceFormState): UpdateResourcePayload {
	return {
		name: form.name.trim(),
		type: form.type,
		config: sanitizeResourceConfig(form.config)
	};
}

function sanitizeResourceConfig(cfg: ResourceConfig): ResourceConfig {
	const out = structuredClone(cfg);
	for (const row of out.recovery_schedules ?? []) {
		if (row.trigger !== RECOVERY_TRIGGER_CUSTOM) {
			row.conditions = row.conditions ?? '';
		}
	}
	return out;
}

export function validateResourceForm(form: ResourceFormState): string | null {
	if (!form.name.trim()) return 'Name is required.';
	if (!form.config.current_max_format) return 'Value bounds format is required.';
	for (const row of form.config.recovery_schedules ?? []) {
		if (!row.amount.trim()) return 'Each recovery event must have a recovery amount.';
		if (row.trigger === RECOVERY_TRIGGER_CUSTOM && !row.amount.trim()) {
			return 'Recovery amount is required.';
		}
	}
	return null;
}

export function resourceTypeLabel(type: string): string {
	const labels: Record<string, string> = {
		pool: 'Pool',
		slot_track: 'Slot track',
		counter: 'Counter',
		currency: 'Currency',
		custom: 'Custom'
	};
	return labels[type] ?? type;
}

export function resourceFormatLabel(format: string): string {
	const labels: Record<string, string> = {
		integer: 'Integer',
		float: 'Float',
		die_pool: 'Die-pool',
		step_die: 'Step-die'
	};
	return labels[format] ?? format;
}

import type { SelectOption } from './resolutionOptions';

export const RESOURCE_TYPE_POOL = 'pool';
export const RESOURCE_TYPE_SLOT_TRACK = 'slot_track';
export const RESOURCE_TYPE_COUNTER = 'counter';
export const RESOURCE_TYPE_CURRENCY = 'currency';
export const RESOURCE_TYPE_CUSTOM = 'custom';

export const RESOURCE_TYPE_OPTIONS: SelectOption[] = [
	{ value: RESOURCE_TYPE_POOL, label: 'Pool' },
	{ value: RESOURCE_TYPE_SLOT_TRACK, label: 'Slot track' },
	{ value: RESOURCE_TYPE_COUNTER, label: 'Counter' },
	{ value: RESOURCE_TYPE_CURRENCY, label: 'Currency' },
	{ value: RESOURCE_TYPE_CUSTOM, label: 'Custom' }
];

export const RESOURCE_FORMAT_INTEGER = 'integer';
export const RESOURCE_FORMAT_FLOAT = 'float';
export const RESOURCE_FORMAT_DIE_POOL = 'die_pool';
export const RESOURCE_FORMAT_STEP_DIE = 'step_die';

export const RESOURCE_FORMAT_OPTIONS: SelectOption[] = [
	{ value: RESOURCE_FORMAT_INTEGER, label: 'Integer' },
	{ value: RESOURCE_FORMAT_FLOAT, label: 'Float' },
	{ value: RESOURCE_FORMAT_DIE_POOL, label: 'Die-pool' },
	{ value: RESOURCE_FORMAT_STEP_DIE, label: 'Step-die' }
];

export const RECOVERY_TRIGGER_CUSTOM = 'custom';

export const RECOVERY_TRIGGER_OPTIONS: SelectOption[] = [
	{ value: 'short_rest', label: 'Short Rest' },
	{ value: 'long_rest', label: 'Long Rest' },
	{ value: 'dawn', label: 'Dawn' },
	{ value: RECOVERY_TRIGGER_CUSTOM, label: 'Custom' }
];

export const ADD_RECOVERY_EVENT_LABEL = 'Add Recovery Event';
export const MIN_VAL_HELPER = 'Can go negative?';

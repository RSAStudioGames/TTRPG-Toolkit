import type { AttributesConfig } from '$lib/types/mechanics';

export function defaultAttributesConfig(): AttributesConfig {
	return { enabled_derived: false };
}

export function normalizeAttributesConfig(raw?: Partial<AttributesConfig> | null): AttributesConfig {
	return {
		enabled_derived: raw?.enabled_derived === true
	};
}

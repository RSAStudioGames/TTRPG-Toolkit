export type SystemStatus = 'draft' | 'published' | 'locked' | 'archived';

export interface GameSystem {
	id: string;
	name: string;
	slug: string;
	edition?: string;
	publisher?: string;
	description?: string;
	license_type?: string;
	version: string;
	playstyle?: string;
	complexity?: number;
	measurement_unit?: string;
	currency_symbol?: string;
	status: SystemStatus;
	is_active: boolean;
	system_family?: string;
	player_count_min?: number;
	player_count_max?: number;
	official_links: string[];
	tags: string[];
	core_rulebooks: string[];
	icon_url?: string;
	cover_url?: string;
	parent_system_id?: string;
	is_core: boolean;
	is_protected: boolean;
	created_at: string;
	updated_at: string;
}

export interface SystemListResponse {
	items: GameSystem[];
	page: number;
	per_page: number;
	total: number;
	total_pages: number;
}

export interface DeletePreview {
	tag_count: number;
	rulebook_count: number;
	link_count: number;
	child_count: number;
	total_associated: number;
}

export interface SystemFormValues {
	name: string;
	slug: string;
	edition: string;
	publisher: string;
	description: string;
	license_type: string;
	version: string;
	playstyle: string;
	complexity: number;
	measurement_unit: string;
	currency_symbol: string;
	is_active: boolean;
	system_family: string;
	player_count_min: string;
	player_count_max: string;
	official_links: string[];
	tags: string[];
	core_rulebooks: string[];
	is_core: boolean;
	parent_system_id: string;
	is_protected: boolean;
}

export function defaultFormValues(): SystemFormValues {
	return {
		name: '',
		slug: '',
		edition: '',
		publisher: '',
		description: '',
		license_type: 'Homebrew',
		version: '0.1.0',
		playstyle: 'Narrative',
		complexity: 3,
		measurement_unit: 'Imperial',
		currency_symbol: '',
		is_active: true,
		system_family: '',
		player_count_min: '',
		player_count_max: '',
		official_links: [''],
		tags: [],
		core_rulebooks: [],
		is_core: true,
		parent_system_id: '',
		is_protected: false
	};
}

export function systemToFormValues(s: GameSystem): SystemFormValues {
	return {
		name: s.name,
		slug: s.slug,
		edition: s.edition ?? '',
		publisher: s.publisher ?? '',
		description: s.description ?? '',
		license_type: s.license_type ?? 'Homebrew',
		version: s.version ?? '0.1.0',
		playstyle: s.playstyle ?? 'Narrative',
		complexity: s.complexity ?? 3,
		measurement_unit: s.measurement_unit ?? 'Imperial',
		currency_symbol: s.currency_symbol ?? '',
		is_active: s.is_active,
		system_family: s.system_family ?? '',
		player_count_min: s.player_count_min != null ? String(s.player_count_min) : '',
		player_count_max: s.player_count_max != null ? String(s.player_count_max) : '',
		official_links: s.official_links?.length ? [...s.official_links] : [''],
		tags: s.tags ?? [],
		core_rulebooks: s.core_rulebooks ?? [],
		is_core: s.is_core,
		parent_system_id: s.parent_system_id ?? '',
		is_protected: s.is_protected
	};
}

export function formValuesToPayload(v: SystemFormValues): Record<string, unknown> {
	const min = v.player_count_min !== '' ? Number(v.player_count_min) : undefined;
	const max = v.player_count_max !== '' ? Number(v.player_count_max) : undefined;
	const links = v.official_links.map((u) => u.trim()).filter(Boolean);
	return {
		name: v.name.trim(),
		slug: v.slug.trim(),
		edition: v.edition.trim() || undefined,
		publisher: v.publisher.trim() || undefined,
		description: v.description.trim() || undefined,
		license_type: v.license_type || undefined,
		version: v.version.trim() || '0.1.0',
		playstyle: v.playstyle || undefined,
		complexity: v.complexity,
		measurement_unit: v.measurement_unit || undefined,
		currency_symbol: v.currency_symbol.trim() || undefined,
		is_active: v.is_active,
		system_family: v.system_family.trim() || undefined,
		player_count_min: min,
		player_count_max: max,
		official_links: links,
		tags: v.tags,
		core_rulebooks: v.core_rulebooks,
		is_core: v.is_core,
		parent_system_id: v.parent_system_id.trim() || undefined,
		is_protected: v.is_protected
	};
}

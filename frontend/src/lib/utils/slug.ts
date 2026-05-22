const SLUG_PATTERN = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;

export function slugFromName(name: string): string {
	return name
		.trim()
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '-')
		.replace(/^-+|-+$/g, '');
}

export function isValidSlug(slug: string): boolean {
	return slug.length > 0 && SLUG_PATTERN.test(slug);
}

export function slugValidationMessage(slug: string): string | null {
	if (!slug.trim()) return 'URL Identifier is required';
	if (!isValidSlug(slug)) {
		return 'URL Identifier must be lowercase alphanumeric with hyphens';
	}
	return null;
}

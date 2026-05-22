import { apiEnvelope } from './client';
import type {
	AttributeResponse,
	CreateAttributePayload,
	CreateSkillPayload,
	ListAttributesResponse,
	ListSkillsResponse,
	MechanicsResponse,
	ResolutionConfig,
	SkillResponse,
	UpdateAttributePayload,
	UpdateSkillPayload,
	ValidateFormulaResponse
} from '$lib/types/mechanics';

export async function validateFormula(
	systemId: string,
	formula: string
): Promise<ValidateFormulaResponse> {
	return apiEnvelope<ValidateFormulaResponse>(`/api/systems/${systemId}/validate-formula`, {
		method: 'POST',
		body: JSON.stringify({ formula })
	});
}

export async function getMechanics(systemId: string): Promise<MechanicsResponse> {
	return apiEnvelope<MechanicsResponse>(`/api/systems/${systemId}/mechanics`);
}

export async function saveResolutionConfig(
	systemId: string,
	body: ResolutionConfig
): Promise<MechanicsResponse> {
	return apiEnvelope<MechanicsResponse>(`/api/systems/${systemId}/mechanics/resolution`, {
		method: 'PUT',
		body: JSON.stringify(body)
	});
}

export async function listAttributes(systemId: string): Promise<ListAttributesResponse> {
	return apiEnvelope<ListAttributesResponse>(`/api/systems/${systemId}/attributes`);
}

export async function createAttribute(
	systemId: string,
	body: CreateAttributePayload
): Promise<AttributeResponse> {
	return apiEnvelope<AttributeResponse>(`/api/systems/${systemId}/attributes`, {
		method: 'POST',
		body: JSON.stringify(body)
	});
}

export async function getAttribute(
	systemId: string,
	attrId: string
): Promise<AttributeResponse> {
	return apiEnvelope<AttributeResponse>(
		`/api/systems/${systemId}/attributes/${attrId}`
	);
}

export async function updateAttribute(
	systemId: string,
	attrId: string,
	body: UpdateAttributePayload
): Promise<AttributeResponse> {
	return apiEnvelope<AttributeResponse>(`/api/systems/${systemId}/attributes/${attrId}`, {
		method: 'PUT',
		body: JSON.stringify(body)
	});
}

export async function deleteAttribute(systemId: string, attrId: string): Promise<void> {
	await apiEnvelope<{ deleted: string }>(`/api/systems/${systemId}/attributes/${attrId}`, {
		method: 'DELETE'
	});
}

export async function listSkills(systemId: string): Promise<ListSkillsResponse> {
	return apiEnvelope<ListSkillsResponse>(`/api/systems/${systemId}/skills`);
}

export async function createSkill(
	systemId: string,
	body: CreateSkillPayload
): Promise<SkillResponse> {
	return apiEnvelope<SkillResponse>(`/api/systems/${systemId}/skills`, {
		method: 'POST',
		body: JSON.stringify(body)
	});
}

export async function updateSkill(
	systemId: string,
	skillId: string,
	body: UpdateSkillPayload
): Promise<SkillResponse> {
	return apiEnvelope<SkillResponse>(`/api/systems/${systemId}/skills/${skillId}`, {
		method: 'PUT',
		body: JSON.stringify(body)
	});
}

export async function deleteSkill(systemId: string, skillId: string): Promise<void> {
	await apiEnvelope<{ deleted: string }>(`/api/systems/${systemId}/skills/${skillId}`, {
		method: 'DELETE'
	});
}

import { apiEnvelope } from './client';
import type {
	AttributeGroupResponse,
	AttributeResponse,
	CreateAttributeGroupPayload,
	CreateAttributePayload,
	CreateResourcePayload,
	CreateSkillPayload,
	ListAttributeGroupsResponse,
	ListAttributesResponse,
	ListResourcesResponse,
	ListSkillsResponse,
	ActionEconomyConfig,
	AttributesConfig,
	MechanicsResponse,
	ProgressionConfig,
	ResolutionConfig,
	ResourceResponse,
	SkillResponse,
	UpdateAttributeGroupPayload,
	UpdateAttributePayload,
	UpdateResourcePayload,
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

export async function listAttributeGroups(
	systemId: string
): Promise<ListAttributeGroupsResponse> {
	return apiEnvelope<ListAttributeGroupsResponse>(`/api/systems/${systemId}/attribute-groups`);
}

export async function createAttributeGroup(
	systemId: string,
	body: CreateAttributeGroupPayload
): Promise<AttributeGroupResponse> {
	return apiEnvelope<AttributeGroupResponse>(`/api/systems/${systemId}/attribute-groups`, {
		method: 'POST',
		body: JSON.stringify(body)
	});
}

export async function updateAttributeGroup(
	systemId: string,
	groupId: string,
	body: UpdateAttributeGroupPayload
): Promise<AttributeGroupResponse> {
	return apiEnvelope<AttributeGroupResponse>(
		`/api/systems/${systemId}/attribute-groups/${groupId}`,
		{ method: 'PUT', body: JSON.stringify(body) }
	);
}

export async function deleteAttributeGroup(systemId: string, groupId: string): Promise<void> {
	await apiEnvelope<{ deleted: string }>(
		`/api/systems/${systemId}/attribute-groups/${groupId}`,
		{ method: 'DELETE' }
	);
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

export async function saveProgressionConfig(
	systemId: string,
	body: ProgressionConfig
): Promise<MechanicsResponse> {
	return apiEnvelope<MechanicsResponse>(`/api/systems/${systemId}/mechanics/progression`, {
		method: 'PUT',
		body: JSON.stringify(body)
	});
}

export async function saveActionEconomyConfig(
	systemId: string,
	body: ActionEconomyConfig
): Promise<MechanicsResponse> {
	return apiEnvelope<MechanicsResponse>(`/api/systems/${systemId}/mechanics/action-economy`, {
		method: 'PUT',
		body: JSON.stringify(body)
	});
}

export async function saveAttributesConfig(
	systemId: string,
	body: AttributesConfig
): Promise<MechanicsResponse> {
	return apiEnvelope<MechanicsResponse>(`/api/systems/${systemId}/mechanics/attributes-config`, {
		method: 'PUT',
		body: JSON.stringify(body)
	});
}

export async function listResources(systemId: string): Promise<ListResourcesResponse> {
	return apiEnvelope<ListResourcesResponse>(`/api/systems/${systemId}/resources`);
}

export async function createResource(
	systemId: string,
	body: CreateResourcePayload
): Promise<ResourceResponse> {
	return apiEnvelope<ResourceResponse>(`/api/systems/${systemId}/resources`, {
		method: 'POST',
		body: JSON.stringify(body)
	});
}

export async function updateResource(
	systemId: string,
	resourceId: string,
	body: UpdateResourcePayload
): Promise<ResourceResponse> {
	return apiEnvelope<ResourceResponse>(
		`/api/systems/${systemId}/resources/${resourceId}`,
		{
			method: 'PUT',
			body: JSON.stringify(body)
		}
	);
}

export async function deleteResource(systemId: string, resourceId: string): Promise<void> {
	await apiEnvelope<{ deleted: string }>(
		`/api/systems/${systemId}/resources/${resourceId}`,
		{ method: 'DELETE' }
	);
}

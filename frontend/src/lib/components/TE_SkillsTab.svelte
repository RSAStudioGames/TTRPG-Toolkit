<script lang="ts">
	import TE_SkillModal from './TE_SkillModal.svelte';
	import { deleteSkill, listAttributes, listSkills } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import { skillTypeLabel } from '$lib/utils/skillDefaults';
	import type { AttributeResponse, SkillResponse } from '$lib/types/mechanics';

	interface Props {
		systemId: string;
		disabled?: boolean;
		onModalOpenChange?: (open: boolean) => void;
	}

	let { systemId, disabled = false, onModalOpenChange }: Props = $props();

	export function closeModalDiscard() {
		closeModal();
	}

	let skills = $state<SkillResponse[]>([]);
	let attributes = $state<AttributeResponse[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let modalOpen = $state(false);
	let editingSkill = $state<SkillResponse | null>(null);

	const attributeNameById = $derived(
		new Map(attributes.map((a) => [a.id, a.name]))
	);

	const sortedSkills = $derived(
		[...skills].sort((a, b) => a.sort_order - b.sort_order || a.name.localeCompare(b.name))
	);

	function linkedAttributeName(skill: SkillResponse): string {
		if (!skill.linked_attribute_id) return '—';
		return attributeNameById.get(skill.linked_attribute_id) ?? '—';
	}

	async function load() {
		loading = true;
		error = null;
		try {
			const [skillsData, attrsData] = await Promise.all([
				listSkills(systemId),
				listAttributes(systemId)
			]);
			skills = skillsData.items ?? [];
			attributes = attrsData.items ?? [];
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to load skills';
		} finally {
			loading = false;
		}
	}

	function openCreate() {
		editingSkill = null;
		modalOpen = true;
	}

	function openEdit(skill: SkillResponse) {
		editingSkill = skill;
		modalOpen = true;
	}

	function closeModal() {
		modalOpen = false;
		editingSkill = null;
	}

	async function handleDelete(skill: SkillResponse) {
		if (!confirm(`Delete skill "${skill.name}"?`)) return;
		error = null;
		try {
			await deleteSkill(systemId, skill.id);
			await load();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to delete skill';
		}
	}

	$effect(() => {
		onModalOpenChange?.(modalOpen);
	});

	$effect(() => {
		if (systemId) load();
	});
</script>

<div class="skills-tab">
	{#if loading}
		<p class="status-msg">Loading skills…</p>
	{:else}
		{#if error}
			<p class="form-error">{error}</p>
		{/if}

		<div class="mechanics-table-wrap">
			<table class="mechanics-table">
				<thead>
					<tr>
						<th>Name</th>
						<th>Linked Attribute</th>
						<th>Rating Type</th>
						<th>Category</th>
						<th class="col-actions">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each sortedSkills as skill (skill.id)}
						<tr>
							<td>{skill.name}</td>
							<td>{linkedAttributeName(skill)}</td>
							<td>{skillTypeLabel(skill.type)}</td>
							<td>{skill.category ?? '—'}</td>
							<td class="col-actions">
								<button
									type="button"
									class="btn-sm"
									{disabled}
									onclick={() => openEdit(skill)}
								>
									Edit
								</button>
								<button
									type="button"
									class="btn-sm btn-danger-outline"
									{disabled}
									onclick={() => handleDelete(skill)}
								>
									Delete
								</button>
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="5" class="empty-row">No skills defined yet.</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div class="tab-footer">
			<button type="button" class="btn-primary" {disabled} onclick={openCreate}>
				Create Skill
			</button>
		</div>
	{/if}
</div>

<TE_SkillModal
	open={modalOpen}
	{systemId}
	skill={editingSkill}
	{attributes}
	siblingSkillCount={skills.length}
	{disabled}
	onclose={closeModal}
	onsaved={load}
/>

<style>
	.skills-tab {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.status-msg {
		margin: 0;
		color: var(--text-muted, #6b7280);
		font-size: 0.9rem;
	}

	.form-error {
		color: #b91c1c;
		margin: 0;
		font-size: 0.875rem;
	}

	.mechanics-table-wrap {
		overflow-x: auto;
	}

	.mechanics-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.9rem;
	}

	.mechanics-table th,
	.mechanics-table td {
		padding: 0.5rem;
		border-bottom: 1px solid #e5e7eb;
		text-align: left;
		vertical-align: middle;
	}

	.col-actions {
		white-space: nowrap;
	}

	.empty-row {
		color: var(--text-muted, #6b7280);
		font-style: italic;
	}

	.tab-footer {
		margin-top: 0.25rem;
	}

	.btn-primary {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		background: var(--accent-gm, #c9a227);
		color: #1a1a1a;
		font-weight: 600;
		cursor: pointer;
		font: inherit;
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-sm {
		padding: 0.25rem 0.55rem;
		font-size: 0.8rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
		margin-right: 0.35rem;
	}

	.btn-sm:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-danger-outline {
		border-color: #fca5a5;
		color: #b91c1c;
	}
</style>

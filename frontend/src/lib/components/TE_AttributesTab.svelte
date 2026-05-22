<script lang="ts">
	import { onMount } from 'svelte';
	import TE_AttributeGroupModal from './TE_AttributeGroupModal.svelte';
	import TE_AttributeModal from './TE_AttributeModal.svelte';
	import TE_DerivedAttributeModal from './TE_DerivedAttributeModal.svelte';
	import TE_DragList, { type DragListItem } from './TE_DragList.svelte';
	import {
		deleteAttribute,
		deleteAttributeGroup,
		getMechanics,
		listAttributeGroups,
		listAttributes,
		saveAttributesConfig,
		updateAttribute,
		updateAttributeGroup
	} from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import { normalizeAttributesConfig } from '$lib/utils/attributesDefaults';
	import { attributeTypeLabel } from '$lib/utils/attributeDefaults';
	import type { AttributeGroupResponse, AttributeResponse } from '$lib/types/mechanics';

	const UNGROUPED_ID = '__ungrouped__';

	interface Props {
		systemId: string;
		disabled?: boolean;
		onModalOpenChange?: (open: boolean) => void;
	}

	let { systemId, disabled = false, onModalOpenChange }: Props = $props();

	export function closeModalDiscard() {
		closeCoreModal();
		closeDerivedModal();
		closeGroupModal();
	}

	let attributes = $state<AttributeResponse[]>([]);
	let groups = $state<AttributeGroupResponse[]>([]);
	let enabledDerived = $state(false);
	let loading = $state(true);
	let reordering = $state(false);
	let toggleSaving = $state(false);
	let error = $state<string | null>(null);
	let coreModalOpen = $state(false);
	let derivedModalOpen = $state(false);
	let groupModalOpen = $state(false);
	let editingAttribute = $state<AttributeResponse | null>(null);
	let editingGroup = $state<AttributeGroupResponse | null>(null);
	let defaultGroupId = $state('');

	const coreAttributes = $derived(
		[...attributes]
			.filter((a) => !a.config.is_derived)
			.sort((a, b) => a.sort_order - b.sort_order)
	);

	const derivedAttributes = $derived(
		[...attributes]
			.filter((a) => a.config.is_derived)
			.sort((a, b) => a.sort_order - b.sort_order)
	);

	const sortedGroups = $derived(
		[...groups].sort((a, b) => a.sort_order - b.sort_order || a.name.localeCompare(b.name))
	);

	const groupDragItems = $derived.by((): DragListItem[] => {
		const items: DragListItem[] = sortedGroups.map((g) => ({
			id: g.id,
			sortOrder: g.sort_order,
			header: g.name,
			subtitle: `${coreInGroup(g.id).length} attributes`
		}));
		const ungroupedCount = coreInGroup(null).length;
		items.push({
			id: UNGROUPED_ID,
			sortOrder: items.length,
			header: 'Ungrouped',
			subtitle: `${ungroupedCount} attributes`
		});
		return items;
	});

	const derivedDragItems = $derived.by((): DragListItem[] =>
		derivedAttributes.map((attr) => ({
			id: attr.id,
			sortOrder: attr.sort_order,
			header: attr.name,
			summary: `${attributeTypeLabel(attr.type)} · Derived`
		}))
	);

	function coreInGroup(groupId: string | null): AttributeResponse[] {
		return coreAttributes
			.filter((a) => {
				const gid = a.attribute_group_id ?? null;
				if (groupId === null) return gid === null;
				return gid === groupId;
			})
			.sort((a, b) => a.sort_order - b.sort_order);
	}

	async function load() {
		loading = true;
		error = null;
		try {
			const [attrData, groupData, mech] = await Promise.all([
				listAttributes(systemId),
				listAttributeGroups(systemId),
				getMechanics(systemId)
			]);
			attributes = attrData.items ?? [];
			groups = groupData.items ?? [];
			enabledDerived = normalizeAttributesConfig(mech.attributes_config).enabled_derived;
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to load attributes';
		} finally {
			loading = false;
		}
	}

	async function onToggleDerived(checked: boolean) {
		const prev = enabledDerived;
		enabledDerived = checked;
		toggleSaving = true;
		error = null;
		try {
			await saveAttributesConfig(systemId, { enabled_derived: checked });
		} catch (e) {
			enabledDerived = prev;
			error = e instanceof ApiError ? e.message : 'Failed to save derived attributes setting';
		} finally {
			toggleSaving = false;
		}
	}

	function openCreateGroup() {
		editingGroup = null;
		groupModalOpen = true;
	}

	function openEditGroup(id: string) {
		if (id === UNGROUPED_ID) return;
		const g = groups.find((x) => x.id === id);
		if (!g) return;
		editingGroup = g;
		groupModalOpen = true;
	}

	function closeGroupModal() {
		groupModalOpen = false;
		editingGroup = null;
	}

	function openCreateCore(groupId: string | null = null) {
		editingAttribute = null;
		defaultGroupId = groupId && groupId !== UNGROUPED_ID ? groupId : '';
		coreModalOpen = true;
	}

	function openEditCore(id: string) {
		const attr = attributes.find((a) => a.id === id);
		if (!attr || attr.config.is_derived) return;
		editingAttribute = attr;
		defaultGroupId = '';
		coreModalOpen = true;
	}

	function openCreateDerived() {
		editingAttribute = null;
		derivedModalOpen = true;
	}

	function openEditDerived(id: string) {
		const attr = attributes.find((a) => a.id === id);
		if (!attr || !attr.config.is_derived) return;
		editingAttribute = attr;
		derivedModalOpen = true;
	}

	function closeCoreModal() {
		coreModalOpen = false;
		if (!derivedModalOpen) editingAttribute = null;
	}

	function closeDerivedModal() {
		derivedModalOpen = false;
		if (!coreModalOpen) editingAttribute = null;
	}

	async function handleGroupReorder(items: DragListItem[]) {
		const real = items.filter((i) => i.id !== UNGROUPED_ID);
		const updates: { id: string; sort_order: number }[] = [];
		for (const item of real) {
			const g = groups.find((x) => x.id === item.id);
			if (g && g.sort_order !== item.sortOrder) {
				updates.push({ id: item.id, sort_order: item.sortOrder });
			}
		}
		if (updates.length === 0) return;
		reordering = true;
		try {
			await Promise.all(
				updates.map((u) => updateAttributeGroup(systemId, u.id, { sort_order: u.sort_order }))
			);
			await load();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to reorder groups';
		} finally {
			reordering = false;
		}
	}

	async function handleReorderSubset(
		items: DragListItem[],
		subset: AttributeResponse[]
	) {
		const idSet = new Set(subset.map((a) => a.id));
		const updates: { id: string; sort_order: number }[] = [];
		for (const item of items) {
			if (!idSet.has(item.id)) continue;
			const prev = subset.find((a) => a.id === item.id)?.sort_order;
			if (prev !== item.sortOrder) {
				updates.push({ id: item.id, sort_order: item.sortOrder });
			}
		}
		if (updates.length === 0) return;
		reordering = true;
		try {
			await Promise.all(
				updates.map((u) => updateAttribute(systemId, u.id, { sort_order: u.sort_order }))
			);
			await load();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to reorder attributes';
			await load();
		} finally {
			reordering = false;
		}
	}

	async function moveAttributeInGroup(
		attrId: string,
		groupKey: string,
		direction: -1 | 1
	) {
		const groupId = groupKey === UNGROUPED_ID ? null : groupKey;
		const list = coreInGroup(groupId);
		const idx = list.findIndex((a) => a.id === attrId);
		const target = idx + direction;
		if (idx < 0 || target < 0 || target >= list.length) return;
		const order = list.map((a) => a.id);
		order.splice(idx, 1);
		order.splice(target, 0, attrId);
		const updates = order.map((id, sort_order) => ({ id, sort_order }));
		reordering = true;
		try {
			await Promise.all(
				updates.map((u) => updateAttribute(systemId, u.id, { sort_order: u.sort_order }))
			);
			await load();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to reorder attributes';
		} finally {
			reordering = false;
		}
	}

	async function handleDeleteGroup(id: string) {
		if (id === UNGROUPED_ID) return;
		const g = groups.find((x) => x.id === id);
		if (!g) return;
		if (!confirm(`Delete group "${g.name}"? Attributes will move to Ungrouped.`)) return;
		error = null;
		try {
			await deleteAttributeGroup(systemId, id);
			await load();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to delete group';
		}
	}

	async function handleDeleteAttribute(id: string) {
		const attr = attributes.find((a) => a.id === id);
		if (!attr) return;
		if (!confirm(`Delete attribute "${attr.name}"?`)) return;
		error = null;
		try {
			await deleteAttribute(systemId, id);
			await load();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to delete attribute';
		}
	}

	$effect(() => {
		onModalOpenChange?.(coreModalOpen || derivedModalOpen || groupModalOpen);
	});

	onMount(() => {
		load();
	});
</script>

<div class="attributes-tab">
	{#if loading}
		<p class="status-msg">Loading attributes…</p>
	{:else}
		{#if error}
			<p class="form-error">{error}</p>
		{/if}
		{#if reordering}
			<p class="status-msg">Saving order…</p>
		{/if}

		<section class="attr-section">
			<div class="section-head">
				<h3 class="section-title">Attribute Groups</h3>
				<button type="button" class="btn-secondary-sm" {disabled} onclick={openCreateGroup}>
					Create Group
				</button>
			</div>
			<TE_DragList
				items={groupDragItems}
				disabled={disabled || reordering}
				onReorder={handleGroupReorder}
				onEdit={openEditGroup}
				onDelete={handleDeleteGroup}
				renderBody={groupBody}
			/>
		</section>

		<label class="toggle-row">
			<input
				type="checkbox"
				checked={enabledDerived}
				disabled={disabled || toggleSaving}
				onchange={(e) => onToggleDerived((e.currentTarget as HTMLInputElement).checked)}
			/>
			<span>Enable Derived Attributes</span>
		</label>

		{#if enabledDerived}
			<section class="attr-section">
				<h3 class="section-title">Derived Attributes</h3>
				<TE_DragList
					items={derivedDragItems}
					disabled={disabled || reordering}
					onReorder={(items) => handleReorderSubset(items, derivedAttributes)}
					onEdit={openEditDerived}
					onDelete={handleDeleteAttribute}
				/>
				<div class="tab-footer">
					<button type="button" class="btn-primary" {disabled} onclick={openCreateDerived}>
						Create Derived Attribute
					</button>
				</div>
			</section>
		{/if}
	{/if}
</div>

{#snippet groupBody(groupKey: string)}
	<ul class="attr-rows">
		{#each coreInGroup(groupKey === UNGROUPED_ID ? null : groupKey) as attr (attr.id)}
			<li class="attr-row">
				<span class="attr-row__name">{attr.name}</span>
				<span class="attr-row__type">{attributeTypeLabel(attr.type)}</span>
				<span class="attr-row__actions">
					{#if !disabled}
						<button
							type="button"
							class="btn-icon"
							title="Move up"
							aria-label="Move up"
							disabled={reordering}
							onclick={() => moveAttributeInGroup(attr.id, groupKey, -1)}
						>
							↑
						</button>
						<button
							type="button"
							class="btn-icon"
							title="Move down"
							aria-label="Move down"
							disabled={reordering}
							onclick={() => moveAttributeInGroup(attr.id, groupKey, 1)}
						>
							↓
						</button>
					{/if}
					<button
						type="button"
						class="btn-sm"
						{disabled}
						onclick={() => openEditCore(attr.id)}
					>
						Edit
					</button>
					<button
						type="button"
						class="btn-sm btn-danger-outline"
						{disabled}
						onclick={() => handleDeleteAttribute(attr.id)}
					>
						Delete
					</button>
				</span>
			</li>
		{:else}
			<li class="attr-row attr-row--empty">No attributes in this group.</li>
		{/each}
	</ul>
	{#if groupKey !== UNGROUPED_ID && !disabled}
		<button
			type="button"
			class="btn-secondary-sm attr-add"
			onclick={() => openCreateCore(groupKey)}
		>
			Create Attribute
		</button>
	{:else if groupKey === UNGROUPED_ID && !disabled}
		<button type="button" class="btn-secondary-sm attr-add" onclick={() => openCreateCore(null)}>
			Create Attribute
		</button>
	{/if}
{/snippet}

<TE_AttributeGroupModal
	open={groupModalOpen}
	{systemId}
	group={editingGroup}
	siblingGroupCount={groups.length}
	{disabled}
	onclose={closeGroupModal}
	onsaved={load}
/>

<TE_AttributeModal
	open={coreModalOpen}
	{systemId}
	attribute={editingAttribute}
	siblingAttributes={attributes}
	attributeGroups={groups}
	{defaultGroupId}
	{disabled}
	onclose={closeCoreModal}
	onsaved={load}
/>

<TE_DerivedAttributeModal
	open={derivedModalOpen}
	{systemId}
	attribute={editingAttribute}
	siblingAttributes={attributes}
	{disabled}
	onclose={closeDerivedModal}
	onsaved={load}
/>

<style>
	.attributes-tab {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.attr-section {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.section-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
	}

	.section-title {
		margin: 0;
		font-size: 0.95rem;
		font-weight: 600;
	}

	.toggle-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
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

	.tab-footer {
		margin-top: 0.25rem;
	}

	.btn-primary,
	.btn-secondary-sm {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		font: inherit;
	}

	.btn-primary {
		border: none;
		background: var(--accent-gm, #c9a227);
		color: #1a1a1a;
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary-sm {
		padding: 0.35rem 0.75rem;
		font-size: 0.85rem;
		border: 1px solid #d1d5db;
		background: #fff;
	}

	.attr-rows {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
	}

	.attr-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.35rem 0.5rem;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		background: #fff;
		font-size: 0.875rem;
	}

	.attr-row--empty {
		color: var(--text-muted, #6b7280);
		font-style: italic;
	}

	.attr-row__name {
		flex: 1;
		font-weight: 600;
	}

	.attr-row__type {
		color: var(--text-muted, #6b7280);
	}

	.attr-row__actions {
		display: flex;
		align-items: center;
		gap: 0.35rem;
	}

	.btn-icon {
		padding: 0.2rem 0.45rem;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		background: #fff;
		cursor: pointer;
		font: inherit;
		line-height: 1;
	}

	.btn-sm {
		padding: 0.25rem 0.55rem;
		font-size: 0.8rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.btn-danger-outline {
		border-color: #fca5a5;
		color: #b91c1c;
	}

	.attr-add {
		margin-top: 0.35rem;
	}
</style>

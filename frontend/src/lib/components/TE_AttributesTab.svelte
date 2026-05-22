<script lang="ts">
	import { onMount } from 'svelte';
	import TE_AttributeModal from './TE_AttributeModal.svelte';
	import TE_DragList, { type DragListItem } from './TE_DragList.svelte';
	import { deleteAttribute, listAttributes, updateAttribute } from '$lib/api/mechanics';
	import { ApiError } from '$lib/api/client';
	import {
		attributeTypeLabel,
		modifierDisplayLabel
	} from '$lib/utils/attributeDefaults';
	import type { AttributeResponse } from '$lib/types/mechanics';

	interface Props {
		systemId: string;
		disabled?: boolean;
		onModalOpenChange?: (open: boolean) => void;
	}

	let { systemId, disabled = false, onModalOpenChange }: Props = $props();

	export function closeModalDiscard() {
		closeModal();
	}

	let attributes = $state<AttributeResponse[]>([]);
	let loading = $state(true);
	let reordering = $state(false);
	let error = $state<string | null>(null);
	let modalOpen = $state(false);
	let editingAttribute = $state<AttributeResponse | null>(null);

	const sortedAttributes = $derived(
		[...attributes].sort((a, b) => a.sort_order - b.sort_order)
	);

	const dragItems = $derived.by((): DragListItem[] =>
		sortedAttributes.map((attr) => ({
			id: attr.id,
			sortOrder: attr.sort_order,
			header: attr.name,
			subtitle: attr.group_name?.trim() || 'Ungrouped',
			summary: attributeSummary(attr)
		}))
	);

	function attributeSummary(attr: AttributeResponse): string {
		const parts = [attributeTypeLabel(attr.type)];
		if (attr.config.is_derived) parts.push('Derived');
		if (attr.config.modifier_display) {
			parts.push(modifierDisplayLabel(attr.config.modifier_display));
		}
		return parts.join(' · ');
	}

	async function load() {
		loading = true;
		error = null;
		try {
			const data = await listAttributes(systemId);
			attributes = data.items ?? [];
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to load attributes';
		} finally {
			loading = false;
		}
	}

	function openCreate() {
		editingAttribute = null;
		modalOpen = true;
	}

	function openEdit(id: string) {
		const attr = attributes.find((a) => a.id === id);
		if (!attr) return;
		editingAttribute = attr;
		modalOpen = true;
	}

	function closeModal() {
		modalOpen = false;
		editingAttribute = null;
	}

	async function handleReorder(items: DragListItem[]) {
		const prev = new Map(attributes.map((a) => [a.id, a.sort_order]));
		const updates: { id: string; sort_order: number }[] = [];
		for (const item of items) {
			if (prev.get(item.id) !== item.sortOrder) {
				updates.push({ id: item.id, sort_order: item.sortOrder });
			}
		}
		if (updates.length === 0) return;

		attributes = items
			.map((item) => {
				const attr = attributes.find((a) => a.id === item.id)!;
				return { ...attr, sort_order: item.sortOrder };
			})
			.sort((a, b) => a.sort_order - b.sort_order);

		reordering = true;
		error = null;
		try {
			await Promise.all(
				updates.map((u) => updateAttribute(systemId, u.id, { sort_order: u.sort_order }))
			);
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to reorder attributes';
			await load();
		} finally {
			reordering = false;
		}
	}

	async function handleDelete(id: string) {
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
		onModalOpenChange?.(modalOpen);
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

		<TE_DragList
			items={dragItems}
			disabled={disabled || reordering}
			onReorder={handleReorder}
			onEdit={openEdit}
			onDelete={handleDelete}
		/>

		<div class="tab-footer">
			<button type="button" class="btn-primary" {disabled} onclick={openCreate}>
				Create Attribute
			</button>
		</div>
	{/if}
</div>

<TE_AttributeModal
	open={modalOpen}
	{systemId}
	attribute={editingAttribute}
	siblingAttributes={attributes}
	{disabled}
	onclose={closeModal}
	onsaved={load}
/>

<style>
	.attributes-tab {
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

	.tab-footer {
		margin-top: 0.5rem;
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
</style>

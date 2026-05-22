<script lang="ts">
	import type { Snippet } from 'svelte';

	export interface DragListItem {
		id: string;
		sortOrder: number;
		header: string;
		subtitle?: string;
		summary?: string;
		body?: Snippet;
	}

	interface Props {
		items: DragListItem[];
		disabled?: boolean;
		onReorder?: (items: DragListItem[]) => void;
		onEdit?: (id: string) => void;
		onDelete?: (id: string) => void;
	}

	let { items, disabled = false, onReorder, onEdit, onDelete }: Props = $props();

	let dragId = $state<string | null>(null);
	let dropTargetId = $state<string | null>(null);

	const sortedItems = $derived([...items].sort((a, b) => a.sortOrder - b.sortOrder));

	function reorderList(fromId: string, toId: string) {
		if (fromId === toId || disabled) return;
		const order = sortedItems.map((i) => i.id);
		const fromIdx = order.indexOf(fromId);
		const toIdx = order.indexOf(toId);
		if (fromIdx < 0 || toIdx < 0) return;
		order.splice(fromIdx, 1);
		order.splice(toIdx, 0, fromId);
		const byId = new Map(items.map((i) => [i.id, i]));
		const reordered = order.map((id, index) => {
			const item = byId.get(id)!;
			return { ...item, sortOrder: index };
		});
		onReorder?.(reordered);
	}

	function moveByKeyboard(id: string, direction: -1 | 1) {
		const order = sortedItems.map((i) => i.id);
		const idx = order.indexOf(id);
		const target = idx + direction;
		if (target < 0 || target >= order.length) return;
		reorderList(id, order[target]);
	}

	function handleDragStart(e: DragEvent, id: string) {
		if (disabled) return;
		dragId = id;
		e.dataTransfer?.setData('text/plain', id);
		if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move';
	}

	function handleDragOver(e: DragEvent, id: string) {
		if (disabled || !dragId) return;
		e.preventDefault();
		if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
		dropTargetId = id;
	}

	function handleDrop(e: DragEvent, id: string) {
		e.preventDefault();
		if (disabled || !dragId) return;
		reorderList(dragId, id);
		dragId = null;
		dropTargetId = null;
	}

	function handleDragEnd() {
		dragId = null;
		dropTargetId = null;
	}
</script>

<ul class="drag-list" role="list">
	{#each sortedItems as item (item.id)}
		<li
			class="drag-list__item"
			class:drag-over={dropTargetId === item.id && dragId !== item.id}
			role="listitem"
			ondragover={(e) => handleDragOver(e, item.id)}
			ondrop={(e) => handleDrop(e, item.id)}
		>
			<details class="drag-card" open>
				<summary class="drag-card__summary">
					<span
						class="drag-handle"
						title="Drag to reorder"
						draggable={!disabled}
						role="button"
						tabindex={disabled ? -1 : 0}
						aria-label="Drag to reorder"
						ondragstart={(e) => handleDragStart(e, item.id)}
						ondragend={handleDragEnd}
					>
						<i class="fas fa-grip-vertical" aria-hidden="true"></i>
					</span>
					<span class="drag-card__titles">
						<span class="drag-card__header">{item.header}</span>
						{#if item.subtitle}
							<span class="drag-card__subtitle">{item.subtitle}</span>
						{/if}
					</span>
					<span class="drag-card__actions" role="group" aria-label="Item actions">
						{#if !disabled}
							<button
								type="button"
								class="btn-icon"
								title="Move up"
								aria-label="Move up"
								onclick={(e) => {
									e.preventDefault();
									moveByKeyboard(item.id, -1);
								}}
							>
								<i class="fas fa-arrow-up" aria-hidden="true"></i>
							</button>
							<button
								type="button"
								class="btn-icon"
								title="Move down"
								aria-label="Move down"
								onclick={(e) => {
									e.preventDefault();
									moveByKeyboard(item.id, 1);
								}}
							>
								<i class="fas fa-arrow-down" aria-hidden="true"></i>
							</button>
						{/if}
						{#if onEdit}
							<button
								type="button"
								class="btn-sm"
								{disabled}
								onclick={(e) => {
									e.preventDefault();
									onEdit(item.id);
								}}
							>
								Edit
							</button>
						{/if}
						{#if onDelete}
							<button
								type="button"
								class="btn-sm btn-danger-outline"
								{disabled}
								onclick={(e) => {
									e.preventDefault();
									onDelete(item.id);
								}}
							>
								Delete
							</button>
						{/if}
					</span>
				</summary>
				{#if item.body}
					<div class="drag-card__body">
						{@render item.body()}
					</div>
				{:else if item.summary}
					<div class="drag-card__body">
						<p class="drag-card__summary-text">{item.summary}</p>
					</div>
				{/if}
			</details>
		</li>
	{/each}
</ul>

{#if sortedItems.length === 0}
	<p class="drag-list__empty">No items yet.</p>
{/if}

<style>
	.drag-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.drag-list__item {
		margin: 0;
	}

	.drag-list__item.drag-over .drag-card {
		outline: 2px dashed var(--accent-gm, #c9a227);
	}

	.drag-card {
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		background: #fafafa;
		overflow: hidden;
	}

	.drag-card__summary {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.6rem 0.75rem;
		cursor: pointer;
		list-style: none;
	}

	.drag-card__summary::-webkit-details-marker {
		display: none;
	}

	.drag-handle {
		display: flex;
		align-items: center;
		padding: 0.25rem;
		color: var(--text-muted, #6b7280);
		cursor: grab;
		flex-shrink: 0;
	}

	.drag-handle:active {
		cursor: grabbing;
	}

	.drag-card__titles {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.drag-card__header {
		font-weight: 600;
		font-size: 0.95rem;
	}

	.drag-card__subtitle {
		font-size: 0.8rem;
		color: var(--text-muted, #6b7280);
	}

	.drag-card__actions {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.35rem;
		flex-shrink: 0;
	}

	.drag-card__body {
		padding: 0 0.75rem 0.75rem;
		font-size: 0.85rem;
		color: var(--text-muted, #6b7280);
	}

	.drag-card__summary-text {
		margin: 0;
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

	.btn-sm:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-danger-outline {
		border-color: #fca5a5;
		color: #b91c1c;
	}

	.btn-icon {
		border: none;
		background: none;
		padding: 0.25rem;
		cursor: pointer;
		color: var(--text-muted, #6b7280);
	}

	.btn-icon:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.drag-list__empty {
		margin: 0;
		color: var(--text-muted, #6b7280);
		font-size: 0.9rem;
	}
</style>

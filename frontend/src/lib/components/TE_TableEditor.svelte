<script lang="ts">
	import type { Snippet } from 'svelte';

	export interface SelectOption {
		value: string;
		label: string;
	}

	export interface TableColumn {
		key: string;
		label: string;
		type?: 'text' | 'number' | 'select';
		options?: SelectOption[];
	}

	interface RowDetailArgs {
		row: Record<string, string | number>;
		index: number;
	}

	interface Props {
		columns: TableColumn[];
		rows: Record<string, string | number>[];
		disabled?: boolean;
		addLabel?: string;
		minRows?: number;
		expandable?: boolean;
		rowDetail?: Snippet<[RowDetailArgs]>;
		onadd?: () => void;
		onremove?: (index: number) => void;
	}

	let {
		columns,
		rows = $bindable([]),
		disabled = false,
		addLabel = 'Add Row',
		minRows = 0,
		expandable = false,
		rowDetail,
		onadd,
		onremove
	}: Props = $props();

	let expanded = $state<Set<number>>(new Set());

	const colSpan = $derived(
		columns.length + (expandable ? 1 : 0) + (onremove && !disabled ? 1 : 0)
	);

	function cellValue(row: Record<string, string | number>, key: string): string | number {
		return row[key] ?? (columns.find((c) => c.key === key)?.type === 'number' ? 0 : '');
	}

	function setCell(index: number, key: string, value: string) {
		const col = columns.find((c) => c.key === key);
		const next = [...rows];
		const row = { ...next[index] };
		if (col?.type === 'number') {
			row[key] = value === '' ? 0 : Number(value);
		} else {
			row[key] = value;
		}
		next[index] = row;
		rows = next;
	}

	function toggleExpanded(index: number) {
		const next = new Set(expanded);
		if (next.has(index)) next.delete(index);
		else next.add(index);
		expanded = next;
	}

	function isExpanded(index: number): boolean {
		return expanded.has(index);
	}

	$effect(() => {
		if (!expandable || rows.length === 0) return;
		const valid = new Set<number>();
		for (const i of expanded) {
			if (i < rows.length) valid.add(i);
		}
		if (valid.size === 0 && rows.length > 0) {
			valid.add(0);
		}
		if (valid.size !== expanded.size || [...valid].some((i) => !expanded.has(i))) {
			expanded = valid;
		}
	});
</script>

<div class="table-editor">
	<div class="mechanics-table-wrap">
		<table class="mechanics-table">
			<thead>
				<tr>
					{#if expandable}
						<th scope="col" class="col-expand"><span class="sr-only">Expand</span></th>
					{/if}
					{#each columns as col}
						<th scope="col">{col.label}</th>
					{/each}
					{#if onremove && !disabled}
						<th scope="col"><span class="sr-only">Actions</span></th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each rows as row, i (i)}
					<tr>
						{#if expandable}
							<td class="col-expand">
								{#if rowDetail}
									<button
										type="button"
										class="btn-expand"
										aria-expanded={isExpanded(i)}
										aria-label={isExpanded(i) ? 'Collapse row details' : 'Expand row details'}
										{disabled}
										onclick={() => toggleExpanded(i)}
									>
										{isExpanded(i) ? '▼' : '▶'}
									</button>
								{/if}
							</td>
						{/if}
						{#each columns as col}
							<td>
								{#if col.type === 'select' && col.options}
									<select
										value={String(cellValue(row, col.key))}
										{disabled}
										aria-label={col.label}
										onchange={(e) =>
											setCell(i, col.key, (e.currentTarget as HTMLSelectElement).value)}
									>
										{#each col.options as opt}
											<option value={opt.value}>{opt.label}</option>
										{/each}
									</select>
								{:else}
									<input
										type={col.type === 'number' ? 'number' : 'text'}
										value={cellValue(row, col.key)}
										{disabled}
										aria-label={col.label}
										oninput={(e) =>
											setCell(i, col.key, (e.currentTarget as HTMLInputElement).value)}
									/>
								{/if}
							</td>
						{/each}
						{#if onremove && !disabled}
							<td>
								<button
									type="button"
									class="btn-row-remove"
									disabled={rows.length <= minRows}
									title={rows.length <= minRows ? 'At least one row is required' : undefined}
									onclick={() => onremove(i)}
								>
									Remove
								</button>
							</td>
						{/if}
					</tr>
					{#if expandable && rowDetail && isExpanded(i)}
						<tr class="slot-details">
							<td colspan={colSpan}>
								{@render rowDetail({ row, index: i })}
							</td>
						</tr>
					{/if}
				{/each}
			</tbody>
		</table>
	</div>
	{#if onadd && !disabled}
		<button type="button" class="btn-add" onclick={onadd}>{addLabel}</button>
	{/if}
</div>

<style>
	.table-editor {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
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

	.mechanics-table input,
	.mechanics-table select {
		width: 100%;
		min-width: 4rem;
		padding: 0.4rem 0.55rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.col-expand {
		width: 2.5rem;
	}

	.btn-expand {
		padding: 0.2rem 0.4rem;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		background: #fff;
		cursor: pointer;
		font: inherit;
		line-height: 1;
	}

	.slot-details td {
		background: #f9fafb;
		padding: 0.75rem 1rem;
	}

	.btn-add {
		align-self: flex-start;
		padding: 0.35rem 0.75rem;
		font-size: 0.85rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.btn-row-remove {
		padding: 0.25rem 0.55rem;
		font-size: 0.8rem;
		border: 1px solid #fca5a5;
		border-radius: 6px;
		background: #fff;
		color: #b91c1c;
		cursor: pointer;
		font: inherit;
	}

	.btn-row-remove:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}
</style>

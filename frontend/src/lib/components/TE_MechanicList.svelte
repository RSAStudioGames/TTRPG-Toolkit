<script lang="ts">
	import { MECHANIC_TYPE_OPTIONS } from '$lib/constants/resolutionOptions';
	import type { AdvantageDisadvantageEntry } from '$lib/types/mechanics';

	interface Props {
		entries: AdvantageDisadvantageEntry[];
		disabled?: boolean;
	}

	let { entries = $bindable([]), disabled = false }: Props = $props();

	function addEntry() {
		entries = [...entries, { name: '', mechanic_type: 'add_modifier' }];
	}

	function removeEntry(index: number) {
		entries = entries.filter((_, i) => i !== index);
	}
</script>

<div class="mechanic-list">
	<div class="mechanics-table-wrap">
		<table class="mechanics-table">
			<thead>
				<tr>
					<th scope="col">Name</th>
					<th scope="col">Mechanic Type</th>
					{#if !disabled}
						<th scope="col"><span class="sr-only">Actions</span></th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each entries as entry, i (i)}
					<tr>
						<td>
							<input
								type="text"
								bind:value={entry.name}
								placeholder="Name"
								{disabled}
								aria-label="Mechanic name"
							/>
						</td>
						<td>
							<select bind:value={entry.mechanic_type} {disabled} aria-label="Mechanic type">
								{#each MECHANIC_TYPE_OPTIONS as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</select>
						</td>
						{#if !disabled}
							<td>
								<button type="button" class="btn-row-remove" onclick={() => removeEntry(i)}>
									Remove
								</button>
							</td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
	{#if !disabled}
		<button type="button" class="btn-add" onclick={addEntry}>Add Mechanic</button>
	{/if}
</div>

<style>
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

	.mechanics-table-wrap {
		overflow-x: auto;
		margin-bottom: 0.75rem;
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
		min-width: 8rem;
		padding: 0.4rem 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font: inherit;
	}

	.btn-add,
	.btn-row-remove {
		padding: 0.35rem 0.75rem;
		font-size: 0.85rem;
		border-radius: 6px;
		cursor: pointer;
		font: inherit;
	}

	.btn-add {
		border: 1px dashed #9ca3af;
		background: #fff;
		color: #374151;
	}

	.btn-row-remove {
		border: 1px solid #e5e7eb;
		background: #fff;
		color: #6b7280;
	}
</style>

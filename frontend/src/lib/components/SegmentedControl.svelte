<script lang="ts">
	interface Props {
		label: string;
		value: number;
		disabled?: boolean;
	}

	let { label, value = $bindable(), disabled = false }: Props = $props();

	const labels: Record<number, string> = {
		1: 'Light',
		3: 'Moderate',
		5: 'Heavy'
	};
</script>

<div class="segmented">
	<span class="seg-label">{label}</span>
	<div class="seg-track" role="group" aria-label={label}>
		{#each [1, 2, 3, 4, 5] as n}
			<button
				type="button"
				class:selected={value === n}
				{disabled}
				aria-pressed={value === n}
				title={labels[n] ?? String(n)}
				onclick={() => (value = n)}
			>
				<span class="num">{n}</span>
				{#if labels[n]}
					<span class="sub">{labels[n]}</span>
				{/if}
			</button>
		{/each}
	</div>
</div>

<style>
	.segmented {
		margin-bottom: 0.85rem;
	}

	.seg-label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		margin-bottom: 0.4rem;
	}

	.seg-track {
		display: flex;
		gap: 0.25rem;
	}

	.seg-track button {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 0.4rem 0.25rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: #fff;
		cursor: pointer;
		font: inherit;
	}

	.seg-track button.selected {
		border-color: var(--accent-gm);
		background: var(--accent-gm-muted);
		font-weight: 600;
	}

	.num {
		font-size: 0.95rem;
	}

	.sub {
		font-size: 0.65rem;
		color: var(--text-muted);
	}
</style>

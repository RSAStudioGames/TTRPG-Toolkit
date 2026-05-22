<script lang="ts">
	interface Props {
		label: string;
		name: string;
		options: string[];
		value: string;
		disabled?: boolean;
		onchange?: (value: string) => void;
	}

	let { label, name, options, value = $bindable(), disabled = false, onchange }: Props = $props();
</script>

<fieldset class="radio-group" {disabled}>
	<legend>{label}</legend>
	<div class="options">
		{#each options as opt}
			<label class="radio-option">
				<input
					type="radio"
					{name}
					value={opt}
					checked={value === opt}
					{disabled}
					onchange={() => {
						value = opt;
						onchange?.(opt);
					}}
				/>
				<span>{opt}</span>
			</label>
		{/each}
	</div>
</fieldset>

<style>
	.radio-group {
		border: none;
		margin: 0 0 0.85rem;
		padding: 0;
	}

	legend {
		font-size: 0.875rem;
		font-weight: 600;
		margin-bottom: 0.4rem;
	}

	.options {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem 1rem;
	}

	.radio-option {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		font-size: 0.9rem;
		cursor: pointer;
	}

	.radio-option input {
		accent-color: var(--accent-gm);
	}
</style>

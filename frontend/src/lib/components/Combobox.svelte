<script lang="ts">
	interface Props {
		label: string;
		placeholder?: string;
		options?: string[];
		value: string;
		disabled?: boolean;
	}

	let {
		label,
		placeholder = '',
		options = [],
		value = $bindable(),
		disabled = false
	}: Props = $props();

	let listId = `cb-${Math.random().toString(36).slice(2, 9)}`;
</script>

<div class="combobox form-field">
	<label for={listId}>{label}</label>
	<input
		id={listId}
		type="text"
		list="{listId}-opts"
		{placeholder}
		bind:value
		{disabled}
	/>
	<datalist id="{listId}-opts">
		{#each options as opt}
			<option value={opt}></option>
		{/each}
	</datalist>
</div>

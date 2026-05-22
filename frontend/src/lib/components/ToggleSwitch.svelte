<script lang="ts">
	interface Props {
		label: string;
		checked: boolean;
		disabled?: boolean;
	}

	let { label, checked = $bindable(), disabled = false }: Props = $props();
</script>

<label class="toggle-row">
	<span>{label}</span>
	<input type="checkbox" role="switch" bind:checked {disabled} />
</label>

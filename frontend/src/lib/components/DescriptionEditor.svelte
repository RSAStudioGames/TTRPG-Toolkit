<script lang="ts">
	interface Props {
		value: string;
		disabled?: boolean;
	}

	let { value = $bindable(), disabled = false }: Props = $props();
</script>

<div class="form-field">
	<label for="system-desc">Description</label>
	<textarea id="system-desc" bind:value {disabled} placeholder="Describe this game system…"></textarea>
</div>

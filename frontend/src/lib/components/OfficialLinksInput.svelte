<script lang="ts">
	interface Props {
		links: string[];
		disabled?: boolean;
	}

	let { links = $bindable(), disabled = false }: Props = $props();

	function addLink() {
		links = [...links, ''];
	}

	function removeLink(i: number) {
		links = links.filter((_, idx) => idx !== i);
		if (links.length === 0) links = [''];
	}
</script>

<div class="form-field">
	<span class="field-label">Official Links</span>
	{#each links as link, i}
		<div class="link-row">
			<input type="url" placeholder="https://..." bind:value={links[i]} {disabled} />
			{#if links.length > 1}
				<button type="button" class="btn-text" {disabled} onclick={() => removeLink(i)}>Remove</button>
			{/if}
		</div>
	{/each}
	<button type="button" class="btn-secondary-sm" {disabled} onclick={addLink}>Add Link</button>
</div>

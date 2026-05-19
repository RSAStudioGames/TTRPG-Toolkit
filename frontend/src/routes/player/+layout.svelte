<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();
</script>

<div class="role-layout" style="--accent-primary: var(--accent-player)">
	{@render children()}
</div>

<style>
	.role-layout {
		min-height: 100%;
	}
</style>

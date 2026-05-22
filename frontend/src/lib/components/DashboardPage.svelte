<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		headerActions?: Snippet;
		children: Snippet;
	}

	let { title, headerActions, children }: Props = $props();
</script>

<div class="dashboard">
	<header class="dashboard-header">
		<h1>{title}</h1>
		{#if headerActions}
			<div class="header-actions">
				{@render headerActions()}
			</div>
		{/if}
	</header>

	<main class="cards-row">
		{@render children()}
	</main>
</div>

<style>
	.dashboard {
		display: flex;
		flex-direction: column;
		align-items: center;
		min-height: 100vh;
		padding: 3.5rem 1.5rem 3rem;
		width: 100%;
		max-width: 1100px;
		margin: 0 auto;
	}

	.dashboard-header {
		display: flex;
		width: 100%;
		align-items: flex-start;
		justify-content: space-between;
		gap: 1rem;
	}

	.dashboard-header h1 {
		margin: 0;
		font-size: clamp(1.75rem, 4vw, 2.25rem);
		font-weight: 700;
		color: var(--text-ink);
	}

	.header-actions {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		flex-shrink: 0;
	}

	.cards-row {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 2rem;
		align-items: stretch;
		width: 100%;
		margin-top: 3rem;
		justify-items: center;
	}

	@media (max-width: 768px) {
		.cards-row {
			grid-template-columns: 1fr;
		}
	}
</style>

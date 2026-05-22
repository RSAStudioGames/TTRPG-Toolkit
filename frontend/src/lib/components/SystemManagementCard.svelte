<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		iconClass: string;
		title: string;
		description: string;
		children: Snippet;
	}

	let { iconClass, title, description, children }: Props = $props();
</script>

<article
	class="system-mgmt-card"
	style="--accent-primary: var(--accent-gm); background: var(--card-tint-gm)"
>
	<div class="card-icon" aria-hidden="true">
		<i class={iconClass}></i>
	</div>
	<h2 class="card-title">{title}</h2>
	<p class="card-description">{description}</p>
	<div class="cta-slot">
		{@render children()}
	</div>
</article>

<style>
	.system-mgmt-card {
		display: flex;
		flex-direction: column;
		align-items: center;
		height: 100%;
		width: 100%;
		min-width: 280px;
		max-width: 520px;
		padding: 2rem;
		border-radius: 12px;
		border: 2px solid transparent;
		box-shadow: var(--shadow-card);
		text-align: center;
		transition: var(--transition-card);
	}

	.system-mgmt-card:hover {
		transform: translateY(-6px);
		border-color: #f0c94a;
		box-shadow: var(--shadow-card-hover-gm);
	}

	.card-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		font-size: 2.5rem;
		margin-bottom: 1rem;
		color: var(--accent-primary);
	}

	.card-title {
		width: 100%;
		margin: 0 0 0.75rem;
		font-size: 1.5rem;
		font-weight: 700;
		text-align: center;
		color: var(--text-ink);
	}

	.card-description {
		width: 100%;
		margin: 0;
		font-size: 0.9375rem;
		line-height: 1.6;
		text-align: center;
		color: var(--text-muted);
	}

	.cta-slot {
		width: 100%;
		margin-top: auto;
		padding-top: 1.75rem;
	}
</style>

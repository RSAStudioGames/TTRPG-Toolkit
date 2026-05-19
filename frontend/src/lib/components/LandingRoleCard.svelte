<script lang="ts">
	import type { Snippet } from 'svelte';

	type Role = 'player' | 'gm';

	interface Props {
		role: Role;
		iconClass: string;
		title: string;
		description: string;
		children: Snippet;
	}

	let { role, iconClass, title, description, children }: Props = $props();

	const accentVar = $derived(
		role === 'player' ? 'var(--accent-player)' : 'var(--accent-gm)'
	);
	const cardTint = $derived(
		role === 'player' ? 'var(--card-tint-player)' : 'var(--card-tint-gm)'
	);
</script>

<article
	class="role-card"
	class:role-player={role === 'player'}
	class:role-gm={role === 'gm'}
	style="--accent-primary: {accentVar}; background: {cardTint}"
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
	.role-card {
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

	.role-card.role-player:hover {
		transform: translateY(-6px);
		border-color: #ef4444;
		box-shadow: var(--shadow-card-hover-player);
	}

	.role-card.role-gm:hover {
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

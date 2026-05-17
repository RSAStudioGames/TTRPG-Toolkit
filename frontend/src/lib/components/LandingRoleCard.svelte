<script lang="ts">
	import { motion } from '@humanspeak/svelte-motion';
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
	const hoverShadow = $derived(
		role === 'player' ? 'var(--shadow-card-hover-player)' : 'var(--shadow-card-hover-gm)'
	);
	const cardBorder = $derived(
		role === 'player' ? 'var(--card-border-player)' : 'var(--card-border-gm)'
	);
	const cardBorderHover = $derived(
		role === 'player' ? 'var(--card-border-hover-player)' : 'var(--card-border-hover-gm)'
	);

	const cardTransition = { type: 'spring' as const, stiffness: 400, damping: 30 };
</script>

<motion.article
	class="role-card"
	style="--accent-primary: {accentVar}; background: {cardTint}"
	initial={{ y: 0, boxShadow: 'var(--shadow-card)', border: cardBorder }}
	whileHover={{ y: -6, boxShadow: hoverShadow, border: cardBorderHover }}
	transition={cardTransition}
>
	<div class="card-icon" aria-hidden="true">
		<i class={iconClass}></i>
	</div>
	<h2 class="card-title">{title}</h2>
	<p class="card-description">{description}</p>
	<div class="cta-slot">
		{@render children()}
	</div>
</motion.article>

<style>
	:global(.role-card) {
		display: flex;
		flex-direction: column;
		align-items: center;
		height: 100%;
		width: 100%;
		min-width: 280px;
		max-width: 520px;
		padding: 2rem;
		border-radius: 12px;
		box-shadow: var(--shadow-card);
		text-align: center;
		transition: var(--transition-card);
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

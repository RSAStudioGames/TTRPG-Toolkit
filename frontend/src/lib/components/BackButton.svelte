<script lang="ts">
	import { goto } from '$app/navigation';
	import { motion } from '@humanspeak/svelte-motion';

	interface Props {
		href: string;
		parentTitle: string;
	}

	let { href, parentTitle }: Props = $props();

	const buttonTransition = { type: 'spring' as const, stiffness: 500, damping: 35 };

	function handleClick() {
		goto(href);
	}
</script>

<motion.button
	type="button"
	class="back-button"
	whileHover={{ filter: 'brightness(0.97)' }}
	whileTap={{ scale: 0.98 }}
	transition={buttonTransition}
	onclick={handleClick}
>
	<i class="fa-solid fa-arrow-left" aria-hidden="true"></i>
	Back to {parentTitle}
</motion.button>

<style>
	:global(.back-button) {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		border: 1px solid var(--text-muted);
		border-radius: 8px;
		background: transparent;
		color: var(--text-ink);
		font-size: 0.9375rem;
		font-weight: 600;
		font-family: inherit;
		cursor: pointer;
		transition:
			filter 0.15s ease,
			border-color 0.15s ease;
	}

	:global(.back-button:hover) {
		border-color: var(--accent-primary);
	}

	:global(.back-button:active) {
		filter: brightness(0.95);
	}
</style>

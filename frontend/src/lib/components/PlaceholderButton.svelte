<script lang="ts">
	import { motion } from '@humanspeak/svelte-motion';

	const buttonTransition = { type: 'spring' as const, stiffness: 500, damping: 35 };
</script>

<motion.button
	type="button"
	class="placeholder-button"
	whileHover={{ filter: 'brightness(1.05)' }}
	whileTap={{ scale: 0.98 }}
	transition={buttonTransition}
>
	<i class="fas fa-arrow-right" aria-hidden="true"></i>
	Placeholder
</motion.button>

<style>
	:global(.placeholder-button) {
		display: flex;
		width: 100%;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.875rem 1.25rem;
		border: none;
		border-radius: 8px;
		background: var(--accent-primary);
		color: var(--text-ink);
		font-size: 1rem;
		font-weight: 600;
		font-family: inherit;
		cursor: pointer;
		transition: filter 0.15s ease;
	}

	:global(.placeholder-button:active) {
		filter: brightness(0.95);
	}
</style>

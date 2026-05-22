<script lang="ts">
	import { motion } from '@humanspeak/svelte-motion';

	const buttonTransition = { type: 'spring' as const, stiffness: 500, damping: 35 };
</script>

<motion.button
	type="button"
	class="open-content-library-button"
	whileHover={{ filter: 'brightness(0.97)' }}
	whileTap={{ scale: 0.98 }}
	transition={buttonTransition}
>
	<i class="fas fa-arrow-right" aria-hidden="true"></i>
	Open Content Library
</motion.button>

<style>
	:global(.open-content-library-button) {
		display: flex;
		width: 100%;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.875rem 1.25rem;
		border: none;
		border-radius: 8px;
		background: var(--accent-gm-muted);
		color: var(--text-ink);
		font-size: 1rem;
		font-weight: 600;
		font-family: inherit;
		cursor: pointer;
		transition: filter 0.15s ease;
	}

	:global(.open-content-library-button:active) {
		filter: brightness(0.92);
	}
</style>

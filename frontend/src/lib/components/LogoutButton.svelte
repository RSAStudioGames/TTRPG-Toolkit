<script lang="ts">
	import { goto } from '$app/navigation';
	import { motion } from '@humanspeak/svelte-motion';

	const buttonTransition = { type: 'spring' as const, stiffness: 500, damping: 35 };

	function handleLogout() {
		goto('/');
	}
</script>

<motion.button
	type="button"
	class="logout-button"
	whileHover={{ filter: 'brightness(0.97)' }}
	whileTap={{ scale: 0.98 }}
	transition={buttonTransition}
	onclick={handleLogout}
>
	Log out
</motion.button>

<style>
	:global(.logout-button) {
		display: inline-flex;
		align-items: center;
		justify-content: center;
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

	:global(.logout-button:hover) {
		border-color: var(--accent-primary);
	}

	:global(.logout-button:active) {
		filter: brightness(0.95);
	}
</style>

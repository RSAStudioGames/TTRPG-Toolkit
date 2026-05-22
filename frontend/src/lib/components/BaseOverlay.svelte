<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open: boolean;
		onclose: () => void;
		title: string;
		headerExtra?: Snippet;
		children: Snippet;
		footer?: Snippet;
	}

	let { open, onclose, title, headerExtra, children, footer }: Props = $props();

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose();
	}
</script>

<svelte:window onkeydown={open ? handleKeydown : undefined} />

{#if open}
	<div class="overlay-backdrop">
		<div class="overlay-panel" role="dialog" aria-modal="true" aria-labelledby="overlay-title">
			<header class="overlay-header">
				<div class="overlay-title-row">
					<h2 id="overlay-title">{title}</h2>
					{#if headerExtra}
						{@render headerExtra()}
					{/if}
				</div>
				<button type="button" class="overlay-close" aria-label="Close" onclick={onclose}>×</button>
			</header>
			<div class="overlay-body">
				{@render children()}
			</div>
			{#if footer}
				<footer class="overlay-footer">
					{@render footer()}
				</footer>
			{/if}
		</div>
	</div>
{/if}

<style>
	.overlay-backdrop {
		position: fixed;
		inset: 0;
		z-index: 900;
		background: rgba(0, 0, 0, 0.35);
	}

	.overlay-panel {
		position: fixed;
		top: 0;
		right: 0;
		width: min(520px, 100%);
		height: 100%;
		display: flex;
		flex-direction: column;
		background: var(--card-surface);
		border-left: 2px solid var(--accent-gm);
		box-shadow: -8px 0 32px rgba(0, 0, 0, 0.12);
	}

	.overlay-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		padding: 1rem 1.25rem;
		border-bottom: 1px solid var(--accent-gm-muted);
		gap: 0.5rem;
	}

	.overlay-title-row {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		flex: 1;
	}

	.overlay-title-row h2 {
		margin: 0;
		font-size: 1.2rem;
	}

	.overlay-close {
		border: none;
		background: none;
		font-size: 1.5rem;
		cursor: pointer;
		color: var(--text-muted);
	}

	.overlay-body {
		flex: 1;
		overflow-y: auto;
		padding: 1.25rem;
	}

	.overlay-footer {
		padding: 1rem 1.25rem;
		border-top: 1px solid var(--accent-gm-muted);
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}
</style>

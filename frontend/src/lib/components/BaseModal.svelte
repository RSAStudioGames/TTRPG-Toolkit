<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open: boolean;
		title: string;
		onclose: () => void;
		children: Snippet;
		footer?: Snippet;
		wide?: boolean;
		/** When false, clicking the dimmed backdrop does not close the modal. Default true. */
		closeOnBackdrop?: boolean;
		/** Raised z-index for nested modals (e.g. confirm dialogs). */
		stacked?: boolean;
		/** When false, Escape does not close the modal. Default true. */
		closeOnEscape?: boolean;
	}

	let {
		open,
		title,
		onclose,
		children,
		footer,
		wide = false,
		closeOnBackdrop = true,
		stacked = false,
		closeOnEscape = true
	}: Props = $props();

	function handleBackdrop(e: MouseEvent) {
		if (!closeOnBackdrop) return;
		if (e.target === e.currentTarget) onclose();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!closeOnEscape || e.key !== 'Escape') return;
		onclose();
	}
</script>

<svelte:window onkeydown={open ? handleKeydown : undefined} />

{#if open}
	<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
	<div
		class="modal-backdrop"
		class:modal-backdrop--stacked={stacked}
		role="presentation"
		onclick={handleBackdrop}
	>
		<div
			class="modal-panel"
			class:modal-panel--wide={wide}
			role="dialog"
			aria-modal="true"
			aria-labelledby="modal-title"
		>
			<header class="modal-header">
				<h2 id="modal-title">{title}</h2>
				<button type="button" class="modal-close" aria-label="Close" onclick={onclose}>×</button>
			</header>
			<div class="modal-body">
				{@render children()}
			</div>
			{#if footer}
				<footer class="modal-footer">
					{@render footer()}
				</footer>
			{/if}
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		z-index: 1000;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		background: rgba(0, 0, 0, 0.45);
	}

	.modal-backdrop--stacked {
		z-index: 1100;
	}

	.modal-panel {
		width: min(640px, 100%);
		max-height: 90vh;
		display: flex;
		flex-direction: column;
		background: var(--card-surface);
		border-radius: 12px;
		border: 2px solid var(--accent-gm);
		box-shadow: var(--shadow-card-hover-gm);
	}

	.modal-panel--wide {
		width: min(720px, 100%);
	}

	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 1.25rem;
		border-bottom: 1px solid var(--accent-gm-muted);
	}

	.modal-header h2 {
		margin: 0;
		font-size: 1.25rem;
	}

	.modal-close {
		border: none;
		background: none;
		font-size: 1.5rem;
		cursor: pointer;
		line-height: 1;
		color: var(--text-muted);
	}

	.modal-body {
		padding: 1.25rem;
		overflow-y: auto;
		flex: 1;
	}

	.modal-footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.5rem;
		padding: 1rem 1.25rem;
		border-top: 1px solid var(--accent-gm-muted);
	}
</style>

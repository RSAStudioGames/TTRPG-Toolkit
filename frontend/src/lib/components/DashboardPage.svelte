<script lang="ts">
	import type { Snippet } from 'svelte';

	type ContentLayout = 'cards' | 'stack';
	type DashboardRole = 'gm' | 'player';

	interface Props {
		title: string;
		role: DashboardRole;
		showRoleIcon?: boolean;
		contentLayout?: ContentLayout;
		headerActions?: Snippet;
		children: Snippet;
	}

	let {
		title,
		role,
		showRoleIcon = false,
		contentLayout = 'cards',
		headerActions,
		children
	}: Props = $props();

	const roleIconClass = $derived(
		role === 'gm' ? 'fa-solid fa-crown' : 'fa-solid fa-users'
	);
</script>

<div class="dashboard">
	<header class="dashboard-header" class:role-gm={role === 'gm'} class:role-player={role === 'player'}>
		<div class="header-row" class:role-dashboard={showRoleIcon}>
			<div class="header-title">
				{#if showRoleIcon}
					<i class={roleIconClass} aria-hidden="true"></i>
				{/if}
				<h1>{title}</h1>
			</div>
			{#if headerActions}
				<div class="header-actions">
					{@render headerActions()}
				</div>
			{/if}
		</div>
	</header>

	<main class="cards-row" class:stacked={contentLayout === 'stack'}>
		{@render children()}
	</main>
</div>

<style>
	.dashboard {
		display: flex;
		flex-direction: column;
		align-items: center;
		min-height: 100vh;
		padding: 2.75rem 1.5rem 3rem;
		width: 100%;
		max-width: 1100px;
		margin: 0 auto;
	}

	.dashboard-header {
		width: 100%;
		padding: 1.25rem 1.5rem;
		border-radius: 12px;
	}

	.dashboard-header.role-gm {
		background: linear-gradient(
			180deg,
			rgba(212, 175, 55, 0.14) 0%,
			rgba(212, 175, 55, 0.04) 100%
		);
	}

	.dashboard-header.role-player {
		background: linear-gradient(
			180deg,
			rgba(220, 38, 38, 0.1) 0%,
			rgba(220, 38, 38, 0.03) 100%
		);
	}

	.header-row {
		display: grid;
		grid-template-columns: 1fr auto 1fr;
		align-items: center;
		width: 100%;
		gap: 1rem;
	}

	.header-row.role-dashboard {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.header-title {
		grid-column: 2;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.625rem;
		min-width: 0;
	}

	.header-row.role-dashboard .header-title {
		justify-content: flex-start;
	}

	.header-title i {
		font-size: clamp(1.25rem, 3vw, 1.5rem);
		color: var(--accent-primary);
		flex-shrink: 0;
	}

	.header-title h1 {
		margin: 0;
		font-size: clamp(1.75rem, 4vw, 2.25rem);
		font-weight: 700;
		color: var(--text-ink);
		text-align: center;
	}

	.header-row.role-dashboard .header-title h1 {
		text-align: left;
	}

	.header-actions {
		grid-column: 3;
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 0.75rem;
		flex-shrink: 0;
		flex-wrap: nowrap;
	}

	.header-row.role-dashboard .header-actions {
		grid-column: unset;
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

	.cards-row.stacked {
		display: block;
		max-width: 720px;
	}

	@media (max-width: 768px) {
		.cards-row {
			grid-template-columns: 1fr;
		}
	}
</style>

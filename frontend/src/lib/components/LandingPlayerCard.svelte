<script lang="ts">
	import JoinGameButton from './JoinGameButton.svelte';
</script>

<article class="role-card" style="--accent-primary: var(--accent-player)">
	<div class="card-icon player-icon">
		<i class="fa-solid fa-users" aria-hidden="true"></i>
	</div>
	<h2 class="card-title">Player</h2>
	<p class="card-description">
		Enter the game world and embark on epic adventures with your fellow players.
	</p>
	<ul class="feature-list">
		<li>
			<i class="fa-solid fa-dice-d20" aria-hidden="true"></i>
			Roll dice and track results
		</li>
		<li>
			<i class="fa-solid fa-user" aria-hidden="true"></i>
			Manage your character
		</li>
		<li class="feature-indent">Track your inventory</li>
		<li>
			<i class="fa-solid fa-comments" aria-hidden="true"></i>
			Real-time collaboration
		</li>
	</ul>
	<JoinGameButton />
</article>

<style>
	.role-card {
		flex: 1;
		min-width: 280px;
		max-width: 520px;
		padding: 2rem;
		background: var(--card-surface);
		border-radius: 12px;
		box-shadow: var(--shadow-card);
		text-align: center;
	}

	.card-icon {
		font-size: 2.5rem;
		margin-bottom: 1rem;
	}

	.player-icon {
		color: var(--accent-primary);
	}

	.card-title {
		margin: 0 0 0.75rem;
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--text-ink);
	}

	.card-description {
		margin: 0 0 1.5rem;
		font-size: 0.9375rem;
		line-height: 1.6;
		color: var(--text-muted);
	}

	.feature-list {
		list-style: none;
		margin: 0 0 1.75rem;
		padding: 0;
		text-align: left;
	}

	.feature-list li {
		display: flex;
		align-items: flex-start;
		gap: 0.625rem;
		margin-bottom: 0.625rem;
		font-size: 0.875rem;
		color: var(--text-muted);
	}

	.feature-list li i {
		margin-top: 0.15rem;
		color: var(--accent-primary);
		width: 1rem;
		flex-shrink: 0;
	}

	.feature-indent {
		padding-left: 1.625rem;
	}
</style>

<script lang="ts">
	import GMLoginButton from './GMLoginButton.svelte';
</script>

<article class="role-card" style="--accent-primary: var(--accent-gm)">
	<div class="card-icon gm-icon">
		<i class="fa-solid fa-crown" aria-hidden="true"></i>
	</div>
	<h2 class="card-title">Game Master</h2>
	<p class="card-description">
		Take control as the Game Master and guide your players through incredible stories.
	</p>
	<ul class="feature-list">
		<li>
			<i class="fa-solid fa-users-gear" aria-hidden="true"></i>
			Manage all players
		</li>
		<li>
			<i class="fa-solid fa-dragon" aria-hidden="true"></i>
			Control NPCs &amp; enemies
		</li>
		<li>
			<i class="fa-solid fa-map" aria-hidden="true"></i>
			Track story progress
		</li>
		<li>
			<i class="fa-solid fa-gears" aria-hidden="true"></i>
			Full system control
		</li>
	</ul>
	<GMLoginButton />
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

	.gm-icon {
		color: var(--accent-gm);
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
</style>

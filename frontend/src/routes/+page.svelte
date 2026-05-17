<script lang="ts">
	import LandingRoleCard from '$lib/components/LandingRoleCard.svelte';
	import JoinGameButton from '$lib/components/JoinGameButton.svelte';
	import GMLoginButton from '$lib/components/GMLoginButton.svelte';
</script>

<div class="landing">
	<header class="landing-header">
		<div class="title-row">
			<i class="fa-solid fa-dice-d20 dice-spin" aria-hidden="true"></i>
			<h1>Tabletop RPG Toolkit</h1>
		</div>
		<p class="tagline">Your ultimate companion for tabletop role-playing adventures</p>
	</header>

	<main class="cards-row">
		<LandingRoleCard
			role="player"
			iconClass="fa-solid fa-users"
			title="Player"
			description="Enter the game world and embark on epic adventures with your fellow players."
		>
			<JoinGameButton />
		</LandingRoleCard>

		<LandingRoleCard
			role="gm"
			iconClass="fa-solid fa-crown"
			title="Game Master"
			description="Take control as the Game Master and guide your players through incredible stories."
		>
			<GMLoginButton />
		</LandingRoleCard>
	</main>
</div>

<style>
	.landing {
		display: flex;
		flex-direction: column;
		align-items: center;
		min-height: 100vh;
		padding: 3.5rem 1.5rem 3rem;
	}

	.landing-header {
		text-align: center;
		margin-bottom: 0;
		max-width: 720px;
	}

	.title-row {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.625rem;
		margin-bottom: 0.75rem;
	}

	.title-row i {
		font-size: 2.25rem;
		color: var(--accent-player);
	}

	.landing-header h1 {
		margin: 0;
		font-size: clamp(1.75rem, 4vw, 2.25rem);
		font-weight: 700;
		color: var(--text-ink);
	}

	.tagline {
		margin: 0;
		font-size: 1.0625rem;
		color: var(--text-muted);
		line-height: 1.5;
	}

	.cards-row {
		display: grid;
		grid-template-columns: repeat(2, minmax(280px, 1fr));
		gap: 2rem;
		align-items: stretch;
		width: 100%;
		max-width: 1100px;
		margin-top: 3.5rem;
	}

	@media (max-width: 768px) {
		.cards-row {
			grid-template-columns: 1fr;
			justify-items: center;
		}
	}
</style>

<script lang="ts">
	import LandingPlayerCard from '$lib/components/LandingPlayerCard.svelte';
	import LandingGMCard from '$lib/components/LandingGMCard.svelte';
	import LandingFooter from '$lib/components/LandingFooter.svelte';
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
		<LandingPlayerCard />
		<LandingGMCard />
	</main>

	<LandingFooter />
</div>

<style>
	.landing {
		display: flex;
		flex-direction: column;
		align-items: center;
		min-height: 100vh;
		padding: 2.5rem 1.5rem 0;
	}

	.landing-header {
		text-align: center;
		margin-bottom: 2.5rem;
		max-width: 720px;
	}

	.title-row {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
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
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 2rem;
		width: 100%;
		max-width: 1100px;
		margin-bottom: 2rem;
	}

	@media (max-width: 768px) {
		.cards-row {
			flex-direction: column;
			align-items: center;
		}
	}
</style>

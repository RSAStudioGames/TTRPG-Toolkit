<script lang="ts">
	import DashboardPage from '$lib/components/DashboardPage.svelte';
	import LandingRoleCard from '$lib/components/LandingRoleCard.svelte';
	import LogoutButton from '$lib/components/LogoutButton.svelte';
	import PlaceholderButton from '$lib/components/PlaceholderButton.svelte';
</script>

<DashboardPage title="Player" role="player" showRoleIcon>
	{#snippet headerActions()}
		<LogoutButton />
	{/snippet}

	<LandingRoleCard
		role="player"
		iconClass="fa-solid fa-layer-group"
		title="Placeholder"
		description="Placeholder"
	>
		<PlaceholderButton />
	</LandingRoleCard>
</DashboardPage>

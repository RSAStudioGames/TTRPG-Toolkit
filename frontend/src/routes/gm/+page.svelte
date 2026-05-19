<script lang="ts">
	import DashboardPage from '$lib/components/DashboardPage.svelte';
	import LandingRoleCard from '$lib/components/LandingRoleCard.svelte';
	import LogoutButton from '$lib/components/LogoutButton.svelte';
	import OpenSystemBuilderButton from '$lib/components/OpenSystemBuilderButton.svelte';
</script>

<DashboardPage title="Game Master Dashboard">
	{#snippet headerActions()}
		<LogoutButton />
	{/snippet}

	<LandingRoleCard
		role="gm"
		iconClass="fa-solid fa-screwdriver-wrench"
		title="System Builder"
		description="TTRPG Management Hub"
	>
		<OpenSystemBuilderButton />
	</LandingRoleCard>
</DashboardPage>

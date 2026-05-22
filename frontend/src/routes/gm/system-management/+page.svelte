<script lang="ts">
	import BackButton from '$lib/components/BackButton.svelte';
	import DashboardPage from '$lib/components/DashboardPage.svelte';
	import LogoutButton from '$lib/components/LogoutButton.svelte';
	import OpenCharacterSystemsButton from '$lib/components/OpenCharacterSystemsButton.svelte';
	import OpenContentLibraryButton from '$lib/components/OpenContentLibraryButton.svelte';
	import OpenGameSystemsButton from '$lib/components/OpenGameSystemsButton.svelte';
	import OpenTemplateManagerButton from '$lib/components/OpenTemplateManagerButton.svelte';
	import SystemManagementCard from '$lib/components/SystemManagementCard.svelte';
</script>

<DashboardPage title="System Management">
	{#snippet headerActions()}
		<BackButton href="/gm" parentTitle="Game Master Dashboard" />
		<LogoutButton />
	{/snippet}

	<SystemManagementCard
		iconClass="fa-solid fa-dice-d20"
		title="Game Systems"
		description="Create and manage game systems."
	>
		<OpenGameSystemsButton />
	</SystemManagementCard>
	<SystemManagementCard
		iconClass="fa-solid fa-table-columns"
		title="Template Manager"
		description="PLACEHOLDER"
	>
		<OpenTemplateManagerButton />
	</SystemManagementCard>
	<SystemManagementCard
		iconClass="fa-solid fa-user-gear"
		title="Character Systems"
		description="PLACEHOLDER"
	>
		<OpenCharacterSystemsButton />
	</SystemManagementCard>
	<SystemManagementCard
		iconClass="fa-solid fa-book-open"
		title="Content Library"
		description="PLACEHOLDER"
	>
		<OpenContentLibraryButton />
	</SystemManagementCard>
</DashboardPage>

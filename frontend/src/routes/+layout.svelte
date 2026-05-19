<script lang="ts">
	import '@fortawesome/fontawesome-free/css/all.min.css';
	import '$lib/styles/theme.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>Tabletop RPG Toolkit</title>
	<meta
		name="description"
		content="Your companion for tabletop role-playing games"
	/>
</svelte:head>

{@render children()}

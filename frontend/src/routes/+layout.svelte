<script lang="ts">
	import '@fortawesome/fontawesome-free/css/all.min.css';
	import '$lib/styles/theme.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>Tabletop RPG Toolkit</title>
	<meta
		name="description"
		content="Your ultimate companion for tabletop role-playing adventures"
	/>
</svelte:head>

{@render children()}

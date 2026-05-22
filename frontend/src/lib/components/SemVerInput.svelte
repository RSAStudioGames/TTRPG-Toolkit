<script lang="ts">
	interface Props {
		value: string;
		disabled?: boolean;
	}

	let { value = $bindable(), disabled = false }: Props = $props();

	let major = $state('0');
	let minor = $state('1');
	let patch = $state('0');

	$effect(() => {
		const parts = value.split('.');
		if (parts.length === 3) {
			major = parts[0];
			minor = parts[1];
			patch = parts[2];
		}
	});

	function sync() {
		value = `${major || '0'}.${minor || '0'}.${patch || '0'}`;
	}
</script>

<div class="form-field">
	<label for="semver-major">Version</label>
	<div class="semver">
		<input id="semver-major" type="number" min="0" bind:value={major} {disabled} oninput={sync} />
		<span aria-hidden="true">.</span>
		<input type="number" min="0" aria-label="Minor version" bind:value={minor} {disabled} oninput={sync} />
		<span aria-hidden="true">.</span>
		<input type="number" min="0" aria-label="Patch version" bind:value={patch} {disabled} oninput={sync} />
	</div>
</div>

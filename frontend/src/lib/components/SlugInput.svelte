<script lang="ts">
	import { slugFromName, slugValidationMessage } from '$lib/utils/slug';

	interface Props {
		name: string;
		value: string;
		disabled?: boolean;
		auto?: boolean;
	}

	let { name, value = $bindable(), disabled = false, auto = true }: Props = $props();

	let slugTouched = $state(false);
	let slugError = $state<string | null>(null);

	$effect(() => {
		if (auto && !slugTouched && name) {
			value = slugFromName(name);
		}
	});

	function onSlugBlur() {
		slugTouched = true;
		slugError = slugValidationMessage(value);
	}

	export function validate(): boolean {
		slugError = slugValidationMessage(value);
		return !slugError;
	}
</script>

<div class="form-field">
	<label for="system-slug">URL Identifier</label>
	<p class="hint">Used in URLs and references. Must be lowercase, alphanumeric with hyphens.</p>
	<input
		id="system-slug"
		type="text"
		bind:value
		{disabled}
		onblur={onSlugBlur}
		oninput={() => {
			slugTouched = true;
			auto = false;
		}}
	/>
	{#if slugError}
		<p class="form-error">{slugError}</p>
	{/if}
</div>

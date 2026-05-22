<script lang="ts">
	interface Props {
		label: string;
		placeholder?: string;
		tokens: string[];
		disabled?: boolean;
	}

	let { label, placeholder = 'Add item…', tokens = $bindable(), disabled = false }: Props = $props();

	let draft = $state('');
	let inputId = `token-${Math.random().toString(36).slice(2, 9)}`;

	function addToken() {
		const t = draft.trim();
		if (!t || tokens.includes(t)) {
			draft = '';
			return;
		}
		tokens = [...tokens, t];
		draft = '';
	}

	function removeToken(i: number) {
		tokens = tokens.filter((_, idx) => idx !== i);
	}

	function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			e.preventDefault();
			addToken();
		}
	}
</script>

<div class="form-field token-input">
	<label for={inputId}>{label}</label>
	<div class="chips">
		{#each tokens as token, i}
			<span class="chip">
				{token}
				<button type="button" aria-label="Remove" {disabled} onclick={() => removeToken(i)}>×</button>
			</span>
		{/each}
	</div>
	<div class="add-row">
		<input id={inputId} type="text" {placeholder} bind:value={draft} {disabled} onkeydown={onKeydown} />
		<button type="button" class="btn-secondary-sm" {disabled} onclick={addToken}>Add</button>
	</div>
</div>

<style>
	.chips {
		display: flex;
		flex-wrap: wrap;
		gap: 0.35rem;
		margin-bottom: 0.35rem;
	}

	.chip {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.2rem 0.5rem;
		background: var(--accent-gm-muted);
		border-radius: 999px;
		font-size: 0.85rem;
	}

	.chip button {
		border: none;
		background: none;
		cursor: pointer;
		padding: 0;
		line-height: 1;
	}
</style>

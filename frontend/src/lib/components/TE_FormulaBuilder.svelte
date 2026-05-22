<script lang="ts">
	import { validateFormula, ApiError } from '$lib/api';
	import {
		DICE_TOKENS,
		FUNCTION_TOKENS,
		LOGIC_TOKENS,
		MODIFIER_TOKENS,
		TOKEN_GROUP_LABELS,
		TOKEN_GROUP_ORDER,
		type FormulaToken,
		type TokenGroup
	} from './formulaTokens';

	interface Props {
		value: string;
		label: string;
		variables: string[];
		systemId: string;
		disabled?: boolean;
	}

	let {
		value = $bindable(''),
		label,
		variables = [],
		systemId,
		disabled = false
	}: Props = $props();

	type ValidationState = 'idle' | 'loading' | 'valid' | 'invalid';

	let textareaEl: HTMLTextAreaElement | undefined = $state();
	let tokensOpen = $state(false);
	let validationState: ValidationState = $state('idle');
	let validationErrors = $state<string[]>([]);
	let validationMessage = $state('');

	const textareaId = `formula-${Math.random().toString(36).slice(2, 9)}`;

	const staticTokensByGroup: Record<Exclude<TokenGroup, 'variables'>, FormulaToken[]> = {
		dice: DICE_TOKENS,
		modifiers: MODIFIER_TOKENS,
		logic: LOGIC_TOKENS,
		functions: FUNCTION_TOKENS
	};

	function insertAtCursor(snippet: string, cursorOffsetAfter = snippet.length) {
		const el = textareaEl;
		if (!el) return;

		const start = el.selectionStart ?? value.length;
		const end = el.selectionEnd ?? value.length;
		const before = value.slice(0, start);
		const after = value.slice(end);
		value = before + snippet + after;

		const pos = start + cursorOffsetAfter;
		requestAnimationFrame(() => {
			el.focus();
			el.setSelectionRange(pos, pos);
		});

		validationState = 'idle';
		validationErrors = [];
		validationMessage = '';
	}

	function insertVariable(name: string) {
		const trimmed = name.trim();
		if (!trimmed) return;
		const snippet = trimmed.startsWith('{') && trimmed.endsWith('}') ? trimmed : `{${trimmed}}`;
		insertAtCursor(snippet);
	}

	function onTokenClick(token: FormulaToken) {
		insertAtCursor(token.insert);
	}

	function onValueInput() {
		if (validationState !== 'idle') {
			validationState = 'idle';
			validationErrors = [];
			validationMessage = '';
		}
	}

	async function handleValidate() {
		if (!systemId.trim()) {
			validationState = 'invalid';
			validationErrors = [];
			validationMessage = 'System ID is required to validate syntax.';
			return;
		}

		validationState = 'loading';
		validationErrors = [];
		validationMessage = '';

		try {
			await validateFormula(systemId, value);
			validationState = 'valid';
		} catch (err) {
			validationState = 'invalid';
			if (err instanceof ApiError) {
				validationErrors = err.errors.length > 0 ? err.errors : [err.message];
				validationMessage = err.errors.length > 0 ? 'Invalid formula' : err.message;
			} else if (err instanceof Error) {
				validationErrors = [err.message];
				validationMessage = 'Validation failed';
			} else {
				validationErrors = ['Validation failed'];
				validationMessage = 'Validation failed';
			}
		}
	}
</script>

<div class="formula-builder">
	<div class="formula-builder__grid">
		<div class="formula-builder__editor">
			<div class="form-field">
				<label for={textareaId}>{label}</label>
				<textarea
					id={textareaId}
					bind:this={textareaEl}
					bind:value
					{disabled}
					class="formula-textarea"
					rows="5"
					spellcheck="false"
					oninput={onValueInput}
				></textarea>
			</div>

			<div class="formula-builder__actions">
				<button
					type="button"
					class="formula-validate-btn"
					disabled={disabled || validationState === 'loading'}
					onclick={handleValidate}
				>
					{validationState === 'loading' ? 'Validating...' : 'Validate Syntax'}
				</button>

				<div
					class="validation-status"
					class:validation-status--valid={validationState === 'valid'}
					class:validation-status--invalid={validationState === 'invalid'}
					aria-live="polite"
				>
					{#if validationState === 'valid'}
						<span class="validation-status__icon" aria-hidden="true">✓</span>
						<span>Valid Syntax</span>
					{:else if validationState === 'invalid'}
						<span class="validation-status__icon validation-status__icon--error" aria-hidden="true"
							>✕</span
						>
						<span>{validationMessage || 'Invalid formula'}</span>
					{/if}
				</div>

				{#if validationState === 'invalid' && validationErrors.length > 0}
					<ul class="validation-errors">
						{#each validationErrors as err}
							<li>{err}</li>
						{/each}
					</ul>
				{/if}
			</div>
		</div>

		<!-- Accepted deviation: spec requests dropdown menus; collapsible details panel is intentional. -->
		<aside class="formula-builder__tokens">
			<details class="token-panel" bind:open={tokensOpen}>
				<summary class="token-panel__toggle">Insert Token</summary>
				<div class="token-panel__body">
					{#each TOKEN_GROUP_ORDER as group}
						<section class="token-group">
							<h4 class="token-group__title">{TOKEN_GROUP_LABELS[group]}</h4>
							<div class="token-group__items">
								{#if group === 'variables'}
									{#if variables.length === 0}
										<p class="token-group__empty">No variables provided.</p>
									{:else}
										{#each variables as variable}
											<button
												type="button"
												class="token-btn"
												{disabled}
												title="Insert {variable}"
												onclick={() => insertVariable(variable)}
											>
												{`{${variable.replace(/^\{|\}$/g, '')}}`}
											</button>
										{/each}
									{/if}
								{:else}
									{#each staticTokensByGroup[group] as token (token.label + token.insert)}
										<button
											type="button"
											class="token-btn"
											{disabled}
											title={token.title ?? token.label}
											onclick={() => onTokenClick(token)}
										>
											{token.label}
										</button>
									{/each}
								{/if}
							</div>
							{#if group === 'modifiers'}
								<p class="token-group__hint">
									Insert a modifier after dice (e.g. 4d6), then type the required number. Some
									modifiers may not validate until rules support is extended.
								</p>
							{/if}
						</section>
					{/each}
				</div>
			</details>
		</aside>
	</div>

	<p class="formula-builder__helper hint">
		Use standard notation (e.g., 2d6, 4d6kh3). Wrap variables in {'{curly_braces}'}.
	</p>
</div>

<style>
	.formula-builder {
		background: var(--bg-paper);
		color: var(--text-ink);
		border-radius: 8px;
		padding: 0.75rem 0;
	}

	.formula-builder__grid {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	@media (min-width: 640px) {
		.formula-builder__grid {
			display: grid;
			grid-template-columns: 2fr 1fr;
			align-items: start;
			gap: 1.25rem;
		}
	}

	.formula-textarea {
		font-family: ui-monospace, 'Cascadia Code', 'Source Code Pro', Menlo, Consolas, monospace;
		font-size: 0.9rem;
	}

	.formula-builder__actions {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}

	.formula-validate-btn {
		align-self: flex-start;
		display: inline-flex;
		align-items: center;
		padding: 0.5rem 1rem;
		border: 1px solid var(--accent-gm);
		border-radius: 6px;
		background: var(--accent-gm-muted);
		color: var(--text-ink);
		font: inherit;
		font-weight: 600;
		cursor: pointer;
		transition: filter 0.15s ease;
	}

	.formula-validate-btn:hover:not(:disabled) {
		filter: brightness(0.97);
	}

	.formula-validate-btn:disabled {
		opacity: 0.65;
		cursor: not-allowed;
	}

	.validation-status {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.875rem;
		min-height: 1.25rem;
	}

	.validation-status--valid {
		color: #15803d;
		font-weight: 600;
	}

	.validation-status--invalid {
		color: #b91c1c;
		font-weight: 600;
	}

	.validation-status__icon {
		font-size: 1rem;
		line-height: 1;
	}

	.validation-errors {
		margin: 0;
		padding-left: 1.25rem;
		color: #b91c1c;
		font-size: 0.8rem;
	}

	.validation-errors li {
		margin-bottom: 0.2rem;
	}

	.token-panel {
		border: 1px solid var(--input-border, #d1d5db);
		border-radius: 8px;
		background: var(--card-surface, #fff);
	}

	.token-panel__toggle {
		padding: 0.6rem 0.85rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		list-style: none;
		user-select: none;
	}

	.token-panel__toggle::-webkit-details-marker {
		display: none;
	}

	.token-panel__body {
		padding: 0 0.85rem 0.85rem;
		max-height: 22rem;
		overflow-y: auto;
	}

	.token-group {
		margin-bottom: 0.85rem;
	}

	.token-group:last-child {
		margin-bottom: 0;
	}

	.token-group__title {
		margin: 0 0 0.4rem;
		font-size: 0.75rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--text-muted);
		border-bottom: 1px solid var(--accent-gm-muted);
		padding-bottom: 0.25rem;
	}

	.token-group__items {
		display: flex;
		flex-wrap: wrap;
		gap: 0.35rem;
	}

	.token-group__empty,
	.token-group__hint {
		margin: 0.25rem 0 0;
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	.token-btn {
		padding: 0.25rem 0.5rem;
		border: 1px solid var(--input-border, #d1d5db);
		border-radius: 4px;
		background: var(--input-bg, #fff);
		color: var(--text-ink);
		font: inherit;
		font-size: 0.8rem;
		cursor: pointer;
		transition:
			border-color 0.15s ease,
			background 0.15s ease;
	}

	.token-btn:hover:not(:disabled) {
		border-color: var(--accent-gm);
		background: var(--accent-gm-muted);
	}

	.token-btn:disabled {
		opacity: 0.55;
		cursor: not-allowed;
	}

	.formula-builder__helper {
		margin: 0.75rem 0 0;
	}
</style>

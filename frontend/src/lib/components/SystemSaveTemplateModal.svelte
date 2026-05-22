<script lang="ts">
	import BaseModal from './BaseModal.svelte';
	import { saveSystemTemplate } from '$lib/api/systems';
	import type { GameSystem } from '$lib/types/system';
	import { ApiError } from '$lib/api/client';

	interface Props {
		open: boolean;
		system: GameSystem | null;
		onclose: () => void;
	}

	let { open, system, onclose }: Props = $props();

	let templateName = $state('');
	let templateDescription = $state('');
	let saving = $state(false);
	let error = $state<string | null>(null);

	$effect(() => {
		if (open && system) {
			templateName = system.name;
			templateDescription = '';
			error = null;
		}
	});

	async function handleSave() {
		if (!system || !templateName.trim()) return;
		saving = true;
		error = null;
		try {
			await saveSystemTemplate(system.id, templateName.trim(), templateDescription.trim());
			onclose();
		} catch (e) {
			error = e instanceof ApiError ? e.message : 'Failed to save template';
		} finally {
			saving = false;
		}
	}
</script>

<BaseModal {open} title="Save as Template" onclose={onclose}>
	<div class="form-field">
		<label for="tpl-name">Template Name</label>
		<input id="tpl-name" type="text" bind:value={templateName} maxlength="120" required />
	</div>
	<div class="form-field">
		<label for="tpl-desc">Template Description</label>
		<textarea id="tpl-desc" bind:value={templateDescription} rows="4"></textarea>
	</div>
	{#if error}
		<p class="form-error">{error}</p>
	{/if}

	{#snippet footer()}
		<button type="button" class="btn-secondary" onclick={onclose}>Cancel</button>
		<button type="button" class="btn-primary" disabled={saving || !templateName.trim()} onclick={handleSave}>
			{saving ? 'Saving…' : 'Save'}
		</button>
	{/snippet}
</BaseModal>

<style>
	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font: inherit;
		font-weight: 600;
		cursor: pointer;
	}

	.btn-primary {
		background: var(--accent-gm-muted);
		border: none;
	}

	.btn-secondary {
		background: #fff;
		border: 1px solid #d1d5db;
	}
</style>

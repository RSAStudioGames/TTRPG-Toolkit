<script lang="ts">
	import Combobox from './Combobox.svelte';
	import DescriptionEditor from './DescriptionEditor.svelte';
	import OfficialLinksInput from './OfficialLinksInput.svelte';
	import RadioGroup from './RadioGroup.svelte';
	import SegmentedControl from './SegmentedControl.svelte';
	import SemVerInput from './SemVerInput.svelte';
	import SlugInput from './SlugInput.svelte';
	import ToggleSwitch from './ToggleSwitch.svelte';
	import TokenInput from './TokenInput.svelte';
	import type { SystemFormValues } from '$lib/types/system';

	interface Props {
		form: SystemFormValues;
		disabled?: boolean;
		showStatusExtras?: boolean;
		showTagsRulebooks?: boolean;
		parentOptions?: { id: string; name: string }[];
		slugRef?: SlugInput;
	}

	let {
		form = $bindable(),
		disabled = false,
		showStatusExtras = false,
		showTagsRulebooks = false,
		parentOptions = []
	}: Props = $props();

	const licenseOptions = ['OGL', 'SRD', 'Proprietary', 'Homebrew', 'Custom'];
	const playstyleOptions = ['Tactical', 'Narrative', 'Simulationist', 'Gamist', 'Custom'];
	const familySuggestions = ['d20 System', 'PbtA', 'OSR', 'Year Zero', 'Powered by the Apocalypse'];
</script>

<div class="metadata-form">
	<section class="form-section">
		<h3>Identity & Attribution</h3>
		<div class="form-field">
			<label for="sys-name">System Name</label>
			<input
				id="sys-name"
				type="text"
				placeholder="e.g., Dungeons & Dragons"
				bind:value={form.name}
				minlength="3"
				maxlength="120"
				required
				{disabled}
			/>
		</div>
		<SlugInput name={form.name} bind:value={form.slug} {disabled} />
		<div class="form-field">
			<label for="sys-edition">Edition</label>
			<input id="sys-edition" type="text" placeholder="e.g., 5th Edition" bind:value={form.edition} {disabled} />
		</div>
		<div class="form-field">
			<label for="sys-publisher">Publisher</label>
			<input id="sys-publisher" type="text" bind:value={form.publisher} {disabled} />
		</div>
		{#if showTagsRulebooks}
			<TokenInput label="Core Rulebooks" placeholder="Add rulebook…" bind:tokens={form.core_rulebooks} {disabled} />
		{/if}
	</section>

	<section class="form-section">
		<h3>Description & Visuals</h3>
		<DescriptionEditor bind:value={form.description} {disabled} />
	</section>

	<section class="form-section">
		<h3>Classification & Mechanics</h3>
		<div class="form-field">
			<label for="sys-license">License Type</label>
			<select id="sys-license" bind:value={form.license_type} {disabled}>
				{#each licenseOptions as opt}
					<option value={opt}>{opt}</option>
				{/each}
			</select>
		</div>
		<SemVerInput bind:value={form.version} {disabled} />
		{#if showTagsRulebooks}
			<TokenInput label="Tags" placeholder="e.g., fantasy, sci-fi, crunchy" bind:tokens={form.tags} {disabled} />
		{/if}
		<Combobox
			label="System Family"
			placeholder="e.g., d20 System, PbtA"
			options={familySuggestions}
			bind:value={form.system_family}
			{disabled}
		/>
		<div class="form-field">
			<label>Recommended Players</label>
			<div class="player-range">
				<input type="number" min="1" placeholder="Min" bind:value={form.player_count_min} {disabled} />
				<span>to</span>
				<input type="number" min="1" placeholder="Max" bind:value={form.player_count_max} {disabled} />
			</div>
		</div>
		<RadioGroup label="Playstyle" name="playstyle" options={playstyleOptions} bind:value={form.playstyle} {disabled} />
		<SegmentedControl label="Complexity" bind:value={form.complexity} {disabled} />
	</section>

	<section class="form-section">
		<h3>Settings</h3>
		<RadioGroup
			label="Default Measurement Unit"
			name="measurement"
			options={['Imperial', 'Metric']}
			bind:value={form.measurement_unit}
			{disabled}
		/>
		<div class="form-field">
			<label for="sys-currency">Currency Symbol</label>
			<input
				id="sys-currency"
				type="text"
				maxlength="8"
				placeholder="e.g., GP, $, ¢"
				bind:value={form.currency_symbol}
				{disabled}
			/>
		</div>
		<OfficialLinksInput bind:links={form.official_links} {disabled} />
	</section>

	<section class="form-section">
		<h3>System Status Controls</h3>
		<ToggleSwitch label="Active Status" bind:checked={form.is_active} {disabled} />
		{#if showStatusExtras}
			<ToggleSwitch label="Core System" bind:checked={form.is_core} {disabled} />
			{#if !form.is_core}
				<div class="form-field">
					<label for="parent-system">Parent System</label>
					<select id="parent-system" bind:value={form.parent_system_id} {disabled}>
						<option value="">Select parent…</option>
						{#each parentOptions as p}
							<option value={p.id}>{p.name}</option>
						{/each}
					</select>
				</div>
			{/if}
			<ToggleSwitch label="Deletion Protection" bind:checked={form.is_protected} {disabled} />
		{/if}
	</section>
</div>

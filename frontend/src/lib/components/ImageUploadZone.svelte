<script lang="ts">
	interface Props {
		variant: 'icon' | 'cover';
		previewUrl?: string;
		disabled?: boolean;
		onfile?: (file: File) => void;
	}

	let { variant, previewUrl, disabled = false, onfile }: Props = $props();

	let error = $state<string | null>(null);
	let warn = $state<string | null>(null);
	let dragOver = $state(false);

	const ICON_MAX = 2 * 1024 * 1024;
	const COVER_MAX = 5 * 1024 * 1024;

	function validateFile(file: File): boolean {
		error = null;
		warn = null;
		const ext = file.name.split('.').pop()?.toLowerCase() ?? '';
		if (variant === 'icon') {
			if (!['png', 'jpg', 'jpeg', 'svg'].includes(ext)) {
				error = 'System Icon must be a PNG, JPG, or SVG file.';
				return false;
			}
			if (file.size > ICON_MAX) {
				error = 'System Icon must be 2 MB or smaller.';
				return false;
			}
		} else {
			if (!['png', 'jpg', 'jpeg'].includes(ext)) {
				error = 'Cover Image must be a PNG or JPG file.';
				return false;
			}
			if (file.size > COVER_MAX) {
				error = 'Cover Image must be 5 MB or smaller.';
				return false;
			}
		}
		checkAspect(file);
		return true;
	}

	function checkAspect(file: File) {
		if (file.type === 'image/svg+xml') return;
		const img = new Image();
		const url = URL.createObjectURL(file);
		img.onload = () => {
			URL.revokeObjectURL(url);
			if (variant === 'icon') {
				const d = Math.abs(img.width - img.height) / Math.max(img.width, img.height);
				if (d > 0.15) warn = '1:1 aspect ratio recommended.';
			} else {
				const ratio = img.width / img.height;
				if (ratio < 1.7 || ratio > 1.85) warn = '16:9 aspect ratio recommended.';
			}
		};
		img.src = url;
	}

	function handleFiles(files: FileList | null) {
		if (!files?.length || disabled) return;
		const file = files[0];
		if (validateFile(file)) onfile?.(file);
	}

	function onDrop(e: DragEvent) {
		e.preventDefault();
		dragOver = false;
		handleFiles(e.dataTransfer?.files ?? null);
	}
</script>

<div class="upload-zone">
	<span class="upload-label">{variant === 'icon' ? 'System Icon' : 'Cover Image'}</span>
	<p class="hint">
		{#if variant === 'icon'}
			PNG, JPG, or SVG · max 2 MB · 1:1 aspect ratio recommended.
		{:else}
			PNG or JPG · max 5 MB · 16:9 aspect ratio recommended.
		{/if}
	</p>
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="drop-area"
		class:drag-over={dragOver}
		role="group"
		ondragover={(e) => {
			e.preventDefault();
			dragOver = true;
		}}
		ondragleave={() => (dragOver = false)}
		ondrop={onDrop}
	>
		{#if previewUrl}
			<img src={previewUrl} alt="" class="preview" />
		{:else}
			<i class="fas fa-cloud-arrow-up" aria-hidden="true"></i>
			<span>Drag and drop or choose a file</span>
		{/if}
		<input
			type="file"
			accept={variant === 'icon' ? '.png,.jpg,.jpeg,.svg' : '.png,.jpg,.jpeg'}
			{disabled}
			onchange={(e) => handleFiles(e.currentTarget.files)}
		/>
	</div>
	{#if error}
		<p class="form-error">{error}</p>
	{/if}
	{#if warn}
		<p class="warn">{warn}</p>
	{/if}
</div>

<style>
	.upload-zone {
		margin-bottom: 0.85rem;
	}

	.upload-label {
		font-size: 0.875rem;
		font-weight: 600;
		display: block;
		margin-bottom: 0.25rem;
	}

	.hint {
		font-size: 0.8rem;
		color: var(--text-muted);
		margin: 0 0 0.5rem;
	}

	.drop-area {
		position: relative;
		border: 2px dashed #d1d5db;
		border-radius: 8px;
		padding: 1.5rem;
		text-align: center;
		background: #fafafa;
		min-height: 120px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
	}

	.drop-area.drag-over {
		border-color: var(--accent-gm);
		background: var(--accent-gm-muted);
	}

	.drop-area input {
		position: absolute;
		inset: 0;
		opacity: 0;
		cursor: pointer;
	}

	.preview {
		max-width: 100%;
		max-height: 140px;
		object-fit: contain;
	}

	.warn {
		color: #b45309;
		font-size: 0.8rem;
		margin: 0.25rem 0 0;
	}
</style>

<script lang="ts">
	import type { SystemStatus } from '$lib/types/system';

	interface Props {
		status: SystemStatus;
	}

	let { status }: Props = $props();

	const labels: Record<SystemStatus, string> = {
		draft: 'Draft',
		published: 'Published',
		locked: 'Locked',
		archived: 'Archived'
	};
</script>

<span class="badge status-{status}">{labels[status]}</span>

<style>
	.badge {
		display: inline-block;
		padding: 0.2rem 0.55rem;
		border-radius: 999px;
		font-size: 0.75rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.03em;
	}

	.status-draft {
		background: #e5e7eb;
		color: #374151;
	}

	.status-published {
		background: #d1fae5;
		color: #065f46;
	}

	.status-locked {
		background: #fef3c7;
		color: #92400e;
	}

	.status-archived {
		background: #fee2e2;
		color: #991b1b;
	}
</style>

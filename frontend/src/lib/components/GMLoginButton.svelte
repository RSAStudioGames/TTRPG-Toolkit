<button type="button" class="gm-login-button">
	<i class="fa-solid fa-key" aria-hidden="true"></i>
	GM Login
</button>

<style>
	.gm-login-button {
		display: flex;
		width: 100%;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.875rem 1.25rem;
		border: none;
		border-radius: 8px;
		background: var(--accent-gm-muted);
		color: var(--text-ink);
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: filter 0.15s ease;
	}

	.gm-login-button:hover {
		filter: brightness(0.97);
	}
</style>

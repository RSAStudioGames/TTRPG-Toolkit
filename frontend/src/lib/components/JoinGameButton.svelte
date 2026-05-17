<button type="button" class="join-game-button">
	<i class="fa-solid fa-right-to-bracket" aria-hidden="true"></i>
	Join Game
</button>

<style>
	.join-game-button {
		display: flex;
		width: 100%;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.875rem 1.25rem;
		border: none;
		border-radius: 8px;
		background: var(--accent-primary);
		color: #ffffff;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: filter 0.15s ease;
	}

	.join-game-button:hover {
		filter: brightness(1.05);
	}
</style>

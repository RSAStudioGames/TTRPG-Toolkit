<script lang="ts">
	const year = new Date().getFullYear();
</script>

<footer class="landing-footer">
	<p class="copyright">© {year} Tabletop RPG Toolkit. Built by RSA.</p>
	<nav class="footer-nav" aria-label="Footer">
		<button type="button" class="footer-link">About</button>
		<button type="button" class="footer-link">Help</button>
	</nav>
</footer>

<style>
	.landing-footer {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		width: 100%;
		max-width: 1100px;
		margin: 0 auto;
		padding: 1.5rem 1rem 2rem;
		border-top: 1px solid #e2e8f0;
	}

	.copyright {
		margin: 0;
		font-size: 0.875rem;
		color: var(--text-muted);
	}

	.footer-nav {
		display: flex;
		gap: 1.5rem;
	}

	.footer-link {
		padding: 0;
		border: none;
		background: none;
		font-size: 0.875rem;
		color: var(--text-muted);
		cursor: pointer;
	}

	.footer-link:hover {
		color: var(--text-ink);
	}
</style>

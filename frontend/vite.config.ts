import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	resolve: {
		dedupe: ['svelte']
	},
	// Optional: only used if you manually run `npm run dev` (not part of prep / production).
	server: {
		host: true,
		port: 3030,
		strictPort: true,
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:8080',
				changeOrigin: true
			}
		}
	},
	preview: {
		host: true
	}
});

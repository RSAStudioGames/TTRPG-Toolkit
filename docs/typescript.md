# TypeScript

## Tooling

| Tool | Version | Role |
|------|---------|------|
| TypeScript | 5.6+ | `strict: true`, `moduleResolution: "bundler"` |
| svelte-check | 4.x | Typecheck Svelte + TS (`npm run check`) |
| SvelteKit | 2.x | File-based routing, adapter-static |

Config: [`frontend/tsconfig.json`](../frontend/tsconfig.json)

## Layout

```
frontend/src/lib/
├── types/          # Interfaces mirroring API JSON (snake_case fields)
│   └── api.ts      # ApiEnvelope, AppConfig, …
├── api/
│   ├── client.ts   # api<T>(), apiEnvelope<T>(), ApiError
│   ├── config.ts   # fetchAppConfig()
│   └── index.ts
├── components/
├── stores/
└── utils/
```

## API client

All HTTP calls go through `$lib/api`, not raw `fetch` in routes:

```typescript
import { apiEnvelope } from '$lib/api/client';
import type { AppConfig } from '$lib/types/api';

const config = await apiEnvelope<AppConfig>('/api/config');
```

- `credentials: 'include'` for future session cookies
- `apiEnvelope` unwraps `{ status, data, message, errors }` from Go Fiber handlers

## Checks in prep

`go run -C tools ./build` runs `npm run check` before `npm run build`.

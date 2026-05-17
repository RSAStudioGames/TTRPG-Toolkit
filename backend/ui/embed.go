package ui

import "embed"

// Rebuild embedded UI from frontend/ (npm run build + copy). Required when Svelte changes.
//go:generate go run ../../tools/buildfrontend

//go:embed all:static
var Static embed.FS

# Prepares frontend embed + ttrpg-toolkit binary at repo root. No npm dev server.
$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot
go run -C tools ./build
exit $LASTEXITCODE

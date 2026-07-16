# figma-cli

[简体中文](README.zh-CN.md)

`figma-cli` is an agent-friendly Figma REST API command-line toolkit.

It is designed for repeatable inspection and asset export:

- Paste a Figma file/design/proto URL directly.
- Resolve file keys and node IDs automatically.
- Read file and node summaries without opening the Figma app.
- Export node image URLs or download assets locally.
- List comments for review context.
- Store the Figma token in the OS keyring; `FIGMA_TOKEN` is supported for CI and one-off runs.

This CLI complements Figma MCP workflows. Use MCP or browser automation for interactive canvas context and write-back workflows. Use this CLI when an agent needs stable, token-based reads and reproducible outputs.

## Install

Local development install:

```bash
go install ./cmd/figma-cli
```

Or build locally:

```bash
make build
./bin/figma-cli help
```

## Configure

Create a Figma personal access token, then store it in the OS keyring:

```bash
printf %s "$FIGMA_TOKEN" | figma-cli init --token-stdin
```

For one-off use or CI, set:

```bash
export FIGMA_TOKEN="figd_..."
```

`FIGMA_TOKEN` takes precedence over the keyring. The token is never printed by the CLI.

For local integration tests against a mock server or proxy:

```bash
export FIGMA_API_BASE_URL="http://127.0.0.1:8080/v1"
```

## Commands

Show the current Figma user:

```bash
figma-cli me
```

Inspect a file:

```bash
figma-cli file get "https://www.figma.com/design/FILE_KEY/name?node-id=1-2"
figma-cli file get FILE_KEY --depth 2
```

Inspect a node:

```bash
figma-cli node inspect "https://www.figma.com/design/FILE_KEY/name?node-id=1-2"
figma-cli node inspect FILE_KEY --node 1:2
```

Export a node image URL:

```bash
figma-cli image export FILE_KEY --node 1:2 --format png --scale 2
```

Download the exported image:

```bash
figma-cli image export FILE_KEY --node 1:2 --format png --scale 2 --out ./figma-assets
```

List comments:

```bash
figma-cli comments list FILE_KEY
```

Use raw JSON when an agent needs the original API response:

```bash
figma-cli node inspect FILE_KEY --node 1:2 --raw
```

## Output Contract

- Default output is concise text for humans and agents.
- `--raw` prints formatted JSON from the Figma REST API.
- `--verbose` prints request diagnostics without token values.
- `FIGMA_API_BASE_URL` can override the API host for local testing.
- Large file payloads should be requested deliberately with `--depth` or `--raw`.

## Roadmap

- Components, styles, and variables inspection.
- Design token export to CSS / JSON.
- Local response cache keyed by file key, node ID, and version.
- GoReleaser and Homebrew packaging.
- Optional skills for Codex and Claude Code.

## Development

```bash
make test
make build
```

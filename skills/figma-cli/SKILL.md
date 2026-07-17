---
name: figma-cli
description: Use when reading Figma files or nodes, exporting node images or assets, listing Figma comments, components, styles, or variables through the figma-cli command line, or setting up Figma token credentials. Token-based reads work headlessly — no Figma app or browser required.
---

# Figma CLI Skill

Use this skill when an agent needs stable, reproducible reads from the Figma REST API through `figma-cli`.

## CLI Setup

- Before using Figma commands, check whether `figma-cli` is available with `command -v figma-cli`.
- If missing, install with `brew install --cask wohsj110/tap/figma-cli`; without Homebrew, `go install github.com/wohsj110/figma_cli/cmd/figma-cli@latest` (requires Go ≥ 1.26).

## Credential Setup

- Verify access first: `figma-cli me`. If it prints the current Figma user, credentials are ready.
- If not configured, ask the user for a Figma personal access token, then store it in the OS keyring:

```bash
printf %s "$FIGMA_TOKEN" | figma-cli init --token-stdin
```

- For CI or one-off runs, `FIGMA_TOKEN` env takes precedence over the keyring.
- Never print, log, or hardcode token values. The CLI never prints the token; keep it that way in your own output too.

## Core Rules

- Paste full Figma URLs directly — file keys and node ids resolve automatically. **Never guess a node id**; if you don't have one, ask the user or list metadata first.
- Default output is concise text. Add `--raw` only when the original API JSON is genuinely needed.
- GET responses are cached for 15 minutes. Use `--no-cache` when freshness matters, `--cache-ttl 1h` to extend.
- Request large files deliberately: start with `--depth 2` before any deep dump.
- `--verbose` prints request diagnostics without token values — use it when debugging API errors.

## Commands

```bash
figma-cli me                                                        # current user / credential check
figma-cli file get "https://www.figma.com/design/KEY/name?node-id=1-2"
figma-cli file get FILE_KEY --depth 2                               # bounded file summary
figma-cli node inspect FILE_KEY --node 1:2 [--raw]                  # node summary / raw API JSON
figma-cli image export FILE_KEY --node 1:2 --format png --scale 2   # export image URL
figma-cli image export FILE_KEY --node 1:2 --format png --out ./dir # download asset locally
figma-cli comments list FILE_KEY                                    # review context
figma-cli components list FILE_KEY
figma-cli styles list FILE_KEY
figma-cli variables list FILE_KEY                                   # design tokens
```

## Boundaries

- This CLI is read-only inspection and asset export by design — it never writes back to Figma.

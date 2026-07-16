# AGENTS.md

This repository builds `figma-cli`, a Figma REST API command-line toolkit for humans and AI agents.

## Language
- Think in English.
- Reply to the user in Simplified Chinese unless they request otherwise.

## Product Goal
Build a small, stable CLI for repeatable Figma inspection and asset export.

The CLI should complement Figma MCP workflows:
- Use MCP or browser automation for interactive canvas context and write-back workflows.
- Use this CLI for token-based, reproducible reads, summaries, and asset downloads.

## Principles
- Default output must be concise, stable, and useful to agents.
- Raw JSON must be opt-in via `--raw`.
- Accept Figma file/design/proto URLs directly; do not force users to manually extract file keys.
- Never log or echo Figma tokens.
- Store tokens in the OS keyring when possible; `FIGMA_TOKEN` is supported for CI and ephemeral use.
- Avoid broad REST mirroring. Add commands around real user workflows.

## MVP Commands
- `figma-cli init`
- `figma-cli me`
- `figma-cli file get URL_OR_KEY`
- `figma-cli node inspect URL_OR_KEY --node NODE_ID`
- `figma-cli image export URL_OR_KEY --node NODE_ID`
- `figma-cli comments list URL_OR_KEY`
- `figma-cli components list URL_OR_KEY`
- `figma-cli styles list URL_OR_KEY`
- `figma-cli variables list URL_OR_KEY`

## Testing
- Unit-test URL parsing and request construction.
- Use mocked HTTP servers for API behavior.
- Do not require real Figma credentials in CI.
- Use `python3 scripts/mock_verify.py` for end-to-end command verification without real credentials.

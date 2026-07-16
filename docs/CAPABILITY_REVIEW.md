# Capability Review

This project intentionally builds a small read-first Figma CLI instead of a full Figma automation surface.

## Official API Boundary

Authoritative Figma surfaces used by this CLI:

- Figma REST API authentication supports OAuth, plan access tokens, and personal access tokens. Local CLI usage is optimized for personal access tokens through `FIGMA_TOKEN` or OS keyring storage.
- `GET /v1/me` identifies the current user.
- `GET /v1/files/:key` returns file metadata, the document tree, and local component/style metadata.
- `GET /v1/files/:key/nodes?ids=...` returns selected nodes for focused inspection.
- `GET /v1/images/:key?ids=...` returns expiring rendered image URLs.
- `GET /v1/files/:key/comments` returns file comments for review context.
- `GET /v1/files/:file_key/variables/local` returns local and referenced variables, but Figma documents this area as seat/plan restricted. The CLI treats it as read-only and surfaces API errors directly.

References:

- Figma REST API overview: https://developers.figma.com/docs/rest-api/
- File endpoints: https://developers.figma.com/docs/rest-api/file-endpoints/
- Comment endpoints: https://developers.figma.com/docs/rest-api/comments-endpoints/
- Variable endpoints: https://developers.figma.com/docs/rest-api/variables-endpoints/
- MCP server: https://developers.figma.com/docs/figma-mcp-server/

## Lessons From Existing Tooling

Useful patterns to preserve:

- Accept full Figma URLs, not only file keys.
- Default to concise inspection output and keep raw JSON opt-in.
- Export image URLs and optionally download assets.
- Cache repeated reads to avoid heavy file fetches during agent loops.
- Make mock-server verification easy so CI does not require real Figma credentials.

Patterns to avoid in v1:

- Mirroring every REST endpoint.
- Canvas write-back or desktop automation.
- OAuth app management, account systems, or cloud sync.
- Default full-file JSON dumps that are too large for agents.
- Logging request headers or token values.

## Product Positioning

`figma-cli` complements Figma MCP:

- MCP is better for interactive design context, current selection, and canvas write-back.
- This CLI is better for reproducible token-based reads, deterministic command output, asset export, and CI/mock verification.

## First-Version Command Surface

- `figma-cli init`
- `figma-cli me`
- `figma-cli file get`
- `figma-cli node inspect`
- `figma-cli image export`
- `figma-cli comments list`
- `figma-cli components list`
- `figma-cli styles list`
- `figma-cli variables list`

Non-goals for v1:

- Updating Figma files.
- Creating comments or reactions.
- Webhook management.
- Organization discovery/activity logs.
- Publishing release artifacts automatically without explicit user approval.

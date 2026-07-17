# Releasing

Releases ship GoReleaser-built binaries on GitHub plus a Homebrew cask published to `wohsj110/homebrew-tap`.

## Preflight

```bash
go test ./...
make build
python3 scripts/mock_verify.py
go run github.com/goreleaser/goreleaser/v2@latest check
```

## Tag Release

```bash
git tag v0.1.0
git push origin v0.1.0
GITHUB_TOKEN="$(gh auth token)" HOMEBREW_TAP_GITHUB_TOKEN="$(gh auth token)" go run github.com/goreleaser/goreleaser/v2@latest release --clean
```

## Homebrew

The `homebrew_casks` section in `.goreleaser.yml` publishes the `figma-cli` cask to `wohsj110/homebrew-tap` on every release. Export `HOMEBREW_TAP_GITHUB_TOKEN` (a token with push access to the tap repo) before running `goreleaser release`. Users install with:

```bash
brew install --cask wohsj110/tap/figma-cli
```

## Safety

- Never include a Figma token in release logs.
- Do not run real-file smoke tests against private files in CI.
- Use `FIGMA_API_BASE_URL` and `scripts/mock_verify.py` for credential-free verification.

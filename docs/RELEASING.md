# Releasing

The first release target is a GitHub release with GoReleaser-built binaries. Homebrew can be enabled by adding a tap repository after the GitHub repository name is final.

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
GITHUB_TOKEN="$(gh auth token)" go run github.com/goreleaser/goreleaser/v2@latest release --clean
```

## Homebrew

Homebrew cask publishing is intentionally not enabled until the repository and tap names are confirmed. After that, add a `homebrew_casks` section to `.goreleaser.yml`.

## Safety

- Never include a Figma token in release logs.
- Do not run real-file smoke tests against private files in CI.
- Use `FIGMA_API_BASE_URL` and `scripts/mock_verify.py` for credential-free verification.

# Releasing blindpay-go

## Version source of truth

`const Version` in [`blindpay.go`](blindpay.go) (`X.Y.Z`, no `v`). It feeds
the `User-Agent` and the release tag `vX.Y.Z`. The module stays on v1 forever
(no `/v2` path): MINOR for new methods, fields, enum values and error codes,
PATCH for fixes (CLAUDE.md section 9). `cmd/apisync -apply` bumps it in sync
PRs; humans bump it in their PR.

## What a merge to `main` does

1. `Main` (`.github/workflows/main.yaml`) runs lint, tests (Go 1.21-1.23,
   oldstable, stable, `-race`), contract check, apisync check and Snyk.
   Nothing gates the merge on it: `main` has no branch protection.
2. `Release` (`.github/workflows/release.yml`) reads `const Version`. If the
   tag does not exist yet, it builds, runs the tests, pushes the tag
   `vX.Y.Z` at the merge commit and creates the GitHub release with
   generated notes. If the tag exists it does nothing, so re-runs are safe.
3. The Go module proxy serves the version as soon as the tag exists; there
   is no registry upload.

`API Sync` (`repository_dispatch` from blindpay-v2) opens an auto-merging
`api-sync` PR; its merge goes through the same two steps.

## Secrets and permissions

| What | Used by | Owner |
|---|---|---|
| `GITHUB_TOKEN` (`contents: write`) | `release.yml` (tag + release) | GitHub, automatic |
| `SDK_SYNC_PAT` | `api-sync.yml`, `api-sync-merged.yml` | BlindPay org admin |
| `SNYK_TOKEN` | `snyk.yaml` | Snyk org `bernardosm` admin |
| `SLACK_WEBHOOK_URL` | `pipeline-alert.yml` | BlindPay org admin |

No registry credentials exist or are needed.

## Re-running a stuck release

Tag-on-push works (push run 43093086939 tagged v1.20.0 at the #51 merge),
but v1.21.0 needed a manual `workflow_dispatch` because its push run stayed
queued during a runner incident. After every merge that bumps the version,
check:

```sh
gh run list -R blindpaylabs/blindpay-go -w Release -L 3
```

If the push run is queued, failed or missing, cancel it and dispatch the
workflow on `main` (idempotent; it skips when the tag exists):

```sh
gh run cancel <run-id> -R blindpaylabs/blindpay-go
gh workflow run release.yml -R blindpaylabs/blindpay-go --ref main
```

Never tag by hand: the workflow tags only after the build and tests pass.

## Verifying the published version

```sh
GOPROXY=https://proxy.golang.org go list -m -versions github.com/blindpaylabs/blindpay-go
GOPROXY=https://proxy.golang.org go list -m github.com/blindpaylabs/blindpay-go@vX.Y.Z
```

Then open `https://pkg.go.dev/github.com/blindpaylabs/blindpay-go@vX.Y.Z`
(it may take a few minutes to appear) and build a scratch module with
`go get github.com/blindpaylabs/blindpay-go@vX.Y.Z`.

## Known traps

- A queued push run (runner incident) leaves the version untagged: use the
  dispatch above.
- `release.yml` runs only build and tests, not lint, contract or apisync
  checks; keep `Main` green before merging.
- `-apply` in a PR that already bumped the version bumps it again; reset it.
- A version that was ever tagged cannot be reused (the proxy caches it).

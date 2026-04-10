# Versioning & Release Guide

This document explains how to properly version and release new versions of the `blindpay-go` SDK.

## Versioning Scheme

We follow **Semantic Versioning** (`v1.MINOR.PATCH`) and stay on the **v1.x** line.

> **Important:** Go requires a `/v2` import path suffix for v2+ modules, which forces all consumers to update their imports. To avoid this, we stay on v1.x indefinitely. See [Go's module version numbering](https://go.dev/doc/modules/version-numbers) for details.

| Version Part | When to Increment | Example |
|---|---|---|
| **PATCH** (`v1.0.X`) | Bug fixes, typo corrections, small non-breaking changes | `v1.3.0` → `v1.3.1` |
| **MINOR** (`v1.X.0`) | New features, new fields, new enum values, breaking changes to existing types | `v1.3.0` → `v1.4.0` |

### Examples

- Adding a new enum type (e.g. `AccountPurpose`) → **MINOR** bump
- Adding a new field to a struct (e.g. `SelfieFile`) → **MINOR** bump
- Renaming a field (e.g. `IndividualHoldingDocFrontFile` → `SelfieFile`) → **MINOR** bump (breaking)
- Changing field types (e.g. `string` → `*string`) → **MINOR** bump (breaking)
- Fixing a wrong enum value (e.g. `"savings"` → `"saving"`) → **PATCH** bump
- Fixing a typo in a JSON tag → **PATCH** bump

## Release Checklist

Every PR that changes SDK behavior **must** include these steps:

### 1. Bump the Version Constant

Update the `Version` constant in `blindpay.go`:

```go
// Before
const Version = "1.3.0"

// After
const Version = "1.4.0"
```

### 2. Make Sure Tests Pass

```bash
go build ./...
go test ./...
```

### 3. Merge the PR

Merge the PR into `main` through GitHub.

### 4. Pull Main Locally

```bash
git checkout main
git pull origin main
```

### 5. Create and Push the Git Tag

The tag **must** match the version in `blindpay.go`:

```bash
git tag v1.4.0
git push origin v1.4.0
```

### 6. Create a GitHub Release

```bash
gh release create v1.4.0 --title "v1.4.0 - Short description" --notes "$(cat <<'EOF'
## What's Changed

- Added new enum types: `AccountPurpose`, `BusinessType`, `BusinessIndustry`
- Added `SelfieFile` field to enhanced KYC params
- Updated `Owner` struct with `OwnershipPercentage` and `Title` fields

## Breaking Changes

- `CreateIndividualStandardParams`: most fields changed from `string` to `*string`
- `CreateIndividualEnhancedParams`: most fields changed from `string` to `*string`
- `CreateBusinessStandardParams`: most fields changed from `string` to `*string`

**Full Changelog**: https://github.com/blindpaylabs/blindpay-go/compare/v1.3.0...v1.4.0
EOF
)"
```

### 7. Verify on pkg.go.dev

Visit `https://pkg.go.dev/github.com/blindpaylabs/blindpay-go@v1.4.0` to trigger indexing.

It may take a few minutes for the new version to appear as "Latest" on the versions tab.

## Common Mistakes to Avoid

### 1. Forgetting to bump the version constant
The `Version` constant in `blindpay.go` **must** be updated in the same PR as the code changes. If you forget, the tag and the constant will be out of sync.

### 2. Creating a v2+ tag
**Never** create tags like `v2.0.0`. Go will not resolve them unless the module path in `go.mod` includes `/v2`. Stick to `v1.x.x`.

### 3. Pushing a tag before merging
Always merge the PR first, then pull main, then create the tag. Otherwise the tag will point to the wrong commit.

### 4. Mismatched tag and constant
The git tag (e.g. `v1.4.0`) and the `Version` constant in `blindpay.go` (e.g. `"1.4.0"`) must always match. The tag has the `v` prefix, the constant does not.

## Quick Reference

```bash
# Full release flow after PR is merged
git checkout main
git pull origin main
git tag v1.X.0
git push origin v1.X.0
gh release create v1.X.0 --title "v1.X.0 - Description" --generate-notes
```

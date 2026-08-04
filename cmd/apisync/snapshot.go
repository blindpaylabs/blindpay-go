package main

import (
	"os"
	"path/filepath"
)

// refreshSnapshot copies the applied spec's raw bytes verbatim over
// .api-sync/spec-snapshot.json. It must never re-marshal the parsed
// document: encoding/json would reorder object keys into map order,
// re-indent, and re-escape, turning every future no-op refresh into an
// unreviewable diff of the entire file (and the committed snapshot would
// stop matching the bytes blindpay-v2 actually ships as spec-current.json).
func refreshSnapshot(repoRoot string, specBytes []byte) error {
	path := filepath.Join(repoRoot, ".api-sync", "spec-snapshot.json")
	return os.WriteFile(path, specBytes, 0o644)
}

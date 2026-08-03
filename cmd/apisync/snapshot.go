package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

// refreshSnapshot re-serializes the applied spec deterministically (Go's
// encoding/json always sorts map keys) and writes it over
// .api-sync/spec-snapshot.json, so code and baseline never drift apart.
func refreshSnapshot(repoRoot string, spec map[string]any) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	if err := enc.Encode(spec); err != nil {
		return err
	}
	path := filepath.Join(repoRoot, ".api-sync", "spec-snapshot.json")
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

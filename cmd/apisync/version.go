package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var versionLineRE = regexp.MustCompile(`const Version = "(\d+)\.(\d+)\.(\d+)"`)

// bumpVersion reads blindpay.go's `const Version = "X.Y.Z"`, bumps it
// according to class ("minor" or "patch"; major is never automatic), and
// writes it back in the exact literal form release.yml's
// `grep -oP 'const Version = "\K[^"]+'` expects.
func bumpVersion(repoRoot, class string) (string, error) {
	path := filepath.Join(repoRoot, "blindpay.go")
	data, err := os.ReadFile(path) //#nosec G304 -- path is developer-supplied (CLI arg or repo-relative), this is a local codegen tool
	if err != nil {
		return "", err
	}

	loc := versionLineRE.FindSubmatchIndex(data)
	if loc == nil {
		return "", fmt.Errorf("blindpay.go: const Version line not found")
	}
	major, err := strconv.Atoi(string(data[loc[2]:loc[3]]))
	if err != nil {
		return "", err
	}
	minor, err := strconv.Atoi(string(data[loc[4]:loc[5]]))
	if err != nil {
		return "", err
	}
	patch, err := strconv.Atoi(string(data[loc[6]:loc[7]]))
	if err != nil {
		return "", err
	}

	switch class {
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		return "", fmt.Errorf("unknown bump class %q", class)
	}

	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	newLine := fmt.Sprintf(`const Version = "%s"`, newVersion)

	out := make([]byte, 0, len(data))
	out = append(out, data[:loc[0]]...)
	out = append(out, newLine...)
	out = append(out, data[loc[1]:]...)

	if err := os.WriteFile(path, out, 0o600); err != nil { //#nosec G703 -- path is repoRoot/blindpay.go, repoRoot comes from the developer CLI arg
		return "", err
	}
	return newVersion, nil
}

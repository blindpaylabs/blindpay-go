package blindpay

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	readmeGoBlockRE = regexp.MustCompile("(?s)```go\n(.*?)```")
	readmeDeclRE    = regexp.MustCompile(`(?m)^((?:\w+, )*\w+) :=`)
)

// TestREADMESnippetsCompile extracts every ```go block from README.md and
// type-checks it against this checkout in a scratch module (go vet). A full
// program (package main) is checked as is; a fragment is wrapped in a
// function whose parameters provide the names the README leaves to the
// reader (client, ctx, params...); a func declaration is kept top level.
func TestREADMESnippetsCompile(t *testing.T) {
	if testing.Short() {
		t.Skip("runs the go command")
	}
	readme, err := os.ReadFile("README.md")
	require.NoError(t, err)
	blocks := readmeGoBlockRE.FindAllStringSubmatch(string(readme), -1)
	require.NotEmpty(t, blocks)

	root, err := filepath.Abs(".")
	require.NoError(t, err)
	dir := t.TempDir()
	write := func(name, content string) {
		require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}
	goMod, err := os.ReadFile("go.mod")
	require.NoError(t, err)
	goVersion := regexp.MustCompile(`(?m)^go (\S+)`).FindStringSubmatch(string(goMod))[1]
	write("go.mod", fmt.Sprintf("module example.com/readmecheck\n\ngo %s\n\nrequire github.com/blindpaylabs/blindpay-go v1.0.0\n\nreplace github.com/blindpaylabs/blindpay-go => %s\n", goVersion, root))
	sum, err := os.ReadFile("go.sum")
	require.NoError(t, err)
	write("go.sum", string(sum))

	var top, funcs []string
	programs := 0
	for i, m := range blocks {
		block := m[1]
		switch {
		case strings.HasPrefix(block, "package main"):
			programs++
			write(fmt.Sprintf("program%d/main.go", i), block)
		case strings.HasPrefix(strings.TrimSpace(block), "func "):
			top = append(top, block)
		default:
			var uses []string
			seen := map[string]bool{}
			for _, d := range readmeDeclRE.FindAllStringSubmatch(block, -1) {
				for _, name := range strings.Split(d[1], ", ") {
					if name != "_" && !seen[name] {
						seen[name] = true
						uses = append(uses, "\t_ = "+name)
					}
				}
			}
			funcs = append(funcs, fmt.Sprintf("func snippet%d(ctx context.Context, client *blindpay.Client, apiKey, instanceID string, httpClient *http.Client, params *payouts.CreateEvmParams, quoteParams *quotes.CreateParams) error {\n%s\n%s\n\treturn nil\n}\n",
				i, block, strings.Join(uses, "\n")))
		}
	}
	write("snippets/main.go", `package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blindpaylabs/blindpay-go"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/quotes"
	"github.com/blindpaylabs/blindpay-go/upload"
)

// The README fragments use different subsets of these imports.
var (
	_ = errors.As
	_ = fmt.Println
	_ = io.ReadAll
	_ = log.Fatal
	_ = os.Open
	_ = time.Second
	_ = upload.BucketOnboarding
)

func main() {}

`+strings.Join(top, "\n")+"\n"+strings.Join(funcs, "\n"))

	goCmd := filepath.Join(runtime.GOROOT(), "bin", "go")
	cmd := exec.Command(goCmd, "vet", "./...") //#nosec G204 -- fixed go binary, test-owned temp dir
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOWORK=off")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "README snippets do not compile:\n%s", out)
	require.Equal(t, 1, programs, "the quick start is the only full program")
	t.Logf("%d README Go blocks type-check", len(blocks))
}

package tol

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// newContext builds a *cli.Context from a plain map of string flag-name → value,
// an optional map of string-slice flag-name → values, and an optional map of
// int flag-name → value.
func newContext(
	t *testing.T,
	strFlags map[string]string,
	sliceFlags map[string][]string,
	intFlags map[string]int,
) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet("test", flag.ContinueOnError)

	for _, name := range []string{
		"csv-path", "document-path", "header", "limiter-left", "limiter-right",
	} {
		set.String(name, "", "")
	}

	for _, name := range []string{"index"} {
		set.Int(name, 0, "")
	}

	for _, name := range []string{"skip"} {
		set.Var(cli.NewStringSlice(), name, "")
	}

	for k, v := range strFlags {
		require.NoError(t, set.Set(k, v), "flag.Set(%q)", k)
	}

	for k, vals := range sliceFlags {
		for _, v := range vals {
			require.NoError(t, set.Set(k, v), "flag.Set(%q, %q)", k, v)
		}
	}

	for k, v := range intFlags {
		require.NoError(t, set.Set(k, strconv.Itoa(v)), "flag.Set(%q)", k)
	}

	cliApp := cli.NewApp()
	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{Name: "csv-path"},
		&cli.StringFlag{Name: "document-path"},
		&cli.StringFlag{Name: "header"},
		&cli.StringFlag{Name: "limiter-left"},
		&cli.StringFlag{Name: "limiter-right"},
		&cli.IntFlag{Name: "index"},
		&cli.StringSliceFlag{Name: "skip"},
	}

	return cli.NewContext(cliApp, set, nil)
}

// writeTempFile creates a temporary file with content inside t.TempDir().
func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "doc-*.md")
	require.NoError(t, err)

	_, err = f.WriteString(content)
	require.NoError(t, err)

	_ = f.Close()

	return f.Name()
}

// writeTempFileNamed creates a file with a specific name inside t.TempDir().
func writeTempFileNamed(t *testing.T, name, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}

// writeTempCSV creates a temporary CSV file with the given rows inside t.TempDir().
func writeTempCSV(t *testing.T, rows []string) string {
	t.Helper()

	return writeTempFileNamed(t, "licenses.csv", strings.Join(rows, "\n")+"\n")
}

// resetGlobals resets all package-level variables mutated by before().
func resetGlobals() {
	csvPath = ""
	documentPath = ""
	header = ""
	limiterL = ""
	limiterR = ""
	skip = nil
	index = 0
	scsv = nil
	body = nil
	rexp = nil
	prefix = ""
	suffix = ""
}

// setupBefore writes content to temp files, runs before(), and returns the
// document and csv paths. The caller owns resetGlobals via t.Cleanup.
func setupBefore(
	t *testing.T,
	docContent, csvContent, hdr, limL, limR string,
	skipPkgs []string,
) (string, string) {
	t.Helper()

	docPath := writeTempFile(t, docContent)
	csvFilePath := writeTempCSV(t, strings.Split(strings.TrimSpace(csvContent), "\n"))

	flags := map[string]string{
		"csv-path":      csvFilePath,
		"document-path": docPath,
		"header":        hdr,
		"limiter-left":  limL,
	}
	if limR != "" {
		flags["limiter-right"] = limR
	}

	sliceFlags := map[string][]string{}
	if len(skipPkgs) > 0 {
		sliceFlags["skip"] = skipPkgs
	}

	require.NoError(
		t,
		before(newContext(t, flags, sliceFlags, nil)),
		"setupBefore: before() failed",
	)

	return docPath, csvFilePath
}

func TestBefore(t *testing.T) {
	tests := []struct {
		name       string
		docContent string
		csvContent string
		hdr        string
		limL       string
		limR       string
		wantErr    error
	}{
		{
			name:       "valid inputs succeed",
			docContent: "## Licenses\n\nold content\n\n",
			csvContent: "pkg/a,https://example.com/LICENSE,MIT",
			hdr:        "Licenses",
			limL:       "##",
		},
		{
			name:       "missing csv file returns errFileOpen",
			docContent: "## Licenses\n\nold content\n\n",
			csvContent: "",
			hdr:        "Licenses",
			limL:       "##",
			wantErr:    errFileOpen,
		},
		{
			name:       "missing document file returns errFileRead",
			docContent: "",
			csvContent: "pkg/a,https://example.com/LICENSE,MIT",
			hdr:        "Licenses",
			limL:       "##",
			wantErr:    errFileRead,
		},
		{
			name:       "section not found returns errExtract",
			docContent: "## Other\n\nsome content\n\n",
			csvContent: "pkg/a,https://example.com/LICENSE,MIT",
			hdr:        "Licenses",
			limL:       "##",
			wantErr:    errExtract,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(resetGlobals)

			var docPath, csvFilePath string

			switch {
			case errors.Is(tt.wantErr, errFileOpen):
				docPath = writeTempFile(t, tt.docContent)
				csvFilePath = filepath.Join(t.TempDir(), "nonexistent.csv")
			case errors.Is(tt.wantErr, errFileRead):
				csvFilePath = writeTempCSV(t, []string{tt.csvContent})
				docPath = filepath.Join(t.TempDir(), "nonexistent.md")
			default:
				docPath = writeTempFile(t, tt.docContent)
				csvFilePath = writeTempCSV(t, []string{tt.csvContent})
			}

			flags := map[string]string{
				"csv-path":      csvFilePath,
				"document-path": docPath,
				"header":        tt.hdr,
				"limiter-left":  tt.limL,
			}
			if tt.limR != "" {
				flags["limiter-right"] = tt.limR
			}

			err := before(newContext(t, flags, nil, nil))

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAction(t *testing.T) {
	const (
		csvContent = `pkg/a,https://example.com/a/LICENSE,MIT
pkg/b,https://example.com/b/LICENSE,Apache-2.0
pkg/skip-me,https://example.com/skip/LICENSE,BSD-2-Clause`
		docHeader = "Licenses"
		limL      = "##"
		limR      = "##"
	)

	tests := []struct {
		name         string
		docContent   string
		skip         []string
		index        int
		wantErr      error
		wantContains []string
		wantAbsent   []string
	}{
		{
			name:       "generates table with all packages",
			docContent: "## Licenses\n\nold content\n\n## Next\n",
			wantContains: []string{
				"| Package ",
				"| Licence ",
				"| Type    ",
				"| pkg/a ",
				"| pkg/b ",
				"| pkg/skip-me ",
			},
		},
		{
			name:       "skip flag excludes specified package",
			docContent: "## Licenses\n\nold content\n\n## Next\n",
			skip:       []string{"pkg/skip-me"},
			wantContains: []string{
				"| pkg/a ",
				"| pkg/b ",
			},
			wantAbsent: []string{
				"pkg/skip-me",
			},
		},
		{
			name:       "multiple skip entries are all excluded",
			docContent: "## Licenses\n\nold content\n\n## Next\n",
			skip:       []string{"pkg/a", "pkg/b"},
			wantContains: []string{
				"| pkg/skip-me ",
			},
			wantAbsent: []string{
				"pkg/a",
				"pkg/b",
			},
		},
		{
			name:       "separator row is correctly formatted",
			docContent: "## Licenses\n\nold content\n\n## Next\n",
			wantContains: []string{
				"|---",
			},
		},
		{
			name:       "table contains correct license types",
			docContent: "## Licenses\n\nold content\n\n## Next\n",
			wantContains: []string{
				"MIT",
				"Apache-2.0",
				"BSD-2-Clause",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(resetGlobals)

			docPath, _ := setupBefore(t, tt.docContent, csvContent, docHeader, limL, limR, tt.skip)

			if tt.index != 0 {
				index = tt.index
			}

			err := action(
				cli.NewContext(cli.NewApp(), flag.NewFlagSet("test", flag.ContinueOnError), nil),
			)

			require.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				updated, readErr := os.ReadFile(docPath) // #nosec G304
				require.NoError(t, readErr)

				got := string(updated)

				for _, sub := range tt.wantContains {
					assert.Contains(t, got, sub)
				}

				for _, sub := range tt.wantAbsent {
					assert.NotContains(t, got, sub)
				}
			}
		})
	}
}

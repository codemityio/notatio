package toc

import (
	"errors"
	"flag"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// newContext builds a *cli.Context from a plain map of string flag-name → value
// and an optional map of string-slice flag-name → values.
//
// StringSliceFlags must be registered using cli.StringSlice as the underlying
// flag.Value — not as a plain string — otherwise c.StringSlice() returns an
// empty slice regardless of what is set on the flag.FlagSet.
func newContext(
	t *testing.T,
	strFlags map[string]string,
	sliceFlags map[string][]string,
) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet("test", flag.ContinueOnError)

	// Plain string flags.
	for _, name := range []string{
		"document", "header", "limiter-left", "limiter-right",
		"summary-header", "summary-limiter-left", "summary-limiter-right",
	} {
		set.String(name, "", "")
	}

	// Int flags.
	for _, name := range []string{"start-from-level", "start-from-item"} {
		set.Int(name, 0, "")
	}

	// StringSlice flags: register each name with a *cli.StringSlice value so
	// that cli.Context.StringSlice() can find and unwrap it correctly.
	sliceNames := []string{"path"}
	for _, name := range sliceNames {
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

	// The app flags list must mirror the set so cli.Context resolves names.
	cliApp := cli.NewApp()
	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{Name: "document"},
		&cli.StringFlag{Name: "header"},
		&cli.StringFlag{Name: "limiter-left"},
		&cli.StringFlag{Name: "limiter-right"},
		&cli.StringFlag{Name: "summary-header"},
		&cli.StringFlag{Name: "summary-limiter-left"},
		&cli.StringFlag{Name: "summary-limiter-right"},
		&cli.IntFlag{Name: "start-from-level"},
		&cli.IntFlag{Name: "start-from-item"},
		&cli.StringSliceFlag{Name: "path"},
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

// resetGlobals resets all package-level variables mutated by before().
func resetGlobals() {
	document = ""
	header = ""
	limiterL = ""
	limiterR = ""
	body = nil
	rexp = nil
	prefix = ""
	suffix = ""
}

// setupBefore writes content to a temp file, runs before(), and returns the
// document path. The caller owns resetGlobals via t.Cleanup.
func setupBefore(t *testing.T, content, hdr, limL, limR string) string {
	t.Helper()

	docPath := writeTempFile(t, content)

	flags := map[string]string{
		"document":     docPath,
		"header":       hdr,
		"limiter-left": limL,
	}
	if limR != "" {
		flags["limiter-right"] = limR
	}

	require.NoError(t, before(newContext(t, flags, nil)), "setupBefore: before() failed")

	return docPath
}

func TestGenerateAnchor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain lowercase title",
			input: "hello world",
			want:  "hello-world",
		},
		{
			name:  "uppercase is lowercased",
			input: "Hello World",
			want:  "hello-world",
		},
		{
			name:  "special characters are replaced with hyphens",
			input: "Hello, World!",
			want:  "hello-world",
		},
		{
			name:  "leading and trailing hyphens are trimmed",
			input: "!hello!",
			want:  "hello",
		},
		{
			name:  "numbers are preserved",
			input: "Step 1 Setup",
			want:  "step-1-setup",
		},
		{
			name:  "consecutive special characters collapse to one hyphen",
			input: "A -- B",
			want:  "a-b",
		},
		{
			name:  "already valid anchor is unchanged",
			input: "valid-anchor",
			want:  "valid-anchor",
		},
		{
			name:  "empty string returns empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, generateAnchor(tt.input))
		})
	}
}

func TestGenerateInternalTOC(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		startFromLevel int
		startFromItem  int
		wantErr        error
		wantContains   []string
		wantAbsent     []string
	}{
		{
			name:    "basic headers produce toc entries",
			content: "# Doc\n\n## Installation\n\n## Usage\n\n## Contributing\n",
			wantContains: []string{
				"- [Installation](#installation)",
				"- [Usage](#usage)",
				"- [Contributing](#contributing)",
			},
		},
		{
			name:    "nested headers are indented",
			content: "# Doc\n\n## Usage\n\n### Basic\n\n### Advanced\n",
			wantContains: []string{
				"- [Usage](#usage)",
				"    - [Basic](#basic)",
				"    - [Advanced](#advanced)",
			},
		},
		{
			name:           "start-from-level skips shallow headers",
			content:        "# Doc\n\n## Usage\n\n### Details\n",
			startFromLevel: 2,
			wantAbsent:     []string{"[Doc]"},
			wantContains:   []string{"[Details]"},
		},
		{
			name:          "start-from-item skips first N headers",
			content:       "# Doc\n\n## First\n\n## Second\n\n## Third\n",
			startFromItem: 2,
			wantAbsent:    []string{"[Doc]", "[First]"},
			wantContains:  []string{"[Third]"},
		},
		{
			name:         "headers inside code blocks are ignored",
			content:      "# Doc\n\n## Real\n\n```\n## FakeInsideBlock\n```\n\n## AfterBlock\n",
			wantContains: []string{"[Real]", "[AfterBlock]"},
			wantAbsent:   []string{"[FakeInsideBlock]"},
		},
		{
			name:         "anchor is correctly generated for titles with spaces",
			content:      "# Doc\n\n## Hello World\n",
			wantContains: []string{"(#hello-world)"},
		},
		{
			name:    "non-existent file returns errDocumentOpen",
			content: "", // file path is overridden below
			wantErr: errDocumentOpen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if errors.Is(tt.wantErr, errDocumentOpen) {
				path = filepath.Join(t.TempDir(), "nonexistent.md")
			} else {
				path = writeTempFile(t, tt.content)
			}

			got, err := generateInternalTOC(tt.startFromLevel, tt.startFromItem, "#", "TOC", path)

			require.ErrorIs(t, err, tt.wantErr)

			for _, sub := range tt.wantContains {
				assert.Contains(t, got, sub)
			}

			for _, sub := range tt.wantAbsent {
				assert.NotContains(t, got, sub)
			}
		})
	}
}

func TestGenerateExternalTOC(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string]string // filename → content
		summaryHeader string
		summaryLimL   string
		summaryLimR   string
		wantErr       error
		wantContains  []string
		wantAbsent    []string
		wantOrdered   []string // substrings that must appear in this relative order
	}{
		{
			name: "single file without summary produces simple list entry",
			files: map[string]string{
				"a.md": "# Alpha\n\nSome content.\n",
			},
			wantContains: []string{"- [Alpha]"},
		},
		{
			name: "multiple files are sorted alphabetically by title",
			files: map[string]string{
				"z.md": "# Zebra\n\ncontent\n",
				"a.md": "# Apple\n\ncontent\n",
				"m.md": "# Mango\n\ncontent\n",
			},
			wantOrdered: []string{"Apple", "Mango", "Zebra"},
		},
		{
			name: "file without a title returns errTitleNotFound",
			files: map[string]string{
				"notitle.md": "no heading here\n",
			},
			wantErr: errTitleNotFound,
		},
		{
			name:  "missing file returns errFileRead",
			files: map[string]string{
				// path injected as nonexistent in test loop
			},
			wantErr: errFileRead,
		},
		{
			name: "summary section is appended after title",
			files: map[string]string{
				"a.md": "# Alpha\n\n## Summary\n\nThis is the summary.\n\n## Next\n",
			},
			summaryHeader: "Summary",
			summaryLimL:   "##",
			summaryLimR:   "##",
			wantContains:  []string{"This is the summary."},
		},
		{
			name: "summary with no right limiter reads to EOF",
			files: map[string]string{
				"a.md": "# Alpha\n\n## Summary\n\nEOF summary content.\n",
			},
			summaryHeader: "Summary",
			summaryLimL:   "##",
			summaryLimR:   "",
			wantContains:  []string{"EOF summary content."},
		},
		{
			name: "invalid summary regex returns errRegexCompile",
			files: map[string]string{
				"a.md": "# Alpha\n\ncontent\n",
			},
			summaryHeader: "Summary",
			summaryLimL:   "[invalid",
			summaryLimR:   "##",
			wantErr:       errRegexCompile,
		},
		{
			name: "toc header and prefix are included in output",
			files: map[string]string{
				"a.md": "# Alpha\n\ncontent\n",
			},
			wantContains: []string{"# TOC"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var paths []string

			if errors.Is(tt.wantErr, errFileRead) {
				paths = []string{filepath.Join(t.TempDir(), "nonexistent.md")}
			} else {
				for name, content := range tt.files {
					paths = append(paths, writeTempFileNamed(t, name, content))
				}
			}

			got, err := generateExternalTOC(
				"#",
				"TOC",
				paths,
				tt.summaryHeader,
				tt.summaryLimL,
				tt.summaryLimR,
			)

			require.ErrorIs(t, err, tt.wantErr)

			for _, sub := range tt.wantContains {
				assert.Contains(t, got, sub)
			}

			for _, sub := range tt.wantAbsent {
				assert.NotContains(t, got, sub)
			}

			if len(tt.wantOrdered) > 1 {
				assertOrder(t, got, tt.wantOrdered)
			}
		})
	}
}

// assertOrder checks that all items in order appear in s in the given sequence.
func assertOrder(t *testing.T, s string, order []string) {
	t.Helper()

	pos := 0
	for _, item := range order {
		idx := strings.Index(s[pos:], item)
		require.GreaterOrEqualf(t, idx, 0,
			"expected %q to appear after previous items in:\n%s", item, s)
		pos += idx + len(item)
	}
}

func TestInternal(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		startFromLevel int
		startFromItem  int
		wantErr        error
		wantContains   []string
	}{
		{
			name:         "writes toc section into document",
			content:      "# Doc\n\n## TOC\n\nold toc\n\n## Next\n\n## Installation\n\n## Usage\n",
			wantContains: []string{"Installation", "Usage"},
		},
		{
			name:           "start-from-level flag is forwarded",
			content:        "# Doc\n\n## TOC\n\nold toc\n\n## Next\n\n## Section\n\n### Sub\n",
			startFromLevel: 2,
			wantContains:   []string{"Sub"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(resetGlobals)

			docPath := setupBefore(t, tt.content, "TOC", "#", "#")

			// internal() only reads start-from-level and start-from-item from
			// the context; build a minimal flagset that carries those int values.
			set := flag.NewFlagSet("test", flag.ContinueOnError)
			set.Int("start-from-level", tt.startFromLevel, "")
			set.Int("start-from-item", tt.startFromItem, "")
			ctx := cli.NewContext(cli.NewApp(), set, nil)

			require.ErrorIs(t, internal(ctx), tt.wantErr)

			updated, err := os.ReadFile(docPath) // #nosec G304
			require.NoError(t, err)

			got := string(updated)
			for _, sub := range tt.wantContains {
				assert.Contains(t, got, sub)
			}
		})
	}
}

func TestExternal(t *testing.T) {
	tests := []struct {
		name         string
		docContent   string
		extFiles     map[string]string
		summaryFlags map[string]string
		wantErr      error
		wantContains []string
	}{
		{
			name:       "writes external toc into document",
			docContent: "# Doc\n\n## TOC\n\nold toc\n\n## Next\n",
			extFiles: map[string]string{
				"alpha.md": "# Alpha\n\ncontent\n",
				"beta.md":  "# Beta\n\ncontent\n",
			},
			wantContains: []string{"Alpha", "Beta"},
		},
		{
			name:       "summary flags are forwarded to generateExternalTOC",
			docContent: "# Doc\n\n## TOC\n\nold toc\n\n## Next\n",
			extFiles: map[string]string{
				"a.md": "# Alpha\n\n## Desc\n\nGreat tool.\n\n## End\n",
			},
			summaryFlags: map[string]string{
				"summary-header":        "Desc",
				"summary-limiter-left":  "##",
				"summary-limiter-right": "##",
			},
			wantContains: []string{"Great tool."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(resetGlobals)

			docPath := setupBefore(t, tt.docContent, "TOC", "#", "#")

			var extPaths []string
			for name, content := range tt.extFiles {
				extPaths = append(extPaths, writeTempFileNamed(t, name, content))
			}

			strFlags := map[string]string{"document": docPath}
			maps.Copy(strFlags, tt.summaryFlags)

			ctx := newContext(t, strFlags, map[string][]string{"path": extPaths})

			require.ErrorIs(t, external(ctx), tt.wantErr)

			updated, err := os.ReadFile(docPath) // #nosec G304
			require.NoError(t, err)

			got := string(updated)
			for _, sub := range tt.wantContains {
				assert.Contains(t, got, sub)
			}
		})
	}
}

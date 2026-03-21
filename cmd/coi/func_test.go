package coi

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// newContext builds a *cli.Context from a plain map of flag-name → value.
func newContext(t *testing.T, flags map[string]string) *cli.Context {
	t.Helper()

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "document"},
		&cli.StringFlag{Name: "header"},
		&cli.StringFlag{Name: "limiter-left"},
		&cli.StringFlag{Name: "limiter-right"},
		&cli.StringFlag{Name: "shell-name"},
		&cli.StringFlag{Name: "shell-prompt"},
		&cli.StringFlag{Name: "command"},
		&cli.StringFlag{Name: "output"},
	}

	set := flag.NewFlagSet("test", flag.ContinueOnError)

	for _, f := range app.Flags {
		sf, ok := f.(*cli.StringFlag)
		require.True(t, ok)

		set.String(sf.Name, sf.Value, sf.Usage)
	}

	for k, v := range flags {
		require.NoError(t, set.Set(k, v), "flag.Set(%q, %q)", k, v)
	}

	return cli.NewContext(app, set, nil)
}

// writeTempFile creates a temporary file with the given content and returns
// its path. The file is removed automatically when the test ends.
func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "doc-*.md")
	require.NoError(t, err, "CreateTemp")

	_, err = f.WriteString(content)
	require.NoError(t, err, "WriteString")

	_ = f.Close()

	return f.Name()
}

// resetGlobals resets all package-level variables that before/action mutate.
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

func TestBefore(t *testing.T) {
	tests := []struct {
		name    string
		content string
		flags   map[string]string

		// error expectation
		wantErr error

		// optional state assertions (only evaluated on success)
		wantLimiterR string
		wantPrefix   bool
		checkGlobals bool
	}{
		{
			name: "missing file returns",
			// content ignored — document path is overridden below
			flags: map[string]string{
				"header":        "Usage",
				"limiter-left":  "#",
				"limiter-right": "#",
			},
			wantErr: errFileRead,
		},
		{
			name:    "header absent from document",
			content: "# Installation\n\nsome content\n\n# Next\n",
			flags: map[string]string{
				"header":        "Usage",
				"limiter-left":  "#",
				"limiter-right": "#",
			},
			wantErr: errDocumentSectionExtract,
		},
		{
			name:    "wrong header name",
			content: "# Usage\n\nbody\n",
			flags: map[string]string{
				"header":        "Installation",
				"limiter-left":  "#",
				"limiter-right": "",
			},
			wantErr: errDocumentSectionExtract,
		},
		{
			name:    "empty document",
			content: "",
			flags: map[string]string{
				"header":        "Usage",
				"limiter-left":  "#",
				"limiter-right": "",
			},
			wantErr: errDocumentSectionExtract,
		},
		{
			name:    "omitted right limiter defaults to \\z",
			content: "# Usage\n\nsome content here\n",
			flags: map[string]string{
				"header":       "Usage",
				"limiter-left": "#",
				// limiter-right intentionally omitted
			},
			wantLimiterR: `\z`,
		},
		{
			name:    "section found between two headers sets prefix",
			content: "# Usage\n\nsome content\n\n# Next\n",
			flags: map[string]string{
				"header":        "Usage",
				"limiter-left":  "#",
				"limiter-right": "#",
			},
			wantPrefix: true,
		},
		{
			name:    "multi-level header is matched",
			content: "## Usage\n\nsome content\n\n## Next\n",
			flags: map[string]string{
				"header":        "Usage",
				"limiter-left":  "#",
				"limiter-right": "#",
			},
		},
		{
			name:    "single section ending at EOF is matched",
			content: "# Usage\n\nsome content\n",
			flags: map[string]string{
				"header":        "Usage",
				"limiter-left":  "#",
				"limiter-right": "",
			},
		},
		{
			name:    "all globals are populated on success",
			content: "# Usage\n\nsome content\n\n# Next\n",
			flags: map[string]string{
				"header":        "Usage",
				"limiter-left":  "#",
				"limiter-right": "#",
			},
			checkGlobals: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(resetGlobals)

			flags := tt.flags

			if errors.Is(tt.wantErr, errFileRead) {
				flags["document"] = filepath.Join(t.TempDir(), "nonexistent.md")
			} else {
				flags["document"] = writeTempFile(t, tt.content)
			}

			require.ErrorIs(t, before(newContext(t, flags)), tt.wantErr)

			if tt.wantLimiterR != "" {
				assert.Equal(t, tt.wantLimiterR, limiterR)
			}

			if tt.wantPrefix {
				assert.NotEmpty(t, prefix, "expected prefix to be populated")
			}

			if tt.checkGlobals {
				assert.NotEmpty(t, document)
				assert.NotEmpty(t, header)
				assert.NotNil(t, rexp)
				assert.NotNil(t, body)
			}
		})
	}
}

// setupAction calls before() against the given document content, registers
// global cleanup, and returns the document path plus a *cli.Context ready
// for action().
func setupAction(
	t *testing.T,
	content string,
	actionFlags map[string]string,
) (string, *cli.Context) {
	t.Helper()
	t.Cleanup(resetGlobals)

	docPath := writeTempFile(t, content)

	require.NoError(t,
		before(newContext(t, map[string]string{
			"document":      docPath,
			"header":        "Usage",
			"limiter-left":  "#",
			"limiter-right": "#",
		})),
		"before() setup failed",
	)

	actionFlags["document"] = docPath

	return docPath, newContext(t, actionFlags)
}

func TestAction(t *testing.T) {
	const baseContent = "# Usage\n\n``` shell\n$ old command\nold output\n```\n\n# Next\n"

	tests := []struct {
		name    string
		content string
		flags   map[string]string

		// error expectation
		wantErr error

		// optional document assertions (only evaluated on success)
		wantChanged  bool     // document must differ from its original content
		wantContains []string // substrings that must appear in the updated document
	}{
		{
			name:    "both --command and --output",
			content: baseContent,
			flags: map[string]string{
				"shell-name":   "shell",
				"shell-prompt": "$",
				"command":      "echo hi",
				"output":       "some output",
			},
			wantErr: nil,
		},
		{
			name:    "neither --command nor --output",
			content: baseContent,
			flags: map[string]string{
				"shell-name":   "shell",
				"shell-prompt": "$",
			},
			wantErr: errExclusiveFlags,
		},
		{
			name:    "non-existent command",
			content: baseContent,
			flags: map[string]string{
				"shell-name":   "shell",
				"shell-prompt": "$",
				"command":      "this-binary-does-not-exist-xyz",
			},
			wantErr: errCommandExecute,
		},
		{
			name:    "--output flag rewrites the matched section",
			content: baseContent,
			flags: map[string]string{
				"shell-name":   "shell",
				"shell-prompt": "$",
				"output":       "new output\n",
			},
			wantChanged: true,
		},
		{
			name:    "--command stdout is written into the document",
			content: baseContent,
			flags: map[string]string{
				"shell-name":   "shell",
				"shell-prompt": "$",
				"command":      "echo hello",
			},
			wantChanged:  true,
			wantContains: []string{"hello"},
		},
		{
			name:    "shell-name appears in the rewritten section",
			content: baseContent,
			flags: map[string]string{
				"shell-name":   "zsh",
				"shell-prompt": "%",
				"output":       "result\n",
			},
			wantContains: []string{"zsh"},
		},
		{
			name:    "content outside the matched section is preserved",
			content: "preamble line\n\n# Usage\n\n``` shell\n$ cmd\nout\n```\n\n# Next\n\nfooter line\n",
			flags: map[string]string{
				"shell-name":   "shell",
				"shell-prompt": "$",
				"output":       "new output\n",
			},
			wantContains: []string{"preamble line", "footer line"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docPath, ctx := setupAction(t, tt.content, tt.flags)

			require.ErrorIs(t, action(ctx), tt.wantErr)

			// #nosec G304
			updated, err := os.ReadFile(docPath)
			require.NoError(t, err)

			got := string(updated)

			if tt.wantChanged {
				assert.NotEqual(t, tt.content, got, "expected document content to change")
			}

			for _, sub := range tt.wantContains {
				assert.Contains(t, got, sub)
			}
		})
	}
}

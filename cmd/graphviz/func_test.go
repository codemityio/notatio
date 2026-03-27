package graphviz

import (
	"bytes"
	"flag"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// installFakeDot writes a stub "dot" binary into a temp dir and prepends that
// dir to PATH for the duration of the test. The stub parses "-o <out>" and
// touches the output file so every code path that reaches runDot succeeds
// without a real Graphviz installation.
func installFakeDot(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("fake-dot setup not implemented for Windows")
	}

	dir := t.TempDir()
	stub := filepath.Join(dir, "dot")

	script := "#!/bin/sh\n" +
		"while [ $# -gt 0 ]; do\n" +
		"  case \"$1\" in\n" +
		"    -o) shift; mkdir -p \"$(dirname \"$1\")\"; touch \"$1\";;\n" +
		"  esac\n" +
		"  shift\n" +
		"done\n"

	// #nosec G306
	require.NoError(t, os.WriteFile(stub, []byte(script), 0o700))
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// installFailingDot writes a stub that always exits non-zero, placed ahead of
// any existing dot on PATH.
func installFailingDot(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("fake-dot setup not implemented for Windows")
	}

	dir := t.TempDir()
	stub := filepath.Join(dir, "dot")

	require.NoError(t, os.WriteFile(stub, []byte("#!/bin/sh\nexit 1\n"), 0o700)) // #nosec G302 G306
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func newContext(
	t *testing.T,
	buf *bytes.Buffer,
	strFlags map[string]string,
	boolFlags map[string]bool,
) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.String("input-path", "", "")
	set.String("output-format", "", "")
	set.Bool("recursive", false, "")

	for k, v := range strFlags {
		require.NoError(t, set.Set(k, v), "flag.Set(%q)", k)
	}

	for k, v := range boolFlags {
		val := "false"
		if v {
			val = "true"
		}

		require.NoError(t, set.Set(k, val), "flag.Set(%q)", k)
	}

	cliApp := cli.NewApp()
	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{Name: "input-path"},
		&cli.StringFlag{Name: "output-format"},
		&cli.BoolFlag{Name: "recursive"},
	}
	cliApp.Writer = buf

	return cli.NewContext(cliApp, set, nil)
}

func writeDotFile(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte("digraph G { A -> B; }\n"), 0o600))

	return path
}

func TestAction(t *testing.T) {
	tests := []struct {
		name         string
		strFlags     map[string]string
		boolFlags    map[string]bool
		setup        func(t *testing.T) string // returns value for input-path; overrides strFlags["input-path"]
		setupDot     func(t *testing.T)        // installs a dot stub; defaults to installFakeDot
		wantErr      error
		wantContains []string // substrings expected in writer output
	}{
		{
			name: "empty input-path",
			strFlags: map[string]string{
				"input-path":    "",
				"output-format": "png",
			},
			wantErr: errInputPathEmpty,
		},
		{
			name: "unsupported output format",
			strFlags: map[string]string{
				"input-path":    "somefile.dot",
				"output-format": "pdf",
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "format check is case-sensitive — PNG is rejected",
			strFlags: map[string]string{
				"input-path":    "somefile.dot",
				"output-format": "PNG",
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "format check is case-sensitive — SVG is rejected",
			strFlags: map[string]string{
				"input-path":    "somefile.dot",
				"output-format": "SVG",
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "non-existent input path",
			strFlags: map[string]string{
				"output-format": "png",
			},
			setup: func(t *testing.T) string {
				t.Helper()

				return filepath.Join(t.TempDir(), "nonexistent.dot")
			},
			wantErr: errInputPath,
		},
		{
			name: "single .dot file is processed without error",
			strFlags: map[string]string{
				"output-format": "png",
			},
			setup: func(t *testing.T) string {
				t.Helper()

				return writeDotFile(t, t.TempDir(), "graph.dot")
			},
		},
		{
			name: "single .gv file is processed without error",
			strFlags: map[string]string{
				"output-format": "svg",
			},
			setup: func(t *testing.T) string {
				t.Helper()

				return writeDotFile(t, t.TempDir(), "graph.gv")
			},
		},
		{
			name: "directory input prints scanning message",
			strFlags: map[string]string{
				"output-format": "png",
			},
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeDotFile(t, dir, "a.dot")

				return dir
			},
			wantContains: []string{"scanning..."},
		},
		{
			name: "directory input logs all .dot and .gv files",
			strFlags: map[string]string{
				"output-format": "png",
			},
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeDotFile(t, dir, "a.dot")
				writeDotFile(t, dir, "b.gv")

				return dir
			},
			wantContains: []string{"a.dot", "b.gv"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installFakeDot(t)

			buf := &bytes.Buffer{}

			flags := make(map[string]string, len(tt.strFlags))
			maps.Copy(flags, tt.strFlags)

			if tt.setup != nil {
				flags["input-path"] = tt.setup(t)
			}

			require.ErrorIs(t, action(newContext(t, buf, flags, tt.boolFlags)), tt.wantErr)

			log := buf.String()
			for _, sub := range tt.wantContains {
				assert.Contains(t, log, sub)
			}
		})
	}
}

func TestIterate(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T) string // returns root directory
		format       string
		recursive    bool
		wantErr      error
		wantContains []string
		wantAbsent   []string
	}{
		{
			name: "processes .dot files in top-level directory",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeDotFile(t, dir, "first.dot")
				writeDotFile(t, dir, "second.dot")

				return dir
			},
			format:       "png",
			wantContains: []string{"first.dot", "second.dot"},
		},
		{
			name: "processes .gv files in top-level directory",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeDotFile(t, dir, "schema.gv")

				return dir
			},
			format:       "svg",
			wantContains: []string{"schema.gv"},
		},
		{
			name: "files with unrecognised extensions are skipped",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				require.NoError(
					t,
					os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hi"), 0o600),
				)
				writeDotFile(t, dir, "real.dot")

				return dir
			},
			format:       "png",
			wantContains: []string{"real.dot"},
			wantAbsent:   []string{"notes.txt"},
		},
		{
			name: "extension matching is case-insensitive",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeDotFile(t, dir, "upper.DOT")

				return dir
			},
			format:       "png",
			wantContains: []string{"upper.DOT"},
		},
		{
			name: "non-recursive does not descend into subdirectories",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				sub := filepath.Join(root, "sub")
				require.NoError(t, os.MkdirAll(sub, 0o700))
				writeDotFile(t, sub, "nested.dot")

				return root
			},
			format:     "png",
			recursive:  false,
			wantAbsent: []string{"nested.dot"},
		},
		{
			name: "recursive descends into subdirectories",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				sub := filepath.Join(root, "sub")
				require.NoError(t, os.MkdirAll(sub, 0o700))
				writeDotFile(t, sub, "nested.dot")

				return root
			},
			format:       "png",
			recursive:    true,
			wantContains: []string{"scanning...", "nested.dot"},
		},
		{
			name: "recursive descends multiple levels",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				deep := filepath.Join(root, "a", "b", "c")
				require.NoError(t, os.MkdirAll(deep, 0o700))
				writeDotFile(t, deep, "deep.dot")

				return root
			},
			format:       "svg",
			recursive:    true,
			wantContains: []string{"deep.dot"},
		},
		{
			name: "unreadable directory",
			setup: func(t *testing.T) string {
				t.Helper()

				if os.Getuid() == 0 {
					t.Skip("running as root — permission checks are ineffective")
				}

				dir := t.TempDir()
				locked := filepath.Join(dir, "locked")
				require.NoError(t, os.MkdirAll(locked, 0o000))
				t.Cleanup(func() { _ = os.Chmod(locked, 0o700) }) // #nosec G302

				return locked
			},
			format:  "png",
			wantErr: errReadDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installFakeDot(t)

			buf := &bytes.Buffer{}
			rootDir := tt.setup(t)

			cliApp := cli.NewApp()
			cliApp.Writer = buf
			ctx := cli.NewContext(cliApp, flag.NewFlagSet("test", flag.ContinueOnError), nil)

			require.ErrorIs(t, iterate(ctx, rootDir, tt.format, tt.recursive), tt.wantErr)

			log := buf.String()
			for _, sub := range tt.wantContains {
				assert.Contains(t, log, sub)
			}

			for _, sub := range tt.wantAbsent {
				assert.NotContains(t, log, sub)
			}
		})
	}
}

func TestRunDot(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string // returns input file path
		setupDot    func(t *testing.T)        // installs the dot stub; defaults to installFakeDot
		format      string
		wantErr     error
		wantLogFrag string
		checkOutput func(t *testing.T, inputPath string)
	}{
		{
			name: "successful run logs a generating message",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeDotFile(t, t.TempDir(), "graph.dot")
			},
			format:      "png",
			wantLogFrag: "generating",
		},
		{
			name: "log line contains the input file name",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeDotFile(t, t.TempDir(), "mydiagram.dot")
			},
			format:      "png",
			wantLogFrag: "mydiagram",
		},
		{
			name: "output file has the correct format extension",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeDotFile(t, t.TempDir(), "check.dot")
			},
			format: "svg",
			checkOutput: func(t *testing.T, inputPath string) {
				t.Helper()

				rel, err := filepath.Rel(".", inputPath)
				if err != nil {
					rel = inputPath
				}

				expected := strings.TrimSuffix(rel, filepath.Ext(rel)) + ".svg"
				assert.FileExists(t, expected)
			},
		},
		{
			name: "dot command failure",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeDotFile(t, t.TempDir(), "fail.dot")
			},
			setupDot: installFailingDot,
			format:   "png",
			wantErr:  errCommandRun,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupDot != nil {
				tt.setupDot(t)
			} else {
				installFakeDot(t)
			}

			buf := &bytes.Buffer{}
			inputPath := tt.setup(t)

			cliApp := cli.NewApp()
			cliApp.Writer = buf
			ctx := cli.NewContext(cliApp, flag.NewFlagSet("test", flag.ContinueOnError), nil)

			require.ErrorIs(t, runDot(ctx, inputPath, tt.format), tt.wantErr)

			if tt.wantLogFrag != "" {
				assert.Contains(t, buf.String(), tt.wantLogFrag)
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, inputPath)
			}
		})
	}
}

func TestRoundFloatsInValue(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		precision int
		want      string
	}{
		{
			name:      "rounds single float to integer",
			input:     "100.25",
			precision: 0,
			want:      "100",
		},
		{
			name:      "rounds multiple floats in path d attribute",
			input:     "M29435.88,-1704C29578.98,-1704",
			precision: 0,
			want:      "M29436,-1704C29579,-1704",
		},
		{
			name:      "preserves integers unchanged",
			input:     "100,-200 300,-400",
			precision: 0,
			want:      "100,-200 300,-400",
		},
		{
			name:      "rounds negative floats correctly",
			input:     "-853.5,-1704.2",
			precision: 0,
			want:      "-854,-1704",
		},
		{
			name:      "respects precision of 1",
			input:     "29435.88,-1704.22",
			precision: 1,
			want:      "29435.9,-1704.2",
		},
		{
			name:      "empty string returns empty string",
			input:     "",
			precision: 0,
			want:      "",
		},
		{
			name:      "string with no floats is unchanged",
			input:     "M100,-200 L300,-400",
			precision: 0,
			want:      "M100,-200 L300,-400",
		},
		{
			name:      "rounds polygon points",
			input:     "42241.43,-855.1 42247.43,-853 42241.43,-850.9",
			precision: 0,
			want:      "42241,-855 42247,-853 42241,-851",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, roundFloatsInValue(tt.input, tt.precision))
		})
	}
}

func writeSVGFile(t *testing.T, dir, name, content string) string { //nolint: unparam
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}

func TestNormalizeSVG(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string // returns path to SVG file
		precision   int
		wantErr     error
		wantContent string // substring expected in normalised output
		wantAbsent  string // substring that must NOT appear in output
	}{
		{
			name: "rounds floats in d attribute",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeSVGFile(t, t.TempDir(), "graph.svg",
					`<path fill="none" stroke="#a8a8a8" d="M29435.88,-1704C29578.98,-1704"/>`,
				)
			},
			precision:   0,
			wantContent: `d="M29436,-1704C29579,-1704"`,
			wantAbsent:  "29435.88",
		},
		{
			name: "rounds floats in points attribute",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeSVGFile(
					t,
					t.TempDir(),
					"graph.svg",
					`<polygon fill="#a8a8a8" stroke="#a8a8a8" points="42241.43,-855.1 42247.43,-853 42241.43,-850.9"/>`,
				)
			},
			precision:   0,
			wantContent: `points="42241,-855 42247,-853 42241,-851"`,
			wantAbsent:  "42241.43",
		},
		{
			name: "does not modify non-coordinate attributes",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeSVGFile(
					t,
					t.TempDir(),
					"graph.svg",
					`<svg version="1.1" width="2519pt" font-size="14.00"><path d="M100.5,-200.5"/></svg>`,
				)
			},
			precision:   0,
			wantContent: `version="1.1"`,
			wantAbsent:  `version="1"`,
		},
		{
			name: "does not modify font-size",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeSVGFile(t, t.TempDir(), "graph.svg",
					`<text font-size="14.00" x="100.5" y="-107.33">label</text>`,
				)
			},
			precision:   0,
			wantContent: `font-size="14.00"`,
		},
		{
			name: "rounds floats in x and y attributes",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeSVGFile(t, t.TempDir(), "graph.svg",
					`<text x="1022.12" y="-107.33">label</text>`,
				)
			},
			precision:   0,
			wantContent: `x="1022" y="-107"`,
			wantAbsent:  "1022.12",
		},
		{
			name: "preserves file content outside coordinate attributes unchanged",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeSVGFile(t, t.TempDir(), "graph.svg",
					`<!-- Generated by graphviz --><path d="M100.5,-200.5"/>`,
				)
			},
			precision:   0,
			wantContent: "<!-- Generated by graphviz -->",
		},
		{
			name: "respects precision of 1",
			setup: func(t *testing.T) string {
				t.Helper()

				return writeSVGFile(t, t.TempDir(), "graph.svg",
					`<path d="M29435.88,-1704.25"/>`,
				)
			},
			precision:   1,
			wantContent: `d="M29435.9,-1704.3"`,
			wantAbsent:  "29435.88",
		},
		{
			name: "returns error for non-existent file",
			setup: func(t *testing.T) string {
				t.Helper()

				return filepath.Join(t.TempDir(), "nonexistent.svg")
			},
			precision: 0,
			wantErr:   os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)

			require.ErrorIs(t, normalizeSVG(path, tt.precision), tt.wantErr)

			if tt.wantErr != nil {
				return
			}

			content, e := os.ReadFile(path) // #nosec G304
			require.NoError(t, e)

			assert.Contains(t, string(content), tt.wantContent)

			if tt.wantAbsent != "" {
				assert.NotContains(t, string(content), tt.wantAbsent)
			}
		})
	}
}

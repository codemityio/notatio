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

			assert.ErrorIs(t, action(newContext(t, buf, flags, tt.boolFlags)), tt.wantErr)

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

			assert.ErrorIs(t, iterate(ctx, rootDir, tt.format, tt.recursive), tt.wantErr)

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

			assert.ErrorIs(t, runDot(ctx, inputPath, tt.format), tt.wantErr)

			if tt.wantLogFrag != "" {
				assert.Contains(t, buf.String(), tt.wantLogFrag)
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, inputPath)
			}
		})
	}
}

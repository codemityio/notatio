package mermaid

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func installFakeBinaries(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("fake binary setup not implemented for Windows")
	}

	dir := t.TempDir()
	for _, name := range []string{"mmdc", "chromium-browser"} {
		require.NoError(t, os.WriteFile( // #nosec G306
			filepath.Join(dir, name),
			[]byte("#!/bin/sh\nexit 0\n"),
			0o700,
		))
	}

	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func installFailingMmdc(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("fake binary setup not implemented for Windows")
	}

	dir := t.TempDir()
	require.NoError(t, os.WriteFile( // #nosec G306
		filepath.Join(dir, "mmdc"),
		[]byte("#!/bin/sh\nexit 1\n"),
		0o700,
	))
	require.NoError(t, os.WriteFile( // #nosec G306
		filepath.Join(dir, "chromium-browser"),
		[]byte("#!/bin/sh\nexit 0\n"),
		0o700,
	))
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
	set.String("puppeteer-config-json-path", "", "")
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
		&cli.StringFlag{Name: "puppeteer-config-json-path"},
		&cli.BoolFlag{Name: "recursive"},
	}
	cliApp.Writer = buf

	return cli.NewContext(cliApp, set, nil)
}

func writeMmdFile(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte("graph TD\n  A --> B\n"), 0o600))

	return path
}

func writePuppeteerConfig(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "puppeteer-config.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"args":["--no-sandbox"]}`), 0o600))

	return path
}

func TestAction(t *testing.T) {
	tests := []struct {
		name         string
		setupBin     func(t *testing.T)
		flags        func(t *testing.T) map[string]string
		boolFlags    map[string]bool
		wantErr      error
		wantContains []string
	}{
		{
			name: "empty input-path",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":                 "",
					"output-format":              "png",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
			wantErr: errInputPathEmpty,
		},
		{
			name: "unsupported output format",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":                 "somefile.mmd",
					"output-format":              "pdf",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "output format case-sensitive — PNG rejected",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":                 "somefile.mmd",
					"output-format":              "PNG",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "output format case-sensitive — SVG rejected",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":                 "somefile.mmd",
					"output-format":              "SVG",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "non-existent input path",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":                 filepath.Join(t.TempDir(), "nonexistent.mmd"),
					"output-format":              "png",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
			wantErr: errInputPath,
		},
		{
			name: "missing puppeteer config",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":    "somefile.mmd",
					"output-format": "png",
					"puppeteer-config-json-path": filepath.Join(
						"nonexistent",
						"puppeteer-config.json",
					),
				}
			},
			wantErr: errWrite,
		},
		{
			name: "single .mmd file is processed",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				dir := t.TempDir()
				path := writeMmdFile(t, dir, "diagram.mmd")

				return map[string]string{
					"input-path":                 path,
					"output-format":              "png",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
			wantContains: []string{"generating"},
		},
		{
			name: "non-.mmd file is skipped",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				dir := t.TempDir()
				path := filepath.Join(dir, "diagram.txt")
				require.NoError(t, os.WriteFile(path, []byte("content"), 0o600))

				return map[string]string{
					"input-path":                 path,
					"output-format":              "png",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
		},
		{
			name: "directory input iterates files",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				dir := t.TempDir()
				writeMmdFile(t, dir, "a.mmd")

				return map[string]string{
					"input-path":                 dir,
					"output-format":              "png",
					"puppeteer-config-json-path": writePuppeteerConfig(t),
				}
			},
			wantContains: []string{"generating"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupBin != nil {
				tt.setupBin(t)
			} else {
				installFakeBinaries(t)
			}

			buf := &bytes.Buffer{}

			assert.ErrorIs(t, action(newContext(t, buf, tt.flags(t), tt.boolFlags)), tt.wantErr)

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
		setup        func(t *testing.T) string
		format       string
		recursive    bool
		setupBin     func(t *testing.T)
		wantErr      error
		wantContains []string
		wantAbsent   []string
	}{
		{
			name: "processes .mmd files in top-level directory",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeMmdFile(t, dir, "a.mmd")
				writeMmdFile(t, dir, "b.mmd")

				return dir
			},
			format:       "png",
			wantContains: []string{"generating"},
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

				return dir
			},
			format:     "png",
			wantAbsent: []string{"generating"},
		},
		{
			name: "non-recursive does not descend into subdirectories",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				sub := filepath.Join(root, "sub")
				require.NoError(t, os.MkdirAll(sub, 0o700))
				writeMmdFile(t, sub, "nested.mmd")

				return root
			},
			format:     "png",
			recursive:  false,
			wantAbsent: []string{"scanning..."},
		},
		{
			name: "recursive descends into subdirectories",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				sub := filepath.Join(root, "sub")
				require.NoError(t, os.MkdirAll(sub, 0o700))
				writeMmdFile(t, sub, "nested.mmd")

				return root
			},
			format:       "png",
			recursive:    true,
			wantContains: []string{"scanning...", "generating"},
		},
		{
			name: "recursive descends multiple levels",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				deep := filepath.Join(root, "a", "b", "c")
				require.NoError(t, os.MkdirAll(deep, 0o700))
				writeMmdFile(t, deep, "deep.mmd")

				return root
			},
			format:       "svg",
			recursive:    true,
			wantContains: []string{"scanning..."},
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
		{
			name: "mmdc failure propagates error",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeMmdFile(t, dir, "a.mmd")

				return dir
			},
			setupBin: installFailingMmdc,
			format:   "png",
			wantErr:  errCommandRun,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupBin != nil {
				tt.setupBin(t)
			} else {
				installFakeBinaries(t)
			}

			buf := &bytes.Buffer{}
			rootDir := tt.setup(t)
			puppeteerConfig := writePuppeteerConfig(t)

			cliApp := cli.NewApp()
			cliApp.Writer = buf
			ctx := cli.NewContext(cliApp, flag.NewFlagSet("test", flag.ContinueOnError), nil)

			assert.ErrorIs(
				t,
				iterate(ctx, rootDir, tt.format, puppeteerConfig, tt.recursive),
				tt.wantErr,
			)

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

func TestGenerate(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T) string
		setupBin     func(t *testing.T)
		format       string
		wantErr      error
		wantContains []string
	}{
		{
			name: "successful png generation",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeMmdFile(t, dir, "diagram.mmd")

				return filepath.Join(dir, "diagram")
			},
			format:       "png",
			wantContains: []string{"generating"},
		},
		{
			name: "successful svg generation",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeMmdFile(t, dir, "diagram.mmd")

				return filepath.Join(dir, "diagram")
			},
			format:       "svg",
			wantContains: []string{"generating"},
		},
		{
			name: "log line contains the output path",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeMmdFile(t, dir, "mydiagram.mmd")

				return filepath.Join(dir, "mydiagram")
			},
			format:       "png",
			wantContains: []string{"mydiagram"},
		},
		{
			name: "mmdc failure",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writeMmdFile(t, dir, "diagram.mmd")

				return filepath.Join(dir, "diagram")
			},
			setupBin: installFailingMmdc,
			format:   "png",
			wantErr:  errCommandRun,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupBin != nil {
				tt.setupBin(t)
			} else {
				installFakeBinaries(t)
			}

			buf := &bytes.Buffer{}
			path := tt.setup(t)
			puppeteerConfig := writePuppeteerConfig(t)

			cliApp := cli.NewApp()
			cliApp.Writer = buf
			ctx := cli.NewContext(cliApp, flag.NewFlagSet("test", flag.ContinueOnError), nil)

			assert.ErrorIs(t, generate(ctx, path, tt.format, puppeteerConfig), tt.wantErr)

			log := buf.String()
			for _, sub := range tt.wantContains {
				assert.Contains(t, log, sub)
			}
		})
	}
}

package plantuml

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

func installFakeJava(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("fake-java setup not implemented for Windows")
	}

	dir := t.TempDir()
	stub := filepath.Join(dir, "java")
	require.NoError(t, os.WriteFile(stub, []byte("#!/bin/sh\nexit 0\n"), 0o700)) // #nosec G306
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func installFailingJava(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("fake-java setup not implemented for Windows")
	}

	dir := t.TempDir()
	stub := filepath.Join(dir, "java")
	require.NoError(t, os.WriteFile(stub, []byte("#!/bin/sh\nexit 1\n"), 0o700)) // #nosec G306
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
	set.String("plantuml-jar-path", "", "")
	set.String("plantuml-limit-size", "", "")
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
		&cli.StringFlag{Name: "plantuml-jar-path"},
		&cli.StringFlag{Name: "plantuml-limit-size"},
		&cli.BoolFlag{Name: "recursive"},
	}
	cliApp.Writer = buf

	return cli.NewContext(cliApp, set, nil)
}

func writePumlFile(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte("@startuml\nA -> B\n@enduml\n"), 0o600))

	return path
}

func writeFakeJar(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "plantuml.jar")
	require.NoError(t, os.WriteFile(path, []byte(""), 0o600))

	return path
}

func TestAction(t *testing.T) {
	tests := []struct {
		name         string
		setupJava    func(t *testing.T)
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
					"input-path":        "",
					"output-format":     "png",
					"plantuml-jar-path": writeFakeJar(t),
				}
			},
			wantErr: errInputPathEmpty,
		},
		{
			name: "unsupported output format",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":        "somefile.puml",
					"output-format":     "pdf",
					"plantuml-jar-path": writeFakeJar(t),
				}
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "output format case-sensitive — PNG rejected",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":        "somefile.puml",
					"output-format":     "PNG",
					"plantuml-jar-path": writeFakeJar(t),
				}
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "output format case-sensitive — SVG rejected",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":        "somefile.puml",
					"output-format":     "SVG",
					"plantuml-jar-path": writeFakeJar(t),
				}
			},
			wantErr: errUnsupportedOutputFormat,
		},
		{
			name: "non-existent input path",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":        filepath.Join(t.TempDir(), "nonexistent.puml"),
					"output-format":     "png",
					"plantuml-jar-path": writeFakeJar(t),
				}
			},
			wantErr: errInputPath,
		},
		{
			name: "missing plantuml jar",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				return map[string]string{
					"input-path":        "somefile.puml",
					"output-format":     "png",
					"plantuml-jar-path": filepath.Join("nonexistent", "plantuml.jar"),
				}
			},
			wantErr: errWrite,
		},
		{
			name: "directory input prints scanning message",
			flags: func(t *testing.T) map[string]string {
				t.Helper()

				dir := t.TempDir()
				writePumlFile(t, dir, "a.puml")

				return map[string]string{
					"input-path":        dir,
					"output-format":     "png",
					"plantuml-jar-path": writeFakeJar(t),
				}
			},
			wantContains: []string{"scanning..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupJava != nil {
				tt.setupJava(t)
			} else {
				installFakeJava(t)
			}

			buf := &bytes.Buffer{}

			require.ErrorIs(t, action(newContext(t, buf, tt.flags(t), tt.boolFlags)), tt.wantErr)

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
		setupJava    func(t *testing.T)
		wantErr      error
		wantContains []string
		wantAbsent   []string
	}{
		{
			name: "processes .puml files in top-level directory",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writePumlFile(t, dir, "a.puml")
				writePumlFile(t, dir, "b.puml")

				return dir
			},
			format: "png",
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
			wantAbsent: []string{"notes.txt"},
		},
		{
			name: "non-recursive does not descend into subdirectories",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				sub := filepath.Join(root, "sub")
				require.NoError(t, os.MkdirAll(sub, 0o700))
				writePumlFile(t, sub, "nested.puml")

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
				writePumlFile(t, sub, "nested.puml")

				return root
			},
			format:       "png",
			recursive:    true,
			wantContains: []string{"scanning..."},
		},
		{
			name: "recursive descends multiple levels",
			setup: func(t *testing.T) string {
				t.Helper()

				root := t.TempDir()
				deep := filepath.Join(root, "a", "b", "c")
				require.NoError(t, os.MkdirAll(deep, 0o700))
				writePumlFile(t, deep, "deep.puml")

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
			name: "java failure propagates error",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				writePumlFile(t, dir, "a.puml")

				return dir
			},
			setupJava: installFailingJava,
			format:    "png",
			wantErr:   errCommandRun,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupJava != nil {
				tt.setupJava(t)
			} else {
				installFakeJava(t)
			}

			buf := &bytes.Buffer{}
			rootDir := tt.setup(t)
			jarPath := writeFakeJar(t)

			cliApp := cli.NewApp()
			cliApp.Writer = buf
			ctx := cli.NewContext(cliApp, flag.NewFlagSet("test", flag.ContinueOnError), nil)

			require.ErrorIs(
				t,
				iterate(ctx, rootDir, tt.format, jarPath, "", tt.recursive),
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
		name      string
		setup     func(t *testing.T) string
		setupJava func(t *testing.T)
		format    string
		wantErr   error
	}{
		{
			name: "successful png generation",
			setup: func(t *testing.T) string {
				t.Helper()

				return writePumlFile(t, t.TempDir(), "diagram.puml")
			},
			format: "png",
		},
		{
			name: "successful svg generation",
			setup: func(t *testing.T) string {
				t.Helper()

				return writePumlFile(t, t.TempDir(), "diagram.puml")
			},
			format: "svg",
		},
		{
			name: "java failure",
			setup: func(t *testing.T) string {
				t.Helper()

				return writePumlFile(t, t.TempDir(), "diagram.puml")
			},
			setupJava: installFailingJava,
			format:    "png",
			wantErr:   errCommandRun,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupJava != nil {
				tt.setupJava(t)
			} else {
				installFakeJava(t)
			}

			buf := &bytes.Buffer{}
			inputPath := tt.setup(t)
			jarPath := writeFakeJar(t)

			cliApp := cli.NewApp()
			cliApp.Writer = buf
			ctx := cli.NewContext(cliApp, flag.NewFlagSet("test", flag.ContinueOnError), nil)

			require.ErrorIs(t, generate(ctx, inputPath, tt.format, jarPath, "8192"), tt.wantErr)
		})
	}
}

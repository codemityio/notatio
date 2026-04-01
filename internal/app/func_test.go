package app

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestWithValues(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		inputDesc   string
		inputVer    string
		inputCopy   string
		authorName  string
		authorEmail string
		buildTime   string
		expectPanic bool
	}{
		{
			name:        "All valid values with build time",
			inputName:   "MyApp",
			inputDesc:   "Test app",
			inputVer:    "1.0",
			inputCopy:   "2026",
			authorName:  "John",
			authorEmail: "john@example.com",
			buildTime:   "2026-03-23T15:04:05Z",
			expectPanic: false,
		},
		{
			name:        "Empty build time",
			inputName:   "App2",
			inputDesc:   "Another app",
			inputVer:    "2.0",
			inputCopy:   "2026",
			authorName:  "Alice",
			authorEmail: "alice@example.com",
			buildTime:   "",
			expectPanic: false,
		},
		{
			name:        "Invalid build time triggers panic",
			inputName:   "App3",
			inputDesc:   "Bad app",
			inputVer:    "3.0",
			inputCopy:   "2026",
			authorName:  "Bob",
			authorEmail: "bob@example.com",
			buildTime:   "invalid-time",
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appInstance := &cli.App{}

			if tt.expectPanic {
				assert.Panics(t, func() {
					opt := WithValues(
						tt.inputName,
						tt.inputDesc,
						tt.inputVer,
						tt.inputCopy,
						tt.authorName,
						tt.authorEmail,
						tt.buildTime,
					)
					opt(appInstance)
				})

				return
			}

			opt := WithValues(
				tt.inputName,
				tt.inputDesc,
				tt.inputVer,
				tt.inputCopy,
				tt.authorName,
				tt.authorEmail,
				tt.buildTime,
			)
			opt(appInstance)

			// Assertions
			assert.Equal(t, tt.inputName, appInstance.Name)
			assert.Equal(t, tt.inputDesc, appInstance.Description)
			assert.Equal(t, tt.inputVer, appInstance.Version)
			assert.Equal(t, tt.inputCopy, appInstance.Copyright)
			assert.Equal(
				t,
				[]*cli.Author{{Name: tt.authorName, Email: tt.authorEmail}},
				appInstance.Authors,
			)
			assert.False(t, appInstance.HideVersion)

			if tt.buildTime != "" {
				parsed, _ := time.Parse(time.RFC3339, tt.buildTime)
				assert.Equal(t, parsed, appInstance.Compiled)
			} else {
				assert.True(t, appInstance.Compiled.IsZero())
			}
		})
	}
}

func TestCheckFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "file.txt")

	require.NoError(t, os.WriteFile(existingFile, []byte("content"), 0o644)) // #nosec G306

	nonExistingFile := filepath.Join(tmpDir, "missing.txt")

	tests := []struct {
		name    string
		path    string
		message string
		wantErr error
	}{
		{
			name:    "Existing file succeeds",
			path:    existingFile,
			message: "",
			wantErr: nil,
		},
		{
			name:    "Non-existing file fails",
			path:    nonExistingFile,
			message: "file missing",
			wantErr: errWrite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ctx := &cli.Context{
				App: &cli.App{
					Writer: buf,
				},
			}

			err := CheckFileExists(ctx, tt.path, tt.message)

			require.ErrorIs(t, err, tt.wantErr)

			assert.Contains(t, buf.String(), tt.message)
		})
	}
}

func TestCheckCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		message string
		wantErr error
	}{
		{
			name:    "Existing command succeeds",
			cmd:     "go",
			message: "",
			wantErr: nil,
		},
		{
			name:    "Non-existing command fails",
			cmd:     "noop",
			message: "command missing",
			wantErr: errWrite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ctx := &cli.Context{
				App: &cli.App{
					Writer: buf,
				},
			}

			err := CheckCommand(ctx, tt.cmd, tt.message)

			require.ErrorIs(t, err, tt.wantErr)

			assert.Contains(t, buf.String(), tt.message)
		})
	}
}

func TestResolveVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fallback string
		bi       *debug.BuildInfo
		want     string
	}{
		{
			name:     "fallback takes precedence",
			fallback: "v1.0.0",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: "v2.0.0"}},
			want:     "v1.0.0",
		},
		{
			name:     "nil build info returns latest",
			fallback: "",
			bi:       nil,
			want:     "latest",
		},
		{
			name:     "devel returns latest",
			fallback: "",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}},
			want:     "latest",
		},
		{
			name:     "empty version returns latest",
			fallback: "",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: ""}},
			want:     "latest",
		},
		{
			name:     "module version used when no fallback",
			fallback: "",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}},
			want:     "v1.2.3",
		},
		{
			name:     "pseudo version used when no fallback",
			fallback: "",
			bi: &debug.BuildInfo{
				Main: debug.Module{Version: "v0.0.4-0.20260401175942-30fd43debbd8"},
			},
			want: "v0.0.4-0.20260401175942-30fd43debbd8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, resolveVersion(tt.fallback, tt.bi))
		})
	}
}

func TestResolveBuildTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fallback string
		bi       *debug.BuildInfo
		want     string
	}{
		{
			name:     "nil build info returns fallback",
			fallback: "2026-04-01T18:00:00Z",
			bi:       nil,
			want:     "2026-04-01T18:00:00Z",
		},
		{
			name:     "vcs.time returned when present",
			fallback: "",
			bi: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.time", Value: "2026-04-01T12:00:00Z"},
				},
			},
			want: "2026-04-01T12:00:00Z",
		},
		{
			name:     "fallback returned when vcs.time absent",
			fallback: "2026-04-01T18:00:00Z",
			bi: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "abc123"},
				},
			},
			want: "2026-04-01T18:00:00Z",
		},
		{
			name:     "empty fallback and no vcs.time returns empty",
			fallback: "",
			bi:       &debug.BuildInfo{Settings: []debug.BuildSetting{}},
			want:     "",
		},
		{
			name:     "vcs.time takes precedence over fallback",
			fallback: "2026-04-01T18:00:00Z",
			bi: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.time", Value: "2026-04-01T12:00:00Z"},
				},
			},
			want: "2026-04-01T12:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, resolveBuildTime(tt.fallback, tt.bi))
		})
	}
}

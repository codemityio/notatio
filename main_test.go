package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const subprocessEnv = "TEST_SUBPROCESS_ARGS"

func TestMainCLI(t *testing.T) {
	if args := os.Getenv(subprocessEnv); args != "" {
		os.Args = strings.Fields(args)

		main()

		return
	}

	tests := []struct {
		name        string
		args        string
		wantErr     bool
		wantContain string
	}{
		{name: "no args", args: "notatio", wantErr: false},
		{name: "help flag", args: "notatio --help", wantErr: false},
		{name: "version flag", args: "notatio --version", wantErr: false},
		{
			name:        "unknown command",
			args:        "notatio nonexistent",
			wantErr:     true,
			wantContain: "unknown command",
		},
		{name: "coi help", args: "notatio coi --help", wantErr: false},
		{name: "graphviz help", args: "notatio graphviz --help", wantErr: false},
		{name: "mermaid help", args: "notatio mermaid --help", wantErr: false},
		{name: "plantuml help", args: "notatio plantuml --help", wantErr: false},
		{name: "toc help", args: "notatio toc --help", wantErr: false},
		{name: "tol help", args: "notatio tol --help", wantErr: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=TestMainCLI")

			cmd.Env = append(os.Environ(), subprocessEnv+"="+tc.args)

			out, err := cmd.CombinedOutput()

			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got none; output: %s", out)
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v; output: %s", err, out)
			}

			if tc.wantContain != "" && !strings.Contains(string(out), tc.wantContain) {
				t.Errorf("expected output to contain %q, got: %s", tc.wantContain, out)
			}
		})
	}
}

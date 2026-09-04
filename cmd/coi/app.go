//nolint:gochecknoglobals
package coi

import (
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{ //nolint:exhaustruct_v5
	Name:         "coi",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "Command output injector.",
	ArgsUsage:    "",
	Category:     "",
	BashComplete: nil,
	Before:       before,
	After:        nil,
	Action:       action,
	OnUsageError: nil,
	Flags: []cli.Flag{
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:     "document-path",
			Usage:    "markdown file path to be updated",
			Required: true,
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:     "header",
			Usage:    "header to use for document lookups and generation",
			Required: true,
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:  "limiter-left",
			Usage: "string to use as a lookup limiter",
			Value: "##",
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:  "limiter-right",
			Usage: "string to use as a lookup limiter - empty will use end of file as a limit",
			Value: "##",
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:  "shell-name",
			Usage: "shell name to use in the output",
			Value: "bash",
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:  "shell-prompt",
			Usage: "shell prompt prefix to use in the output",
			Value: "$",
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:  "command",
			Usage: "command to execute (command execution is skipped if --output is also provided)",
			Value: "",
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:  "output",
			Usage: "output to inject",
			Value: "",
		},
		&cli.IntFlag{ //nolint:exhaustruct_v5
			Name: "index",
			Usage: `index of a section to be used as a placeholder (useful if limiters refer to more than one section,
0 = replace all)`,
			Value: 0,
		},
	},
	Subcommands: []*cli.Command{},
}

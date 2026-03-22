//nolint:gochecknoglobals
package coi

import (
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{ //nolint:exhaustruct
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
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "document",
			Usage:    "markdown file path to be updated",
			Required: true,
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "header",
			Usage:    "header to use for document lookups and generation",
			Required: true,
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "limiter-left",
			Usage: "string to use as a lookup limiter",
			Value: "##",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "limiter-right",
			Usage: "string to use as a lookup limiter - empty will use end of file as a limit",
			Value: "##",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "shell-name",
			Usage: "shell name to use in the output",
			Value: "bash",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "shell-prompt",
			Usage: "shell prompt prefix to use in the output",
			Value: "$",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "command",
			Usage: "command to execute (use it only if the -output flag is not provided)",
			Value: "",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "output",
			Usage: "output to inject (use it only if the --command flag is not provided)",
			Value: "",
		},
	},
	Subcommands: []*cli.Command{},
}

package tol

import (
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{ //nolint:gochecknoglobals,exhaustruct
	Name:         "tol",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "Table of licences generator",
	ArgsUsage:    "",
	Category:     "",
	BashComplete: nil,
	Before:       before,
	After:        nil,
	Action:       action,
	OnUsageError: nil,
	Flags: []cli.Flag{
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "csv-path",
			Usage:    "input csv file (go-licenses output)",
			Required: true,
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "document-path",
			Usage:    "markdown document file path to be updated",
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
		&cli.IntFlag{ //nolint:exhaustruct
			Name: "index",
			Usage: `index of a section to be used as a placeholder (useful if limiters refer to more than one section,
0 = replace all)`,
			Value: 0,
		},
		&cli.StringSliceFlag{ //nolint:exhaustruct
			Name:  "skip",
			Usage: "packages to skip",
		},
	},
	Subcommands: []*cli.Command{},
}

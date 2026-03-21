//nolint:gochecknoglobals
package graphviz

import (
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{ //nolint:exhaustruct
	Name:         "graphviz",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "A tool to convert dot/gv files to svg/png images.",
	ArgsUsage:    "",
	Category:     "",
	BashComplete: nil,
	Before: func(c *cli.Context) error {
		return nil
	},
	After:        nil,
	Action:       action,
	OnUsageError: nil,
	Flags: []cli.Flag{
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "input-path",
			Usage:    "input path, either a file to be converted or a directory to be scanned",
			Value:    ".",
			Required: false,
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "output-format",
			Usage:    "output format (svg or png)",
			Value:    "svg",
			Required: false,
		},
		&cli.BoolFlag{ //nolint:exhaustruct
			Name:     "recursive",
			Usage:    "enable recursive directories scan",
			Value:    false,
			Required: false,
		},
	},
	Subcommands: []*cli.Command{},
}

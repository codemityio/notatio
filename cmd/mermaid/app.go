//nolint:gochecknoglobals
package mermaid

import (
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{ //nolint:exhaustruct_v5
	Name:         "mermaid",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "A tool to convert mmd files to svg/png images.",
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
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:     "input-path",
			Usage:    "input path, either a file to be converted or a directory to be scanned",
			Value:    ".",
			Required: false,
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:     "output-format",
			Usage:    "output format (svg or png)",
			Value:    "svg",
			Required: false,
		},
		&cli.StringFlag{ //nolint:exhaustruct_v5
			Name:     "puppeteer-config-json-path",
			Usage:    "puppeteer-config.json file path",
			Value:    "/usr/local/lib/puppeteer-config.json",
			Required: false,
		},
		&cli.BoolFlag{ //nolint:exhaustruct_v5
			Name:     "recursive",
			Usage:    "enable recursive directories scan",
			Value:    false,
			Required: false,
		},
	},
	Subcommands: []*cli.Command{},
}

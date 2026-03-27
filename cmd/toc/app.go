//nolint:gochecknoglobals
package toc

import (
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{ //nolint:exhaustruct
	Name:         "toc",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "Table of contents generator.",
	ArgsUsage:    "",
	Category:     "",
	BashComplete: nil,
	Before:       before,
	After:        nil,
	Action:       nil,
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
		&cli.IntFlag{ //nolint:exhaustruct
			Name: "index",
			Usage: `index of a section to be used as a placeholder (useful if limiters refer to more than one section,
0 = replace all)`,
			Value: 0,
		},
	},
	Subcommands: []*cli.Command{
		{
			Name: "int",
			Usage: `Generate table of content from headers within a document 
	e.g. toc --document=README.md --header="Table of contents" --limiter-right="##" int.`,
			Flags: []cli.Flag{
				&cli.IntFlag{ //nolint:exhaustruct
					Name:     "start-from-level",
					Usage:    "indicate what level of headers to start from",
					Required: false,
					Value:    0,
				},
				&cli.IntFlag{ //nolint:exhaustruct
					Name:     "start-from-item",
					Usage:    "indicate what item from the list to start from",
					Required: false,
					Value:    0,
				},
			},
			Action: internal,
		},
		{
			Name: "ext",
			Usage: `Generate table of content within a document and use provided paths as a list 
	e.g. toc --document=README.md --header="Table of contents" --limiter-right="##" ext --path=one/document.md --path=two/document.md.`,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{ //nolint:exhaustruct
					Name:  "path",
					Usage: "path to be included in the table of contents",
				},
				&cli.StringFlag{ //nolint:exhaustruct
					Name:     "summary-header",
					Usage:    "summary header to use for document lookups",
					Required: true,
				},
				&cli.StringFlag{ //nolint:exhaustruct
					Name:  "summary-limiter-left",
					Usage: "string to use as a summary lookup limiter",
					Value: "##",
				},
				&cli.StringFlag{ //nolint:exhaustruct
					Name:  "summary-limiter-right",
					Usage: "string to use as a summary lookup limiter - empty will use end of file as a limit",
					Value: "##",
				},
			},
			Action: external,
		},
	},
}

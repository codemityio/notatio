package app

import (
	"time"

	"github.com/urfave/cli/v2"
)

// WithValues set values.
func WithValues(
	name, description, version, copyright, authorName, authorEmail, buildTime string,
) Option {
	return func(app *cli.App) {
		app.Name = name
		app.Description = description
		app.Version = version
		app.Copyright = copyright

		app.Authors = []*cli.Author{
			{
				Name:  authorName,
				Email: authorEmail,
			},
		}

		app.HideVersion = false

		if buildTime == "" {
			return
		}

		parsedBuildTime, err := time.Parse(time.RFC3339, buildTime)
		if err != nil {
			panic(err)
		}

		app.Compiled = parsedBuildTime
	}
}

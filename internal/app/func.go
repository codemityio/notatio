package app

import (
	"fmt"
	"os"
	"os/exec"
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

func CheckFileExists(ctx *cli.Context, path string, message string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if _, ee := fmt.Fprintln(ctx.App.Writer, message); ee != nil {
				return fmt.Errorf("%w: %w", errWrite, ee)
			}

			return fmt.Errorf("%w: %w", errWrite, err)
		}

		if _, ee := fmt.Fprintln(ctx.App.Writer, err); ee != nil {
			return fmt.Errorf("%w: %w", errWrite, ee)
		}

		return fmt.Errorf("%w: %w", errWrite, err)
	}

	return nil
}

func CheckCommand(ctx *cli.Context, cmd string, message string) error {
	if _, e := exec.LookPath(cmd); e != nil {
		if _, ee := fmt.Fprintln(ctx.App.Writer, message); ee != nil {
			return fmt.Errorf("%w: %w", errWrite, ee)
		}

		return fmt.Errorf("%w: %w", errWrite, e)
	}

	return nil
}

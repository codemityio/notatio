package main

import (
	"log"
	"os"

	"github.com/codemity/notatio/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	application := app.New(
		app.WithValues(
			name,
			`A tool to support work with Markdown.`,
			version,
			copyright,
			authorName,
			authorEmail,
			buildTime,
		),
	)

	application.Commands = []*cli.Command{}

	if e := application.Run(os.Args); e != nil {
		log.Fatalf("error occurred during execution")
	}
}

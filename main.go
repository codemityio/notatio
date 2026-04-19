package main

import (
	"log"
	"os"

	"github.com/codemityio/notatio/cmd/coi"
	"github.com/codemityio/notatio/cmd/fs"
	"github.com/codemityio/notatio/cmd/graphviz"
	"github.com/codemityio/notatio/cmd/mermaid"
	"github.com/codemityio/notatio/cmd/plantuml"
	"github.com/codemityio/notatio/cmd/toc"
	"github.com/codemityio/notatio/cmd/tol"
	"github.com/codemityio/notatio/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	application := app.New(
		app.WithValues(
			name,
			`A tool designed to streamline working with documentation and diagrams.`,
			version,
			copyright,
			authorName,
			authorEmail,
			buildTime,
		),
	)

	application.Commands = []*cli.Command{
		&coi.App,
		&fs.App,
		&graphviz.App,
		&mermaid.App,
		&plantuml.App,
		&toc.App,
		&tol.App,
	}

	if e := application.Run(os.Args); e != nil {
		log.Fatalf("error: %v", e)
	}
}

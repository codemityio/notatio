package plantuml

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codemityio/notatio/internal/app"
	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) error {
	if e := app.CheckCommand(ctx, "java", "java not found"); e != nil {
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	inputPath := ctx.String("input-path")
	outputFormat := ctx.String("output-format")
	plantumlLimitSize := ctx.String("plantuml-limit-size")
	plantumlJarPath := ctx.String("plantuml-jar-path")
	recursive := ctx.Bool("recursive")

	if e := app.CheckFileExists(ctx, plantumlJarPath, plantumlJarPath+" not found"); e != nil {
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	if inputPath == "" {
		return errInputPathEmpty
	}

	if outputFormat != "png" && outputFormat != "svg" {
		return fmt.Errorf("%w: %s", errUnsupportedOutputFormat, outputFormat)
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("%w: %s: %w", errInputPath, inputPath, err)
	}

	if info.IsDir() {
		if _, e := fmt.Fprintln(ctx.App.Writer, inputPath+" - scanning..."); e != nil {
			return fmt.Errorf("%w: %w", errWrite, e)
		}

		if e := iterate(
			ctx,
			inputPath,
			outputFormat,
			plantumlJarPath,
			plantumlLimitSize,
			recursive,
		); e != nil {
			return e
		}

		return nil
	}

	if e := generate(ctx, inputPath, outputFormat, plantumlJarPath, plantumlLimitSize); e != nil {
		return e
	}

	return nil
}

func iterate(
	ctx *cli.Context,
	path, format, plantumlJarPath, plantumlLimitSize string,
	recurse bool,
) error {
	items, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("%w: failed to read directory: %w", errReadDir, err)
	}

	for _, item := range items {
		if !item.IsDir() && strings.HasSuffix(item.Name(), ".puml") {
			if e := generate(ctx, path, format, plantumlJarPath, plantumlLimitSize); e != nil {
				return e
			}
		}
	}

	if !recurse {
		return nil
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		if _, e := fmt.Fprintln(
			ctx.App.Writer,
			path+"/"+item.Name()+" - scanning...",
		); e != nil {
			return fmt.Errorf("%w: %w", errWrite, e)
		}

		if e := iterate(
			ctx,
			path+"/"+item.Name(),
			format,
			plantumlJarPath,
			plantumlLimitSize,
			recurse,
		); e != nil {
			return e
		}
	}

	return nil
}

func generate(ctx *cli.Context, path, format, plantumlJarPath, plantumlLimitSize string) error {
	a1 := "-DPLANTUML_LIMIT_SIZE=" + plantumlLimitSize
	a2 := "-jar"
	a3 := plantumlJarPath
	a4 := path
	a5 := "-t" + format

	// #nosec G204
	cmd := exec.CommandContext(ctx.Context, "java", a1, a2, a3, a4, a5)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if e := cmd.Run(); e != nil {
		return fmt.Errorf("%w: %w", errCommandRun, e)
	}

	return nil
}

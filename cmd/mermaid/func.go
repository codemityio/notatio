package mermaid

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codemityio/notatio/internal/app"
	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) error {
	if e := app.CheckCommand(ctx, "mmdc", "mmdc not found"); e != nil {
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	if e := app.CheckCommand(ctx, "chromium-browser", "chromium-browser not found"); e != nil {
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	inputPath := ctx.String("input-path")
	outputFormat := ctx.String("output-format")
	puppeteerConfigJSONPath := ctx.String("puppeteer-config-json-path")
	recursive := ctx.Bool("recursive")

	if e := app.CheckFileExists(
		ctx,
		puppeteerConfigJSONPath,
		puppeteerConfigJSONPath+" not found",
	); e != nil {
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
		if e := iterate(
			ctx,
			inputPath,
			outputFormat,
			puppeteerConfigJSONPath,
			recursive,
		); e != nil {
			return e
		}

		return nil
	}

	index := strings.LastIndex(inputPath, ".")

	if inputPath[index:] == ".mmd" {
		if e := generate(ctx, inputPath[:index], outputFormat, puppeteerConfigJSONPath); e != nil {
			return e
		}
	}

	return nil
}

func iterate(ctx *cli.Context, path, format, puppeteerConfigJSONPath string, recurse bool) error {
	items, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("%w: failed to read directory: %w", errReadDir, err)
	}

	for _, item := range items {
		if item.IsDir() && recurse {
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
				puppeteerConfigJSONPath,
				recurse,
			); e != nil {
				return e
			}
		}

		index := strings.LastIndex(item.Name(), ".")

		if index != -1 {
			if item.Name()[index:] == ".mmd" {
				if e := generate(
					ctx,
					path+"/"+item.Name()[:index],
					format,
					puppeteerConfigJSONPath,
				); e != nil {
					return e
				}
			}
		}
	}

	return nil
}

func generate(ctx *cli.Context, path, format, puppeteerConfigJSONPath string) error {
	relPath, err := filepath.Rel(".", path)
	if err != nil {
		relPath = path // fallback
	}

	outputPath := strings.TrimSuffix(relPath, filepath.Ext(relPath)) + "." + format
	outputDir := filepath.Dir(outputPath)

	// #nosec G301
	if e := os.MkdirAll(outputDir, permsDir); e != nil {
		return fmt.Errorf("%w: failed to create output directory: %w", errMkdir, e)
	}

	if _, e := fmt.Fprintf(ctx.App.Writer, "generating %s from %s\n", outputPath, path); e != nil {
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	a1 := "-p" + puppeteerConfigJSONPath
	a2 := "-e" + format
	a3 := "-i" + path + ".mmd"
	a4 := "-o" + outputPath

	// #nosec G204
	cmd := exec.CommandContext(ctx.Context, "mmdc", a1, a2, a3, a4)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if e := cmd.Run(); e != nil {
		return fmt.Errorf("%w: %w", errCommandRun, e)
	}

	return nil
}

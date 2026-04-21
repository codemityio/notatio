package graphviz

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codemityio/notatio/internal/app"
	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) error {
	if e := app.CheckCommand(ctx, "dot", "dot not found"); e != nil {
		return fmt.Errorf("%w: %w", errDep, e)
	}

	inputPath := ctx.String("input-path")
	outputFormat := ctx.String("output-format")
	recursive := ctx.Bool("recursive")

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

		if e := iterate(ctx, inputPath, outputFormat, recursive); e != nil {
			return e
		}

		return nil
	}

	if e := runDot(ctx, inputPath, outputFormat); e != nil {
		return e
	}

	return nil
}

func iterate(ctx *cli.Context, path, format string, recurse bool) error {
	items, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("%w: failed to read directory: %w", errReadDir, err)
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		name := item.Name()
		ext := strings.ToLower(filepath.Ext(name))

		if ext == ".dot" || ext == ".gv" {
			fullPath := filepath.Join(path, name)

			if e := runDot(ctx, fullPath, format); e != nil {
				return fmt.Errorf("%w: failed to process file %s: %w", errFileRead, fullPath, e)
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

		subPath := filepath.Join(path, item.Name())

		if _, e := fmt.Fprintln(ctx.App.Writer, subPath+" - scanning..."); e != nil {
			return fmt.Errorf("%w: %w", errWrite, e)
		}

		if e := iterate(ctx, subPath, format, recurse); e != nil {
			return e
		}
	}

	return nil
}

func runDot(ctx *cli.Context, path, format string) error {
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

	// #nosec G204
	cmd := exec.CommandContext(ctx.Context, "dot", "-T"+format, path, "-o", outputPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if e := cmd.Run(); e != nil {
		return fmt.Errorf("%w: %w", errCommandRun, e)
	}

	if format == "svg" {
		if e := normalizeSVG(outputPath, 0); e != nil {
			return e
		}
	}

	return nil
}

func roundFloatsInValue(val string, precision int) string {
	scale := math.Pow(powBase, float64(precision))

	return floatRegex.ReplaceAllStringFunc(val, func(match string) string {
		f, err := strconv.ParseFloat(match, 64)
		if err != nil {
			return match
		}

		return strconv.FormatFloat(math.Round(f*scale)/scale, 'f', precision, 64)
	})
}

func normalizeSVG(path string, precision int) error {
	data, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return fmt.Errorf("%w: failed to read svg: %w", errNromalise, err)
	}

	normalised := coordAttrRegex.ReplaceAllFunc(data, func(match []byte) []byte {
		// Extract the three capture groups: prefix, value, suffix
		groups := coordAttrRegex.FindSubmatch(match)
		if groups == nil {
			return match
		}

		return []byte(
			string(
				groups[1],
			) + roundFloatsInValue(
				string(groups[2]),
				precision,
			) + string(
				groups[3],
			),
		)
	})

	if e := os.WriteFile(path, normalised, permsFile); e != nil { // #nosec G306 G703
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	return nil
}

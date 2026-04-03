package tol

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
)

func before(ctx *cli.Context) error {
	csvPath = ctx.String("csv-path")
	documentPath = ctx.String("document-path")
	header = ctx.String("header")
	limiterL = ctx.String("limiter-left")
	limiterR = ctx.String("limiter-right")
	skip = ctx.StringSlice("skip")
	index = ctx.Int("index") // -1 or 0 means "replace all"

	var err error

	scsv, err = os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("%w: with csv path `%s`: %w", errFileOpen, csvPath, err)
	}

	body, err = os.ReadFile(documentPath)
	if err != nil {
		return fmt.Errorf("%w: with document `%s`: %w", errFileRead, documentPath, err)
	}

	if limiterR == "" {
		limiterR = "\\z"
	}

	rexp, err = regexp.Compile(fmt.Sprintf("%s+ %s[\\s\\S]*?%s+", limiterL, header, limiterR))
	if err != nil {
		return fmt.Errorf("unable to compile regex %w", err)
	}

	section := rexp.FindStringSubmatch(string(body))

	regxSuffix, err := regexp.Compile(limiterR)
	if err != nil {
		return fmt.Errorf("%w: suffix: %w", errRegexCompile, err)
	}

	if len(section) == 0 {
		return fmt.Errorf(
			"%w: unable to extract document section, please make sure the header `%s` and limiter `%s` are correct",
			errExtract,
			header,
			limiterR,
		)
	}

	prefix = regxPrefix.FindString(section[0])
	suffix = regxSuffix.FindString(section[len(section)-1])

	return nil
}

func readRecords(reader *csv.Reader) ([][]string, error) {
	records := [][]string{{"Package", "Licence", "Type"}}

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("%w: %w", errRead, err)
		}

		if len(record) > 0 && slices.Contains(skip, record[0]) {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func columnWidths(records [][]string) []int {
	widths := make([]int, columnWidth)

	for _, row := range records {
		for colIdx, cell := range row {
			if l := len(cell); l > widths[colIdx] {
				widths[colIdx] = l
			}
		}
	}

	return widths
}

func buildRow(row []string, widths []int) (string, error) {
	var b strings.Builder

	for colIdx, cell := range row {
		if _, e := fmt.Fprintf(&b, "| %-*s ", widths[colIdx], cell); e != nil {
			return "", fmt.Errorf("%w: %w", errWrite, e)
		}
	}

	return b.String() + "|\n", nil
}

func buildSeparator(widths []int) (string, error) {
	var b strings.Builder

	for _, width := range widths {
		if _, e := fmt.Fprintf(&b, "|-%s-", strings.Repeat("-", width)); e != nil {
			return "", fmt.Errorf("%w: %w", errWrite, e)
		}
	}

	return b.String() + "|\n", nil
}

func buildTable(records [][]string, widths []int) (string, error) {
	heading, err := buildRow(records[0], widths)
	if err != nil {
		return "", err
	}

	separator, err := buildSeparator(widths)
	if err != nil {
		return "", err
	}

	var rowsBuilder strings.Builder

	for _, row := range records[1:] {
		r, e := buildRow(row, widths)
		if e != nil {
			return "", e
		}

		rowsBuilder.WriteString(r)
	}

	return heading + separator + rowsBuilder.String() + "\n", nil
}

func action(_ *cli.Context) error {
	records, err := readRecords(csv.NewReader(scsv))
	if err != nil {
		return err
	}

	widths := columnWidths(records)

	table, err := buildTable(records, widths)
	if err != nil {
		return err
	}

	tol := fmt.Sprintf("%s %s\n\n", prefix, header) + table

	matchCount := 0

	replaced := rexp.ReplaceAllStringFunc(string(body), func(match string) string {
		matchCount++

		if index <= 0 || matchCount == index {
			return tol + suffix
		}

		return match
	})

	if err := os.WriteFile(documentPath, []byte(replaced), permsWrite); err != nil {
		return fmt.Errorf("%w: with a document path `%s`: %w", errFileWrite, documentPath, err)
	}

	return nil
}

package toc

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

func before(ctx *cli.Context) error {
	document = ctx.String("document")
	header = ctx.String("header")
	limiterL = ctx.String("limiter-left")
	limiterR = ctx.String("limiter-right")
	index = ctx.Int("index") // -1 or 0 means "replace all"

	var err error

	body, err = os.ReadFile(document)
	if err != nil {
		return fmt.Errorf("%w: %s: %w", errFileRead, document, err)
	}

	if limiterR == "" {
		limiterR = "\\z"
	}

	rexp, err = regexp.Compile(fmt.Sprintf("%s+ %s[\\s\\S]*?%s+", limiterL, header, limiterR))
	if err != nil {
		return fmt.Errorf("%w: %w", errRegexCompile, err)
	}

	section := rexp.FindStringSubmatch(string(body))

	regxSuffix, err := regexp.Compile(limiterR)
	if err != nil {
		return fmt.Errorf("%w: %w", errRegexCompile, err)
	}

	if len(section) == 0 {
		return fmt.Errorf(
			"%w: please make sure the header `%s` and limiter `%s` are correct",
			errDocumentSectionExtract,
			header,
			limiterR,
		)
	}

	prefix = regxPrefix.FindString(section[0])
	suffix = regxSuffix.FindString(section[len(section)-1])

	return nil
}

func internal(ctx *cli.Context) error {
	content, err := generateInternalTOC(
		ctx.Int("start-from-level"),
		ctx.Int("start-from-item"),
		prefix, header, document,
	)
	if err != nil {
		return err
	}

	matchCount := 0

	replaced := rexp.ReplaceAllStringFunc(string(body), func(match string) string {
		matchCount++

		if index <= 0 || matchCount == index {
			return content + suffix
		}

		return match
	})

	if e := os.WriteFile(document, []byte(replaced), permsWrite); e != nil {
		return fmt.Errorf("%w: `%s`: %w", errFileWrite, document, e)
	}

	return nil
}

func external(c *cli.Context) error {
	path := c.StringSlice("path")
	sh := c.String("summary-header")
	sll := c.String("summary-limiter-left")
	slr := c.String("summary-limiter-right")

	content, err := generateExternalTOC(prefix, header, path, sh, sll, slr)
	if err != nil {
		return err
	}

	matchCount := 0

	replaced := rexp.ReplaceAllStringFunc(string(body), func(match string) string {
		matchCount++

		if index <= 0 || matchCount == index {
			return content + suffix
		}

		return match
	})

	if e := os.WriteFile(document, []byte(replaced), permsWrite); e != nil {
		return fmt.Errorf("%w: `%s`: %w", errFileWrite, document, e)
	}

	return nil
}

func generateExternalTOC(
	prefix, header string,
	filePaths []string,
	summaryHeader, summaryLimiterL, summaryLimiterR string,
) (string, error) {
	toc := fmt.Sprintf("%s %s\n\n", prefix, header)

	list := make([]string, 0, len(filePaths))

	for _, path := range filePaths {
		var err error

		cont, err := os.ReadFile(path) // #nosec G304
		if err != nil {
			return "", fmt.Errorf("%w: `%s`: %w", errFileRead, path, err)
		}

		// Find the first match
		matches := regxTitle.FindStringSubmatch(string(cont))

		if len(matches) <= 1 {
			return "", errTitleNotFound
		}

		title := matches[1] // captured title without the #

		if summaryHeader == "" {
			list = append(list, fmt.Sprintf("- [%s](%s)\n", title, path))

			continue
		}

		// extract summary section
		if summaryLimiterR == "" {
			summaryLimiterR = "\\z"
		}

		// create a regular expression with capturing groups
		rex, err := regexp.Compile(fmt.Sprintf(
			"%s %s([\\s\\S]*?)%s",
			summaryLimiterL, summaryHeader, summaryLimiterR,
		))
		if err != nil {
			return "", fmt.Errorf("%w: %w", errRegexCompile, err)
		}

		var section string

		// find the section and extract the desired content
		matches = rex.FindStringSubmatch(string(cont))
		if len(matches) > 1 {
			section = matches[1] // the captured content
		}

		list = append(list, fmt.Sprintf(
			"- [%s](%s) - %s\n",
			title, path, strings.TrimSpace(strings.ReplaceAll(section, "\n", " "))),
		)
	}

	sort.Strings(list)

	var builder strings.Builder

	for _, line := range list {
		builder.WriteString(line)
	}

	toc += builder.String()

	// patch to write the right amount of new lines for documents without summary right limiter
	if summaryLimiterR != "" {
		toc += "\n"
	}

	return toc, nil
}

func generateInternalTOC(
	startFromLevel, startFromItem int,
	prefix, header, path string,
) (string, error) {
	toc := fmt.Sprintf("%s %s\n\n", prefix, header)

	inputFile, err := os.Open(path) // #nosec G304
	if err != nil {
		return "", fmt.Errorf("%w: `%s`: %w", errDocumentOpen, path, err)
	}

	defer func() {
		_ = inputFile.Close()
	}()

	// initialise variables
	inCodeBlock := false
	headerStack := make([]string, 0)

	// regular expression for identifying headers
	headerRegex := regexp.MustCompile(`^(#+)\s*(.*)`)

	scanner := bufio.NewScanner(inputFile)

	var (
		builder strings.Builder
		counter int
	)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
		}

		if !inCodeBlock && headerRegex.MatchString(line) {
			matches := headerRegex.FindStringSubmatch(line)
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])

			// calculate the depth in the hierarchy
			depth := level - 1

			if depth < startFromLevel {
				continue
			}

			if counter < startFromItem {
				counter++

				continue
			}

			// remove items from the stack that are at the same or deeper depth
			for len(headerStack) > depth {
				headerStack = headerStack[:len(headerStack)-1]
			}

			// add the current header to the stack
			headerStack = append(headerStack, title)

			// print the nested table of contents with indentation
			indentation := strings.Repeat("  ", depth)

			if _, e := fmt.Fprintf(
				&builder,
				"%s- [%s](#%s)\n",
				indentation,
				title,
				generateAnchor(title),
			); e != nil {
				return "", fmt.Errorf("%w: %w", errPrint, e)
			}

			counter++
		}
	}

	toc += builder.String()

	toc += "\n"

	return toc, nil
}

func generateAnchor(title string) string {
	anchor := strings.ToLower(title)
	anchor = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(anchor, "-")
	anchor = strings.Trim(anchor, "-")

	return anchor
}

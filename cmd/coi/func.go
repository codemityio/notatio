package coi

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

func action(ctx *cli.Context) error {
	coi := fmt.Sprintf("%s %s\n\n", prefix, header)

	shellName := ctx.String("shell-name")
	shellPromptPrefix := ctx.String("shell-prompt")
	command := ctx.String("command")
	output := ctx.String("output")

	switch {
	case command != "" && output == "":
		// #nosec G204
		cmd := exec.CommandContext(ctx.Context, "sh", "-c", command)

		var outBuffer bytes.Buffer

		cmd.Stdout = &outBuffer
		cmd.Stderr = &outBuffer

		if e := cmd.Run(); e != nil {
			return fmt.Errorf("%w: `%s`: %w", errCommandExecute, command, e)
		}

		output = stripANSI(outBuffer.String())
	case output != "":
		output = stripANSI(output)
	default:
		return fmt.Errorf(
			"%w: one of the following flags must be provided: --command, --output",
			errExclusiveFlags,
		)
	}

	matchCount := 0

	replaced := rexp.ReplaceAllStringFunc(string(body), func(match string) string {
		matchCount++

		if index <= 0 || matchCount == index {
			return fmt.Sprintf(
				"%s``` %s\n%s %s\n%s\n```\n\n",
				coi, shellName, shellPromptPrefix, command, strings.TrimSpace(output),
			) + suffix
		}

		return match // leave this occurrence unchanged
	})

	if e := os.WriteFile(document, []byte(replaced), permsWrite); e != nil {
		return fmt.Errorf("%w: `%s`: %w", errFileWrite, document, e)
	}

	return nil
}

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

package coi

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/urfave/cli/v2"
	"mvdan.cc/sh/v3/shell"
)

func before(c *cli.Context) error {
	document = c.String("document")
	header = c.String("header")
	limiterL = c.String("limiter-left")
	limiterR = c.String("limiter-right")

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

	if command != "" && output != "" {
		return fmt.Errorf(
			"%w: only one of the following flags is allowed at the same time: --command=%s, --output=%s",
			errExclusiveFlags,
			command,
			output,
		)
	}

	switch {
	case command != "":
		parts, err := shell.Fields(command, nil)
		if err != nil {
			return fmt.Errorf("%w: `%s` command: %w", errCommandParse, command, err)
		}

		name := parts[0]
		args := parts[1:]

		// #nosec G204
		cmd := exec.CommandContext(ctx.Context, name, args...)

		var outBuffer bytes.Buffer

		cmd.Stdout = &outBuffer
		cmd.Stderr = &outBuffer

		if e := cmd.Run(); e != nil {
			return fmt.Errorf("%w: `%s`: %w", errCommandExecute, command, e)
		}
	case output != "":
	default:
		return fmt.Errorf(
			"%w: one of the following flags must be provided: --command, --output",
			errExclusiveFlags,
		)
	}

	if e := os.WriteFile(
		document,
		[]byte(rexp.ReplaceAllString(string(body), fmt.Sprintf(
			"%s``` %s\n%s %s\n%s```\n\n",
			coi, shellName, shellPromptPrefix, command, output,
		)+suffix)),
		permsWrite,
	); e != nil {
		return fmt.Errorf("%w: `%s`: %w", errFileWrite, document, e)
	}

	return nil
}

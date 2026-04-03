//nolint:gochecknoglobals
package coi

import "regexp"

var (
	documentPath, header, limiterL, limiterR, prefix, suffix string
	index                                                    int
	body                                                     []byte
	rexp                                                     *regexp.Regexp

	regxPrefix = regexp.MustCompile(`^#*`)
	ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\[[?][0-9;]*[a-zA-Z]`)
)

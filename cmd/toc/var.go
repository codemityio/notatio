//nolint:gochecknoglobals
package toc

import "regexp"

var (
	documentPath, header, limiterL, limiterR, prefix, suffix string
	index                                                    int
	body                                                     []byte
	rexp                                                     *regexp.Regexp

	regxPrefix = regexp.MustCompile(`^#*`)
	regxTitle  = regexp.MustCompile(`(?m)^# (.+)$`)
)

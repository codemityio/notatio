//nolint:gochecknoglobals
package toc

import "regexp"

var (
	document, header, limiterL, limiterR, prefix, suffix string
	body                                                 []byte
	rexp                                                 *regexp.Regexp

	regxPrefix = regexp.MustCompile(`^#*`)
	regxTitle  = regexp.MustCompile(`(?m)^# (.+)$`)
)

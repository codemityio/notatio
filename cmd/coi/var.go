//nolint:gochecknoglobals
package coi

import "regexp"

var (
	document, header, limiterL, limiterR, prefix, suffix string
	index                                                int
	body                                                 []byte
	rexp                                                 *regexp.Regexp

	regxPrefix = regexp.MustCompile(`^#*`)
)

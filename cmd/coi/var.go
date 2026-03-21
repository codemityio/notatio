//nolint:gochecknoglobals
package coi

import "regexp"

var (
	document, header, limiterL, limiterR, prefix, suffix string
	body                                                 []byte
	rexp                                                 *regexp.Regexp

	regxPrefix = regexp.MustCompile(`^#*`)
)

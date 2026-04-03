//nolint:gochecknoglobals
package tol

import (
	"os"
	"regexp"
)

var (
	csvPath, documentPath, header, limiterL, limiterR, prefix, suffix string
	skip                                                              []string
	index                                                             int
	scsv                                                              *os.File
	body                                                              []byte
	rexp                                                              *regexp.Regexp

	regxPrefix = regexp.MustCompile(`^#*`)
)

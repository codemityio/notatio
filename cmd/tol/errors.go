package tol

import (
	"errors"
	"fmt"
)

var (
	errPkg          = errors.New("tol")
	errFileRead     = fmt.Errorf("%w: unable to read file", errPkg)
	errFileOpen     = fmt.Errorf("%w: unable to open file", errPkg)
	errFileWrite    = fmt.Errorf("%w: unable to write file", errPkg)
	errRead         = fmt.Errorf("%w: unable to read", errPkg)
	errRegexCompile = fmt.Errorf("%w: unable to compile regex", errPkg)
	errExtract      = fmt.Errorf("%w: unable to extract", errPkg)
	errWrite        = fmt.Errorf("%w: unable to write", errPkg)
)

package mermaid

import (
	"errors"
	"fmt"
)

var (
	errPkg                     = errors.New("mermaid")
	errDep                     = fmt.Errorf("%w: unable to find dependency", errPkg)
	errWrite                   = fmt.Errorf("%w: unable to write", errPkg)
	errInputPath               = fmt.Errorf("%w: input path", errPkg)
	errInputPathEmpty          = fmt.Errorf("%w: input path empty", errPkg)
	errUnsupportedOutputFormat = fmt.Errorf("%w: unsupported output format", errPkg)
	errReadDir                 = fmt.Errorf("%w: unable to read a directory", errPkg)
	errMkdir                   = fmt.Errorf("%w: unable to make a directory", errPkg)
	errCommandRun              = fmt.Errorf("%w: unable to run a command", errPkg)
)

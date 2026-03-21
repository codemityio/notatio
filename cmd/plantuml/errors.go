package plantuml

import (
	"errors"
	"fmt"
)

var (
	errPkg                     = errors.New("plantuml")
	errWrite                   = fmt.Errorf("%w: unable to write", errPkg)
	errInputPath               = fmt.Errorf("%w: input path", errPkg)
	errInputPathEmpty          = fmt.Errorf("%w: input path empty", errPkg)
	errUnsupportedOutputFormat = fmt.Errorf("%w: unsupported output format", errPkg)
	errReadDir                 = fmt.Errorf("%w: unable to read a directory", errPkg)
	errCommandRun              = fmt.Errorf("%w: unable to run a command", errPkg)
)

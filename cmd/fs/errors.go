package fs

import (
	"errors"
	"fmt"
)

var (
	errPkg              = errors.New("file")
	errOutputFormat     = fmt.Errorf("%w: output format", errPkg)
	errRegex            = fmt.Errorf("%w: regex", errPkg)
	errPath             = fmt.Errorf("%w: path", errPkg)
	errStdin            = fmt.Errorf("%w: stdin", errPkg)
	errRead             = fmt.Errorf("%w: unable to read", errPkg)
	errParse            = fmt.Errorf("%w: unable to parse", errPkg)
	errWrite            = fmt.Errorf("%w: unable to write", errPkg)
	errHeatMapField     = fmt.Errorf("%w: heat map field", errPkg)
	errExcludeDateField = fmt.Errorf("%w: exclude date field", errPkg)
	errExcludeDateValue = fmt.Errorf("%w: exclude date value", errPkg)
	errShutdown         = fmt.Errorf("%w: shutdown", errPkg)
)

package coi

import (
	"errors"
	"fmt"
)

var (
	errPkg                    = errors.New("toc")
	errCommandParse           = fmt.Errorf("%w: unable to parse a comand", errPkg)
	errCommandExecute         = fmt.Errorf("%w: unable to execute a comand", errPkg)
	errFileRead               = fmt.Errorf("%w: unable to read a file", errPkg)
	errFileWrite              = fmt.Errorf("%w: unable to write a file", errPkg)
	errRegexCompile           = fmt.Errorf("%w: unable to compile a regex", errPkg)
	errDocumentSectionExtract = fmt.Errorf("%w: unable to extract a document sction", errPkg)
)

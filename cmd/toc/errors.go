package toc

import (
	"errors"
	"fmt"
)

var (
	errPkg                    = errors.New("toc")
	errFileRead               = fmt.Errorf("%w: unable to read file", errPkg)
	errFileWrite              = fmt.Errorf("%w: unable to write file", errPkg)
	errRegexCompile           = fmt.Errorf("%w: unable to compile regex", errPkg)
	errDocumentOpen           = fmt.Errorf("%w: unable to open document", errPkg)
	errDocumentSectionExtract = fmt.Errorf("%w: unable to extract document sction", errPkg)
	errTitleNotFound          = fmt.Errorf("%w: unable to find title", errPkg)
	errPrint                  = fmt.Errorf("%w: unable to print", errPkg)
)

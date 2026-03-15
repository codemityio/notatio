package graph

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New(`graph`)
	// ErrFormatNotAllowed error.
	ErrFormatNotAllowed = fmt.Errorf("%w: format not allowed", ErrPkg)
)

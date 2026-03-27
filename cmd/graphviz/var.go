package graphviz

import "regexp"

// Only match floats inside coordinate attributes Graphviz uses.
var coordAttrRegex = regexp.MustCompile(
	`(\b(?:points|d|cx|cy|x|y|x1|y1|x2|y2)\s*=\s*")([^"]*)(")`,
)

var floatRegex = regexp.MustCompile(`-?\d+\.\d+`)

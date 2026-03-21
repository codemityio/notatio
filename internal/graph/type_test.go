package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatOutput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  FormatOutput
		errIs error
	}{
		{"valid png", "png", "png", nil},
		{"valid svg", "svg", "svg", nil},

		{"invalid empty", "", "", ErrFormatNotAllowed},
		{"invalid txt", "txt", "", ErrFormatNotAllowed},
		{"invalid case SVG", "SVG", "", ErrFormatNotAllowed},
		{"invalid mixed", "Svg", "", ErrFormatNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f FormatOutput

			require.ErrorIs(t, f.Set(tt.input), tt.errIs)

			assert.Equal(t, tt.want, f)
		})
	}
}

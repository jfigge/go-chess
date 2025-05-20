package colors

import (
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

func TestColors_Tint(t *testing.T) {
	tests := map[string]struct {
		base color.Color
		tint color.Color
		want color.Color
	}{
		"red tinted yellow": {
			base: &color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
			tint: &color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0x80},
			want: &color.RGBA{R: 0xff, G: 0x7e, B: 0x00, A: 0xff},
		},
		"blue tinted yellow": {
			base: &color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
			tint: &color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0x80},
			want: &color.RGBA{R: 0x7e, G: 0x7e, B: 0x80, A: 0xff},
		},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			actual := tintColor(test.base, test.tint)
			assert.Equal(tt, test.want, actual, "tintColor(%v, %v) = %v, want %v", test.base, test.tint, actual, test.want)
		})
	}
}

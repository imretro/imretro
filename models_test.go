package imretro

import "testing"

// TestModelConvertNoColor tests that a model without enough colors will return
// no color for a color that would be converted to an undefined color.
func TestModelConvertNoColor(t *testing.T) {
	model := ColorModel{Black, Black, Black}

	c := model.Convert(White)
	CompareColors(t, c, NoColor)
}

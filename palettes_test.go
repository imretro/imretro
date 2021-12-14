package imretro

import (
	"image/color"
	"testing"
)

// TestModelBitMode checks that the correct bit mode is interpreted from the
// color model.
func TestModelBitMode(t *testing.T) {
	if mode, err := ModelBitMode(Default1BitColorModel); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	} else if mode != OneBit {
		t.Errorf(
			`mode = %v (%08b), want %v (%08b)`,
			mode, mode,
			OneBit, OneBit,
		)
	}

	if mode, err := ModelBitMode(Default2BitColorModel); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	} else if mode != TwoBit {
		t.Errorf(
			`mode = %v (%08b), want %v (%08b)`,
			mode, mode,
			TwoBit, TwoBit,
		)
	}
}

// Test2BitModelBits checks that the correct bits (ranging [0b00, 0b11]) are
// returned by colors of varying brightness and opacity.
func Test2BitModelBits(t *testing.T) {
	model := TwoBitColorModel{}

	if bits := model.Bits(color.Alpha{0}); bits != 0 {
		t.Errorf(`bits = %02b, want 0`, bits)
	}
	if bits := model.Bits(color.RGBA{0xFF, 0xFF, 0xFF, 0x7F}); bits != 0 {
		t.Errorf(`bits = %02b, want 0`, bits)
	}

	if bits := model.Bits(DarkerGray); bits != 1 {
		t.Errorf(`bits = %02b, want 1`, bits)
	}
	if bits := model.Bits(MediumGray); bits != 2 {
		t.Errorf(`bits = %02b, want 2`, bits)
	}

	if bits := model.Bits(color.Gray{0xE0}); bits != 3 {
		t.Errorf(`bits = %02b, want 3`, bits)
	}
}

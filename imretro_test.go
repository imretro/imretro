package imretro

import (
	"image/color"
	"testing"
)

// TestIsBitCountSupported tests that true is returned when the bit count is
// supported, false when not supported.
func TestIsBitCountSupported(t *testing.T) {
	if v := IsBitCountSupported(0b0000_0001); v {
		t.Errorf(`IsBitCountSupported(0b0000_0001) = %v, want false`, v)
	}
	if v := IsBitCountSupported(0b1000_0000); !v {
		t.Errorf(`IsBitCountSupported(0b1000_0000) = %v, want true`, v)
	}
}

// TestUnsupportedError tests the error message for unsupported number of bits error.
func TestUnsupportedError(t *testing.T) {
	if actual, want := UnsupportedBitsError(0b10).Error(), "Unsupported bit count byte: 0b10"; actual != want {
		t.Fatalf(`err = %q, want %q`, actual, want)
	}
}

// TestColorAsBytes tests that a color would be converted to 4 bytes.
func TestColorAsBytes(t *testing.T) {
	white := color.Gray{255}
	gray := color.Gray{127}
	invisible := color.RGBA{0, 0, 0, 0}

	if r, g, b, a := ColorAsBytes(white); r != 255 || g != 255 || b != 255 || a != 255 {
		t.Fatalf(`r, g, b, a = %d %d %d %d, want 255, 255, 255, 255`, r, g, b, a)
	}
	if r, g, b, a := ColorAsBytes(gray); r != 127 || g != 127 || b != 127 || a != 255 {
		t.Fatalf(`r, g, b, a = %d %d %d %d, want 127, 127, 127, 255`, r, g, b, a)
	}
	if _, _, _, a := ColorAsBytes(invisible); a != 0 {
		t.Fatalf(`a = %d, want 0`, a)
	}
}

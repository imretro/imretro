package imretro

import "testing"

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

package imretro

import (
	"bytes"
	"testing"
)

// TestEncode1Bit checks that an image would be encoded to a 1-bit imretro
// file.
func TestEncode1Bit(t *testing.T) {
	var b bytes.Buffer

	Encode(&b, nil, OneBit)

	wantHeader := []byte{'I', 'M', 'R', 'E', 'T', 'R', 'O', 0b001_00000}

	for i, actual := range b.Next(len(wantHeader)) {

		if want := wantHeader[i]; actual != want {
			t.Errorf(
				`Header byte %d = %#b (%#x %c), want %#b (%#x %c)`,
				i,
				actual, actual, actual,
				want, want, want,
			)
		}
	}
}

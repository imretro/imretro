package imretro

import (
	"bytes"
	"image"
	"testing"
)

// TestEncode1BitHeader checks that image info would be encoded to a 1-bit imretro
// file.
func TestEncode1BitHeader(t *testing.T) {
	var b bytes.Buffer
	m := image.NewRGBA(image.Rect(0, 0, 320, 240))

	Encode(&b, m, OneBit)

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

	FailDimensionHelper(t, &b, "x", "Most", 1)
	FailDimensionHelper(t, &b, "x", "Least", 64)
	FailDimensionHelper(t, &b, "y", "Most", 0)
	FailDimensionHelper(t, &b, "y", "Least", 240)
}

// FailDimensionHelper fails if the dimension is not the wanted value.
func FailDimensionHelper(t *testing.T, b *bytes.Buffer, dimension, byteSignificance string, want byte) {
	t.Helper()
	actual, err := b.ReadByte()
	if err != nil {
		panic(err)
	}

	if actual != want {
		t.Errorf(
			`%s significant byte of %s dimension = %d (%b), want %d (%b)`,
			byteSignificance, dimension,
			actual, actual,
			want, want,
		)
	}
}

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
	Encode1Bit(t, &b)

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

// TestEncode1BitPalette checks that a black & white palette would be encoded
// to a 1-bit imretro file.
func TestEncode1BitPalette(t *testing.T) {
	var b bytes.Buffer
	Encode1Bit(t, &b)

	t.Log("Skipping to palette")
	b.Next(12)

	channels := []string{"r", "g", "b", "a"}
	for i, want := range []byte{0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} {
		t.Logf(`Checking %s channel of color %d`, channels[i%4], i/4)
		FailByteHelper(t, &b, want)
	}
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

// FailByteHelper fails if the next byte does not match the wanted value.
func FailByteHelper(t *testing.T, b *bytes.Buffer, want byte) {
	t.Helper()
	actual, err := b.ReadByte()
	if err != nil {
		panic(err)
	}

	if actual != want {
		t.Errorf(`byte = %d (%b), want %d (%b)`, actual, actual, want, want)
	}
}

// Encode1Bit creates a 1-bit image and encodes it to a buffer.
func Encode1Bit(t *testing.T, b *bytes.Buffer) {
	t.Helper()
	m := image.NewRGBA(image.Rect(0, 0, 320, 240))
	Encode(b, m, OneBit)
}

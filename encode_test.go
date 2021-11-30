package imretro

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

// TestEncode1BitHeader checks that image info would be encoded to a 1-bit imretro
// file.
func TestEncode1BitHeader(t *testing.T) {
	var b bytes.Buffer
	Encode1Bit(t, &b, 320, 240)

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
	Encode1Bit(t, &b, 320, 240)

	t.Log("Skipping to palette")
	b.Next(12)

	channels := []string{"r", "g", "b", "a"}
	for i, want := range []byte{0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} {
		t.Logf(`Checking %s channel of color %d`, channels[i%4], i/4)
		FailByteHelper(t, &b, want)
	}
}

// TestEncode1BitPixels checks that the pixels have been given the proper indices
// to the palette.
func TestEncode1BitPixels(t *testing.T) {
	var b bytes.Buffer
	Encode1Bit(t, &b, 10, 5)

	t.Log("Skipping to pixels")
	b.Next(12)
	b.Next(8)

	FailByteHelper(t, &b, 0b0000_1111)
	FailByteHelper(t, &b, 0b1010_0000)

	remaining := b.Bytes()

	// NOTE 50 pixels, 8 pixels per byte results in 6 complete bytes (48 pixels)
	// and 1 byte for the 2 remaining pixels. Subtract 2 for the bytes we just
	// tested above.
	if l, want := len(remaining), 5; l != want {
		t.Fatalf(
			`%d remaining pixel bytes (%d total pixel bytes), want %d`,
			l, l+2,
			want,
		)
	}

	t.Logf(`Remaining bytes: %v`, remaining)

	if final, want := remaining[len(remaining)-1], byte(0b0100_0000); final != want {
		t.Errorf(`final byte = %d (%08b), want %d (%08b)`, final, final, want, want)
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
			`%s significant byte of %s dimension = %d (%08b), want %d (%08b)`,
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
		t.Errorf(`byte = %d (%08b), want %d (%08b)`, actual, actual, want, want)
	}
}

// Encode1Bit creates a 1-bit image and encodes it to a buffer.
func Encode1Bit(t *testing.T, b *bytes.Buffer, width, height int) {
	t.Helper()

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < 8; i++ {
		var c color.Color = Black
		if i >= 4 {
			c = White
		}
		m.Set(i, i/width, c)
	}
	for i := 8; i < 16; i++ {
		var c color.Color = Black
		if i%2 == 0 && i < 12 {
			c = White
		}
		m.Set(i%width, i/width, c)
	}
	m.Set(width-2, height-1, DarkerGray)
	m.Set(width-1, height-1, LighterGray)

	Encode(b, m, OneBit)
}

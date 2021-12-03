package imretro

import (
	"bytes"
	"image/color"
	"io"
	"testing"
)

// TestPassCheckHeader tests that a reader starting with "IMRETRO" bytes will
// pass.
func TestPassCheckHeader(t *testing.T) {
	buff := make([]byte, 8)
	r := Make1Bit(t, 0b1010_0000, nil, 0, 0, nil)
	mode, err := checkHeader(r, buff)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if pixelMode := mode & (0b1100_0000); pixelMode != EightBit {
		t.Errorf(
			`pixelMode = %d (%08b), want %d (%08b)`,
			pixelMode, pixelMode,
			EightBit, EightBit,
		)
	}
	if hasPalette := mode & (0b0010_0000); hasPalette != 0x20 {
		t.Error("mode does not signify in-file palette")
	}
}

// TestFailCheckHeader tests that a reader with unexpected magic bytes will
// fail.
func TestFailCheckHeader(t *testing.T) {
	buff := make([]byte, 8)
	partialSignature := "IMRET"
	jpgSignature := "\xFF\xD8\xFF\xE0\x00\x10\x4A\x46\x49\x46\x00\x01"

	partialr := bytes.NewBufferString(partialSignature)
	if _, err := checkHeader(partialr, buff); err != io.ErrUnexpectedEOF {
		t.Errorf(`err = %v, want %v`, err, io.ErrUnexpectedEOF)
	}

	jpgr := bytes.NewBufferString(jpgSignature)
	if _, err := checkHeader(jpgr, buff); err != DecodeError("unexpected signature byte") {
		t.Fatalf(`err = %v, want %v`, err, DecodeError("unexpected signature byte"))
	}
}

// TestDecode1BitNoPalette tests that a 1-bit-mode image with no palette can be decoded.
func TestDecode1BitNoPalette(t *testing.T) {
	const width, height int = 320, 240
	var pixels = make([]byte, width*height)
	r := Make1Bit(t, 0x00, [][]byte{}, uint16(320), uint16(240), pixels)

	config, err := DecodeConfig(r)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if config.Width != width {
		t.Errorf(`Width = %v, want %v`, config.Width, width)
	}
	if config.Height != height {
		t.Errorf(`Height = %v, want %v`, config.Height, height)
	}
	if _, ok := config.ColorModel.(OneBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want OneBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{{DarkGray, Black}, {LightGray, White}}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode1BitPalette tests that a 1-bit palette would be properly decoded.
func TestDecode1BitPalette(t *testing.T) {
	palette := [][]byte{
		{0x00, 0xFF, 0x00, 0xFF},
		{0xEF, 0xFF, 0x00, 0xFF},
	}
	r := Make1Bit(t, 0x20, palette, 2, 2, make([]byte, 1))

	config, err := DecodeConfig(r)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if _, ok := config.ColorModel.(OneBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want OneBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{
		{Black, color.RGBA{0x00, 0xFF, 0x00, 0xFF}},
		{White, color.RGBA{0xEF, 0xFF, 0x00, 0xFF}},
	}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// Make1Bit makes a 1-bit imretro reader.
func Make1Bit(t *testing.T, mode byte, palette [][]byte, width, height uint16, pixels []byte) *bytes.Buffer {
	t.Helper()
	b := bytes.NewBuffer([]byte{
		// signature/magic bytes
		'I', 'M', 'R', 'E', 'T', 'R', 'O',
		// Mode byte (8-bit, in-file palette)
		mode,
		byte(width >> 8), byte(width & 0xFF),
		byte(height >> 8), byte(height & 0xFF),
	})
	for _, color := range palette {
		b.Write(color)
	}
	b.Write(pixels)
	return b
}

// CompareColors helps compare colors to each other.
func CompareColors(t *testing.T, actual, want color.Color) {
	t.Helper()
	r, g, b, a := actual.RGBA()
	wr, wg, wb, wa := want.RGBA()
	comparisons := [4]channelComparison{
		{"red", r, wr},
		{"green", g, wg},
		{"blue", b, wb},
		{"alpha", a, wa},
	}

	for _, comparison := range comparisons {
		if comparison.actual != comparison.want {
			t.Errorf(
				`%s channel = %v, want %v`,
				comparison.name, comparison.actual,
				comparison.want,
			)
		}
	}
}

// ChannelComparison is used to compare color channels.
type channelComparison struct {
	name         string
	actual, want uint32
}

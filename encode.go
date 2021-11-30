package imretro

import (
	"errors"
	"image"
	"image/color"
	"io"
)

// Encode writes the image m to w in imretro format.
func Encode(w io.Writer, m image.Image, bits byte) error {
	w.Write([]byte("IMRETRO"))
	w.Write([]byte{bits | WithPalette})

	bounds := m.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	for _, d := range []int{width, height} {
		if d > MaximumDimension {
			return DimensionsTooLargeError(d)
		}
		byte1 := IntLast8(d >> 8)
		byte2 := IntLast8(d)
		w.Write([]byte{byte1, byte2})
	}

	switch bits {
	case OneBit:
		return encodeOneBit(w, m)
	case TwoBit:
		return encodeTwoBit(w, m)
	case EightBit:
		return encodeEightBit(w, m)
	}
	return UnsupportedBitsError(bits)
}

func encodeOneBit(w io.Writer, m image.Image) error {
	// NOTE Write the palette
	if err := writeColor(w, Black); err != nil {
		return err
	}
	if err := writeColor(w, White); err != nil {
		return err
	}
	return errors.New("Not implemented")
}

func encodeTwoBit(w io.Writer, m image.Image) error {
	return errors.New("Not implemented")
}

func encodeEightBit(w io.Writer, m image.Image) error {
	return errors.New("Not implemented")
}

// WriteColor writes a color as 4 bytes to a Writer.
func writeColor(w io.Writer, c color.Color) error {
	r, g, b, a := ColorAsBytes(c)
	_, err := w.Write([]byte{r, g, b, a})
	return err
}

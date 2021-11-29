package imretro

import (
	"errors"
	"image"
	"io"
)

// Encode writes the image m to w in imretro format.
func Encode(w io.Writer, m image.Image, bits byte) error {
	w.Write([]byte("IMRETRO"))
	w.Write([]byte{bits | WithPalette})

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
	return errors.New("Not implemented")
}

func encodeTwoBit(w io.Writer, m image.Image) error {
	return errors.New("Not implemented")
}

func encodeEightBit(w io.Writer, m image.Image) error {
	return errors.New("Not implemented")
}

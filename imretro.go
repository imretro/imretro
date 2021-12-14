// Encode and decode retro-style images in the imretro format.
package imretro

import (
	"fmt"
	"image/color"
)

const (
	OneBit byte = iota << 6
	TwoBit
	EightBit
)

// PaletteIndex is the "index" (from the left) of the bit in the mode byte that
// signifies if there is an in-file palette.
const PaletteIndex byte = 2

// WithPalette can be used with a union with the bit count when setting the
// header.
const WithPalette byte = 1 << (7 - PaletteIndex)

// MaximumDimension is the maximum size of an image's boundary in the imretro
// format.
const MaximumDimension int = 0xFF_FF

// UnsupportedBitModeError should be returned when an unexpected number
// of bits is received.
type UnsupportedBitModeError byte

// DimensionsTooLargeError should be returned when an encoded image would
// have boundaries that are not valid in the encoding.
type DimensionsTooLargeError int

// IsBitCountSupported checks if the bit count is supported by the imretro
// format.
func IsBitCountSupported(count byte) bool {
	for _, bits := range []byte{OneBit, TwoBit, EightBit} {
		if count == bits {
			return true
		}
	}
	return false
}

// Error converts to an error string.
func (e UnsupportedBitModeError) Error() string {
	return fmt.Sprintf("Unsupported bit count byte: %#b", byte(e))
}

// Error makes a string representation of the too-large error.
func (e DimensionsTooLargeError) Error() string {
	return fmt.Sprintf("Dimensions too large for 16-bit number: %d", int(e))
}

// ColorAsBytes converts a color to a 4-byte (one byte for each channel)
// representation.
func ColorAsBytes(c color.Color) (r, g, b, a byte) {
	rchan, gchan, bchan, achan := c.RGBA()
	return ChannelAsByte(rchan), ChannelAsByte(gchan), ChannelAsByte(bchan), ChannelAsByte(achan)
}

// ColorFromBytes converts 4 bytes into a color. Panics if the slice has less
// than 4 bytes.
func ColorFromBytes(bs []byte) color.Color {
	return color.RGBA{bs[0], bs[1], bs[2], bs[3]}
}

// ChannelAsByte converts a uint32 color channel ranging within [0, 0xFFFF] to
// a byte.
func ChannelAsByte(channel uint32) byte {
	return byte(channel >> 8)
}

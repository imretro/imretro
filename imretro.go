// Encode and decode retro-style images in the imretro format.
package imretro

import (
	"fmt"
	"image"
	"image/color"

	"github.com/spenserblack/go-byteutils"
)

// PixelMode is the type for managing the number of bits per pixel.
type PixelMode = byte

const (
	OneBit PixelMode = iota << 6
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
func IsBitCountSupported(count PixelMode) bool {
	for _, bits := range []PixelMode{OneBit, TwoBit, EightBit} {
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

// ImretroImage is an image decoded from the imretro format.
type ImretroImage interface {
	image.Image
	// ColorIndex converts the x/y coordinates of a pixel to the index in the
	// palette.
	ColorIndex(x, y int) int
	// Palette gets the palette of the image.
	Palette() Palette
}

// ImretroImage is the helper struct for imretro images.
type imretroImage struct {
	config image.Config
	pixels []byte
}

// ColorModel returns the Image's color model.
func (i imretroImage) ColorModel() color.Model {
	return i.config.ColorModel
}

// Bounds returns the boundaries of the image.
func (i imretroImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.config.Width, i.config.Height)
}

// Image1Bit is the 1-bit image type.
type image1Bit struct {
	imretroImage
}

// Image2Bit is the 2-bit image type.
type image2Bit struct {
	imretroImage
}

// Image8Bit is the 8-bit image type.
type image8Bit struct {
	imretroImage
}

// ColorIndex converts the x/y coordinates of a pixel to the index in the
// palette.
func (i *image1Bit) ColorIndex(x, y int) int {
	index := (y * i.config.Width) + x
	byteIndex := index / 8
	bitIndex := byte(index % 8)
	b := i.pixels[byteIndex]
	bit := byteutils.GetL(b, bitIndex)
	return int(bit)
}

// Palette gets the 1-bit image palette.
func (i *image1Bit) Palette() Palette {
	return i.ColorModel().(OneBitColorModel).colors
}

// At returns the color at the given pixel.
func (i *image1Bit) At(x, y int) color.Color {
	if !image.Pt(x, y).In(i.Bounds()) {
		return NoColor
	}
	palette := i.Palette()
	return palette[i.ColorIndex(x, y)]
}

// ColorIndex converts the x/y coordinates of a pixel to the index in the
// palette.
func (i *image2Bit) ColorIndex(x, y int) int {
	index := (y * i.config.Width) + x
	byteIndex := index / 4
	bitIndex := byte(index%4) * 2
	bits := byteutils.SliceL(i.pixels[byteIndex], bitIndex, bitIndex+2)
	return int(bits)
}

// Palette gets the 2-bit image palette.
func (i *image2Bit) Palette() Palette {
	return i.ColorModel().(TwoBitColorModel).colors
}

// At returns the color at the given pixel.
func (i *image2Bit) At(x, y int) color.Color {
	if !image.Pt(x, y).In(i.Bounds()) {
		return NoColor
	}
	palette := i.Palette()
	return palette[i.ColorIndex(x, y)]
}

// ColorIndex converts the x/y coordinates of a pixel to the index in the
// palette.
func (i *image8Bit) ColorIndex(x, y int) int {
	index := (y * i.config.Width) + x
	pixel := i.pixels[index]
	return int(pixel)
}

// Palette gets the 8-bit image palette.
func (i *image8Bit) Palette() Palette {
	return i.ColorModel().(EightBitColorModel).colors
}

// At returns the color at the given pixel.
func (i *image8Bit) At(x, y int) color.Color {
	if !image.Pt(x, y).In(i.Bounds()) {
		return NoColor
	}
	palette := i.Palette()
	return palette[i.ColorIndex(x, y)]
}

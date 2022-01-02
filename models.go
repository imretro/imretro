package imretro

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/spenserblack/go-byteutils"
)

// ModelMap maps bit modes to color models.
type ModelMap = map[PixelMode]color.Model

// ErrUnknownModel is raised when an unknown color model is interpreted.
var ErrUnknownModel = errors.New("Color model not recognized")

// MissingModelError is raised when there is no model for the given pixel bit
// mode.
type MissingModelError PixelMode

// Error reports the pixel mode lacking the color model.
func (mode MissingModelError) Error() string {
	return fmt.Sprintf("No model for pixel mode %02b", mode)
}

var (
	Default1BitColorModel = NewOneBitColorModel(Black, White)
	Default2BitColorModel = NewTwoBitColorModel(Black, DarkGray, LightGray, White)
	Default8BitColorModel = NewEightBitColorModel(Default8BitPalette)
)

// DefaultModelMap maps bit modes to the default color models.
var DefaultModelMap = ModelMap{
	OneBit:   Default1BitColorModel,
	TwoBit:   Default2BitColorModel,
	EightBit: Default8BitColorModel,
}

// OneBitColorModel is color model for 1-bit-pixel images.
type OneBitColorModel struct {
	colors Palette
}

// TwoBitColorModel is a color model for 2-bit-pixel images.
type TwoBitColorModel struct {
	colors Palette
}

// EightBitColorModel is a color model for 8-bit-pixel images.
type EightBitColorModel struct {
	colors Palette
}

// ModelBitMode gets the bits-per-pixel according to the color model.
func ModelBitMode(model color.Model) (byte, error) {
	switch model.(type) {
	case OneBitColorModel:
		return OneBit, nil
	case TwoBitColorModel:
		return TwoBit, nil
	case EightBitColorModel:
		return EightBit, nil
	}
	return 0, ErrUnknownModel
}

// NewOneBitColorModel creates a new color model for 1-bit-pixel images.
func NewOneBitColorModel(off color.Color, on color.Color) OneBitColorModel {
	return OneBitColorModel{Palette{off, on}}
}

func (model OneBitColorModel) Convert(c color.Color) color.Color {
	return model.colors[int(model.Bit(c))]
}

// Bit gets the bit that should point to the color index.
func (model OneBitColorModel) Bit(c color.Color) byte {
	r, g, b, a := ColorAsBytes(c)
	brightness := r | g | b
	isBright := (brightness >= 128) && (a >= 128)
	if isBright {
		return 1
	}
	return 0
}

// NewTwoBitColorModel creates a new color model for 2-bit-pixel images.
func NewTwoBitColorModel(off, light, strong, full color.Color) TwoBitColorModel {
	return TwoBitColorModel{Palette{off, light, strong, full}}
}

func (model TwoBitColorModel) Convert(c color.Color) color.Color {
	return model.colors[int(model.Bits(c))]
}

// Bits gets the two bits that should point to the color index.
//
// Possible values are in range [0, 4).
func (model TwoBitColorModel) Bits(c color.Color) byte {
	r, g, b, a := ColorAsBytes(c)
	// NOTE Return the "off" color if <50% opacity
	if a < 0x80 {
		return 0
	}
	// NOTE Two most significant bits of the combined colors.
	return (r | g | b) >> 6
}

// NewEightBitColorModel creates a new color model for 8-bit-pixel images.
func NewEightBitColorModel(colors Palette) EightBitColorModel {
	return EightBitColorModel{colors}
}

func (model EightBitColorModel) Convert(c color.Color) color.Color {
	index := int(model.Bits(c))
	if index >= len(model.colors) {
		return NoColor
	}
	return model.colors[index]
}

// Bits gets the eight bits that should point to the color index.
//
// Possible values are in range [0, 256).
func (model EightBitColorModel) Bits(c color.Color) byte {
	r, g, b, a := ColorAsBytes(c)
	r = byteutils.SliceL(r, 0, 2)
	g = byteutils.SliceL(g, 0, 2) << 2
	b = byteutils.SliceL(b, 0, 2) << 4
	a = byteutils.SliceL(a, 0, 2) << 6
	return r | g | b | a
}

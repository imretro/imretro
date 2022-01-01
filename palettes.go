package imretro

import (
	"errors"
	"image/color"

	"github.com/spenserblack/go-byteutils"
)

// Palette is a palette of colors.
type Palette = color.Palette

// ErrUnknownModel is raised when an unknown color model is interpreted.
var ErrUnknownModel = errors.New("Color model not recognized")

var (
	Default1BitPalette = Palette{Black, White}
	Default2BitPalette = Palette{Black, DarkGray, LightGray, White}
	// Default8BitPalette has 256 possible colors, and is defined on
	// initialization.
	Default8BitPalette = make8BitPalette()
)

// DefaultPaletteMap maps bit modes to the appropriate default palettes.
var DefaultPaletteMap = map[PixelMode]Palette{
	OneBit:   Default1BitPalette,
	TwoBit:   Default2BitPalette,
	EightBit: Default8BitPalette,
}

var (
	// NoColor is "invisible" and signifies a lack of color.
	NoColor color.Color = color.Alpha{0}
	Black   color.Color = color.Gray{0}
	// DarkerGray is 25% light.
	DarkerGray = color.Gray{0x40}
	// DarkGray is 33% light, and can be used for splitting a monochromatic
	// color range into 4 parts (0, 33%, 66%, 100%).
	DarkGray = color.Gray{0x55}
	// MediumGray is the exact middle between black and white.
	MediumGray = color.Gray{0x80}
	// LightGray is 66% light, and can be used for splitting a monochromatic
	// color range into 4 parts (0, 33%, 66%, 100%).
	LightGray = color.Gray{0xAA}
	// LighterGray is 75% light.
	LighterGray = color.Gray{0xC0}
	White       = color.Gray{0xFF}
)

var (
	Default1BitColorModel = NewOneBitColorModel(Black, White)
	Default2BitColorModel = NewTwoBitColorModel(Black, DarkGray, LightGray, White)
	Default8BitColorModel = NewEightBitColorModel(Default8BitPalette)
)

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

// Make8BitPalette creates the default 8-bit palette as described in the format
// documentation.
func make8BitPalette() Palette {
	palette := make(Palette, 0, 256)
	for i := 0; i < cap(palette); i++ {
		rgba := make([]byte, 4)
		for ci := range rgba {
			channelIndex := byte(ci)
			channel := byteutils.SliceR(byte(i), channelIndex*2, (channelIndex*2)+2)
			channel |= (channel << 6) | (channel << 4) | (channel << 2)
			rgba[ci] = channel
		}
		c := ColorFromBytes(rgba)
		palette = append(palette, c)
	}
	return palette
}

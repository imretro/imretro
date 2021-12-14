package imretro

import (
	"errors"
	"image/color"
)

// Palette is a palette of colors.
type Palette = color.Palette

// ErrUnknownModel is raised when an unknown color model is interpreted.
var ErrUnknownModel = errors.New("Color model not recognized")

var (
	Default1BitPalette = Palette{Black, White}
	Default2BitPalette = Palette{Black, DarkGray, LightGray, White}
)

// DefaultPaletteMap maps bit modes to the appropriate default palettes.
var DefaultPaletteMap = map[byte]Palette{
	OneBit: Default1BitPalette,
	TwoBit: Default2BitPalette,
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
)

// OneBitColorModel is color model for 1-bit-pixel images.
type OneBitColorModel struct {
	colors Palette
}

// ModelBitMode gets the bits-per-pixel according to the color model.
func ModelBitMode(model color.Model) (byte, error) {
	switch model.(type) {
	case OneBitColorModel:
		return OneBit, nil
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

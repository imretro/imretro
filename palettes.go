package imretro

import (
	"errors"
	"image/color"
)

// Palette is a palette of colors.
type Palette = color.Palette

// PaletteMap maps a pixel bit-count to a palette.
type PaletteMap = map[byte]Palette

// ErrUnknownModel is raised when an unknown color model is interpreted.
var ErrUnknownModel = errors.New("Color model not recognized")

var (
	// DecodingPalettes maps a byte for the pixel bit-count
	// to an appropriate default palette to be faithful to
	// a "retro" style.
	DecodingPalettes PaletteMap
	// EncodingPalettes maps a byte for the pixel bit-count
	// to an appropriate default palette to keep the most
	// color accuracy possible when encoding an image.
	EncodingPalettes PaletteMap
)

var (
	// NoColor is "invisible" and signifies a lack of color.
	NoColor color.Color = color.Alpha{0}
	Black   color.Color = color.Gray{0}
	// DarkerGray is 25% light.
	DarkerGray = color.Gray{64}
	// DarkGray is 33% light, and can be used for splitting a monochromatic
	// color range into 4 parts (0, 33%, 66%, 100%).
	DarkGray = color.Gray{85}
	// MediumGray is the exact middle between black and white.
	MediumGray = color.Gray{128}
	// LightGray is 66% light, and can be used for splitting a monochromatic
	// color range into 4 parts (0, 33%, 66%, 100%).
	LightGray = color.Gray{170}
	// LighterGray is 75% light.
	LighterGray = color.Gray{192}
	White       = color.Gray{255}
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

func init() {
	DecodingPalettes = make(PaletteMap)
	EncodingPalettes = make(PaletteMap)

	DecodingPalettes[OneBit] = Palette{Black, White}
	EncodingPalettes[OneBit] = Palette{Black, White}
}

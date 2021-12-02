package imretro

import "image/color"

// Palette is a palette of colors.
type Palette = color.Palette

// PaletteMap maps a pixel bit-count to a palette.
type PaletteMap = map[byte]Palette

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
	Black color.Color = color.Gray{0}
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

func init() {
	DecodingPalettes = make(PaletteMap)
	EncodingPalettes = make(PaletteMap)

	DecodingPalettes[OneBit] = Palette{Black, White}
	EncodingPalettes[OneBit] = Palette{Black, White}
}

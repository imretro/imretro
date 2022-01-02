package imretro

import (
	"image/color"

	"github.com/spenserblack/go-byteutils"
)

// Palette is a palette of colors.
type Palette = color.Palette

// PaletteMap maps pixel modes to palettes.
type PaletteMap = map[PixelMode]Palette

var (
	Default1BitPalette = Palette{Black, White}
	Default2BitPalette = Palette{Black, DarkGray, LightGray, White}
	// Default8BitPalette has 256 possible colors, and is defined on
	// initialization.
	Default8BitPalette = make8BitPalette()
)

// DefaultPaletteMap maps bit modes to the appropriate default palettes.
var DefaultPaletteMap = PaletteMap{
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

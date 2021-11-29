// Encode and decode retro-style images in the imretro format.
package imretro

import "fmt"

const (
	OneBit byte = iota << 6
	TwoBit
	EightBit
)

// WithPalette can be used with a union with the bit count when setting the
// header.
const WithPalette byte = 1 << 5

// UnsupportedBitsError should be returned when an unexpected number
// of bits is received.
type UnsupportedBitsError int

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
func (e UnsupportedBitsError) Error() string {
	return fmt.Sprintf("Unsupported number of bits: %d", int(e))
}

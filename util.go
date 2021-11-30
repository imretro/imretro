package imretro

// IntLast8 gets the last 8 bits of an int and converts to a byte.
func IntLast8(i int) byte {
	return byte(i & 0xFF)
}

// IndexByte gets the bit (as a byte) at the given index.
// index is an int to keep consistency, even though the range of values is
// extremely limited.
//
// Index 0 is the most significant bit.
//
// Panics if index is not in [0, 8).
func IndexByte(b byte, index int) (bit byte) {
	if index < 0 || index >= 8 {
		panic("index out of bounds")
	}
	return (b >> (7 - index)) & 1
}

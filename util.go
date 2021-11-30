package imretro

// IntLast8 gets the last 8 bits of an int and converts to a byte.
func IntLast8(i int) byte {
	return byte(i & 0xFF)
}

package imretro

import "testing"

func SanityCheck(t *testing.T) {
	if add(2, 2) != 4 {
		t.Fatal()
	}
}
